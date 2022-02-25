package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ruelala/arconn/pkg/awsClients"
	"github.com/ruelala/arconn/pkg/awsClients/ec2"
	"github.com/ruelala/arconn/pkg/awsClients/ssm"
	"github.com/urfave/cli/v2"
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
	// check_deps()
	app := &cli.App{
		Name:    "arconn",
		Version: print_version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "profile",
				Aliases:  []string{"p"},
				EnvVars:  []string{"AWS_PROFILE"},
				Usage:    "aws profile for account specification",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "ec2",
				Usage:  "target ec2 instance for remote connection",
				Action: run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "target",
						Aliases:  []string{"t"},
						Usage:    "ec2 target to connect to remotely",
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\n %s", err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	aws_config := awsClients.AwsConfig(c.String("profile"))

	var id string
	switch c.Command.Name {
	case "ec2":
		{
			ttype := ec2.TargetType(c.String("target"))
			fmt.Println(fmt.Sprintf("input target type: %s", ttype))
			if ttype != "ID" {
				id = ec2.Lookup(aws_config, c.String("target"), ttype)
			} else {
				id = c.String("target")
			}
		}
		ssm.Lookup(aws_config, id)
		ssm.Connect(aws_config, id)
	}

	// ssm_client := ssm.Lookup(aws_config, id)
	// ssm.Connect(ssm_client, id, c.String("profile"))
	return nil
}
