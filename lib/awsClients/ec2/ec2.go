package ec2

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/manifoldco/promptui"
)

func ec2_client(config aws.Config) *ec2.Client {
	client := ec2.NewFromConfig(config, func(o *ec2.Options) {
		o.Region = "us-east-1"
	})
	return client
}

func TargetType(target string) string {
	// classic ec2 id length
	smatch, _ := regexp.MatchString("i-[0-9a-f]{8}", target)
	if smatch {
		return "ID"
	}

	// new ec2 id length
	lmatch, _ := regexp.MatchString("i-[0-9a-f]{17}", target)
	if lmatch {
		return "ID"
	}

	// IP address
	ip := net.ParseIP(target)
	if ip != nil {
		return "IP"
	}

	return "name"
}

func Lookup(config aws.Config, target string, ttype string) string {
	client := ec2_client(config)

	fmt.Println("searching ec2 for matching instances")
	id := ""
	switch ttype {
	case "IP":
		id = lookup_with_filter(client, target, "network-interface.addresses.private-ip-address")
	case "name":
		id = lookup_with_filter(client, target, "tag:Name")
	}

	return id
}

func lookup_with_filter(client *ec2.Client, target string, filter string) string {
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(filter),
				Values: []string{target},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	}
	resp, _ := client.DescribeInstances(context.TODO(), input)
	instance_id := filter_matches(resp, target)
	return instance_id
}

type Instance struct {
	ID, IP string
}

func filter_matches(output *ec2.DescribeInstancesOutput, target string) string {
	var matches []Instance
	for _, res := range output.Reservations {
		for _, instance := range res.Instances {
			matches = append(matches, Instance{
				ID: *instance.InstanceId,
				IP: *instance.PrivateIpAddress,
			})
		}
	}

	if len(matches) == 0 {
		fmt.Println(fmt.Sprintf("no matching EC2 instances found for %s", target))
		os.Exit(1)
		return ""
	} else if len(matches) == 1 {
		fmt.Println(fmt.Sprintf("found %s currently running in EC2", matches[0].ID))
		return matches[0].ID
	} else {
		instance_id := prompt_for_choice(matches)
		return instance_id
	}
}

func prompt_for_choice(instances []Instance) string {
	templates := &promptui.SelectTemplates{
		Active:   "\U00002713 {{ .ID | green }} ({{ .IP | blue }})",
		Inactive: "  {{ .ID | green }} ({{ .IP | blue }})",
		Selected: "\U00002713 {{ .ID | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select Instance to connect to",
		Items:     instances,
		Templates: templates,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return instances[i].ID
}
