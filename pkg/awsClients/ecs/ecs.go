package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/manifoldco/promptui"
	"github.com/ruelala/arconn/pkg/awsClients"
	"github.com/ruelala/arconn/pkg/utils"
)

func Lookup(args utils.Args, target utils.Target) utils.Target {
	fmt.Println("searching ECS for matching tasks")
	client := awsClients.ECSClient(args.Profile)
	clusters := list_clusters(client)
	tasks := find_matching_tasks(client, clusters, args, target.Type)

	chosen_task := FTask{}
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

func construct_target(task FTask) string {
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

type FTask struct {
	ClusterArn arn.ARN
	TaskArn    arn.ARN
	TaskName   string
	RuntimeId  string
}

func find_matching_tasks(client *ecs.Client, clusters []string, args utils.Args, ttype string) []FTask {
	tasks := List_tasks(client, clusters)
	return filter_tasks(tasks, args, ttype)
}

func List_tasks(client *ecs.Client, clusters []string) []types.Task {
	next_token := ""
	all_tasks := []types.Task{}
	for _, cluster := range clusters {
		paginator := ecs.NewListTasksPaginator(
			client,
			&ecs.ListTasksInput{
				Cluster:    aws.String(cluster),
				MaxResults: aws.Int32(100),
				NextToken:  aws.String(next_token),
			},
		)
		for paginator.HasMorePages() {
			resp, err := paginator.NextPage(context.TODO())
			utils.Panic(err)
			if len(resp.TaskArns) == 0 {
				break
			}
			lresp, err := client.DescribeTasks(
				context.TODO(),
				&ecs.DescribeTasksInput{
					Cluster: aws.String(cluster),
					Tasks:   resp.TaskArns,
					Include: []types.TaskField{types.TaskFieldTags},
				},
			)
			utils.Panic(err)
			for _, task := range lresp.Tasks {
				all_tasks = append(all_tasks, task)
			}
		}
	}
	return all_tasks
}

func filter_tasks(tasks []types.Task, args utils.Args, ttype string) []FTask {
	filtered_tasks := []FTask{}
	for _, task := range tasks {
		if len(task.Containers) == 0 {
			continue
		}
		if *task.LastStatus != "RUNNING" {
			continue
		}
		r, _ := regexp.Compile(fmt.Sprintf(".*%s.*", args.Target))
		if (ttype == "ECS_ID") && (strings.Split(args.Target, "_")[2] != *task.Containers[0].RuntimeId) {
			continue
		} else if (ttype == "NAME") && !(r.Match([]byte(*task.Containers[0].Name))) {
			continue
		}
		if task.EnableExecuteCommand == true {
			carn, _ := arn.Parse(*task.ClusterArn)
			tarn, _ := arn.Parse(*task.TaskArn)
			filtered_tasks = append(filtered_tasks,
				FTask{
					ClusterArn: carn,
					TaskArn:    tarn,
					TaskName:   *task.Containers[0].Name,
					RuntimeId:  *task.Containers[0].RuntimeId,
				})
		}
	}
	return filtered_tasks
}

func prompt_for_choice(tasks []FTask) FTask {
	templates := &promptui.SelectTemplates{
		Active:   "\U00002713 {{ .TaskName | green }} {{ .ClusterArn.Resource | blue }} {{ .TaskArn.Resource | blue }}",
		Inactive: "  {{ .TaskName | green }} {{ .ClusterArn.Resource | blue }} {{ .TaskArn.Resource | blue }}",
		Selected: "\U00002713 {{ .TaskName | green }} {{ .ClusterArn.Resource | blue }} {{ .TaskArn.Resource | blue }}",
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
