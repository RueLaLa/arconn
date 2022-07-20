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
		Task:        aws.String("b44e720f4ad244439e09f0b71f0d7a93"),
		Cluster:     aws.String("bear"),
	}
	out, err := client.ExecuteCommand(context.TODO(), input)
	utils.Panic(err)
	session, _ := json.Marshal(out.Session)
	target = fmt.Sprintf("ecs:%s_%s_%s", "bear", "b44e720f4ad244439e09f0b71f0d7a93", "b44e720f4ad244439e09f0b71f0d7a93-eecc21cbeb3055da3b5fc8923039136b184a71b606b2b4199ca7803072063df4")
	return target, string(session)
}
