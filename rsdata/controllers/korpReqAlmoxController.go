package controllers

import (
	"fmt"
	"strconv"

	rsdatamodel "vroth.rsdata/RSDataModel"
	"vroth.rsdata/RSDataTemplate"
	sqlRoth "vroth.rsdata/database"
	"vroth.rsdata/models"
	"vroth.rsdata/repositories"
	"vroth.rsdata/service"
	"vroth.rsdata/settings"
)

func GetMesAnoDia(data string) string {
	runes := []rune(data)
	ano := runes[:4]
	mes := runes[4:6]
	dia := runes[6:8]

	novadata := string(ano) + "-" + string(mes) + "-" + string(dia)
	return novadata
}

func EstoqueKorpToRSData(lc *models.RSDataLancamento, dest *RSDataTemplate.EstoqueRsDataModel) {

	dest.CodInterno = lc.Item
	dest.Cnpj = lc.Cnpj
	dest.LocalEstoque, _ = strconv.Atoi(lc.LocalEstoque)
	dest.Quantidade = lc.QtdMovimento
	dest.VlUnitario = lc.ValUni
	dest.CodGrade = lc.CodGradeEPI
	dest.TamanhoGrade = lc.TamanhoGrade
	dest.Descricao = lc.Descricao
	dest.TpMovimento = "ENTRADA"
	dest.TpClasseProduto = "EPI"
	dest.DtMovimento = GetMesAnoDia(lc.DataMovimento)

	dest.Usuario = "integracao.db0922mariniind@rsdata.com.br"
	dest.Senha = "21753ef24af23bd6cc64435a1a4b4137"
}

func inserirRegistrosRSData(lc *models.RSDataLancamento) (*rsdatamodel.EnvelopeEstRetorno, error) {
	var lancamento RSDataTemplate.EstoqueRsDataModel
	EstoqueKorpToRSData(lc, &lancamento)
	retorno, erro := MovimentarEstoqueRsData(&lancamento)
	return retorno, erro
}

func ProcessarEntradaNfKorpToRsData() {
	db := sqlRoth.GetDatabase(settings.GetInstange().SelectedDB)

	repo := repositories.NewEntradaNFRepository(db)
	korp := service.NewEntradaNFService(db, repo)

	lista, errogetkorp := korp.GetNextLancRSData("54919")
	if errogetkorp != nil {
		fmt.Println(errogetkorp)
	}

	if len(lista) == 0 {
		fmt.Println("Nenhuma nota de entrada para lançamento no RSData")
	}

	for _, lanc := range lista {
		retorno, erro := inserirRegistrosRSData(lanc)
		if erro != nil {
			fmt.Println(retorno)
			fmt.Println(lanc)
		}
		korp.InserirLogMovimentoRsData(lanc, retorno)
	}
}
