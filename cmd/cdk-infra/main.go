package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"github.com/betofloresbaca/expenses-manager/cmd/cdk-infra/stacks"
)

func main() {
	app := awscdk.NewApp(nil)
	stackProps := awscdk.StackProps{
		Env: &awscdk.Environment{},
	}

	// Create the Permissions stack first
	permissionsResult := stacks.PermissionsStack(app, "PermissionsStack", &stacks.PermissionsStackProps{
		StackProps: stackProps,
	})

	// Create the API stack with dependency on Permissions stack
	apiStack := stacks.ApiStack(app, "ApiStack", &stacks.ApiStackProps{
		StackProps: stackProps,
		Roles:      permissionsResult.Roles,
	})
	apiStack.AddDependency(permissionsResult.Stack, jsii.String("Requires IAM roles from PermissionsStack"))

	app.Synth(nil)
}
