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
	permissionsResult := stacks.PermissionsStack(app, "PermissionsStack", &stacks.PermissionsStackProps{
		StackProps: stackProps,
	})

	stepMachineResult := stacks.StepMachineStack(app, "StepMachineStack", &stacks.LambdasStackProps{
		StackProps: stackProps,
		Roles:      permissionsResult.Roles,
	})
	stepMachineResult.Stack.AddDependency(permissionsResult.Stack, jsii.String("Requires IAM roles from PermissionsStack"))

	apiStack := stacks.ApiStack(app, "ApiStack", &stacks.ApiStackProps{
		StackProps:   stackProps,
		Roles:        permissionsResult.Roles,
		StateMachine: stepMachineResult.StateMachine,
	})
	apiStack.AddDependency(stepMachineResult.Stack, jsii.String("Requires Step Machine Stack to connect the API"))

	app.Synth(nil)
}
