package constructs

import (
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// StateMachineProps defines the properties for creating a State Machine
type StateMachineProps struct {
	StateMachineName *string
	AslFilePath      string
	ARNReplacements  map[string]string
	Role             awsiam.IRole
}

// NewStateMachine creates a new State Machine from an ASL file
// replacing placeholders with the provided ARN values
func NewStateMachine(scope constructs.Construct, id *string, props *StateMachineProps) awsstepfunctions.StateMachine {
	if props == nil {
		panic("StateMachineProps cannot be nil")
	}

	if props.AslFilePath == "" {
		panic("AslFilePath is required")
	}
	aslBytes, err := os.ReadFile(props.AslFilePath)
	if err != nil {
		panic("Failed to read ASL file: " + props.AslFilePath + " - " + err.Error())
	}
	aslString := string(aslBytes)

	definitionWithReplacements := aslString
	for placeholder, arnValue := range props.ARNReplacements {
		definitionWithReplacements = strings.ReplaceAll(definitionWithReplacements, placeholder, arnValue)
	}

	stateMachine := awsstepfunctions.NewStateMachine(scope, id, &awsstepfunctions.StateMachineProps{
		StateMachineName: props.StateMachineName,
		DefinitionBody:   awsstepfunctions.DefinitionBody_FromString(jsii.String(definitionWithReplacements)),
		Role:             props.Role,
	})

	return stateMachine
}
