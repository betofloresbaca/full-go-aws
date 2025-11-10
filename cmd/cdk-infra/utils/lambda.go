package utils

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

func NewFunction(
	name string,
	zipPath string,
	stack awscdk.Stack,
	environment map[string]*string,
	role awsiam.IRole) awslambda.Function {
	return awslambda.NewFunction(stack, jsii.String(name), &awslambda.FunctionProps{
		FunctionName: jsii.String(name),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("main"),
		Code:         awslambda.Code_FromAsset(jsii.String(zipPath), nil),
		Architecture: awslambda.Architecture_ARM_64(),
		Environment:  &environment,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(60)),
		Role:         role,
	})
}
