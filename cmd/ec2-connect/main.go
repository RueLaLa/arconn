package main

import (
	"github.com/ruelala/ec2-connect/lib"
	"github.com/ruelala/ec2-connect/lib/awsClients"
	"github.com/ruelala/ec2-connect/lib/awsClients/ec2"
	"github.com/ruelala/ec2-connect/lib/awsClients/ssm"
)

func main() {
	args := lib.SetupApp()

	aws_config := awsClients.AwsConfig(args.Profile)

	id := ec2.Lookup(aws_config, args.Target)

	ssm.Connect(aws_config, id)
}
