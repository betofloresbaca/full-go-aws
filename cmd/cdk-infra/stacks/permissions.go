package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/cloudformationinclude"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// PermissionsStackProps defines the properties for the Permissions stack
type PermissionsStackProps struct {
	awscdk.StackProps
}

// PermissionsStackResult contains the stack and the roles map
type PermissionsStackResult struct {
	Stack awscdk.Stack
	Roles map[string]awsiam.IRole
}

// PermissionsStack creates IAM roles and policies from a CloudFormation template
func PermissionsStack(scope constructs.Construct, id string, props *PermissionsStackProps) *PermissionsStackResult {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	templateFile := "cmd/cdk-infra/stacks/permissions-cfn.yaml"
	cfnTemplate := cloudformationinclude.NewCfnInclude(stack, jsii.String("PermissionsTemplate"), &cloudformationinclude.CfnIncludeProps{
		TemplateFile: jsii.String(templateFile),
	})
	roles := make(map[string]awsiam.IRole)
	roleLogicalIds := []string{
		"TelegramApiAuthorizerRole",
		"TelegramApiGatewayRole",
	}
	for _, roleLogicalId := range roleLogicalIds {
		resource := cfnTemplate.GetResource(jsii.String(roleLogicalId))
		if resource != nil {
			if cfnRole, ok := resource.(awsiam.CfnRole); ok {
				role := awsiam.Role_FromRoleArn(stack, jsii.String("Imported"+roleLogicalId), cfnRole.AttrArn(), nil)
				roles[roleLogicalId] = role
			}
		}
	}
	return &PermissionsStackResult{
		Stack: stack,
		Roles: roles,
	}
}
