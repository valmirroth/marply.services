package RSDataTemplate

import "github.com/shopspring/decimal"

type EstoqueRsDataModel struct {
	CodInterno      string
	TpClasseProduto string
	CodGrade        string
	TamanhoGrade    string

	LocalEstoque int
	TpMovimento  string
	Descricao    string
	DtMovimento  string
	Usuario      string
	Senha        string
	Quantidade   decimal.Decimal
	VlUnitario   decimal.Decimal
	Cnpj         string
	DataIni      string
	DataFim      string
}
