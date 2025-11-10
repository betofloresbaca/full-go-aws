package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/betofloresbaca/expenses-manager/cmd/cdk-infra/utils"
)

// Define the properties for your stack
type LambdasStackProps struct {
	awscdk.StackProps
}

// LambdasStack creates the CloudFormation Stack
func LambdasStack(scope constructs.Construct, id string, props *LambdasStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Deploy lambda functions
	utils.NewFunction("CreateUser", "bin/create-user.zip", stack, nil, nil)

	return stack
}
