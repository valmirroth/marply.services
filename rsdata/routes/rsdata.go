package routes

import (
	"net/http"

	"vroth.rsdata/controllers"
)

var rsdataRotas = []Rota{
	{
		URI:                "/ListarOrdensPainel",
		Metodo:             http.MethodGet,
		Funcao:             controllers.TesteComunicacaoRSData,
		RequerAutenticacao: true,
	},
	{URI: "/ListarOrdensPainel",
		Metodo:             http.MethodGet,
		Funcao:             controllers.TesteComunicacaoRSData,
		RequerAutenticacao: true,
	},
}
