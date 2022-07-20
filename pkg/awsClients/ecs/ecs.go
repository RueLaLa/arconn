package ecs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetTarget() {
	input := &ecs.ExecuteCommandInput{
		Command:     aws.String("/bin/bash"),
		Interactive: true,
		Task:        aws.String("b44e720f4ad244439e09f0b71f0d7a93"),
		Cluster:     aws.String("bear"),
	}
	out, _ := client.ExecuteCommand(context.TODO(), input)
	session, _ := json.Marshal(out.Session)
	target := fmt.Sprintf("ecs:%s_%s_%s", "bear", "b44e720f4ad244439e09f0b71f0d7a93", "b44e720f4ad244439e09f0b71f0d7a93-eecc21cbeb3055da3b5fc8923039136b184a71b606b2b4199ca7803072063df4")
	ssmTarget := &ssm.StartSessionInput{
		Target: &target,
	}
	targetJSON, _ := json.Marshal(ssmTarget)
	fmt.Println(fmt.Sprintf("session-manager-plugin %s us-east-1 StartSession %s https://ssm.us-east-1.amazonaws.com", string(session), string(targetJSON)))
}
