package rsdatamodel

import "encoding/xml"

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	SOAPENV string   `xml:"SOAP-ENV,attr"`
	Header  string   `xml:"Header"`
	Body    struct {
		Text                string `xml:",chardata"`
		GetEmpresasResponse struct {
			Text   string `xml:",chardata"`
			Ns2    string `xml:"ns2,attr"`
			Rsdata struct {
				Text     string `xml:",chardata"`
				Empresas struct {
					Text    string `xml:",chardata"`
					Empresa struct {
						Text          string `xml:",chardata"`
						IdEmpresa     string `xml:"idEmpresa"`
						CodIntegracao string `xml:"codIntegracao"`
						NrCNPJEmpresa string `xml:"nrCNPJEmpresa"`
						RazaoSocial   string `xml:"razaoSocial"`
						Denominacao   string `xml:"denominacao"`
						Empregados    struct {
							Text      string `xml:",chardata"`
							Empregado struct {
								Text          string `xml:",chardata"`
								IdEmpregado   string `xml:"idEmpregado"`
								NomeEmpregado string `xml:"nomeEmpregado"`
								Admissao      string `xml:"admissao"`
								Demissao      string `xml:"demissao"`
							} `xml:"empregado"`
						} `xml:"empregados"`
					} `xml:"empresa"`
				} `xml:"empresas"`
			} `xml:"rsdata"`
		} `xml:"getEmpresasResponse"`
	} `xml:"Body"`
}
