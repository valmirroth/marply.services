package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	UseSQLOrigin bool
	SrcConn      string
	InputCSV     string
	CSVSep       rune

	DestConn     string

	TblDetalhado string
	TblResumo    string

	ClearDest    bool
	UpsertResumo bool

	HTTPPort     string
}

func Load() Config {
	_ = godotenv.Load()

	useSQL := parseBool(getEnv("USE_SQL_ORIGIN", "true"))
	csvSep := ';'
	if v := os.Getenv("CSV_SEP"); v != "" {
		csvSep = rune(v[0])
	}
	clearDest := parseBool(getEnv("CLEAR_DEST", "true"))
	upsertResumo := parseBool(getEnv("UPSERT_RESUMO", "true"))

	cfg := Config{
		UseSQLOrigin: useSQL,
		SrcConn:      os.Getenv("SRC_SQL_CONNECTION"),
		InputCSV:     getEnv("INPUT_CSV", "./resultado_query.csv"),
		CSVSep:       csvSep,

		DestConn:     os.Getenv("DEST_SQL_CONNECTION"),

		TblDetalhado: getEnv("TBL_DETALHADO", "dbo.BH_DetalhadoPosAbatimento"),
		TblResumo:    getEnv("TBL_RESUMO", "dbo.BH_ResumoMensalSaldo"),

		ClearDest:    clearDest,
		UpsertResumo: upsertResumo,

		HTTPPort:     getEnv("HTTP_PORT", ":8080"),
	}

	// Valida conexões mínimas
	if cfg.DestConn == "" {
		log.Fatal("DEST_SQL_CONNECTION não definido no ambiente (.env)")
	}
	if cfg.UseSQLOrigin && cfg.SrcConn == "" {
		log.Fatal("SRC_SQL_CONNECTION não definido e USE_SQL_ORIGIN=true")
	}
	return cfg
}

func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil { return false }
	return b
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" { return v }
	return d
}
