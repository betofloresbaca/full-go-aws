package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/betofloresbaca/expenses-manager/pkg/quick"
	"github.com/betofloresbaca/expenses-manager/pkg/telegram"
)

type Request struct {
	ChatID    int64  `json:"ChatID"`
	Text      string `json:"Text"`
	ParseMode string `json:"ParseMode,omitempty"`
}

// Response representa la respuesta del Lambda.
type Response struct {
	Ok        bool `json:"Ok"`
	MessageID int  `json:"MessageID,omitempty"`
}

func handleRequest(ctx context.Context, request Request) (Response, error) {

	// Obtener el nombre del parámetro desde la variable de entorno
	telegramToken, err := quick.GetParameter(ctx, os.Getenv("TELEGRAM_TOKEN_PARAM"), true)
	if err != nil {
		return Response{Ok: false}, err
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken)

	// Crear el request usando el modelo de telegram
	sendRequest := telegram.SendMessageRequest{
		ChatID:    request.ChatID,
		Text:      request.Text,
		ParseMode: request.ParseMode,
	}
	jsonBody, _ := json.Marshal(sendRequest)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return Response{Ok: false}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Response{Ok: false}, fmt.Errorf("error telegram: %d %s", resp.StatusCode, resp.Status)
	}

	var telegramResp telegram.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return Response{Ok: false}, err
	}

	if !telegramResp.OK {
		return Response{Ok: false}, fmt.Errorf("telegram api error: %s", telegramResp.Description)
	}

	return Response{Ok: true, MessageID: int(telegramResp.Result.MessageID)}, nil
}

func main() {
	lambda.Start(handleRequest)
}
