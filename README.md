arconn: Aws Remote CONNect
===

## Introduction
`arconn` is a colourless, odourless CLI utility that is totally inert to other remote shell scripts. It enables the ability to connect to ECS containers and EC2 hosts remotely, leveraging the SSM Session Manager.

## Installation
Head on over to the [latest release](https://github.com/RueLaLa/arconn/releases/latest) and download the OS and architecture appropriate zip file. Extract the binary and place it somewhere in your `$PATH`.

## Usage
```
arconn

  Flags:
       --version        Displays the program version string.
    -h --help           Displays help with available flag, subcommand, and positional value parameters.
    -p --profile        aws profile to use (defaults to value of AWS_PROFILE env var)
    -t --target         name of target (required)
    -c --command        command to pass to ecs targets instead of default shell
    -P --port-forward   port forward map (syntax 80 or 80:80 local:remote)
    -r --remote-host    remote host to port forward to
```

## Examples
Connecting to a simple EC2 host:
```
arconn -p myProfile -t i-12345678
```

Port Forwarding to an ECS container:
```
arconn -p myProfile -t myContainer -p 8080:8080
```

Port Forwarding to a remote host through an EC2 host:
```
arconn -p myProfile -t 10.0.1.66 -p 3306:3306 -r myDatabase.domain
```

## Types & Searching
`arconn` attempts to resolve the target input to a real resource running in AWS. It will also check that the target resource is capable of accepting an SSM session for certain target types. If more than one target is found, you are prompted to choose one. If the input target is an EC2 instance ID or SSM managed instance ID, the resolution logic is skipped and `arconn` will simply ensure it exists and that it can receive SSM sessions. If the target is an arbitrary name, `arconn` will search in ECS first, then EC2, and finally SSM. `arconn` also supports multiple session types for each of the resolved targets. Here is a matrix of each target and their supported session types and input formats.

Target type | Sessions | Custom commands | Port forwarding | Port forwarding to remote hosts | Supported Input formats
:---------- | :------- | :-------------- | :-------------- | :------------------------------ | :----------------------
EC2 | :white_check_mark: | :x: | :white_check_mark: | :white_check_mark: | <ul><li>Instance ID</li><li>IP Address</li><li>Name</li></ul>
ECS | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | <ul><li>Name</li></ul>
SSM Managed Instance | :white_check_mark: | :question: | :question: | :question: | <ul><li>Managed Instance Id</li><li>Name</li></ul>

## Permissions
`arconn` uses the below permissions to some extent across the application:
- `ec2:DescribeInstances`
- `ecs:DescribeTasks`
- `ecs:ExecuteCommand`
- `ecs:ListClusters`
- `ecs:ListTasks`
- `ssm:DescribeInstanceInformation`
- `ssm:StartSession`

For more information on all that's required to make SSM sessions work, see [this documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-getting-started.html). Additionally for ECS, AWS provides [a nice script](https://github.com/aws-containers/amazon-ecs-exec-checker) that verifies all the necessary permissions are in place for a given container.

## Development
To develop and contribute to this project, refer to the `go.mod` file for dependencies, and [goreleaser](https://goreleaser.com/) is used for publishing releases.
