package models

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type RSDataLLog struct {
	RecnoHistlise    int             `gorm:"column:RECNO_HISTLISE"`
	Item             sql.NullString  `gorm:"column:ITEM"`
	GradeTamanho     sql.NullString  `gorm:"column:GRADE_TAMANHO"`
	Descricao        sql.NullString  `gorm:"column:DESCRICAO"`
	QtdMovimento     decimal.Decimal `gorm:"column:QTDEMOV"`
	DataMovimento    sql.NullString  `gorm:"column:DATA_MOV"`
	RecnoHF          int             `gorm:"column:RECNO_HISTLISE_FOR"`
	RetornoRSData    sql.NullString  `gorm:"column:RETORNORSDATA"`
	RetornoDetalhado sql.NullString  `gorm:"column:RETORNO_DETALHADO"`
	Recno            sql.NullInt64   `gorm:"primaryKey;column:R_E_C_N_O_"`
	Cnpj             sql.NullString  `gorm:"column:CNPJ"`
	LocalEstoque     sql.NullString  `gorm:"column:LOCAL_ESTOQUE"`
}

func (RSDataLLog) TableName() string {
	return "CST_LANCAMENTOS_RSDATA"
}
