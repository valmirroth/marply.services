package rsdatamodel

import "encoding/xml"

type FuncionarioEnvelope struct {
	XMLName xml.Name  `xml:"Envelope"`
	Text    string    `xml:",chardata"`
	SOAPENV string    `xml:"SOAP-ENV,attr"`
	Header  string    `xml:"Header"`
	Body    RetornoRS `xml:"Body"`
}

type RetornoRS struct {
	Text                  string      `xml:",chardata"`
	GetEmpregadosResponse RSResponseF `xml:"getEmpregadosResponse"`
}
type RSResponseF struct {
	Text   string    `xml:",chardata"`
	Ns2    string    `xml:"ns2,attr"`
	Rsdata Resultado `xml:"rsdata"`
}

type Resultado struct {
	Text   string `xml:",chardata"`
	Versao string `xml:"versao,attr"`
	Config struct {
		Text           string `xml:",chardata"`
		TpVerEmpregado string `xml:"tpVerEmpregado"`
	} `xml:"config"`
	Empregados RSEmpregados `xml:"empregados"`
}

type RSEmpregados struct {
	Text      string      `xml:",chardata"`
	Empregado Funcionario `xml:"empregado"`
}

type Funcionario struct {
	Text                    string `xml:",chardata"`
	IdEmpresa               string `xml:"idEmpresa"`
	CodIntegracaoEmpresa    string `xml:"codIntegracaoEmpresa"`
	NrCNPJEmpresa           string `xml:"nrCNPJEmpresa"`
	IdEmpregado             string `xml:"idEmpregado"`
	CodIntegracaoEmpregado  string `xml:"codIntegracaoEmpregado"`
	NomeEmpregado           string `xml:"nomeEmpregado"`
	CarteiraTrabalhoDigital string `xml:"carteiraTrabalhoDigital"`
	DtNascimento            string `xml:"dtNascimento"`
	TpSexo                  string `xml:"tpSexo"`
	NrCTPS                  string `xml:"nrCTPS"`
	NrSerieCTPS             string `xml:"nrSerieCTPS"`
	DtEmiCTPS               string `xml:"dtEmiCTPS"`
	UfEmiCTPS               string `xml:"ufEmiCTPS"`
	NrIdentidade            string `xml:"nrIdentidade"`
	OrgaoExpedidorRG        string `xml:"orgaoExpedidorRG"`
	DtEmiRG                 string `xml:"dtEmiRG"`
	UfEmiRG                 string `xml:"ufEmiRG"`
	NrNIT                   string `xml:"nrNIT"`
	NrCPF                   string `xml:"nrCPF"`
	NrMatricula             string `xml:"nrMatricula"`
	MatriculaRh             string `xml:"matriculaRh"`
	CategoriaTrabalhador    string `xml:"categoriaTrabalhador"`
	TpVinculo               string `xml:"tpVinculo"`
	BRPDH                   string `xml:"BR_PDH"`
	RegimeRevezamento       string `xml:"regimeRevezamento"`
	DtAdmissao              string `xml:"dtAdmissao"`
	EnderecoEmpregado       string `xml:"enderecoEmpregado"`
	CidadeEmpregado         string `xml:"cidadeEmpregado"`
	CidadeCodIbge           string `xml:"cidadeCodIbge"`
	BairroEmpregado         string `xml:"bairroEmpregado"`
	EstadoEmpregado         string `xml:"estadoEmpregado"`
	NrCEP                   string `xml:"nrCEP"`
	NrTelefone              string `xml:"nrTelefone"`
	NomeMae                 string `xml:"nomeMae"`
	TpFilPrevidencia        string `xml:"tpFilPrevidencia"`
	TpEstadoCivil           string `xml:"tpEstadoCivil"`
	TpAposentado            string `xml:"tpAposentado"`
	NrEleitor               string `xml:"nrEleitor"`
	Anulado                 string `xml:"anulado"`
	SetoresCargos           Setor  `xml:"setoresCargos"`
}

type Setor struct {
	Text       string  `xml:",chardata"`
	SetorCargo []Cargo `xml:"setorCargo"`
}
type Cargo struct {
	Text              string `xml:",chardata"`
	DtInicio          string `xml:"dtInicio"`
	CdSetor           string `xml:"cdSetor"`
	NomeSetor         string `xml:"nomeSetor"`
	CdCargo           string `xml:"cdCargo"`
	NomeCargo         string `xml:"nomeCargo"`
	CargoDesenvolvido string `xml:"cargoDesenvolvido"`
	CargoCBO          string `xml:"cargoCBO"`
	DtSaida           string `xml:"dtSaida"`
}
