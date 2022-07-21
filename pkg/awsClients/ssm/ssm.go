package ssm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/ruelala/arconn/pkg/awsClients"
	"github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session"
	_ "github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session/portsession"
	_ "github.com/ruelala/arconn/pkg/session-manager-plugin/sessionmanagerplugin/session/shellsession"
	"github.com/ruelala/arconn/pkg/utils"
)

func Lookup(profile, target string) *ssm.Client {
	client := awsClients.SSMClient(profile)
	resp := lookup_instance_in_ssm(client, target)
	instance_online(resp, target)
	return client
}

func lookup_instance_in_ssm(client *ssm.Client, target string) []types.InstanceInformation {
	input := &ssm.DescribeInstanceInformationInput{
		Filters: []types.InstanceInformationStringFilter{
			{
				Key:    aws.String("InstanceIds"),
				Values: []string{target},
			},
		},
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
