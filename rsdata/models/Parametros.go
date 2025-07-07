package models

type Parametros struct {
	Vspay struct {
		Servicename string `json:"ServiceName"`
		User        string `json:"User"`
		Key         string `json:"Key"`
	} `json:"VSPay"`
	ContasP struct {
		CtOperacao    string `json:"CtOperacao"`
		CentroDeCusto string `json:"CentroDeCusto"`
		PlanoGr       string `json:"PlanoGr"`
		EmpresaRecno  int64  `json:"EmpresaRecno"`
		TipoPagamento int    `json:"TipoPagamento"`
		TipoDocumento int    `json:"TipoDocumento"`
		SituacaoPG    int    `json:"SituacaoPG"`
	} `json:"ContasPagar"`

	Dbservice struct {
		Sqlserveraddress  string `json:"SQLServerAddress"`
		Sqlserveruser     string `json:"SqlServerUser"`
		Accessport        string `json:"AccessPort"`
		Sqlserverpassword string `json:"SqlServerPassword"`
		Database          string `json:"DataBase"`
	} `json:"DBService"`
	SMTP struct {
		Enviarnotificacaoemail string `json:"EnviarNotificacaoEmail"`
		Emailadress            string `json:"EmailAdress"`
		EmailUserAutentication string `json:"EmailUserAutentication"`
		Emailpassword          string `json:"EmailPassword"`
		Emailsmtpserver        string `json:"EmailSMTPServer"`
		Emailport              int    `json:"EmailPort"`
		Destinatario           string
	} `json:"SMTP"`
	Vstoken string `json:"VSToken"`
	//KorpResourses        KorpService `json:"KorpService"`
	KorpUser             string
	KorpPass             string
	JWT                  string
	SelectedDB           string
	CurrentLogFile       string
	IntervaloSincronismo int
}
