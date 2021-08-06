package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ruelala/ec2-connect/lib/awsClients"
	"github.com/ruelala/ec2-connect/lib/awsClients/ec2"
	"github.com/ruelala/ec2-connect/lib/awsClients/ssm"
	"github.com/urfave/cli/v2"
)

// these get passed in as ldflags by goreleaser
var version string
var commit string
var date string

func print_version() string {
	go_version := runtime.Version()
	return fmt.Sprintf("ec2-connect %s built with %s on commit %s at %s", version, go_version, commit, date)
}

func main() {
	app := &cli.App{
		Name:    "ec2-connect",
		Version: print_version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				EnvVars: []string{"AWS_PROFILE"},
				Usage:   "aws profile for account specification",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "target",
				Aliases:  []string{"t"},
				Usage:    "target to start ssm session for",
				Required: true,
			},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\n %s", err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	ttype := ec2.TargetType(c.String("target"))
	fmt.Println(fmt.Sprintf("input target type: %s", ttype))

	aws_config := awsClients.AwsConfig(c.String("profile"))

	id := ""
	if ttype != "ID" {
		id = ec2.Lookup(aws_config, c.String("target"), ttype)
	} else {
		id = c.String("target")
	}

	ssm_client := ssm.Lookup(aws_config, id)
	ssm.Connect(ssm_client, id, c.String("profile"))
	return nil
}
