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
	"github.com/ruelala/arconn/pkg/awsClients"
	"github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session"
	_ "github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session/portsession"
	_ "github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session/shellsession"
	"github.com/ruelala/arconn/pkg/utils"
)

func Lookup(profile, target string, use_filter bool) string {
	client := awsClients.SSMClient(profile)
	resp := lookup_instance_in_ssm(client, target, use_filter)
	instance := ""
	if use_filter {
		instance = *resp[0].InstanceId
	} else {
		instance = filter_matches(resp, target)
	}
	instance_online(resp, target)
	return instance
}

type Instance struct {
	Name, ID string
}

func filter_matches(instances []types.InstanceInformation, target string) string {
	var matches []Instance
	for _, instance := range instances {
		if *instance.Name == target {
			matches = append(matches, Instance{
				Name: *instance.Name,
				ID:   *instance.InstanceId,
			})
		}
	}

	if len(matches) == 0 {
		fmt.Println(fmt.Sprintf("no matching SSM instances found for %s", target))
		return ""
	} else if len(matches) == 1 {
		fmt.Println(fmt.Sprintf("found %s currently running in SSM", matches[0].ID))
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

func lookup_instance_in_ssm(client *ssm.Client, target string, use_filter bool) []types.InstanceInformation {
	input := &ssm.DescribeInstanceInformationInput{}
	if use_filter {
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
		fmt.Println(fmt.Sprintf("%s is not currently registered with SSM, make sure agent is configured and online", target))
		os.Exit(1)
	}
	return resp.InstanceInformationList
}

func instance_online(resp []types.InstanceInformation, target string) bool {
	if resp[0].PingStatus == types.PingStatusOnline {
		return true
	} else {
		fmt.Println(fmt.Sprintf("%s is registered with SSM, but the agent is offline", target))
		os.Exit(1)
		return false
	}
}

func Connect(profile, session_json, target string) {
	client := awsClients.SSMClient(profile)

	input := &ssm.StartSessionInput{Target: &target}
	target_json, _ := json.Marshal(input)

	if session_json == "" {
		resp, err := client.StartSession(context.TODO(), input)
		utils.Panic(err)
		session_raw, _ := json.Marshal(resp)
		session_json = string(session_raw)
	}

	args := []string{"session-manager-plugin", string(session_json), "us-east-1", "StartSession", profile, string(target_json), "https://ssm.us-east-1.amazonaws.com"}
	session.ValidateInputAndStartSession(args, os.Stdout)
}
