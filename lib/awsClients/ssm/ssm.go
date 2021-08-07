package ssm

import (
	"context"
	//"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/mmmorris1975/ssm-session-client/ssmclient"
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

// type RespJSON struct {
// 	SessionId, StreamUrl, TokenValue string
// }

// type Target struct {
// 	Target string
// }

func Connect(config aws.Config, target string) {
	// when https://github.com/aws/session-manager-plugin/issues/1 is resolved, switch to using
	// code below and structs above
	//
	// input := &ssm.StartSessionInput{Target: &target}
	// resp, err := client.StartSession(context.TODO(), input)
	// if err != nil {
	//    fmt.Println(err)
	//    os.Exit(1)
	// }
	// session_info := &RespJSON{
	//   SessionId:  *resp.SessionId,
	//   StreamUrl:  *resp.StreamUrl,
	//   TokenValue: *resp.TokenValue,
	// }
	// session_json, _ := json.Marshal(session_info)
	// target_struct := &Target{Target: target}
	// target_json, _ := json.Marshal(target_struct)
	// args := []string{"session-manager-plugin", string(session_json), "us-east-1", "StartSession", profile, string(target_json), "https://ssm.us-east-1.amazonaws.com"}
	// session.ValidateInputAndStartSession(args)

	ssmclient.ShellSession(config, target)
}
