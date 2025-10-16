package controllers

import (
	"fmt"
	"time"

	"vroth.rsdata/RSDataTemplate"
	sqlRoth "vroth.rsdata/database"
	"vroth.rsdata/repositories"
	"vroth.rsdata/service"
	"vroth.rsdata/settings"
)

func limparBaixasDiaAnterior(korp *service.EntradaNFService) error {
	// Limpa os dados de baixa do dia anterior
	hojes := time.Now()
	for i := -2; i <= 0; i++ {
		z := i
		diaAnterior := hojes.AddDate(0, 0, z)
		fmt.Println("Limpando registro referentes a: " + diaAnterior.Format("02/01/2006"))
		erro := korp.LimparDadosBaixaEstoqueSetor(diaAnterior.Format("02/01/2006"))
		if erro != nil {
			fmt.Printf("Erro ao limpar dados de baixa do dia %s: %v\n", diaAnterior.Format("02/01/2006"), erro)
			return erro
		} else {
			fmt.Printf("Dados de baixa do dia %s limpos com sucesso.\n", diaAnterior.Format("02/01/2006"))
		}
	}
	return nil
}

func ProcessarBaixasRsDataToKorp() {
	fmt.Println("Iniciando processamento de baixas do RSData para Korp..." + settings.GetInstange().SelectedDB)
	db := sqlRoth.GetDatabase(settings.GetInstange().SelectedDB)
	repo := repositories.NewEntradaNFRepository(db)
	korp := service.NewEntradaNFService(db, repo)
	errolimparbaixa := limparBaixasDiaAnterior(korp)
	if errolimparbaixa != nil {
		logg.Error("Erro ao limpar dados de baixa: %v", errolimparbaixa)
		return
	}

	// Processar lista de itens para consultar baixa
	processarListaItensBaixa(korp)
	fmt.Println("Processamento concluído.")
	return
}

func setParamsGetBaixas(cnpj string, dest *RSDataTemplate.EstoqueRsDataModel) {
	dest.Cnpj = cnpj
	dest.TpMovimento = "ENTRADA"
	dest.TpClasseProduto = "EPI"
	dest.Usuario = "integracao.db0922mariniind@rsdata.com.br"
	dest.Senha = "21753ef24af23bd6cc64435a1a4b4137"
}

func processarListaItensBaixa(korp *service.EntradaNFService) {
	var cnpjs = []string{
		"05.552.102/0001-33", // Marini Indústria e Comércio de EPI Ltda
		"05.552.102/0005-67", // Marini Indústria e Comércio de EPI Ltda - Filial
	}
	for _, cnpj := range cnpjs {
		var lancamento RSDataTemplate.EstoqueRsDataModel
		setParamsGetBaixas(cnpj, &lancamento)
		lancamento.Cnpj = cnpj
		hoje := time.Now()
		diaAnterior := hoje.AddDate(0, 0, -1)
		fmt.Printf("Buscando movimentos de estoque do dia %s\n", diaAnterior.Format("02/01/2006"))
		for i := -2; i <= 0; i++ {
			z := i
			diaAnterior := hoje.AddDate(0, 0, z)
			diaAnteriorMaisUm := hoje.AddDate(0, 0, (z + 1))
			lancamento.DataIni = diaAnterior.Format("02/01/2006")
			lancamento.DataFim = diaAnteriorMaisUm.Format("02/01/2006")
			fmt.Println(lancamento.DataIni)
			GetBaixasEstoqueSetorRsData(&lancamento, korp)
		}
	}
}
