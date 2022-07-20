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
	if len(resp) == 0 {
		fmt.Println(fmt.Sprintf("%s is not currently registered with SSM, make sure agent is configured and online", target))
		os.Exit(1)
	}
	if instance_online(resp) == false {
		fmt.Println(fmt.Sprintf("%s is registered with SSM, but the agent is offline", target))
		os.Exit(1)
	}
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
	resp, _ := client.DescribeInstanceInformation(context.TODO(), input)
	return resp.InstanceInformationList
}

func instance_online(resp []types.InstanceInformation) bool {
	if resp[0].PingStatus == types.PingStatusOnline {
		return true
	} else {
		return false
	}
}

type SessionInfo struct {
	SessionId, StreamUrl, TokenValue string
}

type SessionTarget struct {
	Target string
}

func Connect(client *ssm.Client, profile string, target string) {
	input := &ssm.StartSessionInput{Target: &target}
	resp, err := client.StartSession(context.TODO(), input)
	utils.Panic(err)

	session_info := &SessionInfo{
		SessionId:  *resp.SessionId,
		StreamUrl:  *resp.StreamUrl,
		TokenValue: *resp.TokenValue,
	}
	session_json, _ := json.Marshal(session_info)

	target_struct := &SessionTarget{Target: target}
	target_json, _ := json.Marshal(target_struct)

	args := []string{"session-manager-plugin", string(session_json), "us-east-1", "StartSession", profile, string(target_json), "https://ssm.us-east-1.amazonaws.com"}
	session.ValidateInputAndStartSession(args, os.Stdout)
}
