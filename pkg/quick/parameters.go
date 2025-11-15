package quick

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/betofloresbaca/expenses-manager/pkg/clients"
)

func GetParameter(ctx context.Context, name string, withDecryption bool) (string, error) {
	ssmClient := clients.GetClient(func(cfg aws.Config) *ssm.Client {
		return ssm.NewFromConfig(cfg)
	})

	parameter, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(withDecryption),
	})

	if err != nil {
		return "", err
	}
	return *parameter.Parameter.Value, nil
}
