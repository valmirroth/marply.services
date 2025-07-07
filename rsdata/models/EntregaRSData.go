package models

type EntregaEPISetor struct {
	IdSetor          string `gorm:"column:IdSetor"`
	NomeSetor        string `gorm:"column:NomeSetor"`
	NomeEpi          string `gorm:"column:NomeEpi"`
	TipoEpi          string `gorm:"column:TipoEpi"`
	Quantidade       string `gorm:"column:Quantidade"`
	CodIntegracaoEpi string `gorm:"column:CodIntegracaoEpi"`
	DataMovimento    string `gorm:"column:DataMovimento"`
	Cnpj             string `gorm:"column:Cnpj"`
	Local            string `gorm:"column:Local"`
	EmpresaRecno     int    `gorm:"column:EmpresaRecno"`
	Tamanho          string `gorm:"column:Tamanho"` // Adicionado campo Tamanho
}

func (EntregaEPISetor) TableName() string {
	return "CST_RSDATA_ENTREGA_SETOR"
}

type SaldosRSData struct {
	IdSetor          string `gorm:"column:Local"`
	CodIntegracaoEpi string `gorm:"column:Codigo"`
	Quantidade       string `gorm:"column:Quantidade"`
}

func (SaldosRSData) TableName() string {
	return "CST_SaldosRSData"
}
