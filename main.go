package main

import (
	"api-persycoins/driver"
	"api-persycoins/logar"
	"api-persycoins/routers"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var (
	dynamoClient *dynamodb.Client
	logs         logar.Logfile
)

func inLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func setupRouter() *gin.Engine {
	apiRouter := gin.Default()

	apiRouter.GET("/", func(ctx *gin.Context) {
		logs.InfoLogger.Println("Servidor Ok")
		routers.ResponseOK(ctx, logs)
	})

	apiRouter.GET("/:nome", func(ctx *gin.Context) {
		Nome := ctx.Param("nome")
		routers.GetSaldo(ctx, logs, dynamoClient, Nome)
	})

	apiRouter.POST("/newclient", func(ctx *gin.Context) {
		routers.PostCliente(ctx, logs, dynamoClient)
	})

	return apiRouter
}

// Para compilar o binario do sistema usamos:
//
//	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o api-persycoins .
//
// para criar o zip do projeto comando:
//
// zip lambda.zip api-persycoins
//
// main.go
func main() {
	InfoLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)

	logs.InfoLogger = InfoLogger
	logs.ErrorLogger = ErrorLogger
	var err error
	// chamada de função para a criação da sessao de login com o banco
	dynamoClient, err = driver.ConfigAws()
	//chamada da função para revificar o erro retornado
	logar.Check(err, logs)

	if inLambda() {
		log.Fatal(gateway.ListenAndServe(":8080", setupRouter()))
	} else {
		log.Fatal(http.ListenAndServe(":8080", setupRouter()))
	}
}
