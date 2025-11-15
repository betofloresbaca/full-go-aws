package constructs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// LambdaFunctionProps defines the properties for the Lambda Function construct
type LambdaFunctionProps struct {
	// FunctionName is the name of the Lambda function
	FunctionName string
	// ZipPath is the path to the deployment package zip file
	ZipPath string
	// Environment contains environment variables for the function
	Environment map[string]*string
	// Role is the IAM role for the Lambda function
	Role awsiam.IRole
	// Timeout is the function execution timeout (default: 60 seconds)
	Timeout awscdk.Duration
	// Runtime is the Lambda runtime (default: PROVIDED_AL2023)
	Runtime awslambda.Runtime
	// Architecture is the instruction set architecture (default: ARM_64)
	Architecture awslambda.Architecture
	// Handler is the function entry point (default: "main")
	Handler string
}

// LambdaFunction represents a Lambda function construct
type LambdaFunction struct {
	constructs.Construct
	Function awslambda.Function
}

// NewLambdaFunction creates a new Lambda function construct
func NewLambdaFunction(scope constructs.Construct, id *string, props *LambdaFunctionProps) *LambdaFunction {
	construct := constructs.NewConstruct(scope, id)

	// Set defaults
	timeout := props.Timeout
	if timeout == nil {
		timeout = awscdk.Duration_Seconds(jsii.Number(60))
	}
	runtime := props.Runtime
	if runtime == nil {
		runtime = awslambda.Runtime_PROVIDED_AL2023()
	}
	architecture := props.Architecture
	if architecture == nil {
		architecture = awslambda.Architecture_ARM_64()
	}
	handler := props.Handler
	if handler == "" {
		handler = "main"
	}

	// Create the Lambda function
	function := awslambda.NewFunction(construct, jsii.String("Function"), &awslambda.FunctionProps{
		FunctionName: jsii.String(props.FunctionName),
		Runtime:      runtime,
		Handler:      jsii.String(handler),
		Code:         awslambda.Code_FromAsset(jsii.String(props.ZipPath), nil),
		Architecture: architecture,
		Environment:  &props.Environment,
		Timeout:      timeout,
		Role:         props.Role,
	})

	return &LambdaFunction{
		Construct: construct,
		Function:  function,
	}
}
