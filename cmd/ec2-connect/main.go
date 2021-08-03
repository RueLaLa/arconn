package main

import (
	"fmt"
	"github.com/ruelala/ec2-connect/lib"
	"github.com/ruelala/ec2-connect/lib/awsClients"
	"github.com/ruelala/ec2-connect/lib/awsClients/ec2"
	"github.com/ruelala/ec2-connect/lib/awsClients/ssm"
)

func main() {
	args := lib.SetupApp()

	ttype := ec2.TargetType(args.Target)
	fmt.Println(fmt.Sprintf("input target type: %s", ttype))

	aws_config := awsClients.AwsConfig(args.Profile)

	id := ""
	if ttype != "ID" {
		id = ec2.Lookup(aws_config, args.Target, ttype)
	} else {
		id = args.Target
	}

	ssm_client := ssm.Lookup(aws_config, id)
	ssm.Connect(ssm_client, id, args.Profile)
}
