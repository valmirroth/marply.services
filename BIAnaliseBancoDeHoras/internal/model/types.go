package model

import "time"

// Linha lida da origem (SQL ou CSV)
type RawRow struct {
	NumEmp           int
	CodBh            int
	NumCad           int
	CodCcu           string
	CodFil           int
	SitAfa           string
	CodCal           int
	PerRef           string
	DtApuracao       *time.Time
	Cracha           string
	Colaborador      string
	TipCol           int
	DesSit           string
	CodSit           int
	ValHoraMes       float64
	ValHoraCalculado float64
	Horas            float64
	ValorReais       float64
}

// Linha processada para escrita no detalhado
type ProcessedRow struct {
	NumEmp                    int
	CodBh                     int
	NumCad                    int
	SitAfa                    string
	CodCcu                    string
	CodFil                    int
	CodCal                    int
	PerRef                    string
	DtApuracao                *time.Time
	Cracha                    string
	Colaborador               string
	TipCol                    int
	DesSit                    string
	CodSit                    int
	ValHoraMes                float64
	ValHoraCalculado          float64
	HorasOriginal             float64
	BancoUsadoNaLinha         float64
	HorasSaldo                float64
	ValorSaldo                float64
	ValorReais                float64
	BancoTotalAplicadoNoGrupo float64
}
