package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"inventory-api/models"

	"github.com/gorilla/mux"
)

type ContagemController struct {
	DB *models.Database
}

// NewContagemController cria um novo controller
func NewContagemController(db *models.Database) *ContagemController {
	return &ContagemController{DB: db}
}

// CriarContagem handler para POST /api/contagens
func (cc *ContagemController) CriarContagem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var contagemReq models.ContagemRequest

	// Decodificar JSON do body
	if err := json.NewDecoder(r.Body).Decode(&contagemReq); err != nil {
		response := models.ContagemResponse{
			Success: false,
			Message: "Erro ao decodificar JSON",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validar campos obrigatórios
	if err := cc.validarContagemRequest(contagemReq); err != nil {
		response := models.ContagemResponse{
			Success: false,
			Message: "Dados inválidos",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Inserir no banco
	contagem, err := cc.DB.InserirContagem(contagemReq)
	if err != nil {
		log.Printf("Erro ao inserir contagem: %v", err)
		response := models.ContagemResponse{
			Success: false,
			Message: "Erro interno do servidor",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Resposta de sucesso
	response := models.ContagemResponse{
		Success: true,
		Message: "Contagem criada com sucesso",
		Data:    contagem,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListarContagens handler para GET /api/contagens
func (cc *ContagemController) ListarContagens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	id := vars["Empresa"]
	if id == "" {
		response := models.ContagemListResponse{
			Success: false,
			Message: "Selecionar uma empresa",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	contagens, err := cc.DB.ListarContagens(id)
	if err != nil {
		log.Printf("Erro ao listar contagens: %v", err)
		response := models.ContagemListResponse{
			Success: false,
			Message: "Erro ao buscar contagens",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ContagemListResponse{
		Success: true,
		Message: "Contagens listadas com sucesso",
		Data:    contagens,
		Total:   len(contagens),
	}
	json.NewEncoder(w).Encode(response)
}

// BuscarContagem handler para GET /api/contagens/{id}
func (cc *ContagemController) BuscarContagem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		response := models.ContagemResponse{
			Success: false,
			Message: "ID é obrigatório",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	contagem, err := cc.DB.BuscarContagemPorID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "não encontrada") {
			status = http.StatusNotFound
		}

		response := models.ContagemResponse{
			Success: false,
			Message: "Erro ao buscar contagem",
			Error:   err.Error(),
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ContagemResponse{
		Success: true,
		Message: "Contagem encontrada",
		Data:    contagem,
	}
	json.NewEncoder(w).Encode(response)
}

// AtualizarContagem handler para PUT /api/contagens/{id}
func (cc *ContagemController) AtualizarContagem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		response := models.ContagemResponse{
			Success: false,
			Message: "ID é obrigatório",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var contagemReq models.ContagemRequest

	if err := json.NewDecoder(r.Body).Decode(&contagemReq); err != nil {
		response := models.ContagemResponse{
			Success: false,
			Message: "Erro ao decodificar JSON",
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	contagem, err := cc.DB.AtualizarContagem(id, contagemReq)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "não encontrada") {
			status = http.StatusNotFound
		}

		response := models.ContagemResponse{
			Success: false,
			Message: "Erro ao atualizar contagem",
			Error:   err.Error(),
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ContagemResponse{
		Success: true,
		Message: "Contagem atualizada com sucesso",
		Data:    contagem,
	}
	json.NewEncoder(w).Encode(response)
}

// ExcluirContagem handler para DELETE /api/contagens/{id}
func (cc *ContagemController) ExcluirContagem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		response := models.ContagemResponse{
			Success: false,
			Message: "ID é obrigatório",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := cc.DB.ExcluirContagem(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "não encontrada") {
			status = http.StatusNotFound
		}

		response := models.ContagemResponse{
			Success: false,
			Message: "Erro ao excluir contagem",
			Error:   err.Error(),
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ContagemResponse{
		Success: true,
		Message: "Contagem excluída com sucesso",
	}
	json.NewEncoder(w).Encode(response)
}

// validarContagemRequest valida os campos obrigatórios
func (cc *ContagemController) validarContagemRequest(req models.ContagemRequest) error {
	if strings.TrimSpace(req.CodigoItem) == "" {
		return fmt.Errorf("código do item é obrigatório")
	}
	if strings.TrimSpace(req.DescricaoItem) == "" {
		return fmt.Errorf("descrição do item é obrigatória")
	}
	if strings.TrimSpace(req.Local) == "" {
		return fmt.Errorf("local é obrigatório")
	}
	if strings.TrimSpace(req.Quantidade) == "" {
		return fmt.Errorf("quantidade é obrigatória")
	}
	volumes, errv := strconv.ParseFloat(req.Volumes, 64)
	if errv != nil {
		return fmt.Errorf("erro ao converter volumes: %w", errv)
	}

	if volumes < 0 {
		return fmt.Errorf("volumes deve ser maior ou igual a zero")
	}
	return nil
}

// BuscarDescricaoItem handler para GET /api/items/{codigo}/description
func (cc *ContagemController) GetMetragem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	codigo := vars["codigo"]

	if codigo == "" {
		response := map[string]interface{}{
			"success": false,
			"message": "Código do item é obrigatório",
			"error":   "Código não fornecido",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	item, err := cc.DB.BuscarDescricaoItem(codigo)
	if err != nil {
		log.Printf("Erro ao buscar descrição do item %s: %v", codigo, err)
		response := map[string]interface{}{
			"success": false,
			"message": "Erro ao buscar descrição do item",
			"error":   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Descrição encontrada com sucesso",
		"data":    item,
	}
	json.NewEncoder(w).Encode(response)
}

// BuscarDescricaoItem handler para GET /api/items/{codigo}/description
func (cc *ContagemController) BuscarDescricaoItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	codigo := vars["codigo"]

	if codigo == "" {
		response := map[string]interface{}{
			"success": false,
			"message": "Código do item é obrigatório",
			"error":   "Código não fornecido",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	item, err := cc.DB.BuscarDescricaoItem(codigo)
	if err != nil {
		log.Printf("Erro ao buscar descrição do item %s: %v", codigo, err)
		response := map[string]interface{}{
			"success": false,
			"message": "Erro ao buscar descrição do item",
			"error":   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Descrição encontrada com sucesso",
		"data":    item,
	}
	json.NewEncoder(w).Encode(response)
}

// ValidarLocal handler para GET /api/locations/{sigla}/validate
func (cc *ContagemController) ValidarLocal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	sigla := vars["sigla"]
	empresa := vars["Empresa"]

	if sigla == "" {
		response := map[string]interface{}{
			"success": false,
			"message": "Sigla do local é obrigatória",
			"error":   "Sigla não fornecida",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	local, err := cc.DB.ValidarLocal(sigla, empresa)
	if err != nil {
		log.Printf("Erro ao validar local %s: %v", sigla, err)

		status := http.StatusNotFound
		if !strings.Contains(err.Error(), "não encontrado") {
			status = http.StatusInternalServerError
		}

		response := map[string]interface{}{
			"success": false,
			"message": "Local não encontrado",
			"error":   err.Error(),
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Local válido",
		"data":    local,
	}
	json.NewEncoder(w).Encode(response)
}

// FinalizarContagem handler para POST /api/contagens/finalizar
func (cc *ContagemController) FinalizarContagem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	vars := mux.Vars(r)
	empresa := vars["Empresa"]
	err := cc.DB.FinalizarContagem(empresa)
	if err != nil {
		log.Printf("Erro ao finalizar contagem: %v", err)

		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "nenhuma contagem ativa") {
			status = http.StatusBadRequest
		}

		response := map[string]interface{}{
			"success": false,
			"message": "Erro ao finalizar contagem",
			"error":   err.Error(),
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Contagem finalizada com sucesso",
	}
	json.NewEncoder(w).Encode(response)
}

// VerificarContagemFinalizada handler para GET /api/contagens/status
func (cc *ContagemController) VerificarContagemFinalizada(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	empresa := vars["Empresa"]

	isFinalizada, err := cc.DB.VerificarContagemFinalizada(empresa)
	if err != nil {
		log.Printf("Erro ao verificar status da contagem: %v", err)
		response := map[string]interface{}{
			"success": false,
			"message": "Erro ao verificar status da contagem",
			"error":   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]bool{
			"finalizada": isFinalizada,
		},
	}
	json.NewEncoder(w).Encode(response)
}
