package utils

import (
	"net"
	"regexp"
)

func TargetType(target string) string {
	switch {
	case regex_match("^mi-[[:xdigit:]]{17}", target):
		return "SSM_MI_ID"
	case regex_match("^i-[[:xdigit:]]{8}", target):
		return "EC2_ID"
	case regex_match("^i-[[:xdigit:]]{17}", target):
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