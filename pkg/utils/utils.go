package utils

import (
	"fmt"
	"os"
	"os/user"

	"github.com/integrii/flaggy"
)

func Panic(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetUser() string {
	user, err := user.Current()
	if err != nil {
		return "unknown user"
	}
	return user.Username
}

func ParseFlags(version_string string) Args {
	args := Args{
		Command:     "/bin/bash",
		Profile:     os.Getenv("AWS_PROFILE"),
		PortForward: "",
		RemoteHost:  "",
	}

	flaggy.String(&args.Profile, "p", "profile", "aws profile to use")
	flaggy.String(&args.Target, "t", "target", "name of target")
	flaggy.String(&args.Command, "c", "command", "command to pass to ecs targets")
	flaggy.String(&args.PortForward, "P", "port-forward", "port forward map to use with ec2 targets (syntax 80 or 80:80 local:remote)")
	flaggy.String(&args.RemoteHost, "r", "remote-host", "remote host to port forward to")

	flaggy.SetVersion(version_string)
	flaggy.Parse()
	return args
}

func ChunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
