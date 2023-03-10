package utils

import (
	"fmt"
	"os"
	"os/user"
	"runtime"

	"github.com/integrii/flaggy"
)

// these get passed in as ldflags by goreleaser
var version string
var commit string
var date string

func print_version() string {
	go_version := runtime.Version()
	return fmt.Sprintf("arconn %s built with %s on commit %s at %s", version, go_version, commit, date)
}

func Panic(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetSessionName() string {
	user, err := user.Current()
	if err != nil {
		return fmt.Sprintf("arconn-%s-unknown-user", version)
	}
	return fmt.Sprintf("arconn-%s-%s", version, user.Username)
}

func ParseFlags() Args {
	args := Args{
		Command:     "",
		Profile:     os.Getenv("AWS_PROFILE"),
		PortForward: "",
		RemoteHost:  "",
	}

	flaggy.String(&args.Profile, "p", "profile", "aws profile to use")
	flaggy.String(&args.Target, "t", "target", "name of target")
	flaggy.String(&args.Command, "c", "command", "command to pass to ecs targets instead of default shell")
	flaggy.String(&args.PortForward, "P", "port-forward", "port forward map to use with ec2 targets (syntax 80 or 80:80 local:remote)")
	flaggy.String(&args.RemoteHost, "r", "remote-host", "remote host to port forward to")

	flaggy.SetVersion(print_version())
	flaggy.Parse()
	return args
}

func ChunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
