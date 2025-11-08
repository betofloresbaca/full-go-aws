package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// Define the properties for your stack
type LambdasStackProps struct {
	awscdk.StackProps
}

func newFunction(name string, zipPath string, stack awscdk.Stack) awslambda.Function {
	return awslambda.NewFunction(stack, jsii.String(name), &awslambda.FunctionProps{
		FunctionName: jsii.String(name),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("main"),
		Code:         awslambda.Code_FromAsset(jsii.String(zipPath), nil),
		Architecture: awslambda.Architecture_ARM_64(),
	})
}

// LambdasStack creates the CloudFormation Stack
func LambdasStack(scope constructs.Construct, id string, props *LambdasStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Deploy lambda functions
	newFunction("CreateUser", "bin/create-user.zip", stack)

	return stack
}
