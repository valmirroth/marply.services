package rsdatamodel

import "encoding/xml"

type RsGetEstoqueEnvelope struct {
	XMLName xml.Name     `xml:"Envelope"`
	Text    string       `xml:",chardata"`
	SOAPENV string       `xml:"SOAP-ENV,attr"`
	Header  string       `xml:"Header"`
	Body    BodyResponse `xml:"Body"`
}

type Mov struct {
	Text         string `xml:",chardata"`
	LocalEstoque string `xml:"localEstoque"`
	CodEmpregado string `xml:"codEmpregado"`
	TpMovimento  string `xml:"tpMovimento"`
	DtMovimento  string `xml:"dtMovimento"`
	CDI          string `xml:"CDI"`
	Quantidade   string `xml:"quantidade"`
	TxObs        string `xml:"txObs"`
	VlUnitario   string `xml:"vlUnitario"`
}

type ListMov struct {
	Text      string `xml:",chardata"`
	Movimento []Mov  `xml:"movimento"`
}

type Produto struct {
	Text            string  `xml:",chardata"`
	ID              string  `xml:"ID"`
	NrCNPJEmpresa   string  `xml:"nrCNPJEmpresa"`
	Descricao       string  `xml:"descricao"`
	Marca           string  `xml:"marca"`
	Modelo          string  `xml:"modelo"`
	Referencia      string  `xml:"referencia"`
	CodInterno      string  `xml:"codInterno"`
	CodFornecedor   string  `xml:"codFornecedor"`
	TpClasseProduto string  `xml:"tpClasseProduto"`
	TpEstoque       string  `xml:"tpEstoque"`
	TpSituacao      string  `xml:"tpSituacao"`
	CA              string  `xml:"CA"`
	DtValCA         string  `xml:"dtValCA"`
	VidaUtil        string  `xml:"vidaUtil"`
	TpVidaUtil      string  `xml:"tpVidaUtil"`
	CodBarra        string  `xml:"codBarra"`
	NRR             string  `xml:"NRR"`
	NRRSF           string  `xml:"NRRSF"`
	TpAuditivo      string  `xml:"tpAuditivo"`
	Movimentos      ListMov `xml:"movimentos"`
}

type RsDataRet struct {
	Text   string `xml:",chardata"`
	Versao string `xml:"versao,attr"`
	Config struct {
		Text         string `xml:",chardata"`
		TpVerEstoque string `xml:"tpVerEstoque"`
	} `xml:"config"`
	Estoques struct {
		Text    string    `xml:",chardata"`
		Estoque []Produto `xml:"estoque"`
	} `xml:"estoques"`
}

type RsResponse struct {
	Text   string    `xml:",chardata"`
	Ns2    string    `xml:"ns2,attr"`
	Rsdata RsDataRet `xml:"rsdata"`
}

type BodyResponse struct {
	Text                string     `xml:",chardata"`
	GetEstoquesResponse RsResponse `xml:"getEstoquesResponse"`
}
