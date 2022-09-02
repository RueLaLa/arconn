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
	profile, target := utils.ParseFlags(print_version())
	ttype := utils.TargetType(target)
	fmt.Println(fmt.Sprintf("input target type: %s", ttype))

	resolved_target := ""
	session := ""
	if (ttype == "EC2_ID") || (ttype == "SSM_MI_ID") {
		resolved_target = target
	} else {
		resolved_target, session = ecs.Lookup(profile, target)
		if session == "" {
			resolved_target = ec2.Lookup(profile, target, ttype)
		}
		if resolved_target == "" {
			resolved_target = ssm.Lookup(profile, target, false)
		}
	}

	if resolved_target == "" {
		os.Exit(1)
	}

	fmt.Println(fmt.Sprintf("connecting to %s", resolved_target))
	ssm.Connect(profile, session, resolved_target)
}
