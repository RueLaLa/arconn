package ssm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/manifoldco/promptui"
	"github.com/ruelala/arconn/pkg/awsClients/AwsConfig"
	"github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session"
	_ "github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session/portsession"
	_ "github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session/shellsession"
	"github.com/ruelala/arconn/pkg/utils"
)

func Lookup(args utils.Args, target utils.Target) utils.Target {
	client := ssm.NewFromConfig(AwsConfig.BuildConfig(args))
	ssm_target := ""
	if target.Resolved {
		ssm_target = target.ResolvedName
	} else {
		ssm_target = args.Target
	}
	resp := lookup_instance_in_ssm(client, ssm_target)
	if len(resp) > 1 {
		filtered := filter_matches(resp, args.Target)
		if filtered == "" {
			return target
		} else {
			target.ResolvedName = filtered
		}
	} else {
		target.ResolvedName = *resp[0].InstanceId
	}
	instance_online(resp, target.ResolvedName)
	target.Resolved = true
	return target
}

type Instance struct {
	Name, ID string
}

func filter_matches(instances []types.InstanceInformation, target string) string {
	var matches []Instance
	for _, instance := range instances {
		if instance.Name == nil {
			continue
		}
		if *instance.Name == target {
			matches = append(matches, Instance{
				Name: *instance.Name,
				ID:   *instance.InstanceId,
			})
		}
	}

	if len(matches) == 0 {
		fmt.Printf("no matching SSM instances found for %s\n", target)
		return ""
	} else if len(matches) == 1 {
		fmt.Printf("found %s currently running in SSM\n", matches[0].ID)
		return matches[0].ID
	} else {
		instance_id := prompt_for_choice(matches)
		return instance_id
	}
}

func prompt_for_choice(instances []Instance) string {
	templates := &promptui.SelectTemplates{
		Active:   "\U00002713 {{ .Name | green }} {{ .ID | blue }}",
		Inactive: "  {{ .Name | green }} {{ .ID | blue }}",
		Selected: "\U00002713 {{ .ID | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select Instance to connect to",
		Items:     instances,
		Size:      10,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	utils.Panic(err)

	return instances[i].ID
}

func lookup_instance_in_ssm(client *ssm.Client, target string) []types.InstanceInformation {
	input := &ssm.DescribeInstanceInformationInput{}
	if utils.TargetType(target) == "EC2_ID" {
		input = &ssm.DescribeInstanceInformationInput{
			Filters: []types.InstanceInformationStringFilter{
				{
					Key:    aws.String("InstanceIds"),
					Values: []string{target},
				},
			},
		}
	}
	resp, err := client.DescribeInstanceInformation(context.TODO(), input)
	utils.Panic(err)

	if len(resp.InstanceInformationList) == 0 {
		utils.Panic(fmt.Errorf("%s is not currently registered with SSM, make sure agent is configured and online", target))
	}
	return resp.InstanceInformationList
}

func instance_online(resp []types.InstanceInformation, target string) bool {
	if resp[0].PingStatus == types.PingStatusOnline {
		return true
	} else {
		utils.Panic(fmt.Errorf("%s is registered with SSM, but the agent is offline", target))
		return false
	}
}

func Connect(args utils.Args, target utils.Target) {
	config := AwsConfig.BuildConfig(args)
	client := ssm.NewFromConfig(config)

	input := &ssm.StartSessionInput{}
	input.Reason = aws.String(fmt.Sprintf("%s %s session", utils.BinaryName(), utils.BinaryVersion()))
	input.Target = &target.ResolvedName
	if len(target.PortForwarding) > 0 {
		input.Reason = aws.String(fmt.Sprintf("%s %s port forward session", utils.BinaryName(), utils.BinaryVersion()))
		type param map[string][]string
		p := make(param)
		p["portNumber"] = []string{target.PortForwarding[1]}
		p["localPortNumber"] = []string{target.PortForwarding[0]}
		if target.RemoteHost != "" {
			p["host"] = []string{target.RemoteHost}
			input.DocumentName = aws.String("AWS-StartPortForwardingSessionToRemoteHost")
		} else {
			input.DocumentName = aws.String("AWS-StartPortForwardingSession")
		}
		input.Parameters = p
	}
	target_json, _ := json.Marshal(input)

	if target.SessionInfo == "" {
		resp, err := client.StartSession(context.TODO(), input)
		utils.Panic(err)
		session_raw, _ := json.Marshal(resp)
		target.SessionInfo = string(session_raw)
	}

	connect_args := []string{
		"session-manager-plugin",
		target.SessionInfo,
		config.Region,
		"StartSession",
		args.Profile,
		string(target_json),
		fmt.Sprintf("https://ssm.%s.amazonaws.com", config.Region),
	}
	session.ValidateInputAndStartSession(connect_args, os.Stdout)
}
