package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
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

const authCookieName = "bh_session"

func main() {
	cfg := config.Load()
	Agenda()
	mux := http.NewServeMux()

	// token de sessão gerado a cada start do processo: reiniciar o serviço
	// derruba todas as sessões logadas. Não depende de armazenar a senha em lugar nenhum.
	authToken := randomToken()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			_ = r.ParseForm()
			pw := r.FormValue("password")
			if cfg.UIPassword != "" && subtle.ConstantTimeCompare([]byte(pw), []byte(cfg.UIPassword)) == 1 {
				http.SetCookie(w, &http.Cookie{
					Name:     authCookieName,
					Value:    authToken,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					MaxAge:   12 * 3600,
				})
				http.Redirect(w, r, "/ui", http.StatusFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(loginPage(true)))
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(loginPage(false)))
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     authCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	// Run ETL end-to-end
	mux.HandleFunc("/run", requireAuth(cfg, authToken, func(w http.ResponseWriter, r *http.Request) {
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
	}))

	// Query summary (JSON)
	mux.HandleFunc("/summary", requireAuth(cfg, authToken, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		perrefs := parsePerrefs(r.URL.Query()["perref"]) // um ou mais YYYYMM/YYYY-MM; vazio = todos os períodos
		numemp := parseInt(r.URL.Query().Get("numemp"))
		codccu := strings.TrimSpace(r.URL.Query().Get("codccu"))
		colab := strings.TrimSpace(r.URL.Query().Get("colab"))

		rows, err := querySummary(ctx, cfg, perrefs, numemp, codccu, colab)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rows)
	}))

	// Lista os períodos (perref) já processados, para alimentar o seletor da UI
	mux.HandleFunc("/periods", requireAuth(cfg, authToken, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		periods, err := listPeriods(ctx, cfg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(periods)
	}))

	// Detalhe (linhas do detalhado) de um colaborador em um período, para drill-down na UI
	mux.HandleFunc("/detail", requireAuth(cfg, authToken, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		numemp := parseInt(r.URL.Query().Get("numemp"))
		numcad := parseInt(r.URL.Query().Get("numcad"))
		perref := strings.TrimSpace(r.URL.Query().Get("perref"))
		if numemp == nil || numcad == nil || perref == "" {
			http.Error(w, "parâmetros numemp, numcad e perref são obrigatórios", http.StatusBadRequest)
			return
		}
		rows, err := queryDetail(ctx, cfg, *numemp, *numcad, perref)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rows)
	}))

	// UI
	mux.HandleFunc("/ui", requireAuth(cfg, authToken, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl := template.Must(template.New("ui").Parse(uiHTML))
		_ = tmpl.Execute(w, nil)
	}))

	log.Printf("HTTP ouvindo em %s\n", cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, mux); err != nil {
		log.Fatal(err)
	}
}

func randomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("falha ao gerar token de sessão: ", err)
	}
	return hex.EncodeToString(b)
}

// requireAuth protege uma rota com a sessão criada em /login. Se UIPassword
// estiver vazio no config, a proteção fica desativada (comportamento anterior).
func requireAuth(cfg config.Config, authToken string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.UIPassword == "" {
			next(w, r)
			return
		}
		c, err := r.Cookie(authCookieName)
		if err != nil || subtle.ConstantTimeCompare([]byte(c.Value), []byte(authToken)) != 1 {
			if r.Header.Get("Accept") != "" && strings.Contains(r.Header.Get("Accept"), "application/json") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if r.URL.Path == "/ui" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func loginPage(wrongPassword bool) string {
	errBlock := ""
	if wrongPassword {
		errBlock = `<div class="err">Senha incorreta.</div>`
	}
	return `<!doctype html>
<html lang="pt-br">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>Login — BH Resumo Mensal</title>
<style>
body { font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif; background:#f1f5f9; margin:0; display:flex; align-items:center; justify-content:center; min-height:100vh; }
.card { background:#fff; border-radius:12px; box-shadow:0 2px 12px rgba(0,0,0,.08); padding:32px; width:320px; }
h1 { font-size:18px; margin:0 0 16px; }
label { display:block; font-size:12px; color:#555; margin-bottom:4px; }
input { padding:10px; border:1px solid #ddd; border-radius:8px; width:100%; font-size:14px; box-sizing:border-box; }
button { margin-top:16px; padding:10px 14px; border:none; border-radius:8px; background:#111827; color:#fff; cursor:pointer; width:100%; font-size:14px; }
.err { color:#b91c1c; font-size:13px; margin-bottom:12px; }
</style>
</head>
<body>
  <form class="card" method="post" action="/login">
    <h1>Resumo Mensal — Banco de Horas</h1>
    ` + errBlock + `
    <label>Senha</label>
    <input type="password" name="password" autofocus required />
    <button type="submit">Entrar</button>
  </form>
</body></html>`
}

// perrefToDate converte "YYYYMM" ou "YYYY-MM" para "YYYY-MM-01", formato
// em que a coluna perref (tipo date) é gravada pelo /run. Não altera nenhuma
// regra de cálculo — é usado apenas para consultas de leitura (/summary, /detail).
func perrefToDate(s string) (string, bool) {
	s = strings.ReplaceAll(strings.TrimSpace(s), "-", "")
	if len(s) != 6 {
		return "", false
	}
	if _, err := strconv.Atoi(s); err != nil {
		return "", false
	}
	return s[:4] + "-" + s[4:6] + "-01", true
}

// parsePerrefs limpa e deduplica a lista de períodos vinda da query string (?perref=A&perref=B...).
func parsePerrefs(raw []string) []string {
	var out []string
	seen := map[string]bool{}
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p != "" && !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
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

// querySummary busca do **Resumo** juntando nome e codccu do detalhado mais recente do mês (para exibir filtros por CCU/Colaborador).
// perrefs pode conter um ou mais períodos (YYYYMM/YYYY-MM); vazio = sem filtro de período (todos).
func querySummary(ctx context.Context, cfg config.Config, perrefs []string, numemp *int, codccu, colab string) ([]SummaryOut, error) {
	db, err := sql.Open("sqlserver", cfg.DestConn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var where []string
	var args []any
	// cada filtro usa seu próprio @pN — no código anterior todos reusavam @p1,
	// o que fazia o SQL Server ignorar os filtros extras quando combinados.
	param := func() string {
		return fmt.Sprintf("@p%d", len(args)+1)
	}

	if len(perrefs) > 0 {
		var placeholders []string
		for _, p := range perrefs {
			if d, ok := perrefToDate(p); ok {
				placeholders = append(placeholders, param())
				args = append(args, d)
			}
		}
		if len(placeholders) > 0 {
			where = append(where, "r.perref IN ("+strings.Join(placeholders, ",")+")")
		}
	}
	if numemp != nil {
		where = append(where, "r.numemp = "+param())
		args = append(args, *numemp)
	}
	// codccu/colab pelo detalhado (pega max dtapuracao por chave dentro do mês)
	if codccu != "" {
		where = append(where, "d.codccu LIKE "+param()+" + '%' ")
		args = append(args, codccu)
	}
	if colab != "" {
		where = append(where, "d.colaborador LIKE '%' + "+param()+" + '%' ")
		args = append(args, colab)
	}

	w := ""
	if len(where) > 0 {
		w = "WHERE " + strings.Join(where, " AND ")
	}

	// ISNULL evita erro de scan quando não há linha correspondente no detalhado (OUTER APPLY sem match).
	// CONVERT(varchar(6), perref, 112) devolve o período como "YYYYMM", já que a coluna é do tipo date.
	q := fmt.Sprintf(`
	SELECT r.numemp, r.numcad, CONVERT(varchar(6), r.perref, 112) AS perref,
	       ISNULL(d.colaborador, '') AS colaborador, ISNULL(d.codccu, '') AS codccu,
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

// listPeriods devolve os períodos (YYYYMM) já existentes no resumo, para alimentar o seletor da UI.
func listPeriods(ctx context.Context, cfg config.Config) ([]string, error) {
	db, err := sql.Open("sqlserver", cfg.DestConn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q := fmt.Sprintf(`SELECT DISTINCT CONVERT(varchar(6), perref, 112) AS perref FROM %s ORDER BY perref DESC`, cfg.TblResumo)
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []string{}
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

type DetailOut struct {
	DtApuracao        string  `json:"dtapuracao"`
	CodCcu            string  `json:"codccu"`
	DesSit            string  `json:"dessit"`
	CodSit            int     `json:"codsit"`
	HorasOriginal     float64 `json:"horas_original"`
	BancoUsadoNaLinha float64 `json:"banco_usado_na_linha"`
	HorasSaldo        float64 `json:"horas_saldo"`
	ValorSaldo        float64 `json:"valor_saldo"`
	ValHoraCalculado  float64 `json:"valhoracalculado"`
}

// queryDetail lista as linhas do detalhado (já processadas pelo /run) de um colaborador em um período,
// para permitir drill-down na UI a partir de uma linha do resumo.
func queryDetail(ctx context.Context, cfg config.Config, numemp, numcad int, perref string) ([]DetailOut, error) {
	d, ok := perrefToDate(perref)
	if !ok {
		return nil, fmt.Errorf("perref inválido: %q", perref)
	}

	db, err := sql.Open("sqlserver", cfg.DestConn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q := fmt.Sprintf(`
	SELECT CONVERT(varchar(10), dtapuracao, 23) AS dtapuracao, ISNULL(codccu,'') AS codccu,
	       ISNULL(dessit,'') AS dessit, codsit,
	       horas_original, banco_usado_na_linha, horas_saldo, valor_saldo, valhoracalculado
	FROM %s
	WHERE numemp = @p1 AND numcad = @p2 AND perref = @p3
	ORDER BY dtapuracao ASC, codsit ASC
	`, cfg.TblDetalhado)

	rows, err := db.QueryContext(ctx, q, numemp, numcad, d)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []DetailOut{}
	for rows.Next() {
		var o DetailOut
		if err := rows.Scan(
			&o.DtApuracao, &o.CodCcu, &o.DesSit, &o.CodSit,
			&o.HorasOriginal, &o.BancoUsadoNaLinha, &o.HorasSaldo, &o.ValorSaldo, &o.ValHoraCalculado,
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
:root { color-scheme: light; }
* { box-sizing: border-box; }
body { font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif; margin: 0; background:#f1f5f9; color:#0f172a; }
.wrap { max-width: 1280px; margin: 0 auto; padding: 24px; }
.card { background: #fff; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,.06); padding: 16px; margin-bottom: 16px; }
h1 { margin: 0 0 4px; font-size: 22px; }
.subtitle { color:#64748b; font-size:13px; margin: 0 0 16px; }
label { display:block; font-size: 12px; color:#555; margin-bottom:4px; }
input, select { padding:8px; border:1px solid #ddd; border-radius:8px; width:100%; font-size:14px; }
.grid { display:grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
button { padding:10px 14px; border:none; border-radius:8px; background:#111827; color:#fff; cursor:pointer; font-size:14px; }
button:disabled { opacity:.6; cursor:default; }
button.secondary { background:#0f766e; }
button.ghost { background:#fff; color:#334155; border:1px solid #cbd5e1; }
.actions { margin-top:12px; display:flex; gap:8px; flex-wrap:wrap; align-items:center; }
.small { color:#64748b; font-size:12px; }
.kpis { display:grid; grid-template-columns: repeat(5, minmax(0,1fr)); gap:12px; margin-bottom:16px; }
.kpi { background:#fff; border-radius:12px; box-shadow: 0 2px 12px rgba(0,0,0,.06); padding:14px 16px; }
.kpi .lbl { font-size:12px; color:#64748b; margin-bottom:4px; }
.kpi .val { font-size:19px; font-weight:600; }
.kpi .val.neg { color:#b91c1c; }
.kpi .val.pos { color:#0f766e; }
.table-scroll { width:100%; overflow-x:auto; }
table { width:100%; border-collapse: collapse; white-space:nowrap; }
th, td { padding: 8px 10px; border-bottom:1px solid #eee; text-align:left; font-size: 13px; }
th { background:#f8fafc; cursor:pointer; user-select:none; position:sticky; top:0; }
th .arrow { font-size:10px; color:#94a3b8; margin-left:4px; }
tbody tr:hover { background:#f8fafc; cursor:pointer; }
tbody tr:nth-child(even) { background:#fcfdfe; }
td.num { text-align:right; font-variant-numeric: tabular-nums; }
td.neg { color:#b91c1c; }
.detail-row td { background:#f8fafc; padding:0; }
.detail-box { padding:12px 16px; }
.detail-box table { background:#fff; border-radius:8px; overflow:hidden; }
.banner { padding:10px 14px; border-radius:8px; margin-bottom:12px; font-size:13px; }
.banner.error { background:#fef2f2; color:#991b1b; border:1px solid #fecaca; }
.banner.info { background:#eff6ff; color:#1e40af; border:1px solid #bfdbfe; }
.empty { padding:24px; text-align:center; color:#94a3b8; }
.spinner { display:inline-block; width:14px; height:14px; border:2px solid rgba(255,255,255,.4); border-top-color:#fff; border-radius:50%; animation: spin .7s linear infinite; vertical-align:-2px; margin-right:6px; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
</head>
<body>
<div class="wrap">
  <div class="card">
    <h1>Resumo Mensal — Banco de Horas</h1>
    <p class="subtitle">Consulta e análise do resumo gerado pelo processamento (/run). A regra de cálculo não é alterada por esta tela.</p>
    <div class="grid">
      <div>
        <label>Período(s) <span class="small">(Ctrl/Cmd+clique para vários; nenhum = todos)</span></label>
        <select id="perref" multiple size="5"></select>
        <div style="margin-top:6px; display:flex; gap:6px;">
          <button type="button" class="ghost" style="padding:4px 8px; font-size:12px;" onclick="selectAllPeriods()">Selecionar todos</button>
          <button type="button" class="ghost" style="padding:4px 8px; font-size:12px;" onclick="clearPeriods()">Limpar</button>
        </div>
      </div>
      <div>
        <label>Crachá/Núm. Cad. (numemp)</label>
        <input id="numemp" placeholder="361" inputmode="numeric" />
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
    <div class="actions">
      <button id="btnBuscar" onclick="buscar()">Buscar</button>
      <button id="btnRun" class="secondary" onclick="rodar()">Processar /run</button>
      <button class="ghost" onclick="exportCsv()">Exportar CSV</button>
      <span id="status" class="small"></span>
      <a href="/logout" class="ghost" style="padding:10px 14px; border-radius:8px; border:1px solid #cbd5e1; color:#334155; text-decoration:none; font-size:14px; margin-left:auto;">Sair</a>
    </div>
  </div>

  <div id="banner"></div>

  <p class="small" id="kpisLabel" style="margin: 0 0 6px;"></p>
  <div class="kpis" id="kpis"></div>

  <div class="card">
    <div class="table-scroll">
      <table id="t">
        <thead>
          <tr>
            <th data-k="perref">Período<span class="arrow"></span></th>
            <th data-k="numcad">Empresa<span class="arrow"></span></th>
            <th data-k="numemp">Crachá/Núm. Cad.<span class="arrow"></span></th>
            <th data-k="colaborador">Colaborador<span class="arrow"></span></th>
            <th data-k="codccu">CCU<span class="arrow"></span></th>
            <th data-k="horas_positivas_original" class="num">Horas +<span class="arrow"></span></th>
            <th data-k="banco_230_consumido_no_mes" class="num">Banco 230 consumido<span class="arrow"></span></th>
            <th data-k="horas_saldo_mes" class="num">Horas saldo<span class="arrow"></span></th>
            <th data-k="valor_saldo_mes" class="num">R$ saldo<span class="arrow"></span></th>
            <th data-k="banco_total_aplicado_no_grupo" class="num">Banco Total (grupo)<span class="arrow"></span></th>
          </tr>
        </thead>
        <tbody></tbody>
      </table>
      <div id="empty" class="empty" hidden>Nenhum registro encontrado para os filtros informados.</div>
    </div>
  </div>
</div>

<script>
let data = [];
let sortKey = 'perref';
let sortDir = 'desc';
let openDetailKey = null;

function fmt(v){
  const n = Number(v);
  const safe = Number.isFinite(n) ? n : 0;
  return safe.toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}
function fmtPerref(p){
  if(!p || p.length !== 6) return p || '';
  return p.slice(0,4) + '-' + p.slice(4,6);
}
function showBanner(msg, kind){
  const b = document.getElementById('banner');
  if(!msg){ b.innerHTML = ''; return; }
  b.innerHTML = '<div class="banner ' + (kind||'info') + '">' + msg + '</div>';
}
function setBusy(busy, btn, label){
  btn.disabled = busy;
  btn.innerHTML = busy ? '<span class="spinner"></span>' + label : label;
}

async function loadPeriods(){
  try{
    const res = await fetch('/periods');
    if(!res.ok) return;
    const periods = await res.json();
    const sel = document.getElementById('perref');
    sel.innerHTML = '';
    for(const p of (periods || [])){
      const opt = document.createElement('option');
      opt.value = p;
      opt.textContent = fmtPerref(p);
      sel.appendChild(opt);
    }
    // por padrão seleciona só o período mais recente; o usuário pode marcar mais.
    if(sel.options.length) sel.options[0].selected = true;
  }catch(e){ /* seletor fica vazio (equivale a "todos") */ }
}

function selectAllPeriods(){
  const sel = document.getElementById('perref');
  for(const o of sel.options) o.selected = true;
}
function clearPeriods(){
  const sel = document.getElementById('perref');
  for(const o of sel.options) o.selected = false;
}
function selectedPeriods(){
  return Array.from(document.getElementById('perref').selectedOptions).map(o => o.value);
}

async function buscar(){
  const btn = document.getElementById('btnBuscar');
  setBusy(true, btn, 'Buscar');
  showBanner('', null);
  const perrefs = selectedPeriods();
  const numemp = document.getElementById('numemp').value.trim();
  const codccu = document.getElementById('codccu').value.trim();
  const colab  = document.getElementById('colab').value.trim();
  const p = new URLSearchParams();
  for(const pr of perrefs) p.append('perref', pr);
  if(numemp) p.append('numemp', numemp);
  if(codccu) p.append('codccu', codccu);
  if(colab)  p.append('colab', colab);

  try{
    const res = await fetch('/summary?'+p.toString());
    if(!res.ok){
      showBanner('Erro ao buscar dados: ' + (await res.text()), 'error');
      data = [];
    } else {
      data = await res.json() || [];
    }
  }catch(e){
    showBanner('Falha de rede ao buscar dados: ' + e, 'error');
    data = [];
  }
  openDetailKey = null;
  const periodCount = new Set(data.map(r => r.perref)).size;
  document.getElementById('status').textContent =
    data.length + ' registro(s)' + (periodCount > 1 ? ' em ' + periodCount + ' períodos (total geral somado)' : '');
  renderKpis();
  render();
  setBusy(false, btn, 'Buscar');
}

function renderKpis(){
  const el = document.getElementById('kpis');
  const periods = Array.from(new Set(data.map(r => r.perref))).sort();
  const lbl = document.getElementById('kpisLabel');
  lbl.textContent = periods.length > 1
    ? 'Total geral somando ' + periods.length + ' períodos: ' + periods.map(fmtPerref).join(', ')
    : (periods.length === 1 ? 'Período: ' + fmtPerref(periods[0]) : '');
  const colabs = new Set(data.map(r => r.numemp + '|' + r.numcad));
  const sum = (k) => data.reduce((a,r) => a + (Number(r[k]) || 0), 0);
  const horasSaldo = sum('horas_saldo_mes');
  const valorSaldo = sum('valor_saldo_mes');
  const cards = [
    ['Colaboradores', colabs.size.toLocaleString('pt-BR'), ''],
    ['Horas positivas', fmt(sum('horas_positivas_original')), ''],
    ['Banco 230 consumido', fmt(sum('banco_230_consumido_no_mes')), ''],
    ['Horas saldo', fmt(horasSaldo), horasSaldo < 0 ? 'neg' : 'pos'],
    ['R$ saldo', fmt(valorSaldo), valorSaldo < 0 ? 'neg' : 'pos'],
  ];
  el.innerHTML = cards.map(c =>
    '<div class="kpi"><div class="lbl">' + c[0] + '</div><div class="val ' + c[2] + '">' + c[1] + '</div></div>'
  ).join('');
}

function sortedData(){
  const arr = data.slice();
  arr.sort((a,b) => {
    let va = a[sortKey], vb = b[sortKey];
    if(typeof va === 'string'){ va = (va||'').toLowerCase(); vb = (vb||'').toLowerCase(); }
    else { va = Number(va)||0; vb = Number(vb)||0; }
    if(va < vb) return sortDir === 'asc' ? -1 : 1;
    if(va > vb) return sortDir === 'asc' ? 1 : -1;
    return 0;
  });
  return arr;
}

function render(){
  const tb = document.querySelector('#t tbody');
  tb.innerHTML = '';
  document.getElementById('empty').hidden = data.length !== 0;

  document.querySelectorAll('#t th').forEach(th => {
    const arrow = th.querySelector('.arrow');
    arrow.textContent = th.dataset.k === sortKey ? (sortDir === 'asc' ? '▲' : '▼') : '';
  });

  for(const r of sortedData()){
    const key = r.numemp + '|' + r.numcad + '|' + r.perref;
    const tr = document.createElement('tr');
    tr.innerHTML =
      '<td>' + fmtPerref(r.perref) + '</td>' +
      '<td>' + (r.numcad ?? '') + '</td>' +
      '<td>' + (r.numemp ?? '') + '</td>' +
      '<td>' + (r.colaborador || '<span class="small">—</span>') + '</td>' +
      '<td>' + (r.codccu || '<span class="small">—</span>') + '</td>' +
      '<td class="num">' + fmt(r.horas_positivas_original) + '</td>' +
      '<td class="num">' + fmt(r.banco_230_consumido_no_mes) + '</td>' +
      '<td class="num ' + (r.horas_saldo_mes < 0 ? 'neg' : '') + '">' + fmt(r.horas_saldo_mes) + '</td>' +
      '<td class="num ' + (r.valor_saldo_mes < 0 ? 'neg' : '') + '">' + fmt(r.valor_saldo_mes) + '</td>' +
      '<td class="num">' + fmt(r.banco_total_aplicado_no_grupo) + '</td>';
    tr.title = 'Clique para ver o detalhado deste colaborador/período';
    tr.addEventListener('click', () => toggleDetail(tr, r, key));
    tb.appendChild(tr);
    if(openDetailKey === key){
      renderDetailRow(tr, r, key);
    }
  }
}

async function toggleDetail(tr, r, key){
  const existing = tr.nextElementSibling;
  if(existing && existing.classList.contains('detail-row')){
    existing.remove();
    openDetailKey = null;
    return;
  }
  document.querySelectorAll('.detail-row').forEach(e => e.remove());
  openDetailKey = key;
  await renderDetailRow(tr, r, key);
}

async function renderDetailRow(tr, r, key){
  const dtr = document.createElement('tr');
  dtr.className = 'detail-row';
  const td = document.createElement('td');
  td.colSpan = 10;
  td.innerHTML = '<div class="detail-box small">Carregando detalhado...</div>';
  dtr.appendChild(td);
  tr.after(dtr);

  try{
    const p = new URLSearchParams({ numemp: r.numemp, numcad: r.numcad, perref: r.perref });
    const res = await fetch('/detail?' + p.toString());
    if(!res.ok){
      td.innerHTML = '<div class="detail-box banner error">Erro ao buscar detalhado: ' + (await res.text()) + '</div>';
      return;
    }
    const rows = await res.json() || [];
    if(!rows.length){
      td.innerHTML = '<div class="detail-box small">Sem linhas de detalhado para este colaborador/período.</div>';
      return;
    }
    let html = '<div class="detail-box"><table><thead><tr>' +
      '<th>Data apuração</th><th>CCU</th><th>Situação</th><th class="num">Horas orig.</th>' +
      '<th class="num">Banco usado</th><th class="num">Horas saldo</th><th class="num">R$ saldo</th>' +
      '</tr></thead><tbody>';
    for(const d of rows){
      html += '<tr>' +
        '<td>' + (d.dtapuracao || '') + '</td>' +
        '<td>' + (d.codccu || '') + '</td>' +
        '<td>' + (d.dessit || '') + ' (' + d.codsit + ')</td>' +
        '<td class="num">' + fmt(d.horas_original) + '</td>' +
        '<td class="num">' + fmt(d.banco_usado_na_linha) + '</td>' +
        '<td class="num ' + (d.horas_saldo < 0 ? 'neg' : '') + '">' + fmt(d.horas_saldo) + '</td>' +
        '<td class="num ' + (d.valor_saldo < 0 ? 'neg' : '') + '">' + fmt(d.valor_saldo) + '</td>' +
        '</tr>';
    }
    html += '</tbody></table></div>';
    td.innerHTML = html;
  }catch(e){
    td.innerHTML = '<div class="detail-box banner error">Falha de rede: ' + e + '</div>';
  }
}

document.querySelectorAll('#t th').forEach(th => {
  th.addEventListener('click', () => {
    const k = th.dataset.k;
    if(sortKey === k){ sortDir = sortDir === 'asc' ? 'desc' : 'asc'; }
    else { sortKey = k; sortDir = 'asc'; }
    render();
  });
});

function exportCsv(){
  if(!data.length){ showBanner('Nada para exportar — faça uma busca primeiro.', 'error'); return; }
  // labels seguem a mesma ordem exibida na tabela (empresa <- numcad, cracha/numcad <- numemp)
  const cols = [
    {label:'perref', key:'perref'},
    {label:'empresa', key:'numcad'},
    {label:'cracha_numcad', key:'numemp'},
    {label:'colaborador', key:'colaborador'},
    {label:'codccu', key:'codccu'},
    {label:'horas_positivas_original', key:'horas_positivas_original'},
    {label:'banco_230_consumido_no_mes', key:'banco_230_consumido_no_mes'},
    {label:'horas_saldo_mes', key:'horas_saldo_mes'},
    {label:'valor_saldo_mes', key:'valor_saldo_mes'},
    {label:'banco_total_aplicado_no_grupo', key:'banco_total_aplicado_no_grupo'},
  ];
  const header = cols.map(c => c.label).join(';');
  const lines = sortedData().map(r => cols.map(c => {
    const v = r[c.key];
    if(typeof v === 'string') return '"' + v.replace(/"/g,'""') + '"';
    return String(v ?? '');
  }).join(';'));
  const csv = '\uFEFF' + [header, ...lines].join('\r\n');
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'bh_resumo.csv';
  document.body.appendChild(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(url);
}

async function rodar(){
  if(!confirm('Processar /run agora? Isso lê a origem e regrava o detalhado/resumo de destino.')) return;
  const btn = document.getElementById('btnRun');
  setBusy(true, btn, 'Processando...');
  showBanner('Processando ETL, isso pode levar alguns segundos...', 'info');
  try{
    const res = await fetch('/run', {method:'POST'});
    if(!res.ok){
      showBanner('Falha no processamento: ' + (await res.text()), 'error');
    } else {
      const j = await res.json();
      showBanner('Processado com sucesso: ' + (j.processed || 0) + ' linhas detalhadas, ' + (j.summary || 0) + ' registros de resumo.', 'info');
      await loadPeriods();
      await buscar();
    }
  }catch(e){
    showBanner('Falha de rede ao processar: ' + e, 'error');
  }
  setBusy(false, btn, 'Processar /run');
}

loadPeriods().then(buscar);
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
