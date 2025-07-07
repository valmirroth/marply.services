package rsdatamodel

import "encoding/xml"

type EnvelopeEstRetorno struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Bodys    `xml:"Body"`
}

type Bodys struct {
	XMLName           xml.Name          `xml:"Body"`
	InsertEstoqueResp InsertEstoqueResp `xml:"insertEstoqueResponse"`
}

type InsertEstoqueResp struct {
	XMLName      xml.Name      `xml:"insertEstoqueResponse"`
	CdMsg        string        `xml:"cdMsg"`
	RetornoMsg   string        `xml:"retornoMsg"`
	RsdataReturn RsdataReturns `xml:"rsdataReturn"`
}

type RsdataReturns struct {
	XMLName   xml.Name  `xml:"rsdataReturn"`
	Mensagens Mensagens `xml:"mensagens"`
}

type Mensagens struct {
	XMLName  xml.Name `xml:"mensagens"`
	Mensagem Mensagem `xml:"mensagem"`
}

type Mensagem struct {
	XMLName     xml.Name `xml:"mensagem"`
	TpRetorno   string   `xml:"tpRetorno"`
	TxDescricao string   `xml:"txDescricao"`
}
