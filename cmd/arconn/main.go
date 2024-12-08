package main

import (
	"fmt"

	"github.com/aws/session-manager-plugin/pkg/log"
	"github.com/ruelala/arconn/pkg/awsClients/ec2"
	"github.com/ruelala/arconn/pkg/awsClients/ecs"
	"github.com/ruelala/arconn/pkg/awsClients/ssm"
	"github.com/ruelala/arconn/pkg/utils"
)

func main() {
	utils.IsLatest()
	args := utils.ParseFlags()
	target := utils.Target{}

	if args.PortForward != "" {
		target.PortForwarding = utils.ResolvePortForwarding(args.PortForward)
	}

	if args.RemoteHost != "" {
		target.RemoteHost = args.RemoteHost
	}

	target.Type = utils.TargetType(args.Target)
	log.Alwaysf("computed target type: %s", target.Type)

	switch target.Type {
	case "EC2_ID", "SSM_MI_ID":
		target.ResolvedName = args.Target
	case "ECS_ID":
		target = ecs.Lookup(args, target)
	case "IP":
		target = ec2.Lookup(args, target)
	default:
		target = ecs.Lookup(args, target)

		if !target.Resolved {
			target = ec2.Lookup(args, target)
		}
	}

	// ECS does not register containers in SSM as inventory
	// so we skip the lookup and assume the task is available
	// for sessions.
	if target.Type != "ECS" {
		target = ssm.Lookup(args, target)
	}

	if !target.Resolved {
		utils.Panic(fmt.Errorf("target %s couldnt be found in ECS, EC2, or SSM", args.Target))
	}

	log.Alwaysf("connecting to %s", target.ResolvedName)
	ssm.Connect(args, target)
}
