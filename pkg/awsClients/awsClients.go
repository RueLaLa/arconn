package awsClients

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func AwsConfig(profile string) aws.Config {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func ECSClient(profile string) *ecs.Client {
	config := AwsConfig(profile)
	client := ecs.NewFromConfig(config, func(o *ecs.Options) {
		o.Region = "us-east-1"
	})
	return client
}

func EC2Client(profile string) *ec2.Client {
	config := AwsConfig(profile)
	client := ec2.NewFromConfig(config, func(o *ec2.Options) {
		o.Region = "us-east-1"
	})
	return client
}

func SSMClient(profile string) *ssm.Client {
	config := AwsConfig(profile)
	client := ssm.NewFromConfig(config, func(o *ssm.Options) {
		o.Region = "us-east-1"
	})
	return client
}
