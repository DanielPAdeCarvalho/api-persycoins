package logar

import (
	"log"
)

type Logfile struct {
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
}

// Logar é uma função que escreve uma mensagem de log
func Check(erro error, logar Logfile) {
	if erro != nil {
		logar.ErrorLogger.Fatal(erro)
	}
}
