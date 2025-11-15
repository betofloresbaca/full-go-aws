package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	customConstructs "github.com/betofloresbaca/expenses-manager/cmd/cdk-infra/constructs"
)

// Define the properties for your stack
type ApiStackProps struct {
	awscdk.StackProps
	Roles        map[string]awsiam.IRole
	StateMachine awsstepfunctions.StateMachine
}

// ApiStack creates the CloudFormation Stack
func ApiStack(scope constructs.Construct, id string, props *ApiStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Create lambda authorizer with the role from permissions stack
	authorizerLambdaConstruct := customConstructs.NewLambdaFunction(
		stack,
		jsii.String("TelegramApiAuthorizer"),
		&customConstructs.LambdaFunctionProps{
			FunctionName: "EM-TelegramApiAuthorizer",
			ZipPath:      "bin/telegram-api-authorizer.zip",
			Environment: map[string]*string{
				"TELEGRAM_SECRET_PARAM": jsii.String("/em/TelegramSecret"),
			},
			Role: props.Roles["TelegramApiAuthorizerRole"],
		},
	)

	// Create API Gateway with webhook using the ApiWebhook construct
	customConstructs.NewApiWebhook(stack, "TelegramBotApi", &customConstructs.ApiWebhookProps{
		ApiName:            "telegram-bot-api",
		ApiDescription:     "Telegram Bot HTTP API Gateway",
		AuthorizerLambda:   authorizerLambdaConstruct.Function,
		IdentityHeaderName: "X-Telegram-Bot-Api-Secret-Token",
		IntegrationRole:    props.Roles["TelegramApiGatewayRole"],
		StateMachineArn:    *props.StateMachine.StateMachineArn(),
		RouteKey:           "POST /webhook",
	})

	return stack
}
