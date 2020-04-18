package globalVars

import (
	"fmt"
)

// General
var AppSpecVersionErr = fmt.Errorf("\nERROR CAUSE: Version not supported. Version should be in this format: JSON - \"version\": 0.0, YAML - version: 0.0. The only versions supported are: %v", AppSpecVersions)

var EmptyAppSpecFileErr = fmt.Errorf("\nERROR CAUSE: AppSpec file is empty")

var EmptyFilePathErr = fmt.Errorf("\nERROR CAUSE: Empty filePath is not allowed")
var ComputePlatformErr = fmt.Errorf("\nERROR CAUSE: computePlatform must be server, lambda, or ecs")
var InvalidFileNameOrExtensionErr = fmt.Errorf("\nERROR CAUSE: File must be named appspec and file extension must be .json or .yml (appspec.json or appspec.yml)")

// ECS
var InvalidECSResourcesErr = fmt.Errorf("ERROR: ECS resources are required and need to be valid")
var InvalidECSHooksAndFunctionsErr = fmt.Errorf("ERROR: The hooks and their function values must be valid")

var UnsupportedNumberOfECSResourcesErr = fmt.Errorf("\nERROR CAUSE: Only 1 ECS resource (TargetService) supported in a deployment")
var InvalidECSTargetServiceTypeErr = fmt.Errorf("\nERROR CAUSE: TargetService Type must be AWS::ECS::Service")
var InvalidECSTargetServicePropsErr = fmt.Errorf("ERROR: Invalid TargetService properties")

var EmptyECSTaskDefErr = fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> TaskDefinition must not be empty (ECS TaskDefinition)")
var InvalidECSLoadBalancerInfoErr = fmt.Errorf("ERROR: Resources -> TargetService -> Properties -> LoadBalancerInfo invalid")
var InvalidECSNetworkConfigurationErr = fmt.Errorf("ERROR: Resources -> TargetService -> Properties -> NetworkConfiguration invalid")

var MissingECSContainerNameErr = fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> LoadBalancerInfo ... ContainerName missing for: ")
var ZeroECSContainerPortWarn = fmt.Errorf("WARNING: Resources -> TargetService -> Properties -> LoadBalancerInfo ... ContainerPort is 0. Please check this was on purpose for: ")

var MissingECSSubnetsErr = fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... Subnets missing for: ")
var EmptyECSSubnetStrsErr = fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... Subnets cannot be empty strings for: ")
var MissingECSSecurityGroupsErr = fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... SecurityGroups missing for: ")
var EmptyECSSecurityGroupStrsErr = fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... SecurityGroups cannot be empty strings for: ")
var MissingECSAssignPublicIpErr = fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... AssignPublicIp missing for: ")
var InvalidECSAssignPublicIpErr = fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... AssignPublicIp invalid (should be ENABLED or DISABLED) for: ")

var EmptyEcsHookValErr = fmt.Errorf("\nERROR CAUSE: Value cannot be empty for hook: ")
var InvalidEcsHookStrErr = fmt.Errorf("\nERROR CAUSE: The hooks must be one of the ECS supported hooks: %v", AppSpecSupportedEcsHooks)

// Lambda
var InvalidLambdaResourcesErr = fmt.Errorf("ERROR: Lambda resources are required and need to be valid")
var InvalidLambdaHooksErr = fmt.Errorf("ERROR: The hooks must be one of the Lambda supported hooks: %v", AppSpecSupportedLambdaHooks)

var UnsupportedNumberOfLambdaResourceErr = fmt.Errorf("\nERROR CAUSE: Only 1 Lambda resource (Function) supported in a deployment")
var EmptyLambdaResourceFunctionNameErr = fmt.Errorf("\nERROR CAUSE: Value should not be empty for FunctionName of Resource")
var InvalidLambdaFunctionTypeErr = fmt.Errorf("\nERROR CAUSE: Function Type must be AWS::Lambda::Function")
var InvalidLambdaFunctionPropsErr = fmt.Errorf("ERROR: Invalid Function properties")

var EmptyLambdaFunctionNameErr = fmt.Errorf("\nERROR CAUSE: Resources -> <Function> -> Properties -> Name must not be empty (Lambda Function Name) : ")
var EmptyLambdaFunctionAliasErr = fmt.Errorf("\nERROR CAUSE: Resources -> <Function> -> Properties -> Alias must not be empty (Lambda Function Alias) : ")
var EmptyLambdaFunctionCurrVersionErr = fmt.Errorf("\nERROR CAUSE: Resources -> <Function> -> Properties -> CurrentVersion must not be empty (Lambda Function current version, ex: 1) : ")
var EmptyLambdaFunctionTargetVersionErr = fmt.Errorf("\nERROR CAUSE: Resources -> <Function> -> Properties -> TargetVersion must not be empty (Lambda Function target version to flip to, ex: 2) : ")

var EmptyLambdaHookValErr = fmt.Errorf("\nERROR CAUSE: Value cannot be empty for hook: ")

// Server (EC2/On-Prem)
var UnsupportedServerOSErr = fmt.Errorf("\nERROR CAUSE: OS not supported. Only 1 OS supported at a time. The only OSs supported are: %v", AppSpecSupportedServerOSs)
var MissingServerFileSpecErr = fmt.Errorf("\nERROR CAUSE: There must be at least 1 File(source, destination) specification")
var InvalidServerFileSpecsErr = fmt.Errorf("ERROR: The Files specifications are invalid")
var InvalidServerPermissionsErr = fmt.Errorf("ERROR: The Permissions are invalid")
var InvalidServerHooksErr = fmt.Errorf("ERROR: The hooks are invalid")

var MissingServerFileSourceErr = fmt.Errorf("\nERROR CAUSE: Missing File Source (Source-Destination pairs must be set together)")
var MissingServerFileDestinationErr = fmt.Errorf("\nERROR CAUSE: Missing File Destination (Source-Destination pairs must be set together)")

var EmptyServerPermissionObjErr = fmt.Errorf("\nERROR CAUSE: Object cannot be empty for permission: ")
var InvalidServerPermissionTypeErr = fmt.Errorf("\nERROR CAUSE: If Permission Type is specified, it must be `file` or `direcory` for permission: ")

var UnsupportedServerHooksErr = fmt.Errorf("\nERROR CAUSE: The hooks must be one of the Ec2/OnPrem supported hooks.\n Deployments without a LoadBalancer: %v \n Deployments with a LoadBalancer: %v", AppSpecSupportedServerHooksWithoutLB, AppSpecSupportedServerHooksWithLB)

var MissingServerHookScriptLocationErr = fmt.Errorf("\nERROR CAUSE: The hook must have a script location: ")
var InvalidServerScriptTimeoutErr = fmt.Errorf("\nERROR CAUSE: Total timeout for all scripts within a single LifecycleEvent added up must not exceed 3600 seconds. : ")
