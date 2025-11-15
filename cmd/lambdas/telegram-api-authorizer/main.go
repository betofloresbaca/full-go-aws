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

func handleRequest(ctx context.Context, event events.APIGatewayV2CustomAuthorizerV2Request) (
	events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	jsonEvent, _ := json.Marshal(event)
	log.Println("Received event:", string(jsonEvent))
	if len(event.IdentitySource) == 0 || event.IdentitySource[0] == "" {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context:      map[string]interface{}{"error": "missing or invalid token"},
		}, nil
	}
	requestToken := event.IdentitySource[0]
	log.Println("Got header token:", requestToken)

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
	} else {
		paramBytes, _ := json.Marshal(parameter)
		log.Println("Got parameter token:", string(paramBytes))
	}

	if err != nil || requestToken != *parameter.Parameter.Value {
		log.Println("Invalid Token")
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context:      map[string]interface{}{"error": "invalid token"},
		}, nil
	}
	log.Println("User Authorized")
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context:      map[string]interface{}{"user": "authorized"},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
