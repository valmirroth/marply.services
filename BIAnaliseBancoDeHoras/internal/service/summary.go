package service

import (
	"bh-mvc/internal/model"
)

type SummaryRow struct {
	NumEmp                    int
	NumCad                    int
	PerRef                    string
	HorasPositivasOriginal    float64
	Banco230ConsumidoMes      float64
	HorasSaldoMes             float64
	ValorSaldoMes             float64
	BancoTotalAplicadoNoGrupo float64
}

// BuildMonthlySummary agrega o slice processado (exclui codsit 230 para as métricas de mês)
func BuildMonthlySummary(processed []model.ProcessedRow) []SummaryRow {
	// key: numemp|numcad|perref
	agg := map[string]*SummaryRow{}
	bancoMax := map[string]float64{} // key: numemp|numcad -> max banco

	key3 := func(n1, n2 int, per string) string { return fmtKey(n1, n2) + "|" + per }

	for _, r := range processed {
		k2 := fmtKey(r.NumEmp, r.NumCad)
		if r.BancoTotalAplicadoNoGrupo > bancoMax[k2] {
			bancoMax[k2] = r.BancoTotalAplicadoNoGrupo
		}

		// Sempre garanta a entrada do mês, mesmo que só existam linhas 230
		k := key3(r.NumEmp, r.NumCad, r.PerRef)
		if _, ok := agg[k]; !ok {
			agg[k] = &SummaryRow{NumEmp: r.NumEmp, NumCad: r.NumCad, PerRef: r.PerRef}
		}
		a := agg[k]

		// Para as linhas "normais"
		if r.CodSit != 230 {
			a.HorasPositivasOriginal += r.HorasOriginal
			a.Banco230ConsumidoMes += r.BancoUsadoNaLinha
		}

		// Sempre refletir o resultado final (positivo ou negativo)
		a.HorasSaldoMes += r.HorasSaldo
		a.ValorSaldoMes += r.ValorSaldo
	}

	// espalha banco_total_aplicado_no_grupo (max por (numemp,numcad))
	for _, a := range agg {
		k2 := fmtKey(a.NumEmp, a.NumCad)
		a.BancoTotalAplicadoNoGrupo = bancoMax[k2]
	}

	// flatten
	out := make([]SummaryRow, 0, len(agg))
	for _, v := range agg {
		out = append(out, *v)
	}
	return out
}
