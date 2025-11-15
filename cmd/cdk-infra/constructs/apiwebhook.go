package constructs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// ApiWebhookProps defines the properties for the ApiWebhook construct
type ApiWebhookProps struct {
	// ApiName is the name of the HTTP API
	ApiName string
	// ApiDescription is the description of the HTTP API
	ApiDescription string
	// AuthorizerLambda is the Lambda function used for authorization
	AuthorizerLambda awslambda.IFunction
	// IdentityHeaderName is the name of the header used for identity validation
	IdentityHeaderName string
	// IntegrationRole is the IAM role used by API Gateway to invoke the integration
	IntegrationRole awsiam.IRole
	// StateMachineArn is the ARN of the Step Functions state machine to integrate
	StateMachineArn string
	// RouteKey is the route key for the webhook (e.g., "POST /webhook")
	RouteKey string
	// AuthorizerTTL is the TTL in seconds for the authorizer cache (default: 300)
	AuthorizerTTL *float64
}

// ApiWebhookOutputs contains the outputs from the ApiWebhook construct
type ApiWebhookOutputs struct {
	// HttpApi is the created HTTP API Gateway
	HttpApi awsapigatewayv2.HttpApi
	// ApiEndpoint is the endpoint URL of the API
	ApiEndpoint *string
	// Authorizer is the CloudFormation authorizer resource
	Authorizer awsapigatewayv2.CfnAuthorizer
	// Integration is the CloudFormation integration resource
	Integration awsapigatewayv2.CfnIntegration
	// Route is the CloudFormation route resource
	Route awsapigatewayv2.CfnRoute
}

// NewApiWebhook creates a new API Gateway HTTP API with Lambda authorizer and Step Functions integration
func NewApiWebhook(scope constructs.Construct, id string, props *ApiWebhookProps) *ApiWebhookOutputs {
	// Set defaults
	if props.AuthorizerTTL == nil {
		props.AuthorizerTTL = jsii.Number(300)
	}

	stack := awscdk.Stack_Of(scope)

	// Create HTTP API Gateway (v2)
	httpApi := awsapigatewayv2.NewHttpApi(scope, jsii.String(id+"HttpApi"), &awsapigatewayv2.HttpApiProps{
		ApiName:     jsii.String(props.ApiName),
		Description: jsii.String(props.ApiDescription),
	})

	// Create Lambda Authorizer using raw CloudFormation
	cfnAuthorizer := awsapigatewayv2.NewCfnAuthorizer(scope, jsii.String(id+"Authorizer"), &awsapigatewayv2.CfnAuthorizerProps{
		ApiId:          httpApi.HttpApiId(),
		AuthorizerType: jsii.String("REQUEST"),
		Name:           jsii.String(props.ApiName + "Authorizer"),
		AuthorizerUri: jsii.String("arn:aws:apigateway:" + *stack.Region() +
			":lambda:path/2015-03-31/functions/" + *props.AuthorizerLambda.FunctionArn() + "/invocations"),
		AuthorizerPayloadFormatVersion: jsii.String("2.0"),
		IdentitySource:                 &[]*string{jsii.String("$request.header." + props.IdentityHeaderName)},
		EnableSimpleResponses:          jsii.Bool(true),
		AuthorizerResultTtlInSeconds:   props.AuthorizerTTL,
	})

	// Grant API Gateway permission to invoke the authorizer Lambda
	props.AuthorizerLambda.GrantInvoke(awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), nil))

	// Create Step Functions integration using raw CloudFormation
	integration := awsapigatewayv2.NewCfnIntegration(scope, jsii.String(id+"Integration"), &awsapigatewayv2.CfnIntegrationProps{
		ApiId:              httpApi.HttpApiId(),
		IntegrationType:    jsii.String("AWS_PROXY"),
		IntegrationSubtype: jsii.String("StepFunctions-StartExecution"),
		CredentialsArn:     props.IntegrationRole.RoleArn(),
		RequestParameters: &map[string]*string{
			"StateMachineArn": jsii.String(props.StateMachineArn),
			"Input":           jsii.String("$request.body"),
		},
		PayloadFormatVersion: jsii.String("1.0"),
	})

	// Create the webhook route with raw CloudFormation
	route := awsapigatewayv2.NewCfnRoute(scope, jsii.String(id+"Route"), &awsapigatewayv2.CfnRouteProps{
		ApiId:             httpApi.HttpApiId(),
		RouteKey:          jsii.String(props.RouteKey),
		Target:            jsii.String("integrations/" + *integration.Ref()),
		AuthorizationType: jsii.String("CUSTOM"),
		AuthorizerId:      cfnAuthorizer.Ref(),
	})

	// Output the API URL
	awscdk.NewCfnOutput(scope, jsii.String(id+"ApiUrl"), &awscdk.CfnOutputProps{
		Value:       httpApi.ApiEndpoint(),
		Description: jsii.String(props.ApiName + " API Gateway URL"),
		ExportName:  jsii.String(props.ApiName + "Url"),
	})

	return &ApiWebhookOutputs{
		HttpApi:     httpApi,
		ApiEndpoint: httpApi.ApiEndpoint(),
		Authorizer:  cfnAuthorizer,
		Integration: integration,
		Route:       route,
	}
}
