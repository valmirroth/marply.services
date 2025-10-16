package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"bh-mvc/internal/config"

	"github.com/robfig/cron/v3"
)

// dispara um POST para /run no próprio servidor
func striggerRunOnce(ctx context.Context, baseURL string) error {
	reqBody, _ := json.Marshal(map[string]any{}) // body vazio
	fmt.Println("vai chamar o agendador de tarefas...")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://192.168.1.28:8093/run", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	fmt.Println("rodando pelo agendador...")
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("POST /run retornou status %d", resp.StatusCode)
	}
	return nil
}

func sAgenda() {
	// Carrega cfg só para saber a porta; não altera o main existente
	cfg := config.Load()

	// Monta a baseURL do próprio servidor (ex.: http://127.0.0.1:8080)
	// cfg.HTTPPort costuma vir no formato ":8080"
	baseURL := "http://127.0.0.1" + cfg.HTTPPort

	// Usa timezone de São Paulo
	var loc *time.Location

	loc = time.FixedZone("BRT", -3*60*60)

	fmt.Println(loc.String())
	log.Writer().Write([]byte(loc.String()))

	c := cron.New(
		cron.WithLocation(loc), // agenda no fuso correto
		cron.WithSeconds(),     // permite campo de segundos no spec
		cron.WithChain( // logs básicos de erro
			cron.Recover(cron.DefaultLogger),
		),
	)

	// “0 0 3 * * *” => todos os dias às 03:00:00
	_, err := c.AddFunc("0 28 17 * * *", func() {
		log.Println("[SCHED] Disparando /run (03:00)…")
		// contexto com timeout generoso para o ETL
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := triggerRunOnce(ctx, baseURL); err != nil {
			log.Printf("[SCHED] Erro ao chamar /run: %v", err)
			return
		}
		log.Println("[SCHED] /run concluído com sucesso.")
	})
	if err != nil {
		log.Printf("[SCHED] Erro ao registrar cron: %v", err)
		return
	}

	// Inicia o cron em background. Não bloqueia o main existente.
	c.Start()

	log.Println("[SCHED] Agendador ativo: todo dia às 03:00 (America/Sao_Paulo)")
}
