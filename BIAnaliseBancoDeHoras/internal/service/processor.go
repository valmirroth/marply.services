package service

import (
	"sort"
	"strconv"
	"time"

	"bh-mvc/internal/model"
)

// ApplyBankOffset executa a mesma lógica do Python para cada grupo (numemp,numcad):
// - Soma total do banco (codsit==230) como positivo
// - Abate das horas positivas em ordem (mais antigo -> mais recente)
// - Sobra do banco vira horas_saldo negativa na linha 230 mais recente
func ApplyBankOffset(rows []model.RawRow) []model.ProcessedRow {
	// agrupa por (numemp,numcad)
	groups := map[string][]model.RawRow{}
	key := func(r model.RawRow) string {
		return fmtKey(r.NumEmp, r.NumCad)
	}
	for _, r := range rows {
		groups[key(r)] = append(groups[key(r)], r)
	}
	var out []model.ProcessedRow
	for _, g := range groups {
		out = append(out, applyGroup(g)...)
	}
	return out
}

func fmtKey(a, b int) string { return fmtS(a) + "|" + fmtS(b) }
func fmtS(i int) string      { return strconv.Itoa(i) }

// applyGroup trata um grupo (numemp,numcad)
func applyGroup(g []model.RawRow) []model.ProcessedRow {
	// copia para não mexer no original
	rows := append([]model.RawRow(nil), g...)

	// total banco (negativos de 230 viram positivo)
	bancoTotal := 0.0
	for _, r := range rows {
		if r.CodSit == 230 {
			bancoTotal += -r.Horas
		}
	}
	bancoRestante := bancoTotal

	// separa positivas (codsit != 230)
	var pos []int
	var bank []int
	for i, r := range rows {
		if r.CodSit != 230 {
			pos = append(pos, i)
		} else {
			bank = append(bank, i)
		}
	}
	// ordena pos por dtapuracao asc, perref asc (nil last)
	sort.SliceStable(pos, func(i, j int) bool {
		ri := rows[pos[i]]
		rj := rows[pos[j]]
		// dt
		var di, dj time.Time
		if ri.DtApuracao != nil {
			di = *ri.DtApuracao
		}
		if rj.DtApuracao != nil {
			dj = *rj.DtApuracao
		}
		if !di.Equal(dj) {
			return di.Before(dj)
		}
		// perref
		return ri.PerRef < rj.PerRef
	})

	processed := make([]model.ProcessedRow, len(rows))

	// inicializa cópia base
	for i, r := range rows {
		processed[i] = model.ProcessedRow{
			NumEmp: r.NumEmp, CodBh: r.CodBh, NumCad: r.NumCad,
			CodCcu: r.CodCcu, CodFil: r.CodFil, CodCal: r.CodCal,
			SitAfa: r.SitAfa,
			PerRef: r.PerRef, DtApuracao: r.DtApuracao,
			Cracha: r.Cracha, Colaborador: r.Colaborador,
			TipCol: r.TipCol, DesSit: r.DesSit, CodSit: r.CodSit,
			ValHoraMes: r.ValHoraMes, ValHoraCalculado: r.ValHoraCalculado,
			HorasOriginal:             r.Horas,
			ValorReais:                r.ValorReais,
			BancoTotalAplicadoNoGrupo: bancoTotal,
		}
	}

	// processa positivas
	for _, idx := range pos {
		h := rows[idx].Horas
		if h <= 0 {
			processed[idx].HorasSaldo = 0
			continue
		}
		if bancoRestante <= 0 {
			processed[idx].HorasSaldo = h
			continue
		}
		abat := h
		if bancoRestante < h {
			abat = bancoRestante
		}
		processed[idx].BancoUsadoNaLinha = abat
		processed[idx].HorasSaldo = h - abat
		bancoRestante -= abat
		if processed[idx].ValHoraCalculado != 0 {
			processed[idx].ValorSaldo = processed[idx].HorasSaldo * processed[idx].ValHoraCalculado
		}
	}

	// iniciais para banco
	for _, idx := range bank {
		processed[idx].BancoUsadoNaLinha = 0
		processed[idx].HorasSaldo = 0
	}

	// Se sobrou banco, coloca negativo na linha 230 mais recente
	if bancoRestante > 0 && len(bank) > 0 {
		// ordenar bank por dtapuracao desc, perref desc
		sort.SliceStable(bank, func(i, j int) bool {
			ri := rows[bank[i]]
			rj := rows[bank[j]]
			var di, dj time.Time
			if ri.DtApuracao != nil {
				di = *ri.DtApuracao
			}
			if rj.DtApuracao != nil {
				dj = *rj.DtApuracao
			}
			if !di.Equal(dj) {
				return dj.Before(di)
			} // desc
			return ri.PerRef > rj.PerRef
		})
		idx := bank[0]

		if processed[idx].SitAfa == "7" {
			processed[idx].HorasSaldo = 0
			if processed[idx].ValHoraCalculado != 0 {
				processed[idx].ValorSaldo = processed[idx].HorasSaldo * processed[idx].ValHoraCalculado
			}
		} else {
			processed[idx].HorasSaldo = -bancoRestante
			if processed[idx].ValHoraCalculado != 0 {
				processed[idx].ValorSaldo = processed[idx].HorasSaldo * processed[idx].ValHoraCalculado
			}
		}
	}

	// calcula ValorSaldo das demais (que ficaram zeradas acima)
	for i := range processed {
		if processed[i].ValorSaldo == 0 && processed[i].ValHoraCalculado != 0 {
			processed[i].ValorSaldo = processed[i].HorasSaldo * processed[i].ValHoraCalculado
		}
	}
	return processed
}
