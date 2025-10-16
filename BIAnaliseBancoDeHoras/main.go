package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bh-mvc/internal/config"
	"bh-mvc/internal/model"
	"bh-mvc/internal/repository"
	"bh-mvc/internal/service"

	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.Load()
	Agenda()
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Run ETL end-to-end
	mux.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		ctx := r.Context()

		if cfg.ClearDest {
			log.Println("[RUN] Limpando destino...")
			_, _, err := repository.DeleteAll(ctx, cfg.DestConn, cfg.TblDetalhado, cfg.TblResumo)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			_, _, errs := repository.DeleteAllBhMonth(ctx, cfg.DestConn, cfg.TblDetalhado, cfg.TblResumo)
			if errs != nil {
				http.Error(w, errs.Error(), 500)
				return
			}
		}

		log.Println("[RUN] Limpando destino...")
		_, _, errd := repository.DeleteAll(ctx, cfg.DestConn, cfg.TblDetalhado, cfg.TblResumo)
		if errd != nil {
			http.Error(w, errd.Error(), 500)
			return
		}

		// Carrega origem
		var rows []model.RawRow
		var err error
		if cfg.UseSQLOrigin {
			log.Println("[RUN] Lendo ORIGEM (SQL)...")
			rows, err = repository.LoadFromSQL(ctx, cfg.SrcConn)
		} else {
			log.Println("[RUN] Lendo ORIGEM (CSV)...")
			rows, err = repository.LoadFromCSV(cfg.InputCSV, cfg.CSVSep)
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("[RUN] Linhas lidas: %d\n", len(rows))

		start := time.Now()
		processed := service.ApplyBankOffset(rows)
		log.Printf("[RUN] Processado em %s. Linhas: %d\n", time.Since(start), len(processed))

		// Detalhado
		if strings.TrimSpace(cfg.TblDetalhado) != "" {
			if _, err := repository.InsertDetalhado(ctx, cfg.DestConn, cfg.TblDetalhado, processed); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		// Resumo
		summary := service.BuildMonthlySummary(processed)
		if cfg.UpsertResumo && strings.TrimSpace(cfg.TblResumo) != "" {
			if _, err := repository.UpsertResumo(ctx, cfg.DestConn, cfg.TblResumo, summary); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		// calculo mensal
		var rowss []model.RawRow
		if cfg.UseSQLOrigin {
			log.Println("[RUN] Lendo ORIGEM (SQL)...")
			rowss, err = repository.LoadFromSQLMensal(ctx, cfg.SrcConn)
		} else {
			log.Println("[RUN] Lendo ORIGEM (CSV)...")
			rowss, err = repository.LoadFromCSV(cfg.InputCSV, cfg.CSVSep)
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Printf("[RUN] Linhas lidas: %d\n", len(rowss))

		starts := time.Now()
		processeds := service.ApplyBankOffset(rowss)
		log.Printf("[RUN] Processado em %s. Linhas: %d\n", time.Since(starts), len(processeds))

		// Detalhado
		if strings.TrimSpace(cfg.TblDetalhado) != "" {
			if _, err := repository.InsertDetalhado(ctx, cfg.DestConn, cfg.TblDetalhado, processeds); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		// Resumo
		service.BuildMonthlySummary(processeds)
		if cfg.UpsertResumo && strings.TrimSpace(cfg.TblResumo) != "" {
			if _, err := repository.UpsertResumo(ctx, cfg.DestConn, cfg.TblResumo, summary); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "processed": len(processeds), "summary": len(summary)})
	})

	// Query summary (JSON)
	mux.HandleFunc("/summary", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		perref := normalizePerref(r.URL.Query().Get("perref")) // YYYYMM or empty
		numemp := parseInt(r.URL.Query().Get("numemp"))
		codccu := strings.TrimSpace(r.URL.Query().Get("codccu"))
		colab := strings.TrimSpace(r.URL.Query().Get("colab"))

		rows, err := querySummary(ctx, cfg, perref, numemp, codccu, colab)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rows)
	})

	// UI
	mux.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl := template.Must(template.New("ui").Parse(uiHTML))
		_ = tmpl.Execute(w, nil)
	})

	log.Printf("HTTP ouvindo em %s\n", cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, mux); err != nil {
		log.Fatal(err)
	}
}

func normalizePerref(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if len(s) == 7 && s[4] == '-' { // YYYY-MM -> YYYYMM
		return s[:4] + s[5:7]
	}
	return s
}

func parseInt(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

type SummaryOut struct {
	NumEmp                    int     `json:"numemp"`
	NumCad                    int     `json:"numcad"`
	PerRef                    string  `json:"perref"`
	Colaborador               string  `json:"colaborador"`
	CodCcu                    string  `json:"codccu"`
	HorasPositivasOriginal    float64 `json:"horas_positivas_original"`
	Banco230ConsumidoMes      float64 `json:"banco_230_consumido_no_mes"`
	HorasSaldoMes             float64 `json:"horas_saldo_mes"`
	ValorSaldoMes             float64 `json:"valor_saldo_mes"`
	BancoTotalAplicadoNoGrupo float64 `json:"banco_total_aplicado_no_grupo"`
}

// querySummary busca do **Resumo** juntando nome e codccu do detalhado mais recente do mês (para exibir filtros por CCU/Colaborador)
func querySummary(ctx context.Context, cfg config.Config, perref string, numemp *int, codccu, colab string) ([]SummaryOut, error) {
	db, err := sql.Open("sqlserver", cfg.DestConn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var where []string
	if perref != "" {
		where = append(where, "r.perref = @p1")
	}
	if numemp != nil {
		where = append(where, "r.numemp = @p1")
	}
	// codccu/colab pelo detalhado (pega max dtapuracao por chave dentro do mês)
	if codccu != "" {
		where = append(where, "d.codccu LIKE @p1 + '%' ")
	}
	if colab != "" {
		where = append(where, "d.colaborador LIKE '%' + @p1 + '%' ")
	}

	w := ""
	if len(where) > 0 {
		w = "WHERE " + strings.Join(where, " AND ")
	}

	q := fmt.Sprintf(`
	SELECT r.numemp, r.numcad, r.perref,
	       d.colaborador, d.codccu,
	       r.horas_positivas_original, r.banco_230_consumido_no_mes,
	       r.horas_saldo_mes, r.valor_saldo_mes, r.banco_total_aplicado_no_grupo
	FROM %s r
	OUTER APPLY (
	  SELECT TOP 1 codccu, colaborador
	  FROM %s d
	  WHERE d.numemp = r.numemp AND d.numcad = r.numcad AND d.perref = r.perref
	  ORDER BY d.dtapuracao DESC
	) d
	%s
	ORDER BY r.perref DESC, r.numemp, r.numcad
	`, cfg.TblResumo, cfg.TblDetalhado, w)

	args := []any{}
	if perref != "" {
		args = append(args, perref)
	}
	if numemp != nil {
		args = append(args, *numemp)
	}
	if codccu != "" {
		args = append(args, codccu)
	}
	if colab != "" {
		args = append(args, colab)
	}

	rows, err := db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SummaryOut
	for rows.Next() {
		var o SummaryOut
		if err := rows.Scan(
			&o.NumEmp, &o.NumCad, &o.PerRef, &o.Colaborador, &o.CodCcu,
			&o.HorasPositivasOriginal, &o.Banco230ConsumidoMes,
			&o.HorasSaldoMes, &o.ValorSaldoMes, &o.BancoTotalAplicadoNoGrupo,
		); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

var uiHTML = `<!doctype html> 
<html lang="pt-br">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>BH Resumo Mensal</title>
<style>
body { font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif; margin: 24px; }
.card { background: #fff; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,.08); padding: 16px; margin-bottom: 16px; }
h1 { margin: 0 0 12px; }
label { display:block; font-size: 12px; color:#555; margin-bottom:4px; }
input { padding:8px; border:1px solid #ddd; border-radius:8px; width:100%; }
.grid { display:grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
button { padding:10px 14px; border:none; border-radius:8px; background:#111827; color:#fff; cursor:pointer; }
table { width:100%; border-collapse: collapse; }
th, td { padding: 8px 10px; border-bottom:1px solid #eee; text-align:left; font-size: 14px; }
th { background:#f8fafc; }
.small { color:#64748b; font-size:12px; }
</style>
</head>
<body>
  <div class="card">
    <h1>Resumo Mensal - Banco de Horas</h1>
    <div class="grid">
      <div>
        <label>Mês/Ano (YYYYMM ou YYYY-MM)</label>
        <input id="perref" placeholder="202509" />
      </div>
      <div>
        <label>Empresa (numemp)</label>
        <input id="numemp" placeholder="1" />
      </div>
      <div>
        <label>Centro de Custo (prefixo)</label>
        <input id="codccu" placeholder="CC" />
      </div>
      <div>
        <label>Colaborador (nome contém)</label>
        <input id="colab" placeholder="Ana" />
      </div>
    </div>
    <div style="margin-top:12px; display:flex; gap:8px;">
      <button onclick="buscar()">Buscar</button>
      <button onclick="rodar()" style="background:#0f766e">Processar /run</button>
    </div>
    <p class="small">Dica: clique em “Processar /run” para atualizar detalhado e resumo antes de consultar.</p>
  </div>

  <div class="card">
    <table id="t">
      <thead>
        <tr>
          <th>Período</th><th>Empresa</th><th>Crachá/Núm. Cad.</th><th>Colaborador</th><th>CCU</th>
          <th>Horas +</th><th>Banco 230 consumido</th><th>Horas saldo</th><th>R$ saldo</th><th>Banco Total (grupo)</th>
        </tr>
      </thead>
      <tbody></tbody>
    </table>
  </div>

<script>
async function buscar(){
  const perref = document.getElementById('perref').value.trim();
  const numemp = document.getElementById('numemp').value.trim();
  const codccu = document.getElementById('codccu').value.trim();
  const colab  = document.getElementById('colab').value.trim();
  const p = new URLSearchParams();
  if(perref) p.append('perref', perref);
  if(numemp) p.append('numemp', numemp);
  if(codccu) p.append('codccu', codccu);
  if(colab)  p.append('colab', colab);

  const res = await fetch('/summary?'+p.toString());
  if(!res.ok){
    alert('Erro ao buscar dados');
    return;
  }
  const data = await res.json();
  const tb = document.querySelector('#t tbody');
  tb.innerHTML = '';
  for(const r of data){
    const tr = document.createElement('tr');
    tr.innerHTML =
      '<td>' + (r.perref || '') + '</td>' +
      '<td>' + (r.numemp || '') + '</td>' +
      '<td>' + (r.numcad || '') + '</td>' +
      '<td>' + (r.colaborador || '') + '</td>' +
      '<td>' + (r.codccu || '') + '</td>' +
      '<td>' + fmt(r.horas_positivas_original) + '</td>' +
      '<td>' + fmt(r.banco_230_consumido_no_mes) + '</td>' +
      '<td>' + fmt(r.horas_saldo_mes) + '</td>' +
      '<td>' + fmt(r.valor_saldo_mes) + '</td>' +
      '<td>' + fmt(r.banco_total_aplicado_no_grupo) + '</td>';
    tb.appendChild(tr);
  }
}

function fmt(v){
  const n = Number(v);
  const safe = Number.isFinite(n) ? n : 0;
  return safe.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

async function rodar(){
  const res = await fetch('/run', {method:'POST'});
  if(!res.ok){
    alert('Falha no processamento');
    return;
  }
  const j = await res.json();
  alert('Processado: linhas=' + (j.processed || 0) + ' | resumo=' + (j.summary || 0));
  buscar();
}
</script>
</body></html>`

// dispara um POST para /run no próprio servidor
func triggerRunOnce(ctx context.Context, baseURL string) error {
	reqBody, _ := json.Marshal(map[string]any{}) // body vazio
	fmt.Println("vai chamar o agendador de tarefas...")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://192.168.1.28:8093/run", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("rodando pelo agendador...")
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("POST /run retornou status %d", resp.StatusCode)
	}
	return nil
}

func Agenda() {
	// Carrega cfg só para saber a porta; não altera o main existente
	cfg := config.Load()

	// Monta a baseURL do próprio servidor (ex.: http://127.0.0.1:8080)
	// cfg.HTTPPort costuma vir no formato ":8080"
	baseURL := "http://127.0.0.1" + cfg.HTTPPort

	// Usa timezone de São Paulo
	var loc *time.Location

	loc = time.FixedZone("BRT", -3*60*60)

	fmt.Println(loc.String())
	log.Writer().Write([]byte(loc.String()))

	c := cron.New(
		cron.WithLocation(loc), // agenda no fuso correto
		cron.WithSeconds(),     // permite campo de segundos no spec
		cron.WithChain( // logs básicos de erro
			cron.Recover(cron.DefaultLogger),
		),
	)

	// “0 0 3 * * *” => todos os dias às 03:00:00
	_, err := c.AddFunc("0 0 3 * * *", func() {
		log.Println("[SCHED] Disparando /run (03:00)…")
		// contexto com timeout generoso para o ETL
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := triggerRunOnce(ctx, baseURL); err != nil {
			log.Printf("[SCHED] Erro ao chamar /run: %v", err)
			return
		}
		log.Println("[SCHED] /run concluído com sucesso.")
	})
	if err != nil {
		log.Printf("[SCHED] Erro ao registrar cron: %v", err)
		return
	}

	// Inicia o cron em background. Não bloqueia o main existente.
	c.Start()

	log.Println("[SCHED] Agendador ativo: todo dia às 03:00 (America/Sao_Paulo)")
}
