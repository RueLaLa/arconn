arconn: **A**ws **R**emote **CONN**ect
===

## Introduction
`arconn` is a colourless, odourless CLI utility that is totally inert to other remote shell scripts. It enables the ability to connect to ECS containers and EC2 hosts remotely, leveraging the SSM Session Manager.

## Installation
Head on over to the [latest release](https://github.com/RueLaLa/arconn/releases/latest) and download the OS appropriate binary archive. Extract the binary and place it somewhere in your `$PATH`

## Usage
```
arconn

  Flags:
       --version        Displays the program version string.
    -h --help           Displays help with available flag, subcommand, and positional value parameters.
    -p --profile        aws profile to use
    -t --target         name of target
    -c --command        command to pass to ecs targets (default: /bin/bash)
    -P --port-forward   port forward map to use with ec2 targets (syntax 80 or 80:80 local:remote)
    -r --remote-host    remote host to port forward to
```

`arconn` accepts multiple flags but only one is required: the `-t` target flag. The target argument accepts a few types of inputs. These inputs could be the name of an ECS container or EC2 instance, The IP of an EC2 instance, or an EC2 instance ID. More information about acceptable inputs and the search logic is described below. The remaining flags are optional but can be useful. `-p` is for the AWS profile name which the CLI uses for target resolution and SSM connections but the `AWS_PROFILE` environment variable can also be used.


## Types & Searching
`arconn` can accept many types of inputs to attempt to find a target for starting a session. Depending on which type you define, the search logic behaves. If it requires searching, searching will be done in a specific order. Target type computing is done via regex.

- for EC2 instance IDs, no searching is done, and `arconn` simply checks that its online and connected to SSM.
- for an IP address, EC2 is searched to find an associated instance ID, then moves on to start a session.
- for an SSM managed instance not running in AWS, the name is looked up in SSM and then moves on to start a session.
- if it is any other string, it is considered a name. Names are then searched across ECS, EC2, and SSM, __in that order__.

`arconn` will then attempt to resolve the target input to a real resource running in AWS. It will also check that the target resource is capable of accepting an SSM session. If one or more are found, you are prompted to choose which resource to connect to, and then an SSM session will be started. If an IP or EC2 instance ID are specified, a lot of the resolution logic is skipped and `arconn` will simply ensure it exists and that it can receive SSM sessions. If the target it an arbitrary name, the order of precedence for resolution is ECS first, then EC2.

## Permissions
There are many permissions required to make SSM sessions work. Some exist on the ECS container or EC2 instance, and some exist on the user starting an SSM session. For more information on all that's required to make SSM sessions work, see [this documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-getting-started.html). Additionally for ECS, AWS provides [a nice script](https://github.com/aws-containers/amazon-ecs-exec-checker) that verifies all the necessary permissions are in place for a given container.

## Development
To develop and contribute to this project, refer to the `go.mod` file for dependencies, and [goreleaser](https://goreleaser.com/) is used for publishing releases.
