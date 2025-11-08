package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

var client *lambda.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	client = lambda.NewFromConfig(cfg)
}

func callLambda(ctx context.Context) (string, error) {
	return "Example response", nil
}

func handleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	// event
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("EVENT: %s", eventJson)
	// environment variables
	log.Printf("REGION: %s", os.Getenv("AWS_REGION"))
	log.Println("ALL ENV VARS:")
	for _, element := range os.Environ() {
		log.Println(element)
	}
	// request context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("REQUEST ID: %s", lc.AwsRequestID)
	// global variable
	log.Printf("FUNCTION NAME: %s", lambdacontext.FunctionName)
	// context method
	deadline, _ := ctx.Deadline()
	log.Printf("DEADLINE: %s", deadline)
	// AWS SDK call
	usage, err := callLambda(ctx)
	if err != nil {
		return "ERROR", err
	}
	return usage, nil
}

func main() {
	runtime.Start(handleRequest)
}
