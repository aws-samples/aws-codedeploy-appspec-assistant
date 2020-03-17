package assistant

// ECS AppSpec strings for Unit Tests
var ecsOutputStr = `{0 [{{AWS::ECS::Service {[Your task definition arn] {[Your container Name] 8000} [Version number, ex: 1.3.0] {{[SubnetId1 SubnetId2] [ecs-security-group-1] ENABLED-or-DISABLED}}}}}] [map[BeforeInstall:BeforeInstallHookLambdaFunctionName] map[AfterInstall:AfterInstallHookLambdaFunctionName] map[AfterAllowTestTraffic:AfterAllowTestTrafficHookLambdaFunctionName] map[BeforeAllowTraffic:SanityTestHookLambdaFunctionName] map[AfterAllowTraffic:ValidationTestHookLambdaFunctionName]]}`

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
              "AssignPublicIp": "ENABLED-or-DISABLED"
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
            AssignPublicIp: "ENABLED-or-DISABLED"
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

// EC2/OnPrem AppSpec strings for Unit Tests
var ec2OnPremOutputStr = `{0 linux [{source-file-location destination-file-location}] [{object-specification pattern-specification exception-specification owner-account-name group-name mode-specification [acls-specification] {user-specification type-specification range-specification} [object-type]}] map[deployment-lifecycle-event-name:[{script-location timeout-in-seconds user-name} {script-location timeout-in-seconds user-name}] deployment-lifecycle-event-name2:[{script-location timeout-in-seconds user-name}]]}`

var ec2OnPremJsonString = `{
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
        "object-type"
      ]
    }
  ],
  "hooks": {
    "deployment-lifecycle-event-name": [
      {
        "location": "script-location",
        "timeout": "timeout-in-seconds",
        "runas": "user-name"
      },
      {
        "location": "script-location",
        "timeout": "timeout-in-seconds",
        "runas": "user-name"
      }
    ],
    "deployment-lifecycle-event-name2": [
      {
        "location": "script-location",
        "timeout": "timeout-in-seconds",
        "runas": "user-name"
      }
    ]
  }
}`

var ec2OnPremYamlString = `version: 0.0
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
      - object-type
# https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file-structure-hooks.html#appspec-hooks-server
hooks:
  deployment-lifecycle-event-name:
    - location: script-location
      timeout: timeout-in-seconds
      runas: user-name
    - location: script-location
      timeout: timeout-in-seconds
      runas: user-name
  deployment-lifecycle-event-name2:
    - location: script-location
      timeout: timeout-in-seconds
      runas: user-name`
