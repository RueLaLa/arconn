package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/integrii/flaggy"
	"github.com/ruelala/arconn/pkg/awsClients"
	"github.com/ruelala/arconn/pkg/awsClients/ec2"
	"github.com/ruelala/arconn/pkg/awsClients/ssm"
)

// these get passed in as ldflags by goreleaser
var version string
var commit string
var date string

func print_version() string {
	go_version := runtime.Version()
	return fmt.Sprintf("arconn %s built with %s on commit %s at %s", version, go_version, commit, date)
}

func main() {
	var profile = os.Getenv("AWS_PROFILE")
	var target string
	flaggy.String(&profile, "p", "profile", "aws profile to use")
	flaggy.String(&target, "t", "target", "resource to target for session connection")
	flaggy.SetVersion(print_version())
	flaggy.Parse()

	aws_config := awsClients.AwsConfig(profile)
	ttype := ec2.TargetType(target)
	fmt.Println(fmt.Sprintf("input target type: %s", ttype))

	var id string
	if ttype != "ID" {
		id = ec2.Lookup(aws_config, target, ttype)
	} else {
		id = target
	}
	ssm.Lookup(aws_config, id)
	ssm.Connect(aws_config, id)

	// ssm_client := ssm.Lookup(aws_config, id)
	// ssm.Connect(ssm_client, id, c.String("profile"))
}
