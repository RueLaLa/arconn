package main

import (
	"fmt"
	"runtime"

	"github.com/ruelala/arconn/pkg/awsClients/ec2"
	"github.com/ruelala/arconn/pkg/awsClients/ecs"
	"github.com/ruelala/arconn/pkg/awsClients/ssm"
	"github.com/ruelala/arconn/pkg/utils"
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
	args := utils.ParseFlags(print_version())
	target := utils.Target{}

	target.Type = utils.TargetType(args.Target)
	fmt.Println(fmt.Sprintf("computed target type: %s", target.Type))

	if (target.Type == "EC2_ID") || (target.Type == "SSM_MI_ID") {
		target.ResolvedName = args.Target
	} else if target.Type == "IP" {
		target = ec2.Lookup(args, target)
	} else {
		target = ecs.Lookup(args, target)

		if target.SessionInfo == "" {
			target = ec2.Lookup(args, target)
		}
	}

	if target.SessionInfo == "" {
		target = ssm.Lookup(args, target)
	}
	fmt.Println(fmt.Sprintf("connecting to %s", target.ResolvedName))
	ssm.Connect(args, target)
}
