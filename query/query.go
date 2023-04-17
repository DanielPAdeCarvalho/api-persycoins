package query

import (
	"api-persycoins/logar"
	"api-persycoins/models"
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// InsertCliente insere um novo cliente no banco
func InsertCliente(client *dynamodb.Client, PersyCoins models.PersyCoins, log logar.Logfile) {
	coins, err := attributevalue.MarshalMap(PersyCoins)
	logar.Check(err, log)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("PersyCoins"),
		Item:      coins,
	}
	_, err = client.PutItem(context.Background(), input)
	logar.Check(err, log)
}

// GetSaldo retorna o saldo de um cliente
func GetSaldo(client *dynamodb.Client, nome string, log logar.Logfile) float64 {
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"Nome": nome,
	})
	logar.Check(err, log)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("PersyCoins"),
		Key:       key,
	}

	// Call the GetItem method with the input
	resp, err := client.GetItem(context.TODO(), input)
	logar.Check(err, log)
	if resp.Item == nil {
		log.ErrorLogger.Println("Cliente não encontrado")
	}

	// Get the Saldo attribute value
	saldoValue, ok := resp.Item["Saldo"]
	if !ok {
		log.ErrorLogger.Println("Saldo não encontrado")
	}

	saldoN := saldoValue.(*types.AttributeValueMemberN)
	saldo, err := strconv.ParseFloat(saldoN.Value, 64)
	logar.Check(err, log)

	return saldo
}
