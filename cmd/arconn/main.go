package main

import (
	"fmt"
	"os"
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

	if args.PortForward != "" {
		target.PortForwarding = utils.ResolvePortForwarding(args.PortForward)
	}

	if args.RemoteHost != "" {
		target.RemoteHost = args.RemoteHost
	}

	target.Type = utils.TargetType(args.Target)
	fmt.Println(fmt.Sprintf("computed target type: %s", target.Type))

	switch target.Type {
	case "EC2_ID", "SSM_MI_ID":
		target.ResolvedName = args.Target
	case "IP":
		target = ec2.Lookup(args, target)
	default:
		target = ecs.Lookup(args, target)

		// this value gets set if a valid ECS target is found
		// if thats found, then we dont need to look in EC2.
		if target.SessionInfo == "" {
			target = ec2.Lookup(args, target)
		}
	}

	// because ECS exec command generates its own session object
	// if that is set in the target struct, it doesnt need to be looked up in ssm.
	// if its not however, the target needs to be identified in ssm and
	// a session has to be created.
	if target.SessionInfo == "" {
		target = ssm.Lookup(args, target)
	}

	if !target.Resolved {
		fmt.Println(fmt.Sprintf("target %s couldnt be found in ECS, EC2, or SSM", args.Target))
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("connecting to %s", target.ResolvedName))
	ssm.Connect(args, target)
}
