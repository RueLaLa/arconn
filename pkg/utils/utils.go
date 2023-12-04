package utils

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/integrii/flaggy"
)

// these get passed in as ldflags by goreleaser
var version string
var commit string
var date string
var binary string

func print_version() string {
	go_version := runtime.Version()
	return fmt.Sprintf("arconn %s built with %s on commit %s at %s", version, go_version, commit, date)
}

func Panic(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func GetSessionName() string {
	user, err := user.Current()
	if err != nil {
		return fmt.Sprintf("%s-%s-unknown-user", binary, version)
	}
	cleanName := strings.Replace(user.Username, "\\", "-", -1)
	return fmt.Sprintf("%s-%s-%s", binary, version, cleanName)
}

func BinaryName() string {
	return binary
}

func ParseFlags() Args {
	args := Args{
		Profile: os.Getenv("AWS_PROFILE"),
		Vault:   os.Getenv("AWS_VAULT"),
	}

	flaggy.String(&args.Profile, "p", "profile", "aws profile to use (defaults to value of AWS_PROFILE env var)")
	flaggy.String(&args.Target, "t", "target", "name of target (required)")
	flaggy.String(&args.Command, "c", "command", "command to pass to ecs targets instead of default shell")
	flaggy.String(&args.PortForward, "P", "port-forward", "port forward map (syntax 80 or 80:80 local:remote)")
	flaggy.String(&args.RemoteHost, "r", "remote-host", "remote host to port forward to")

	flaggy.SetVersion(print_version())
	flaggy.Parse()

	if args.Target == "" {
		flaggy.ShowHelpAndExit("required flag (-t --target) missing")
	}
	if args.Profile == "" && args.Vault == "" {
		flaggy.ShowHelpAndExit("required flag (-p --profile) missing and AWS profile could not be inferred from AWS_PROFILE env var and AWS_VAULT env var is not set")
	}

	return args
}
