package routers

import (
	"api-persycoins/logar"
	"api-persycoins/models"
	"api-persycoins/query"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// ResponseOK retorna uma mensagem de ok
func ResponseOK(c *gin.Context, log logar.Logfile) {
	c.IndentedJSON(http.StatusOK, "Servidor up")
}

// PostCliente insere um novo cliente no banco
func PostCliente(c *gin.Context, log logar.Logfile, client *dynamodb.Client) {
	var PersyCoins models.PersyCoins
	err := c.BindJSON(&PersyCoins)
	logar.Check(err, log)
	query.InsertCliente(client, PersyCoins, log)
	c.IndentedJSON(http.StatusOK, "Cliente cadastrado com sucesso")
}

// GetSaldo retorna o saldo de um cliente
func GetSaldo(c *gin.Context, log logar.Logfile, client *dynamodb.Client, nome string) {
	saldo := query.GetSaldo(client, nome, log)
	retorno := fmt.Sprintf("%.2f", saldo)
	c.IndentedJSON(http.StatusOK, retorno)
}

func GetSaldoByMail(c *gin.Context, log logar.Logfile, client *dynamodb.Client, mail string) {
	saldo := query.GetSaldoByMail(client, mail, log)
	c.IndentedJSON(http.StatusOK, saldo)
}

// PutSaldo atualiza o saldo de um cliente
func PutSaldo(c *gin.Context, log logar.Logfile, client *dynamodb.Client, nome string, operation string, valor float64) {
	query.UpdateSaldo(log, client, nome, operation, valor)
	c.IndentedJSON(http.StatusOK, "Saldo atualizado com sucesso")
}
