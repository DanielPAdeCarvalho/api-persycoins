package models

type PersyCoins struct {
	Nome      string  `json:"Nome"`
	Sobrenome string  `json:"sobrenome"`
	Email     string  `json:"email"`
	Saldo     float64 `json:"saldo"`
}
