package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"vroth.rsdata/controllers"
	"vroth.rsdata/routes"
	"vroth.rsdata/service"
	"vroth.rsdata/settings"
)

func main() {
	logger, err := service.NewLogger("ResultSync", "app", service.DEBUG)
	if err != nil {
		panic(err)
	}
	fmt.Println("Iniciando processamento de baixas do RSData para Korp..." + settings.GetInstange().SelectedDB)

	logger.Info("Aplicação iniciada às %s", time.Now().Format(time.RFC3339))
	logger.Debug("Valor de x: %d", 42)
	logger.Warn("Cuidado: %s", "algo inesperado")
	logger.Error("Erro inesperado: %v", err)
	rotas := mux.NewRouter()
	routes.ConfigureRoutes(rotas)
	controllers.ProcessarBaixasRsDataToKorp()
	//controllers.ProcessarEntradaNfKorpToRsData()
	http.ListenAndServe(":8095", rotas)
}
