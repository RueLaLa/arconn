package utils

import (
	"fmt"
	"os"

	"github.com/integrii/flaggy"
)

func Panic(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ParseFlags(version_string string) Args {
	args := Args{}

	args.Profile = os.Getenv("AWS_PROFILE")
	flaggy.String(&args.Profile, "p", "profile", "aws profile to use")

	args.Command = "/bin/bash"
	flaggy.String(&args.Command, "c", "command", "command to pass to ecs targets")

  flaggy.String(&args.PortForward, "P", "portforward", "port forward map to use with ec2 targets (syntax 80 or 80:80 local:remote)")

	flaggy.AddPositionalValue(&args.Target, "target", 1, true, "name of target")
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
