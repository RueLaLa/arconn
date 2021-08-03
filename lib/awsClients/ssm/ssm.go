package ssm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func ssm_client(config aws.Config) *ssm.Client {
	client := ssm.NewFromConfig(config, func(o *ssm.Options) {
		o.Region = "us-east-1"
	})
	return client
}

func Lookup(config aws.Config, target string) *ssm.Client {
	client := ssm_client(config)
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

type RespJSON struct {
	SessionId, StreamUrl, TokenValue string
}

type Target struct {
	Target string
}

func Connect(client *ssm.Client, target string, profile string) {
	input := &ssm.StartSessionInput{Target: &target}
	resp, err := client.StartSession(context.TODO(), input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	j := &RespJSON{
		SessionId:  *resp.SessionId,
		StreamUrl:  *resp.StreamUrl,
		TokenValue: *resp.TokenValue,
	}
	s, _ := json.Marshal(j)

	t := &Target{Target: target}
	st, _ := json.Marshal(t)

	cmd := exec.Command("session-manager-plugin", string(s), "us-east-1", "StartSession", profile, string(st), "https://ssm.us-east-1.amazonaws.com")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
