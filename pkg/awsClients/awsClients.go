package awsClients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/ruelala/arconn/pkg/utils"
)

func AwsConfig(profile string) aws.Config {
	scp, err := config.LoadSharedConfigProfile(context.TODO(), profile)
	utils.Panic(err)
	if scp.Region == "" {
		scp.Region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithDefaultRegion(scp.Region),
		config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
			options.RoleSessionName = utils.GetSessionName()
		}),
	)
	utils.Panic(err)

	return cfg
}

func ECSClient(profile string) *ecs.Client {
	client := ecs.NewFromConfig(AwsConfig(profile))
	return client
}

func EC2Client(profile string) *ec2.Client {
	client := ec2.NewFromConfig(AwsConfig(profile))
	return client
}

func SSMClient(profile string) *ssm.Client {
	client := ssm.NewFromConfig(AwsConfig(profile))
	return client
}

func CurrentRegion(profile string) string {
	cfg := AwsConfig(profile)
	return cfg.Region
}
