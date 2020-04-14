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

## Validation Assistant Script

### **WARNING**

This script validates basic expected syntax for the AppSpec. It does not guarantee that your deployment will succeed nor that your AppSpec will do what you expect it to.

### Run validator script without having to deal with Golang

You can download a pre-built binary from this repo

Linux:
```
$ ./appSpecAssistant validate --filePath <FILE_PATH> --computePlatform <[server, lambda, or ecs]>
```

Windows:
```
./appSpecAssistantForWindows.exe validate --filePath <FILE_PATH> --computePlatform <[server, lambda, or ecs]>
```

### Capabilities of the Validation Assistant Script

#### March 2020

* Validates file extension
* Validates YAML and JSON syntax
* Does basic validation for ECS, Lambda, and Server (EC2/On-Prem) AppSpec file content
  * Validates type of values
  * Validates values that are required
  * Validates values that require specific strings
  * Validates Hooks used
  * And more basic validations

### Building the Golang Validation Assistant Script Locally

You must have Go `1.13` or later installed.

#### General workflow:

```
# Validate an AppSPec file
$ go run main.go validate --filePath <FILE_PATH> --computePlatform <[server, lambda, or ecs]>
```

Run `go run main.go --help` for full usage.

#### Install Script on your machine

You can install appSpecAssistant on your machine to run the commands without having to build the binary every time.

For Linux on Linux:
```
# Install the binary to your GOPATH (or you can just build it in the current directory $ go build -o appSpecAssistant)
$ go build -o $GOPATH/bin/appSpecAssistant
# Call appSpecAssistant
$ appSpecAssistant validate --filePath <FILE_PATH> --computePlatform <[server, lambda, or ecs]>
```

For Windows on Linux:
https://github.com/golang/go/wiki/WindowsCrossCompiling
```
GOOS=windows GOARCH=386 go build -o appSpecAssistantForWindows.exe main.go
```

For Windows on Windows:
https://github.com/golang/go/wiki/WindowsCrossCompiling
```
In cmd.exe instead of PowerShell:

$ set GOOS=windows
$ set GOARCH=386
$ go build -o hello.exe hello.go
```

### Development (adding changes to the script)

You must have Go `1.13` or later installed.

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
