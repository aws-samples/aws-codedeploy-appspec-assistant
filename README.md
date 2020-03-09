# aws-codedeploy-appspec-assistant
This repo contains Default AppSpec examples and a script to validate an AppSpec file

## CodeDeploy AppSpec examples

Lambda

* lambda-default-appspec-template.yml
 * Contains links and more info
* lambda-default-appspec-template-basic.json
* lambda-default-appspec-template-advanced.json

ECS

* ecs-default-appspec-template.yml
 * Contains links and more info
* ecs-default-appspec-template-basic.json
* ecs-default-appspec-template-advanced.json

EC2/OnPrem

* ec2-and-onPrem-default-appspec-template.yml
 * Contains links and more info
* ec2-and-onPrem-default-appspec-template-basic.json
* ec2-and-onPrem-default-appspec-template-advanced.json

#### CodeDeploy AppSpec Documentation links

Lambda
* [appspec-reference-lambda](https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file.html#appspec-reference-lambda)

ECS
* [appspec-reference-ecs](https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file.html#appspec-reference-ecs)

EC2/OnPrem
* [appspec-reference-server](https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file.html#appspec-reference-server)

## Getting started with the Validation Assistant Script

You must have Go `1.13` or later installed.

## Releases

#### V1

* Basic JSON and YAML parsing and file extension validation

#### V2 - Upcoming

* AppSpec file structure validation and limited value validation

## Usage

General workflow:

```
# Validate an AppSPec file
$ go run main.go validate --filePath <FILE_NAME>
```

Run `go run main.go --help` for full usage.

# Install on your machine

You can install appSpecAssistant on your machine to run the commands without having to build the binary every time.

```
# Install the binary to your GOPATH
$ go build -o $GOPATH/bin/appSpecAssistant
# Call appSpecAssistant
$ appSpecAssistant validate -f <FILE_NAME>
```

You can also download the pre-built binary from this repo

```
$ ./appSpecAssistant validate -f <FILE_NAME>
```

## Development

This CLI uses [cobra](https://github.com/spf13/cobra). See documentation for more information. Install the cobra CLI to auto-generate code for new CLI commands.
* Used to auto-generate code for CLI
* root.go is the default file for Cobra
* Cobra dependencies:
 * github.com/mitchellh/go-homedir
 * github.com/spf13/viper

```
# Add a new command
$ cobra add commandName
```

The assistant uses YAML support: [go-yaml V3](https://github.com/go-yaml/yaml/tree/v3). See documentation for more information. Install the YAML package for development.
* Used to process YAML configuration file syntax

Build the project:

```
$ go build
```

Build and run in one shot:

```
$ go run main.go
```

Run unit tests

```
$ go test -v ./pkg/*
```

Format code:

```
$ ./gofmt
```
