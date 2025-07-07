package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"

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

	// --- Router HTTP
	rotas := mux.NewRouter()
	routes.ConfigureRoutes(rotas)

	// --- Timezone BR (evita rodar 03:00 em UTC por engano)
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		logger.Warn("Não foi possível carregar timezone America/Sao_Paulo; usando local padrão: %v", err)
		loc = time.Local
	}

	// --- Scheduler diário às 03:00
	// Especificação do cron (minuto hora dia-do-mês mês dia-da-semana)
	// "0 3 * * *" => 03:00 todos os dias
	c := cron.New(cron.WithLocation(loc))
	_, err = c.AddFunc("0 3 * * *", func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Pânico durante job de 03:00: %v", r)
			}
		}()

		agora := time.Now().In(loc).Format(time.RFC3339)
		logger.Info("Job diário iniciado às %s", agora)

		// chame aqui as rotinas necessárias
		controllers.ProcessarBaixasRsDataToKorp()
		controllers.ProcessarEntradaNfKorpToRsData()

		logger.Info("Job diário finalizado às %s", time.Now().In(loc).Format(time.RFC3339))
	})
	if err != nil {
		logger.Error("Falha ao registrar job no cron: %v", err)
	}

	// (Opcional) Executa uma vez na inicialização
	// comente se não quiser rodar na subida
	go func() {
		logger.Info("Execução inicial ao subir a aplicação")
		//	controllers.ProcessarBaixasRsDataToKorp()
		//controllers.ProcessarEntradaNfKorpToRsData()
	}()
	logger.Info("Execução inicial ao subir a aplicação")
	controllers.ProcessarBaixasRsDataToKorp()
	// inicia o agendador
	c.Start()
	logger.Info("Scheduler iniciado (rodará diariamente às 03:00 em %s)", loc)

	// --- Sobe o servidor HTTP
	if err := http.ListenAndServe(":8095", rotas); err != nil {
		logger.Error("Erro ao subir servidor HTTP: %v", err)
	}
}
