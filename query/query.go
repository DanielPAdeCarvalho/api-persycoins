package query

import (
	"api-persycoins/logar"
	"api-persycoins/models"
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// InsertCliente insere um novo cliente no banco
func InsertCliente(client *dynamodb.Client, PersyCoins models.PersyCoins, log logar.Logfile) {
	nomeCompleto := PersyCoins.Nome + " " + PersyCoins.Sobrenome
	PersyCoins.Nome = nomeCompleto
	if GetSaldo(client, PersyCoins.Nome, log) > 0 {
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":e": &types.AttributeValueMemberS{Value: PersyCoins.Email},
			},
			TableName: aws.String("PersyCoins"),
			Key: map[string]types.AttributeValue{
				"Nome": &types.AttributeValueMemberS{Value: PersyCoins.Nome},
			},
			ReturnValues:     types.ReturnValueUpdatedNew,
			UpdateExpression: aws.String("set Email = :e"),
		}

		_, err := client.UpdateItem(context.Background(), input)
		logar.Check(err, log)
	} else {
		coins, err := attributevalue.MarshalMap(PersyCoins)
		logar.Check(err, log)

		input := &dynamodb.PutItemInput{
			TableName: aws.String("PersyCoins"),
			Item:      coins,
		}
		_, err = client.PutItem(context.Background(), input)
		logar.Check(err, log)
	}
}

// GetSaldo retorna o saldo de um cliente
func GetSaldo(client *dynamodb.Client, nome string, log logar.Logfile) float64 {
	key, err := attributevalue.MarshalMap(map[string]interface{}{
		"Nome": nome,
	})
	logar.CheckInfo(err, log)

	input := &dynamodb.GetItemInput{
		TableName: aws.String("PersyCoins"),
		Key:       key,
	}

	// Call the GetItem method with the input
	resp, err := client.GetItem(context.TODO(), input)
	logar.CheckInfo(err, log)
	if resp.Item == nil {
		log.InfoLogger.Println("Cliente não encontrado")
	}

	// Get the Saldo attribute value
	saldoValue, ok := resp.Item["Saldo"]
	if !ok {
		log.InfoLogger.Println("Saldo não encontrado")
	}

	saldoN := saldoValue.(*types.AttributeValueMemberN)
	saldo, err := strconv.ParseFloat(saldoN.Value, 64)
	logar.CheckInfo(err, log)

	return saldo
}

// UpdateSaldo atualiza o saldo de um cliente
func UpdateSaldo(log logar.Logfile, client *dynamodb.Client, nome string, operation string, valor float64) {

	// negativar o valor caso seja para subtrair
	if operation == "sub" {
		valor = -valor
	}
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("PersyCoins"),
		Key: map[string]types.AttributeValue{
			"Nome": &types.AttributeValueMemberS{
				Value: nome,
			},
		},
		UpdateExpression:          aws.String(fmt.Sprintf("%s #attr :val", "ADD")),
		ExpressionAttributeNames:  map[string]string{"#attr": "Saldo"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":val": &types.AttributeValueMemberN{Value: fmt.Sprint(valor)}},
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	_, err := client.UpdateItem(context.TODO(), input)
	logar.Check(err, log)
}

// GetSaldoByMail retorna o saldo de um cliente
func GetSaldoByMail(client *dynamodb.Client, mail string, log logar.Logfile) string {
	nome := GetNomeSobrenomeByEmail(client, mail, log)
	fmt.Println(nome)
	saldo := GetSaldo(client, nome, log)
	retorno := nome + " " + fmt.Sprintf("%.2f", saldo)
	return retorno
}

func GetNomeSobrenomeByEmail(client *dynamodb.Client, email string, log logar.Logfile) string {
	input := &dynamodb.ScanInput{
		TableName: aws.String("LoginCliente"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":e": &types.AttributeValueMemberS{
				Value: email,
			},
		},
		FilterExpression: aws.String("Email = :e"),
	}

	pager := dynamodb.NewScanPaginator(client, input)

	for pager.HasMorePages() {
		page, err := pager.NextPage(context.Background())
		logar.Check(err, log)

		for _, item := range page.Items {
			nomeValue, ok := item["Nome"].(*types.AttributeValueMemberS)
			if !ok {
				log.ErrorLogger.Println("Nome não encontrado")
				continue
			}
			return nomeValue.Value
		}
	}
	return "nao encontrado"
}
