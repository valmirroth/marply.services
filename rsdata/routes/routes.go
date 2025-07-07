package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Rota struct {
	URI                string
	Metodo             string
	Funcao             func(http.ResponseWriter, *http.Request)
	RequerAutenticacao bool
}

func ConfigureRoutes(router *mux.Router) *mux.Router {
	var rotas = rsdataRotas
	for _, rota := range rotas {
		router.HandleFunc(rota.URI, rota.Funcao).Methods(rota.Metodo)
	}
	return router
}
