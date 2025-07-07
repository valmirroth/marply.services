package models

import (
	"time"
)

// ContagemLocal representa um registro de contagem de inventário
type ContagemLocal struct {
	ID              string    `json:"id" db:"ID"`
	CodigoItem      string    `json:"codigo_item" db:"CODIGO_ITEM"`
	DescricaoItem   string    `json:"descricao_item" db:"DESCRICAO_ITEM"`
	Local           string    `json:"local" db:"LOCAL"`
	Quantidade      float64   `json:"quantidade" db:"QUANTIDADE"`
	Volumes         float64   `json:"volumes" db:"VOLUMES"`
	DataContagem    time.Time `json:"data_contagem" db:"DATA_CONTAGEM"`
	UsuarioContagem string    `json:"usuario_contagem" db:"USUARIO_CONTAGEM"`
	Status          string    `json:"status" db:"STATUS"`
	QtdeEmbalagem   string    `json:"Metragem" db:"Metragem"`
	NroContagem     string    `json:"nro_contagem" db:"NRO_CONTAGEM"`
	Empresa         int       `json:"empresa" db:"EMPRESA"`
}

// ContagemRequest representa os dados recebidos para criar uma contagem
type ContagemRequest struct {
	CodigoItem      string `json:"itemCode" validate:"required"`
	DescricaoItem   string `json:"itemDescription" validate:"required"`
	Local           string `json:"location" validate:"required"`
	Quantidade      string `json:"quantity" validate:"required"`
	Volumes         string `json:"volumeCount" validate:"required,min=0"`
	UsuarioContagem string `json:"usuario_contagem,omitempty"`
	Empresa         int    `json:"empresa" db:"EMPRESA"`
}

// ContagemResponse representa a resposta padronizada da API
type ContagemResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    *ContagemLocal `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// ContagemListResponse representa a resposta com lista de contagens
type ContagemListResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    []ContagemLocal `json:"data,omitempty"`
	Total   int             `json:"total"`
	Error   string          `json:"error,omitempty"`
}

type LocalInfo struct {
	Codigo    string `json:"codigo"`
	Sigla     string `json:"sigla"`
	Descricao string `json:"descricao"`
}

// ItemInfo representa as informações de um item
type ItemInfo struct {
	Codigo             string  `json:"codigo"`
	Descricao          string  `json:"descricao"`
	QtdeEmbalagem      float64 `json:"qtde_embalagem"`
	CalcularAutomatico bool    `json:"calcular_automatico"`
}
