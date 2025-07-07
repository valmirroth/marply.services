package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"

	"bh-mvc/internal/model"
	"bh-mvc/internal/service"
)

func DeleteAll(ctx context.Context, connStr, tblDetalhado, tblResumo string) (int64, int64, error) {
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return 0, 0, err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	tables := []string{tblDetalhado, tblResumo}
	var aff [2]int64
	for i, t := range tables {
		if strings.TrimSpace(t) == "" {
			continue
		}
		res, e := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s where perref >= '20250902' ;", t))
		if e != nil {
			err = e
			break
		}
		a, _ := res.RowsAffected()
		aff[i] = a
	}
	if err != nil {
		_ = tx.Rollback()
		return aff[0], aff[1], err
	}
	if err = tx.Commit(); err != nil {
		return aff[0], aff[1], err
	}
	return aff[0], aff[1], nil
}
func normalizePerRef(input string, dtApuracao *time.Time) (string, error) {
	// Tenta formatos comuns: RFC3339, YYYY-MM-DD, YYYYMM, YYYY-MM
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
	}

	var t time.Time
	var parseErr error
	for _, layout := range layouts {
		if tt, err := time.Parse(layout, input); err == nil {
			t = tt
			parseErr = nil
			break
		} else {
			parseErr = err
		}
	}
	// Fallback: usa dtApuracao
	if !t.IsZero() {
		// Se quiser "primeiro dia do MESMO mês":
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

		// Se a sua regra for "primeiro dia do MÊS ANTERIOR", descomente estas duas linhas:
		// t = t.AddDate(0, -1, 0)
		// t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)

		return t.Format("2006-01-02"), nil
	}
	if dtApuracao != nil {
		tt := dtApuracao.UTC()
		tt = time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, time.UTC)
		return tt.Format("2006-01-02"), nil
	}
	return "", parseErr
}
func InsertDetalhado(ctx context.Context, connStr, table string, rows []model.ProcessedRow) (int64, error) {
	if len(rows) == 0 || strings.TrimSpace(table) == "" {
		return 0, nil
	}
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	cols := []string{
		"numemp", "codbh", "numcad", "codccu", "codfil", "codcal", "perref", "dtapuracao", "cracha", "colaborador",
		"tipcol", "dessit", "codsit", "valhorames", "valhoracalculado", "horas_original", "banco_usado_na_linha",
		"horas_saldo", "valor_saldo", "valor_reais", "banco_total_aplicado_no_grupo",
	}
	placeholders := "(@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8,@p9,@p10,@p11,@p12,@p13,@p14,@p15,@p16,@p17,@p18,@p19,@p20,@p21)"
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", table, join(cols, ","), placeholders)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	stmt, err := tx.PrepareContext(ctx, sqlStr)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	var total int64
	for _, r := range rows {
		x := r.DtApuracao.Format("2006-01-02")
		perApuStr, erres := normalizePerRef(x, r.DtApuracao)
		if erres != nil {
			_ = tx.Rollback()
			return total, fmt.Errorf("perref inválido (%q): %w", r.PerRef, erres)
		}

		// perref como "YYYY-MM-DD" (dia 1 do mês)
		perrefStr, erre := normalizePerRef(r.PerRef, r.DtApuracao)
		if erre != nil {
			_ = tx.Rollback()
			return total, fmt.Errorf("perref inválido (%q): %w", r.PerRef, erre)
		}
		_, err := stmt.ExecContext(ctx,
			r.NumEmp, r.CodBh, r.NumCad, r.CodCcu, r.CodFil, r.CodCal, perrefStr, perApuStr, r.Cracha, r.Colaborador,
			r.TipCol, r.DesSit, r.CodSit, strconv.FormatFloat(r.ValHoraMes, 'f', 3, 64), strconv.FormatFloat(r.ValHoraCalculado, 'f', 3, 64), strconv.FormatFloat(r.HorasOriginal, 'f', 3, 64), strconv.FormatFloat(r.BancoUsadoNaLinha, 'f', 3, 64),
			strconv.FormatFloat(r.HorasSaldo, 'f', 3, 64), strconv.FormatFloat(r.ValorSaldo, 'f', 3, 64), strconv.FormatFloat(r.ValorReais, 'f', 3, 64), strconv.FormatFloat(r.BancoTotalAplicadoNoGrupo, 'f', 3, 64),
		)
		if err != nil {
			_ = tx.Rollback()
			return total, err
		}
		total++
	}
	if err := tx.Commit(); err != nil {
		return total, err
	}
	return total, nil
}

func UpsertResumo(ctx context.Context, connStr, table string, rows []service.SummaryRow) (int64, error) {
	if len(rows) == 0 || strings.TrimSpace(table) == "" {
		return 0, nil
	}
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// MERGE por linha (chave: numemp,numcad,perref)
	merge := fmt.Sprintf(`MERGE %s AS tgt
							USING (SELECT @p1 AS numemp, @p2 AS numcad, @p3 AS perref) AS src
							ON (tgt.numemp = src.numemp AND tgt.numcad = src.numcad AND tgt.perref = src.perref)
							WHEN MATCHED THEN UPDATE SET
								horas_positivas_original = @p4,
								banco_230_consumido_no_mes = @p5,
								horas_saldo_mes = @p6,
								valor_saldo_mes = @p7,
								banco_total_aplicado_no_grupo = @p8
							WHEN NOT MATCHED THEN INSERT
								(numemp, numcad, perref, horas_positivas_original, banco_230_consumido_no_mes, horas_saldo_mes, valor_saldo_mes, banco_total_aplicado_no_grupo)
								VALUES (@p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16);`, table)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	stmt, err := tx.PrepareContext(ctx, merge)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	var total int64
	for _, r := range rows {
		_, err := stmt.ExecContext(ctx,
			// src
			r.NumEmp, r.NumCad, r.PerRef,
			// update
			r.HorasPositivasOriginal, r.Banco230ConsumidoMes, r.HorasSaldoMes, r.ValorSaldoMes, r.BancoTotalAplicadoNoGrupo,
			// insert
			r.NumEmp, r.NumCad, r.PerRef, r.HorasPositivasOriginal, r.Banco230ConsumidoMes, r.HorasSaldoMes, r.ValorSaldoMes, r.BancoTotalAplicadoNoGrupo,
		)
		if err != nil {
			_ = tx.Rollback()
			return total, err
		}
		total++
	}
	if err := tx.Commit(); err != nil {
		return total, err
	}
	return total, nil
}

func makeQ(n int) string {
	s := "?"
	for i := 1; i < n; i++ {
		s += ",?"
	}
	return s
}
func join(a []string, sep string) string {
	if len(a) == 0 {
		return ""
	}
	s := a[0]
	for i := 1; i < len(a); i++ {
		s += sep + a[i]
	}
	return s
}
