package rsdatamodel

import (
	"encoding/xml"
)

type EstoqueEnvelope struct {
	XMLName xml.Name `xml:"SOAP-ENV:Envelope"`
	Header  string   `xml:"SOAP-ENV:Header"`
	Body    eBody    `xml:"SOAP-ENV:Body"`
}

type eBody struct {
	GetEstoqueResponse GetEstoqueResponse `xml:"ns2:getEstoqueResponse"`
}

type GetEstoqueResponse struct {
	RsdataReturn RsdataReturn `xml:"ns2:rsdataReturn"`
}

type RsdataReturn struct {
	Estoques []Estoque `xml:"ns2:estoques>ns2:estoque"`
}

type Estoque struct {
	LocalEstoque LocalEstoque `xml:"ns2:localEstoque"`
	Epi          Epi          `xml:"ns2:epi"`
}

type LocalEstoque struct {
	IdLocalEstoque            int    `xml:"ns2:idLocalEstoque"`
	NomeLocalEstoque          string `xml:"ns2:nomeLocalEstoque"`
	CodIntegracaoLocalEstoque int    `xml:"ns2:codIntegracaoLocalEstoque"`
}

type Epi struct {
	EpiID            int    `xml:"ns2:epiID"`
	CodIntegracaoEpi int    `xml:"ns2:codIntegracaoEpi"`
	NomeEpi          string `xml:"ns2:nomeEpi"`
	Quantidade       int    `xml:"ns2:quantidade"`
	Grade            Grade  `xml:"ns2:grade"`
}

type Grade struct {
	CodIntegracaoGrade int            `xml:"ns2:codIntegracaoGrade"`
	NomeGrade          string         `xml:"ns2:nomeGrade"`
	GradeTamanhos      []GradeTamanho `xml:"ns2:gradeTamanhos>ns2:gradeTamanho"`
}

type GradeTamanho struct {
	CodIntegracaoEpiGradeTamanho int    `xml:"ns2:codIntegracaoEpiGradeTamanho"`
	Tamanho                      string `xml:"ns2:tamanho"`
	Quantidade                   int    `xml:"ns2:quantidade"`
}
