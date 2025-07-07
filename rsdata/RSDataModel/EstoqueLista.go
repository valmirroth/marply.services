package rsdatamodel

import "encoding/xml"

type EnvelopeEst struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soapenv string   `xml:"soapenv,attr"`
	Est     string   `xml:"est,attr"`
	Header  string   `xml:"Header"`
	Body    struct {
		Text               string `xml:",chardata"`
		GetEstoquesRequest struct {
			Text   string `xml:",chardata"`
			Rsdata struct {
				Text   string `xml:",chardata"`
				Versao string `xml:"versao,attr"`
				Config struct {
					Text         string `xml:",chardata"`
					TpVerEstoque string `xml:"tpVerEstoque"`
				} `xml:"config"`
				Estoques struct {
					Text    string `xml:",chardata"`
					Estoque struct {
						Text       string `xml:",chardata"`
						Descricao  string `xml:"descricao"`
						Movimentos struct {
							Text      string `xml:",chardata"`
							Movimento struct {
								Text         string `xml:",chardata"`
								LocalEstoque string `xml:"localEstoque"`
							} `xml:"movimento"`
						} `xml:"movimentos"`
					} `xml:"estoque"`
				} `xml:"estoques"`
			} `xml:"rsdata"`
		} `xml:"getEstoquesRequest"`
	} `xml:"Body"`
}
