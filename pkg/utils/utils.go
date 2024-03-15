package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/integrii/flaggy"
	"golang.org/x/mod/semver"
)

// these get passed in as ldflags by goreleaser
var version string
var commit string
var date string
var binary string

func IsLatest() {
	var client http.Client
	url := fmt.Sprintf("https://api.github.com/repos/RueLaLa/%s/releases/latest", binary)
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	latest, err := jsonparser.GetString(bodyBytes, "tag_name")
	if err != nil {
		return
	}

	// this function returns 0 for equality, 1 if left value is greater, and -1 if right value is greater
	comparison := semver.Compare(latest, version)
	html_url, err := jsonparser.GetString(bodyBytes, "html_url")
	if err != nil {
		return
	}
	if comparison == 1 {
		fmt.Printf("A new version is available (%s), head to %s to get the latest binary\n", latest, html_url)
	}
}

func Panic(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func UserName() string {
	user, err := user.Current()
	if err != nil {
		return "unknown user"
	}
	return strings.Replace(user.Username, "\\", "-", -1)
}

func BinaryName() string {
	return binary
}

func BinaryVersion() string {
	return version
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

	flaggy.SetVersion(fmt.Sprintf("arconn %s built with %s on commit %s at %s", version, runtime.Version(), commit, date))
	flaggy.Parse()

	if args.Target == "" {
		flaggy.ShowHelpAndExit("required flag (-t --target) missing")
	}
	if args.Profile == "" && args.Vault == "" {
		flaggy.ShowHelpAndExit("required flag (-p --profile) missing and AWS profile could not be inferred from AWS_PROFILE env var and AWS_VAULT env var is not set")
	}

	return args
}
