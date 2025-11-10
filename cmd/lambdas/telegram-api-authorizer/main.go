package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/betofloresbaca/expenses-manager/pkg/clients"
)

const HeaderName = "X-Telegram-Bot-Api-Secret-Token"

func handleRequest(ctx context.Context, event events.APIGatewayV2CustomAuthorizerV2Request) (
	events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	jsonEvent, _ := json.Marshal(event)
	log.Println("Received event:", string(jsonEvent))
	headerToken, ok := event.Headers[HeaderName]
	log.Println("Got header token:", headerToken)
	if !ok || headerToken == "" {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context:      map[string]interface{}{"error": "missing or invalid token"},
		}, nil
	}

	// Obtener el nombre del parámetro desde la variable de entorno
	parameterName := os.Getenv("TELEGRAM_SECRET_PARAM")
	if parameterName == "" {
		log.Println("Error: TELEGRAM_SECRET_PARAM environment variable not set")
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context:      map[string]interface{}{"error": "configuration error"},
		}, nil
	}

	// Obtener el parámetro del SSM (con cache)
	ssmClient := clients.GetClient(func(cfg aws.Config) *ssm.Client {
		return ssm.NewFromConfig(cfg)
	})

	parameter, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(parameterName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Println("Authorization failed:", err)
	}

	if err != nil || headerToken != *parameter.Parameter.Value {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context:      map[string]interface{}{"error": "invalid token"},
		}, nil
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context:      map[string]interface{}{"user": "authorized"},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
