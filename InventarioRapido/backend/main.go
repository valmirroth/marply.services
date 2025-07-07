package main

import (
	"fmt"
	"inventory-api/models"
	"inventory-api/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado, usando variáveis do sistema")
	}

	// Conectar ao banco de dados
	db, err := models.NewDatabase()
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer db.Close()

	// Configurar rotas
	router := routes.SetupRoutes(db)

	// Configurar servidor
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	routes.PrintRoutes(router)

	/*
		log.Printf("🚀 API iniciada na porta %s", port)
		log.Printf("📋 Documentação: http://localhost:%s/health", port)
		log.Printf("🔗 Endpoints disponíveis:")
		log.Printf("   GET    /api/inventory-counts - Listar contagens")
		log.Printf("   POST   /api/inventory-counts - Criar contagem")
		log.Printf("   GET    /api/contagens/{id} - Buscar contagem por ID")
		log.Printf("   PUT    /api/contagens/{id} - Atualizar contagem")
		log.Printf("   DELETE /api/contagens/{id} - Excluir contagem")
	*/
	// Iniciar servidor
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
