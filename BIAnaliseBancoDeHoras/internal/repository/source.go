package repository

import (
	"context"
	"database/sql"

	"strconv"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"

	"bh-mvc/internal/model"
	"bh-mvc/internal/util"
)

const srcQuery = `
	select codcal, codbh, perref, sitafa,
		codccu, dtApuracao, codfil, 
		resumo.CRACHA, resumo.Colaborador, 
		resumo.numcad, resumo.numemp, resumo.tipcol, 
		resumo.ValHoraMes, ValHoraCalculado, resumo.dessit, 
		resumo.codsit, SUM(qtdhor) as horas, (ValHoraCalculado * SUM( resumo.qtdhor ) ) as [R$ valor]
		from (
		SELECT ( select codbhr from r038hsi where numcad = r034fun.numcad and numemp = r034fun.numemp  and tipcol = r034fun.tipcol
				and datalt = (
							select max(datalt) from r038hsi 
							where 
								numcad = r034fun.numcad 
								and numemp = r034fun.numemp 
								and tipcol = r034fun.tipcol
								and datalt <= r066apu.datapu
							) ) codbh,
			r044cal.codcal, r044cal.perref, r034fun.codccu, convert(varchar,(r066sit.datapu),112) dtApuracao,
			r034fun.codfil, getdate() as DataAtual, 
			cast(r034fun.numcad  as varchar) +' - '+ r034fun.nomfun as Colaborador,
			CAST(r034fun.numcad AS VARCHAR) AS CRACHA,               
			r066apu.numemp,       
			r034fun.sitafa,       
			r066apu.numcad,               
			r066apu.datapu,               
			r066apu.tipcol,              
			nomesc AS TURNO ,  
			r034fun.valsal / 
			cast(r034fun.valsal / 
			(hormes/60) as decimal(19,6))as ValHoraMes,
			cast( case when sit.codsit = 331 then ( r034fun.valsal / 
			(hormes/60))*1.6
			when sit.codsit = 330 then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit = 332 then ( r034fun.valsal / 
			(hormes/60))*2
			when sit.codsit in (230,109,336) then ( r034fun.valsal / 
			(hormes/60))	  
			when sit.codsit = 313 then ( r034fun.valsal / 
			(hormes/60))*1.5
			when sit.codsit = 302 then (( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5
			when sit.codsit = 301 then (( r034fun.valsal / 
			(hormes/60))*1.5)  
				when sit.codsit = 305 then (( r034fun.valsal / 
			(hormes/60))*2) 
				when sit.codsit = 311 then (( r034fun.valsal / 
			(hormes/60))*1.6) 
				when sit.codsit = 306 then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 2)
				when sit.codsit = 314 then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.5)
				when sit.codsit = 312 then ((( r034fun.valsal / 
			(hormes/60))*1.2) * 1.6)

			end as decimal(19,6)) as ValHoraCalculado,
			case when sit.codsit in (109,230,336) then 'Banco de Horas Negativo' else 
			replace(
			replace(
			replace(   
			replace(  
			replace(replace(replace(sit.dessit,'Banco de Horas','B.H.'),'H. Extra','H.E.'),'Horas Extras','H.E.') 
			,'Intrajornada','Intr')
			,'Noturnas','Not')
			,'Intraj','Intr')
			,'Apurad','Apu.')
			end as dessit,    
			case when sit.codsit in (230,109,336) then 230 else sit.codsit end as codsit,
			case when sit.codsit in (230,109,336) then (r066sit.qtdhor/60.000) * (-1) else (r066sit.qtdhor/60.000)  end qtdhor 
		FROM r066sit
			left join r066apu on  r066sit.numemp = r066apu.numemp and r066sit.tipcol = r066apu.tipcol and r066sit.numcad = r066apu.numcad           
			and r066sit.datapu = r066apu.datapu
			left join r034fun (nolock) on r034fun.numcad = r066apu.numcad   
			LEFT JOIN r010sit sit on sit.codsit = r066sit.codsit
			INNER JOIN r006esc (NOLOCK) ON r066apu.codesc = r006esc.codesc
			inner join r010sit on r010sit.codsit = r034fun.sitafa
			inner join r044cal on r044cal.tipcal ='11' and  r066apu.datapu between r044cal.iniapu and r044cal.fimapu and r044cal.numemp = r066apu.numemp 
		where 1=1 
			--r010sit.codsit <> '7'
			--and r034fun.codccu in ('cc103','cc789','cc252')
			and convert(varchar,( r066sit.datapu),112) >= '20250326' 
			and convert(varchar,(r066sit.datapu),112) < '20250926'
			and   sit.codsit in ( select  codsit from r011eve where codbhr in (8) )
			and codfil in (1,5)
		) as resumo
		where CRACHA <> 6065 
				and numemp = 1
		--	and CRACHA = '7113'
			and codbh = 8
	group by codcal, codbh, perref, codccu, sitafa,
			dtApuracao, codfil, resumo.CRACHA, 
			ValHoraCalculado, Colaborador, resumo.numcad, 
			resumo.numemp,  resumo.tipcol, resumo.ValHoraMes, 
			resumo.dessit, resumo.codsit
`

// select * from r011lan where numcad = 7094
// From SQL Server
func LoadFromSQL(ctx context.Context, connStr string) ([]model.RawRow, error) {
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, qrySemestral)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.RawRow
	for rows.Next() {
		var (
			codcal      sql.NullInt64
			codbh       sql.NullInt64
			perref      sql.NullString
			sitafa      sql.NullString // --- IGNORE ---
			codccu      sql.NullString
			dtStr       sql.NullString
			codfil      sql.NullInt64
			cracha      sql.NullString
			colab       sql.NullString
			numemp      sql.NullInt64
			numcad      sql.NullInt64
			tipcol      sql.NullInt64
			valHoraMes  sql.NullFloat64
			valHoraCalc sql.NullFloat64
			dessit      sql.NullString
			codsit      sql.NullInt64
			horas       sql.NullFloat64
			valorReais  sql.NullFloat64
		)
		if err := rows.Scan(&codcal, &codbh, &perref, &sitafa, &codccu, &dtStr, &codfil, &cracha, &colab, &numemp, &numcad, &tipcol, &valHoraMes, &valHoraCalc, &dessit, &codsit, &horas, &valorReais); err != nil {
			return nil, err
		}
		var dt *time.Time
		if dtStr.Valid && dtStr.String != "" {
			// dtApuracao vem como yyyyMMdd
			if len(dtStr.String) == 8 {
				t, err := time.Parse("20060102", dtStr.String)
				if err == nil {
					dt = &t
				}
			} else {
				t, err := time.Parse("2006-01-02", dtStr.String)
				if err == nil {
					dt = &t
				}
			}
		}
		out = append(out, model.RawRow{
			NumEmp:           int(numemp.Int64),
			CodBh:            int(codbh.Int64),
			NumCad:           int(numcad.Int64),
			CodCcu:           strings.TrimSpace(codccu.String),
			SitAfa:           strings.TrimSpace(sitafa.String),
			CodFil:           int(codfil.Int64),
			CodCal:           int(codcal.Int64),
			PerRef:           strings.TrimSpace(perref.String),
			DtApuracao:       dt,
			Cracha:           strings.TrimSpace(cracha.String),
			Colaborador:      strings.TrimSpace(colab.String),
			TipCol:           int(tipcol.Int64),
			DesSit:           strings.TrimSpace(dessit.String),
			CodSit:           int(codsit.Int64),
			ValHoraMes:       valHoraMes.Float64,
			ValHoraCalculado: valHoraCalc.Float64,
			Horas:            horas.Float64,
			ValorReais:       valorReais.Float64,
		})
	}
	return out, rows.Err()
}

// From CSV (fallback)
func LoadFromCSV(path string, sep rune) ([]model.RawRow, error) {
	head, rows, err := util.ReadCSVAll(path, sep)
	if err != nil {
		return nil, err
	}

	// normaliza headers
	for i := range head {
		head[i] = strings.ToLower(strings.TrimSpace(head[i]))
	}
	idx := func(name string) int {
		for i, h := range head {
			if h == name {
				return i
			}
		}
		return -1
	}
	get := func(r []string, i int) string {
		if i < 0 || i >= len(r) {
			return ""
		}
		return strings.TrimSpace(r[i])
	}

	// mapeia índices
	iNumEmp := idx("numemp")
	iCodBh := idx("codbh")
	iNumCad := idx("numcad")
	iCodCcu := idx("codccu")
	iCodFil := idx("codfil")
	iCodCal := idx("codcal")
	iPerRef := idx("perref")
	iDt := idx("dtapuracao")
	iCracha := idx("cracha")
	iColab := idx("colaborador")
	iTipCol := idx("tipcol")
	iDesSit := idx("dessit")
	iCodSit := idx("codsit")
	iValHM := idx("valhorames")
	iValHC := idx("valhoracalculado")
	iHoras := idx("horas")
	iValR := idx("valor_reais")

	var out []model.RawRow
	for _, r := range rows {
		// parse
		num := func(s string) float64 {
			f, _ := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", "."), 64)
			return f
		}
		numi := func(s string) int {
			v, _ := strconv.Atoi(strings.Split(s, ".")[0])
			return v
		}
		var dt *time.Time
		ds := get(r, iDt)
		if ds != "" {
			// tenta yyyyMMdd, yyyy-MM-dd, dd/MM/yyyy
			if len(ds) == 8 {
				if t, err := time.Parse("20060102", ds); err == nil {
					dt = &t
				}
			} else if strings.Contains(ds, "/") {
				if t, err := time.Parse("02/01/2006", ds); err == nil {
					dt = &t
				}
			} else if strings.Contains(ds, "-") {
				if t, err := time.Parse("2006-01-02", ds); err == nil {
					dt = &t
				}
			}
		}
		pref := get(r, iPerRef)
		if pref == "" && dt != nil {
			pref = dt.Format("200601")
		}

		out = append(out, model.RawRow{
			NumEmp:           numi(get(r, iNumEmp)),
			CodBh:            numi(get(r, iCodBh)),
			NumCad:           numi(get(r, iNumCad)),
			CodCcu:           get(r, iCodCcu),
			CodFil:           numi(get(r, iCodFil)),
			CodCal:           numi(get(r, iCodCal)),
			PerRef:           pref,
			DtApuracao:       dt,
			Cracha:           get(r, iCracha),
			Colaborador:      get(r, iColab),
			TipCol:           numi(get(r, iTipCol)),
			DesSit:           get(r, iDesSit),
			CodSit:           numi(get(r, iCodSit)),
			ValHoraMes:       num(get(r, iValHM)),
			ValHoraCalculado: num(get(r, iValHC)),
			Horas:            num(get(r, iHoras)),
			ValorReais:       num(get(r, iValR)),
		})
	}
	return out, nil
}

// select * from r011lan where numcad = 7094
// From SQL Server
func LoadFromSQLMensal(ctx context.Context, connStr string) ([]model.RawRow, error) {
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, qryMensal)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.RawRow
	for rows.Next() {
		var (
			codcal      sql.NullInt64
			codbh       sql.NullInt64
			perref      sql.NullString
			sitafa      sql.NullString // --- IGNORE ---
			codccu      sql.NullString
			dtStr       sql.NullString
			codfil      sql.NullInt64
			cracha      sql.NullString
			colab       sql.NullString
			numemp      sql.NullInt64
			numcad      sql.NullInt64
			tipcol      sql.NullInt64
			valHoraMes  sql.NullFloat64
			valHoraCalc sql.NullFloat64
			dessit      sql.NullString
			codsit      sql.NullInt64
			horas       sql.NullFloat64
			valorReais  sql.NullFloat64
		)
		if err := rows.Scan(&codcal, &codbh, &perref, &sitafa, &codccu, &dtStr, &codfil, &cracha, &colab, &numemp, &numcad, &tipcol, &valHoraMes, &valHoraCalc, &dessit, &codsit, &horas, &valorReais); err != nil {
			return nil, err
		}
		var dt *time.Time
		if dtStr.Valid && dtStr.String != "" {
			// dtApuracao vem como yyyyMMdd
			if len(dtStr.String) == 8 {
				t, err := time.Parse("20060102", dtStr.String)
				if err == nil {
					dt = &t
				}
			} else {
				t, err := time.Parse("2006-01-02", dtStr.String)
				if err == nil {
					dt = &t
				}
			}
		}
		out = append(out, model.RawRow{
			NumEmp:           int(numemp.Int64),
			CodBh:            int(codbh.Int64),
			NumCad:           int(numcad.Int64),
			CodCcu:           strings.TrimSpace(codccu.String),
			SitAfa:           strings.TrimSpace(sitafa.String),
			CodFil:           int(codfil.Int64),
			CodCal:           int(codcal.Int64),
			PerRef:           strings.TrimSpace(perref.String),
			DtApuracao:       dt,
			Cracha:           strings.TrimSpace(cracha.String),
			Colaborador:      strings.TrimSpace(colab.String),
			TipCol:           int(tipcol.Int64),
			DesSit:           strings.TrimSpace(dessit.String),
			CodSit:           int(codsit.Int64),
			ValHoraMes:       valHoraMes.Float64,
			ValHoraCalculado: valHoraCalc.Float64,
			Horas:            horas.Float64,
			ValorReais:       valorReais.Float64,
		})
	}
	return out, rows.Err()
}
