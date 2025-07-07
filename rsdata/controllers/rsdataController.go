package controllers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	rsdatamodel "vroth.rsdata/RSDataModel"
	"vroth.rsdata/RSDataTemplate"
	"vroth.rsdata/models"
	"vroth.rsdata/service"
)

var logg *service.Logger

func init() {
	var err error
	logg, err = service.NewLogger("ResultSync", "settings", service.ERROR)
	if err != nil {
		panic(fmt.Errorf("não foi possível inicializar o logger: %w", err))
	}
}

func Index() string {
	return "Hello World"
}

func TesteComunicacaoRSData(w http.ResponseWriter, r *http.Request) {
	logg.Info("TesteComunicacaoRSData", "Validando comunicação com API do RSData")
	w.Write([]byte("teste"))
}

func EfetuarEntradaEstoqueRSData(w http.ResponseWriter, r *http.Request) {
	logg.Info("EfetuarEntradaEstoqueRSData", "Validando comunicação com API do RSData")

	w.Write([]byte("EfetuarEntradaEstoqueRSData"))
}

func MovimentarEstoqueRsData(dadosMov *RSDataTemplate.EstoqueRsDataModel) (*rsdatamodel.EnvelopeEstRetorno, error) {
	var temp *template.Template
	temp = template.Must(template.ParseFiles("./../../RSDataTemplate/InsertMovEstoque.tmpl"))
	var param bytes.Buffer

	err := temp.Execute(&param, dadosMov)
	log.Printf("Tamanho gerado: %d bytes\n", param.Len())

	if err := os.WriteFile("output.txt", param.Bytes(), 0644); err != nil {
		log.Fatalf("gravação falhou: %v", err)
	}
	fmt.Println("gravado em output.txt")
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	fmt.Println("inicio...")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("inicio...")
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest("POST", "https://api.rsdata.com.br/epiEstoqueService/soapws", bytes.NewBuffer(param.Bytes()))
	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", "")
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer response.Body.Close()

	rs := &rsdatamodel.EnvelopeEstRetorno{}

	fmt.Print(string(bodyBytes))
	xml.Unmarshal(bodyBytes, rs)

	return rs, err
}

func GetBaixasEstoqueSetorRsData(dadosMov *RSDataTemplate.EstoqueRsDataModel, korp *service.EntradaNFService) error {
	var temp *template.Template
	temp = template.Must(template.ParseFiles("./../../RSDataTemplate//GetMovEstoque.tmpl"))
	var param bytes.Buffer

	err := temp.Execute(&param, dadosMov)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(string(param.Bytes()))
	//fmt.Print(string(param.Bytes()))
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest("POST", "https://api.rsdata.com.br/epiEstoqueService/soapws", bytes.NewBuffer(param.Bytes()))
	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", "")
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(bodyBytes))
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	//	rs := &rsdatamodel.RsGetEstoqueEnvelope{}
	x := &rsdatamodel.NewEnvelope{}
	//fmt.Print(string(bodyBytes))
	xml.Unmarshal(bodyBytes, x)
	processarListaItensBaixaSetorRSData(x, dadosMov, korp)
	return nil
}

func processarListaItensBaixaSetorRSData(x *rsdatamodel.NewEnvelope, dadosMov *RSDataTemplate.EstoqueRsDataModel, korp *service.EntradaNFService) {
	for _, setor := range x.Body.GetEntregasBySetorResponse.Dados.Empresa.Setores.Setor {
		for _, mov := range setor.Entregas.Entrega {

			var his models.EntregaEPISetor
			if mov.CodIntegracaoEpi == "80751" {
				logg.Error("CodIntegracaoEpi vazio para o EPI: %s", mov.NomeEpi)
				//	continue
			}
			his.DataMovimento = dadosMov.DataIni
			if dadosMov.Cnpj == "05.552.102/0005-67" {
				his.EmpresaRecno = 5
				his.Local = "20"
			} else {
				his.EmpresaRecno = 1
				his.Local = "5"
			}
			his.Quantidade = strconv.Itoa(mov.Quantidade)
			his.CodIntegracaoEpi = mov.CodIntegracaoEpi
			his.NomeEpi = mov.NomeEpi
			his.NomeSetor = setor.NomeSetor
			his.IdSetor = setor.CdSetor
			his.TipoEpi = mov.TipoEpi
			his.Cnpj = dadosMov.Cnpj
			his.Tamanho = mov.Tamanho
			his.NomeEpi = mov.NomeEpi
			mov, errom := korp.GravarMoviementacaoEPISetor(his)
			if errom != nil {
				fmt.Println(errom)
				fmt.Println(mov)
			}
		}
	}
}
