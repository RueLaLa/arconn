package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/manifoldco/promptui"
	"github.com/ruelala/arconn/pkg/awsClients"
	"github.com/ruelala/arconn/pkg/utils"
)

func Lookup(args utils.Args, target utils.Target) utils.Target {
	fmt.Println("searching ECS for matching tasks")
	client := awsClients.ECSClient(args.Profile)
	clusters := list_clusters(client)
	tasks := find_matching_tasks(client, clusters, args.Target)

	chosen_task := Task{}
	if len(tasks) == 0 {
		fmt.Println("no tasks matching target running in ECS with execute command capabilities")
		return target
	} else if len(tasks) == 1 {
		chosen_task = tasks[0]
	} else {
		chosen_task = prompt_for_choice(tasks)
	}

	if args.Command != "" {
		input := &ecs.ExecuteCommandInput{
			Command:     aws.String(args.Command),
			Interactive: true,
			Task:        aws.String(chosen_task.TaskArn.String()),
			Cluster:     aws.String(chosen_task.ClusterArn.String()),
		}
		out, err := client.ExecuteCommand(context.TODO(), input)
		utils.Panic(err)
		session_info, _ := json.Marshal(out.Session)
		target.SessionInfo = string(session_info)
	}

	target.ResolvedName = construct_target(chosen_task)
	target.Type = "ECS"
	target.Resolved = true
	return target
}

func construct_target(task Task) string {
	cluster_name := strings.TrimPrefix(task.ClusterArn.Resource, "cluster/")
	task_id := strings.TrimPrefix(
		task.TaskArn.Resource,
		fmt.Sprintf("task/%s/", cluster_name),
	)
	target := fmt.Sprintf(
		"ecs:%s_%s_%s",
		cluster_name,
		task_id,
		task.RuntimeId,
	)
	return target
}

func clean(nested [][]Task) []Task {
	ret := []Task{}
	for _, nest := range nested {
		for _, elem := range nest {
			ret = append(ret, elem)
		}
	}
	return ret
}

func list_clusters(client *ecs.Client) []string {
	next_token := ""
	clusters := []string{}
	input := &ecs.ListClustersInput{
		MaxResults: aws.Int32(100),
		NextToken:  aws.String(next_token),
	}
	paginator := ecs.NewListClustersPaginator(client, input)
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(context.TODO())
		utils.Panic(err)
		for _, arn := range resp.ClusterArns {
			clusters = append(clusters, arn)
		}
	}
	return clusters
}

type Task struct {
	ClusterArn arn.ARN
	TaskArn    arn.ARN
	TaskName   string
	RuntimeId  string
}

func find_matching_tasks(client *ecs.Client, clusters []string, target string) []Task {
	next_token := ""
	all_tasks := [][]Task{}
	for _, cluster := range clusters {
		cluster_tasks := []string{}
		input := &ecs.ListTasksInput{
			Cluster:    aws.String(cluster),
			MaxResults: aws.Int32(100),
			NextToken:  aws.String(next_token),
		}
		paginator := ecs.NewListTasksPaginator(client, input)
		for paginator.HasMorePages() {
			resp, err := paginator.NextPage(context.TODO())
			utils.Panic(err)
			for _, arn := range resp.TaskArns {
				cluster_tasks = append(cluster_tasks, arn)
			}
		}
		all_tasks = append(all_tasks,
			describe_tasks(client, cluster, cluster_tasks, target))
	}
	return clean(all_tasks)
}

func describe_tasks(client *ecs.Client, cluster string, tasks []string, target string) []Task {
	split_tasks := utils.ChunkBy(tasks, 100)
	task_info := []Task{}
	for _, set := range split_tasks {
		if len(set) == 0 {
			continue
		}
		input := &ecs.DescribeTasksInput{
			Cluster: aws.String(cluster),
			Tasks:   set,
		}
		resp, err := client.DescribeTasks(context.TODO(), input)
		utils.Panic(err)
		for _, task := range resp.Tasks {
			if len(task.Containers) == 0 {
				continue
			}
			if (*task.Containers[0].Name == target) && (task.EnableExecuteCommand == true) {
				carn, _ := arn.Parse(cluster)
				tarn, _ := arn.Parse(*task.TaskArn)
				task_info = append(task_info,
					Task{
						ClusterArn: carn,
						TaskArn:    tarn,
						TaskName:   *task.Containers[0].Name,
						RuntimeId:  *task.Containers[0].RuntimeId,
					})
			}
		}
	}
	return task_info
}

func prompt_for_choice(tasks []Task) Task {
	templates := &promptui.SelectTemplates{
		Active:   "\U00002713 {{ .TaskName | green }} {{ .TaskArn | blue }}",
		Inactive: "  {{ .TaskName | green }} {{ .TaskArn | blue }}",
		Selected: "\U00002713 {{ .TaskName | green }} {{ .TaskArn | blue }}",
	}

	prompt := promptui.Select{
		Label:     "Select ECS task to connect to",
		Items:     tasks,
		Size:      10,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	utils.Panic(err)

	return tasks[i]
}
