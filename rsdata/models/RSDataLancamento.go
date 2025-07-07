package models

import "github.com/shopspring/decimal"

type RSDataLancamento struct {
	RecnoHistlise int             `gorm:"column:RECNO_HISTLISE"`
	Item          string          `gorm:"column:ITEM"`
	Descricao     string          `gorm:"column:DESCRICAO"`
	QtdMovimento  decimal.Decimal `gorm:"column:QTDEMOV"`
	ValUni        decimal.Decimal `gorm:"column:UNI"`
	DataMovimento string          `gorm:"column:DATA_MOV"`
	RecnoHF       int             `gorm:"column:RECNO_HISTLISE_FOR"`
	CodGradeEPI   string          `gorm:"column:COD_GRADE"`
	TamanhoGrade  string          `gorm:"column:TAMANHO_GRADE"`
	RetornoRSData string          `gorm:"column:RETORNORSDATA"`
	Cnpj          string          `gorm:"column:Cnpj"`
	LocalEstoque  string          `gorm:"column:LocalEstoque"`
}

/*
CREATE TABLE CST_LANCAMENTOS_RSDATA (
		RECNO_HISTLISE INT NULL,
		ITEM VARCHAR(10) NULL,
		DESCRICAO VARCHAR(200) NULL,
		QTDEMOV DECIMAL(19,6) NULL,
		DATA_MOV DATETIME DEFAULT(GETDATE()),
		RECNO_HISTLISE_FOR INT NULL,
		RETORNORSDATA VARCHAR(1000) NULL,
		R_E_C_N_O_ INT IDENTITY NOT NULL
)
SELECT * FROM CST_LANCAMENTOS_RSDATA
*/
