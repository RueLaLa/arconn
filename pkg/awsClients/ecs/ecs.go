package ecs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/ruelala/arconn/pkg/awsClients"
	"github.com/ruelala/arconn/pkg/utils"
)

func Lookup(profile, target string) (string, string) {
	client := awsClients.ECSClient(profile)
	input := &ecs.ExecuteCommandInput{
		Command:     aws.String("/bin/bash"),
		Interactive: true,
		Task:        aws.String("4a21cb9e457b4c2bd027f7cc928962da"),
		Cluster:     aws.String("test"),
	}
	out, err := client.ExecuteCommand(context.TODO(), input)
	utils.Panic(err)
	session, _ := json.Marshal(out.Session)

	target = fmt.Sprintf("ecs:%s_%s_%s", "test", "4a21cb9e457b4c2bd027f7cc928962da", "4a21cb9e457b4c2bd027f7cc928962da-2e01cfa350c14834c8b5d78d07fcbd519f7d552b091b3280b4e1953cc6233cad")
	return target, string(session)
}
