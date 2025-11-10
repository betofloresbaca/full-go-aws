package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
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

	// Create lambda with the role from permissions stack
	utils.NewFunction(
		"em-TelegramApiAuthorizer",
		"bin/telegram-api-authorizer.zip",
		stack,
		map[string]*string{
			"TELEGRAM_SECRET_PARAM": jsii.String("/em/TelegramSecretToken"),
		},
		props.Roles["TelegramApiAuthorizerRole"],
	)

	return stack
}
