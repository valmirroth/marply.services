package entities

import "database/sql"

type HisMovEPI struct {
	Text         sql.NullString `gorm:"column:Text"`
	LocalEstoque sql.NullString `gorm:"column:LocalEstoque"`
	CodEmpregado sql.NullString `gorm:"column:CodEmpregado"`
	TpMovimento  sql.NullString `gorm:"column:TpMovimento"`
	DtMovimento  sql.NullString `gorm:"column:DtMovimento"`
	CDI          sql.NullString `gorm:"column:CDI"`
	Quantidade   sql.NullString `gorm:"column:Quantidade"`
	TxObs        sql.NullString `gorm:"column:TxObs"`
	VlUnitario   sql.NullString `gorm:"column:VlUnitario"`
	CodItem      sql.NullString `gorm:"column:CODITEM"`
	EmpresaRecno sql.NullString `gorm:"column:EMPRESA_RECNO"`
	Recno        sql.NullInt64  `gorm:"primaryKey;column:R_E_C_N_O_"`
}

func (HisMovEPI) TableName() string {
	return "CST_ENTREGA_EPI_RSDATA"
}

type EntregaEPISetor struct {
	IdSetor          sql.NullString `gorm:"column:IdSetor"`
	NomeSetor        sql.NullString `gorm:"column:NomeSetor"`
	NomeEpi          sql.NullString `gorm:"column:NomeEpi"`
	TipoEpi          sql.NullString `gorm:"column:TipoEpi"`
	Quantidade       sql.NullString `gorm:"column:Quantidade"`
	CodIntegracaoEpi sql.NullString `gorm:"column:CodIntegracaoEpi"`
	DataMovimento    sql.NullString `gorm:"column:DataMovimento"`
	Tamanho          sql.NullString `gorm:"column:Tamanho"` // Adicionado campo Tamanho
	Cnpj             sql.NullString `gorm:"column:Cnpj"`
	Local            sql.NullString `gorm:"column:Local"`
	EmpresaRecno     sql.NullInt64  `gorm:"column:EmpresaRecno"`
}

func (EntregaEPISetor) TableName() string {
	return "CST_RSDATA_ENTREGA_SETOR"
}

type SaldosRSData struct {
	IdSetor          sql.NullString `gorm:"column:Local"`
	CodIntegracaoEpi sql.NullString `gorm:"column:Codigo"`
	Quantidade       sql.NullString `gorm:"column:Quantidade"`
}

func (SaldosRSData) TableName() string {
	return "CST_SaldosRSData"
}
