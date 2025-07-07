package rsdatamodel

import (
	"encoding/xml"
)

type NewEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	XMLName                    xml.Name                   `xml:"Body"`
	GetEntregasBySetorResponse GetEntregasBySetorResponse `xml:"getEntregasBySetorResponse"`
}

type GetEntregasBySetorResponse struct {
	XMLName xml.Name `xml:"getEntregasBySetorResponse"`
	Dados   Dados    `xml:"dados"`
}

type Dados struct {
	XMLName xml.Name `xml:"dados"`
	Empresa Empresa  `xml:"empresa"`
}

type Empresa struct {
	XMLName              xml.Name `xml:"empresa"`
	IdEmpresa            int      `xml:"idEmpresa"`
	CodIntegracaoEmpresa string   `xml:"codIntegracaoEmpresa"`
	NrCNPJEmpresa        string   `xml:"nrCNPJEmpresa"`
	RazaoSocialEmpresa   string   `xml:"razaoSocialEmpresa"`
	DenominacaoEmpresa   string   `xml:"denominacaoEmpresa"`
	Setores              Setores  `xml:"setores"`
}

type Setores struct {
	XMLName xml.Name   `xml:"setores"`
	Setor   []NewSetor `xml:"setor"`
}

type NewSetor struct {
	XMLName   xml.Name `xml:"setor"`
	CdSetor   string   `xml:"cdSetor"`
	NomeSetor string   `xml:"nomeSetor"`
	Entregas  Entregas `xml:"entregas"`
}

type Entregas struct {
	XMLName xml.Name  `xml:"entregas"`
	Entrega []Entrega `xml:"entrega"`
}

type Entrega struct {
	XMLName          xml.Name `xml:"entrega"`
	CodIntegracaoEpi string   `xml:"codIntegracaoEpi"`
	NomeEpi          string   `xml:"nomeEpi"`
	Tamanho          string   `xml:"tamanho"`
	TipoEpi          string   `xml:"tipoEpi"`
	Quantidade       int      `xml:"quantidade"`
}
