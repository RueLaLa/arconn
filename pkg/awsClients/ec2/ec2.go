package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/manifoldco/promptui"
	"github.com/ruelala/arconn/pkg/awsClients/AwsConfig"
	"github.com/ruelala/arconn/pkg/utils"
)

func Lookup(args utils.Args, target utils.Target) utils.Target {
	fmt.Println("searching EC2 for matching instances")
	client := ec2.NewFromConfig(AwsConfig.BuildConfig(args))

	filter := ""
	switch target.Type {
	case "IP":
		filter = "network-interface.addresses.private-ip-address"
	case "NAME":
		filter = "tag:Name"
	}

	id := lookup_with_filter(client, args.Target, filter)
	if id != "" {
		target.ResolvedName = id
		target.Resolved = true
	}

	return target
}

func lookup_with_filter(client *ec2.Client, target string, filter string) string {
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(filter),
				Values: []string{fmt.Sprintf("*%s*", target)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	}
	resp, err := client.DescribeInstances(context.TODO(), input)
	utils.Panic(err)
	instance_id := filter_matches(resp, target)
	return instance_id
}

type Instance struct {
	Name, ID, IP string
}

func filter_matches(output *ec2.DescribeInstancesOutput, target string) string {
	var matches []Instance
	for _, res := range output.Reservations {
		for _, instance := range res.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					matches = append(matches, Instance{
						Name: *tag.Value,
						ID:   *instance.InstanceId,
						IP:   *instance.PrivateIpAddress,
					})
				}
			}
		}
	}

	if len(matches) == 0 {
		fmt.Printf("no matching EC2 instances found for %s\n", target)
		return ""
	} else if len(matches) == 1 {
		fmt.Printf("found %s currently running in EC2\n", matches[0].ID)
		return matches[0].ID
	} else {
		instance_id := prompt_for_choice(matches)
		return instance_id
	}
}

func prompt_for_choice(instances []Instance) string {
	templates := &promptui.SelectTemplates{
		Active:   "\U00002713 {{ .Name | green }} {{ .ID | green }} ({{ .IP | blue }})",
		Inactive: "  {{ .Name | green }} {{ .ID | green }} ({{ .IP | blue }})",
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
