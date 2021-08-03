package lib

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

type Args struct {
	Profile, Target string
}

func SetupApp() Args {
	args := Args{}
	app := &cli.App{
		Name:    "ec2-connect",
		Version: print_version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				EnvVars: []string{"AWS_PROFILE"},
				Usage:   "aws profile for account specification",
			},
			&cli.StringFlag{
				Name:     "target",
				Aliases:  []string{"t"},
				Usage:    "target to start ssm session for",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			args.Profile = c.String("profile")
			args.Target = c.String("target")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\n %s", err)
		os.Exit(1)
	}
	return args
}

// these get passed in as ldflags by goreleaser
var version string
var commit string
var date string

func print_version() string {
	go_version := runtime.Version()
	return fmt.Sprintf("ec2-connect %s built with %s on commit %s at %s\n", version, go_version, commit, date)
}
