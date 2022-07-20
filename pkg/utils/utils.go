package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/integrii/flaggy"
)

func Panic(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ParseFlags(version_string string) (string, string) {
	var profile = os.Getenv("AWS_PROFILE")
	var target string
	flaggy.String(&profile, "p", "profile", "aws profile to use")
	flaggy.String(&target, "t", "target", "resource to target for session connection")
	flaggy.SetVersion(version_string)
	flaggy.Parse()
	return profile, target
}

func TargetType(target string) string {
	switch {
	case regex_match("[[:xdigit:]]{32}", target):
		return "ECS_ID"
	case regex_match("i-[[:xdigit:]]{8}", target):
		return "EC2_ID"
	case regex_match("i-[[:xdigit:]]{17}", target):
		return "EC2_ID"
	case net.ParseIP(target) != nil:
		return "IP"
	default:
		return "NAME"
	}
}

func regex_match(pattern, target string) bool {
	match, err := regexp.MatchString(pattern, target)
	Panic(err)
	if match {
		return true
	} else {
		return false
	}
}
