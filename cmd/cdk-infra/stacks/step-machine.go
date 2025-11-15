package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	customConstructs "github.com/betofloresbaca/expenses-manager/cmd/cdk-infra/constructs"
)

// Define the properties for your stack
type LambdasStackProps struct {
	awscdk.StackProps
	Roles map[string]awsiam.IRole
}

// StepMachineStackResult contains the stack and its resources
type StepMachineStackResult struct {
	Stack        awscdk.Stack
	StateMachine awsstepfunctions.StateMachine
}

// StepMachineStack creates the CloudFormation Stack
func StepMachineStack(scope constructs.Construct, id string, props *LambdasStackProps) *StepMachineStackResult {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Deploy lambda functions
	telegramSendMessage := customConstructs.NewLambdaFunction(
		stack,
		jsii.String("TelegramSendMessage"),
		&customConstructs.LambdaFunctionProps{
			FunctionName: "EM-TelegramSendMessage",
			ZipPath:      "bin/telegram-send-message.zip",
			Environment: map[string]*string{
				"TELEGRAM_TOKEN_PARAM": jsii.String("/em/TelegramToken"),
			},
			Role: props.Roles["TelegramSendMessageRole"],
		},
	)

	// Create the state machine
	stateMachine := customConstructs.NewStateMachine(
		stack,
		jsii.String("TelegramBotStateMachine"),
		&customConstructs.StateMachineProps{
			StateMachineName: jsii.String("EM-TelegramBotStateMachine"),
			AslFilePath:      "cmd/cdk-infra/resources/telegram-bot-state-machine.asl.json",
			ARNReplacements: map[string]string{
				"{% <TELEGRAM_SEND_MESSAGE> %}": *telegramSendMessage.Function.FunctionArn(),
			},

			Role: props.Roles["TelegramBotStateMachineRole"],
		},
	)

	return &StepMachineStackResult{
		Stack:        stack,
		StateMachine: stateMachine,
	}
}
