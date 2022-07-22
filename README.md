arconn: Aws Remote CONNect
===

## Introduction

`arconn` is a colourless, odourless CLI utility that is totally inert to other remote shell scripts. It enables the ability to connect to ECS containers and EC2 hosts remotely, leveraging the SSM Session Manager.

## Installation

Head on over to the [latest release](https://github.com/RueLaLa/arconn/releases/latest) and download the OS appropriate binary archive. Extract the binary and place it somewhere in your `$PATH`

## Usage
```
arconn

  Flags:
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.
    -p --profile   aws profile to use
    -t --target    resource to target for session connection
```

`arconn` accepts two flags, one for the AWS profile name (`-p`) to which the CLI uses for target resolution and SSM connections. The second is the target flag (`-t`). The target flag accepts a few types of inputs. These inputs could be the name of an ECS container or EC2 instance, The IP of an EC2 instance, or an EC2 instance ID.

`arconn` will then attempt to resolve the target input to a real resource running in AWS. It will also check that the target resource is capable of accepting an SSM session. If one or more are found, you are prompted to choose which resource to connect to, and then an SSM session will be started. If an IP or EC2 instance ID are specified, a lot of the resolution logic is skipped and `arconn` will simply ensure it exists and that it can receive SSM sessions. If the target it an arbitrary name, the order of precedence for resolution is ECS first, then EC2.

## Permissions

There are many permissions required to make SSM sessions work. Some exist on the ECS container or EC2 instance, and some exist on the user starting an SSM session. For more information on all thats required to make SSM sessions work, see [this documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-getting-started.html). Additionally for ECS, AWS provides [a nice script](https://github.com/aws-containers/amazon-ecs-exec-checker) that verifies all the necessary permissions are in place for a given container.

## Development

To develop and contribute to this project, go 1.18 or greater is required and [goreleaser](https://goreleaser.com/) is used for publishing releases.
