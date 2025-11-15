package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/betofloresbaca/expenses-manager/cmd/cdk-infra/utils"
)

// Define the properties for your stack
type ApiStackProps struct {
	awscdk.StackProps
	Roles map[string]awsiam.IRole
}

// ApiStack creates the CloudFormation Stack
func ApiStack(scope constructs.Construct, id string, props *ApiStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// // Create lambda authorizer with the role from permissions stack
	authorizerLambda := utils.NewFunction(
		"em-TelegramApiAuthorizer",
		"bin/telegram-api-authorizer.zip",
		stack,
		map[string]*string{
			"TELEGRAM_SECRET_PARAM": jsii.String("/em/TelegramSecretToken"),
		},
		props.Roles["TelegramApiAuthorizerRole"],
	)

	// Create HTTP API Gateway (v2)
	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("TelegramBotApi"), &awsapigatewayv2.HttpApiProps{
		ApiName:     jsii.String("telegram-bot-api"),
		Description: jsii.String("Telegram Bot HTTP API Gateway"),
	})

	// Create Lambda Authorizer using raw CloudFormation
	cfnAuthorizer := awsapigatewayv2.NewCfnAuthorizer(stack, jsii.String("TelegramApiAuthorizer"), &awsapigatewayv2.CfnAuthorizerProps{
		ApiId:                          httpApi.HttpApiId(),
		AuthorizerType:                 jsii.String("REQUEST"),
		Name:                           jsii.String("TelegramApiAuthorizer"),
		AuthorizerUri:                  jsii.String("arn:aws:apigateway:" + *stack.Region() + ":lambda:path/2015-03-31/functions/" + *authorizerLambda.FunctionArn() + "/invocations"),
		AuthorizerPayloadFormatVersion: jsii.String("2.0"),
		IdentitySource:                 &[]*string{jsii.String("$request.header.X-Telegram-Bot-Api-Secret-Token")},
		EnableSimpleResponses:          jsii.Bool(true),
		AuthorizerResultTtlInSeconds:   jsii.Number(300),
	})

	// Grant API Gateway permission to invoke the authorizer Lambda
	authorizerLambda.GrantInvoke(awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), nil))

	// Create Step Functions integration using raw CloudFormation
	integration := awsapigatewayv2.NewCfnIntegration(stack, jsii.String("StepFunctionsIntegration"), &awsapigatewayv2.CfnIntegrationProps{
		ApiId:              httpApi.HttpApiId(),
		IntegrationType:    jsii.String("AWS_PROXY"),
		IntegrationSubtype: jsii.String("StepFunctions-StartExecution"),
		CredentialsArn:     props.Roles["TelegramApiGatewayRole"].RoleArn(),
		RequestParameters: &map[string]*string{
			"StateMachineArn": jsii.String("arn:aws:states:us-east-1:742377680347:stateMachine:ExampleTelegramApiStateMachine"),
			"Input":           jsii.String("$request.body"),
		},
		PayloadFormatVersion: jsii.String("1.0"),
	})

	// Create the /webhook route with raw CloudFormation
	awsapigatewayv2.NewCfnRoute(stack, jsii.String("WebhookRoute"), &awsapigatewayv2.CfnRouteProps{
		ApiId:             httpApi.HttpApiId(),
		RouteKey:          jsii.String("POST /webhook"),
		Target:            jsii.String(*jsii.String("integrations/" + *integration.Ref())),
		AuthorizationType: jsii.String("CUSTOM"),
		AuthorizerId:      cfnAuthorizer.Ref(),
	})

	// Output the API URL
	awscdk.NewCfnOutput(stack, jsii.String("ApiUrl"), &awscdk.CfnOutputProps{
		Value:       httpApi.ApiEndpoint(),
		Description: jsii.String("Telegram Bot HTTP API Gateway URL"),
		ExportName:  jsii.String("TelegramBotApiUrl"),
	})

	return stack
}
