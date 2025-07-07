package routes

import (
	"fmt"
	"inventory-api/controllers"
	"inventory-api/models"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configura todas as rotas da API
func SetupRoutes(db *models.Database) *mux.Router {
	router := mux.NewRouter()

	// Controllers
	contagemController := controllers.NewContagemController(db)

	// Middleware CORS
	router.Use(corsMiddleware)

	// Rotas da API
	api := router.PathPrefix("/api").Subrouter()

	// Rotas de contagem
	api.HandleFunc("/contagens", contagemController.CriarContagem).Methods("POST", "OPTIONS")
	api.HandleFunc("/inventory-counts/{Empresa}", contagemController.ListarContagens).Methods("GET")
	api.HandleFunc("/contagens/{id}", contagemController.BuscarContagem).Methods("GET")
	api.HandleFunc("/inventory-counts/{id}", contagemController.AtualizarContagem).Methods("PUT", "OPTIONS")
	api.HandleFunc("/inventory-counts/{id}", contagemController.ExcluirContagem).Methods("DELETE", "OPTIONS")
	// Rota para buscar descrição de item
	api.HandleFunc("/items/{codigo}/description", contagemController.BuscarDescricaoItem).Methods("GET", "OPTIONS")
	// Rota para validar local
	api.HandleFunc("/locations/{sigla}/{Empresa}/validate", contagemController.ValidarLocal).Methods("GET", "OPTIONS")
	api.HandleFunc("/contagens/finalizar/{Empresa}", contagemController.FinalizarContagem).Methods("POST", "OPTIONS")
	api.HandleFunc("/contagensstate/status/{Empresa}", contagemController.VerificarContagemFinalizada).Methods("GET")
	api.HandleFunc("/getmetragem/{codigo}", contagemController.GetMetragem).Methods("GET", "OPTIONS")

	// Rota de health check
	router.HandleFunc("/health", healthCheck).Methods("GET")

	return router
}

// corsMiddleware adiciona headers CORS para qualquer origem
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir qualquer origem
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Permitir métodos mais comuns
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Permitir todos os headers comuns e personalizados
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// Permitir credenciais se necessário (cookies, etc.)
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Responde imediatamente para preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Segue para o próximo handler
		next.ServeHTTP(w, r)
	})
}

// healthCheck endpoint para verificar se a API está funcionando
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"status": "ok",
		"message": "API de Contagem de Inventário funcionando",
		"version": "1.0.0"
	}`))
}

func PrintRoutes(router *mux.Router) {
	fmt.Println("Rotas registradas:")
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		// Obter métodos
		methods, err := route.GetMethods()
		if err != nil {
			methods = []string{"*"} // Caso não definido, mostra *
		}

		// Obter caminho
		path, err := route.GetPathTemplate()
		if err != nil {
			path = "<sem caminho definido>"
		}

		fmt.Printf("- %s\t%s\n", methods, path)
		return nil
	})
	if err != nil {
		fmt.Printf("Erro ao listar rotas: %v\n", err)
	}
}
