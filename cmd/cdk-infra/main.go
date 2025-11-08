package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/betofloresbaca/full-go-aws/cmd/cdk-infra/stacks"
)

func main() {
	app := awscdk.NewApp(nil)
	env := &awscdk.Environment{}
	stacks.DbStack(app, "DbStack", &stacks.DbStackProps{
		StackProps: awscdk.StackProps{
			Env: env,
		},
	})
	stacks.LambdasStack(app, "LambdasStack", &stacks.LambdasStackProps{
		StackProps: awscdk.StackProps{
			Env: env,
		},
	})

	app.Synth(nil)
}
