package assistant

// ECS AppSpec strings for Unit Tests
var ecsOutputStr = `{0 [{{AWS::ECS::Service {[Your task definition arn] {[Your container Name] 8000} [Version number, ex: 1.3.0] {{[SubnetId1 SubnetId2] [ecs-security-group-1] DISABLED}}}}}] [map[BeforeInstall:BeforeInstallHookLambdaFunctionName] map[AfterInstall:AfterInstallHookLambdaFunctionName] map[AfterAllowTestTraffic:AfterAllowTestTrafficHookLambdaFunctionName] map[BeforeAllowTraffic:SanityTestHookLambdaFunctionName] map[AfterAllowTraffic:ValidationTestHookLambdaFunctionName]]}`

var ecsJsonString = `{
  "version": 0.0,
  "Resources": [
    {
      "TargetService": {
        "Type": "AWS::ECS::Service",
        "Properties": {
          "TaskDefinition": "[Your task definition arn]",
          "LoadBalancerInfo": {
            "ContainerName": "[Your container Name]",
            "ContainerPort": 8000
          },
          "PlatformVersion": "[Version number, ex: 1.3.0]",
          "NetworkConfiguration": {
            "AwsvpcConfiguration": {
              "Subnets": [
                "SubnetId1",
                "SubnetId2"
              ],
              "SecurityGroups": [
                "ecs-security-group-1"
              ],
              "AssignPublicIp": "DISABLED"
            }
          }
        }
      }
    }
  ],
  "Hooks": [
    {
      "BeforeInstall": "BeforeInstallHookLambdaFunctionName"
    },
    {
      "AfterInstall": "AfterInstallHookLambdaFunctionName"
    },
    {
      "AfterAllowTestTraffic": "AfterAllowTestTrafficHookLambdaFunctionName"
    },
    {
      "BeforeAllowTraffic": "SanityTestHookLambdaFunctionName"
    },
    {
      "AfterAllowTraffic": "ValidationTestHookLambdaFunctionName"
    }
  ]
}`

var ecsYamlString = `version: 0.0
Resources:
  - TargetService:
      Type: AWS::ECS::Service
      Properties:
        TaskDefinition: "[Your task definition arn]"
        LoadBalancerInfo:
          ContainerName: "[Your container Name]"
          ContainerPort: 8000
        PlatformVersion: "[Version number, ex: 1.3.0]"
        NetworkConfiguration:
          AwsvpcConfiguration:
            Subnets: ["SubnetId1","SubnetId2"]
            SecurityGroups: ["ecs-security-group-1"]
            AssignPublicIp: "DISABLED"
Hooks:
  - BeforeInstall: "BeforeInstallHookLambdaFunctionName"
  - AfterInstall: "AfterInstallHookLambdaFunctionName"
  - AfterAllowTestTraffic: "AfterAllowTestTrafficHookLambdaFunctionName"
  - BeforeAllowTraffic: "SanityTestHookLambdaFunctionName"
  - AfterAllowTraffic: "ValidationTestHookLambdaFunctionName"`

// Lambda AppSpec strings for Unit Tests
var lambdaOutputStr = `{0 [map[myLambdaFunction:{AWS::Lambda::Function {myLambdaFunction myLambdaFunctionAlias 1 2}}]] [map[BeforeAllowTraffic:<SanityTestHookLambdaFunctionName>] map[AfterAllowTraffic:<ValidationTestHookLambdaFunctionName>]]}`

var lambdaJsonString = `{
  "version": 0.0,
  "Resources": [
    {
      "myLambdaFunction": {
        "Type": "AWS::Lambda::Function",
        "Properties": {
          "Name": "myLambdaFunction",
          "Alias": "myLambdaFunctionAlias",
          "CurrentVersion": "1",
          "TargetVersion": "2"
        }
      }
    }
  ],
  "Hooks": [
    {
      "BeforeAllowTraffic": "<SanityTestHookLambdaFunctionName>"
    },
    {
      "AfterAllowTraffic": "<ValidationTestHookLambdaFunctionName>"
    }
  ]
}`

var lambdaYamlString = `version: 0.0
Resources:
  - myLambdaFunction:
      Type: AWS::Lambda::Function
      Properties:
        Name: "myLambdaFunction"
        Alias: "myLambdaFunctionAlias"
        CurrentVersion: "1"
        TargetVersion: "2"
Hooks:
  - BeforeAllowTraffic: "<SanityTestHookLambdaFunctionName>"
  - AfterAllowTraffic: "<ValidationTestHookLambdaFunctionName>"`

// EC2/OnPrem (Server) AppSpec strings for Unit Tests
var serverOutputStr = `{0 linux [{source-file-location destination-file-location}] [{object-specification pattern-specification exception-specification owner-account-name group-name mode-specification [acls-specification] {user-specification type-specification range-specification} [file]}] map[ApplicationStop:[{script-location 10 user-name} {script-location 10 user-name}] BeforeInstall:[{script-location 10 user-name}]]}`

var serverJsonString = `{
  "version": 0.0,
  "os": "linux",
  "files": [
    {
      "source": "source-file-location",
      "destination": "destination-file-location"
    }
  ],
  "permissions": [
    {
      "object": "object-specification",
      "pattern": "pattern-specification",
      "except": "exception-specification",
      "owner": "owner-account-name",
      "group": "group-name",
      "mode": "mode-specification",
      "acls": [
        "acls-specification"
      ],
      "context": {
        "user": "user-specification",
        "type": "type-specification",
        "range": "range-specification"
      },
      "type": [
        "file"
      ]
    }
  ],
  "hooks": {
    "ApplicationStop": [
      {
        "location": "script-location",
        "timeout": "10",
        "runas": "user-name"
      },
      {
        "location": "script-location",
        "timeout": "10",
        "runas": "user-name"
      }
    ],
    "BeforeInstall": [
      {
        "location": "script-location",
        "timeout": "10",
        "runas": "user-name"
      }
    ]
  }
}`

var serverYamlString = `version: 0.0
os: linux
# https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file-structure-files.html
files:
 - source: source-file-location
   destination: destination-file-location
# https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file-structure-permissions.html
permissions:
  - object: object-specification
    pattern: pattern-specification
    except: exception-specification
    owner: owner-account-name
    group: group-name
    mode: mode-specification
    acls: 
      - acls-specification 
    context:
      user: user-specification
      type: type-specification
      range: range-specification
    type:
      - file
# https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file-structure-hooks.html#appspec-hooks-server
hooks:
  ApplicationStop:
    - location: script-location
      timeout: 10
      runas: user-name
    - location: script-location
      timeout: 10
      runas: user-name
  BeforeInstall:
    - location: script-location
      timeout: 10
      runas: user-name`
