// File: cmd/cdk-infra/main.go
package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/betofloresbaca/full-go-aws/cmd/cdk-infra/stacks"
)

func main() {
	app := awscdk.NewApp(nil)
	env := &awscdk.Environment{
		// Configure your AWS Account and Region
		// Or leave nil to use default CLI profile
		// Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		// Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
	stacks.LambdasStack(app, "LambdasStack", &stacks.LambdasStackProps{
		StackProps: awscdk.StackProps{
			Env: env,
		},
	})

	app.Synth(nil)
}
