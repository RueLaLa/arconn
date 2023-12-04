package AwsConfig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/ruelala/arconn/pkg/utils"
)

func BuildConfig(args utils.Args) aws.Config {
	if args.Vault != "" {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		utils.Panic(err)
		return cfg
	}

	scp, err := config.LoadSharedConfigProfile(context.TODO(), args.Profile)
	utils.Panic(err)
	if scp.Region == "" {
		scp.Region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(args.Profile),
		config.WithDefaultRegion(scp.Region),
		config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
			options.RoleSessionName = utils.GetSessionName()
		}),
	)
	utils.Panic(err)

	return cfg
}
