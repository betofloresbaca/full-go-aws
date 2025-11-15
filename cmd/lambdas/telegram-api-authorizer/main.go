package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/betofloresbaca/expenses-manager/pkg/quick"
)

func handleRequest(ctx context.Context, event events.APIGatewayV2CustomAuthorizerV2Request) (
	events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	if len(event.IdentitySource) == 0 || event.IdentitySource[0] == "" {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context:      map[string]interface{}{"error": "missing or invalid token"},
		}, nil
	}
	requestSecret := event.IdentitySource[0]

	// Obtener el nombre del parámetro desde la variable de entorno
	paramsSecret, err := quick.GetParameter(ctx, os.Getenv("TELEGRAM_SECRET_PARAM"), true)
	if err != nil {
		log.Println("Error getting Telegram secret:", err)
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
			Context:      map[string]interface{}{"error": "error getting Telegram secret"},
		}, fmt.Errorf("error getting Telegram secret: %w", err)
	}

	if requestSecret != paramsSecret {
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
