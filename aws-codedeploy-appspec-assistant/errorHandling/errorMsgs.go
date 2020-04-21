package errorHandling

const (
	//
	// General
	//

	// Top level errors that will be handled by ErrorHandler
	AppSpecVersionErr = "Version not supported. Version should be in this format: JSON - \"version\": 0.0, YAML - version: 0.0. The only versions supported are: %v"

	EmptyAppSpecFileErr = "AppSpec file is empty"

	EmptyFilePathErr              = "Empty filePath is not allowed"
	ComputePlatformErr            = "computePlatform must be server, lambda, or ecs"
	InvalidFileNameOrExtensionErr = "File must be named appspec and file extension must be .json or .yml (appspec.json or appspec.yml)"

	//
	// ECS
	//

	// Top level errors that will be handled by ErrorHandler
	InvalidECSResourcesErr         = "ECS resources are required and need to be valid"
	InvalidECSHooksAndFunctionsErr = "The hooks and their function values must be valid"

	// Lower-level errors that should be noted as ERROR or ERROR CAUSE
	UnsupportedNumberOfECSResourcesErr = "\nERROR CAUSE: Only 1 ECS resource (TargetService) supported in a deployment"
	InvalidECSTargetServiceTypeErr     = "\nERROR CAUSE: TargetService Type must be AWS::ECS::Service"
	InvalidECSTargetServicePropsErr    = "ERROR: Invalid TargetService properties"

	EmptyECSTaskDefErr                = "\nERROR CAUSE: Resources -> TargetService -> Properties -> TaskDefinition must not be empty (ECS TaskDefinition)"
	InvalidECSLoadBalancerInfoErr     = "ERROR: Resources -> TargetService -> Properties -> LoadBalancerInfo invalid"
	InvalidECSNetworkConfigurationErr = "ERROR: Resources -> TargetService -> Properties -> NetworkConfiguration invalid"

	MissingECSContainerNameErr = "\nERROR CAUSE: Resources -> TargetService -> Properties -> LoadBalancerInfo ... ContainerName missing for:"
	ZeroECSContainerPortWarn   = "WARNING: Resources -> TargetService -> Properties -> LoadBalancerInfo ... ContainerPort is 0. Please check this was on purpose for:"

	MissingECSSubnetsErr         = "\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... Subnets missing for:"
	EmptyECSSubnetStrsErr        = "\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... Subnets cannot be empty strings for:"
	MissingECSSecurityGroupsErr  = "\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... SecurityGroups missing for:"
	EmptyECSSecurityGroupStrsErr = "\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... SecurityGroups cannot be empty strings for:"
	MissingECSAssignPublicIpErr  = "\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... AssignPublicIp missing for:"
	InvalidECSAssignPublicIpErr  = "\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... AssignPublicIp invalid (should be ENABLED or DISABLED) for:"

	EmptyEcsHookValErr   = "\nERROR CAUSE: Value cannot be empty for hook:"
	InvalidEcsHookStrErr = "\nERROR CAUSE: The hooks must be one of the ECS supported hooks:"

	//
	// Lambda
	//

	// Top level errors that will be handled by ErrorHandler
	InvalidLambdaResourcesErr = "Lambda resources are required and need to be valid"
	InvalidLambdaHooksErr     = "The hooks must be one of the Lambda supported hooks: %v"

	// Lower-level errors that should be noted as ERROR or ERROR CAUSE
	UnsupportedNumberOfLambdaResourceErr = "\nERROR CAUSE: Only 1 Lambda resource (Function) supported in a deployment"
	EmptyLambdaResourceFunctionNameErr   = "\nERROR CAUSE: Value should not be empty for FunctionName of Resource"
	InvalidLambdaFunctionTypeErr         = "\nERROR CAUSE: Function Type must be AWS::Lambda::Function"
	InvalidLambdaFunctionPropsErr        = "ERROR: Invalid Function properties"

	EmptyLambdaFunctionNameErr          = "\nERROR CAUSE: Resources -> <Function> -> Properties -> Name must not be empty (Lambda Function Name) :"
	EmptyLambdaFunctionAliasErr         = "\nERROR CAUSE: Resources -> <Function> -> Properties -> Alias must not be empty (Lambda Function Alias) :"
	EmptyLambdaFunctionCurrVersionErr   = "\nERROR CAUSE: Resources -> <Function> -> Properties -> CurrentVersion must not be empty (Lambda Function current version, ex: 1) :"
	EmptyLambdaFunctionTargetVersionErr = "\nERROR CAUSE: Resources -> <Function> -> Properties -> TargetVersion must not be empty (Lambda Function target version to flip to, ex: 2) :"

	EmptyLambdaHookValErr = "\nERROR CAUSE: Value cannot be empty for hook:"

	//
	// Server (EC2/On-Prem)
	//

	// Because of all the optional sections in Server, the error that will be handled by ErrorHandler can vary
	// Lower-level errors that should be noted as ERROR or ERROR CAUSE
	UnsupportedServerOSErr      = "\nERROR CAUSE: OS not supported. Only 1 OS supported at a time. The only OSs supported are: %v"
	MissingServerFileSpecErr    = "\nERROR CAUSE: There must be at least 1 File(source, destination) specification"
	InvalidServerFileSpecsErr   = "ERROR: The Files specifications are invalid"
	InvalidServerPermissionsErr = "ERROR: The Permissions are invalid"
	InvalidServerHooksErr       = "ERROR: The hooks are invalid"

	MissingServerFileSourceErr      = "\nERROR CAUSE: Missing File Source (Source-Destination pairs must be set together)"
	MissingServerFileDestinationErr = "\nERROR CAUSE: Missing File Destination (Source-Destination pairs must be set together)"

	EmptyServerPermissionObjErr    = "\nERROR CAUSE: Object cannot be empty for permission:"
	InvalidServerPermissionTypeErr = "\nERROR CAUSE: If Permission Type is specified, it must be `file` or `direcory` for permission:"

	UnsupportedServerHooksErr        = "\nERROR CAUSE: The hooks must be one of the Ec2/OnPrem supported hooks."
	SupportedServerHooksWithoutLBStr = "\nDeployments without a LoadBalancer: %v"
	SupportedServerHooksWithLBStr    = "\nDeployments with a LoadBalancer: %v"

	MissingServerHookScriptLocationErr = "\nERROR CAUSE: The hook must have a script location:"
	InvalidServerScriptTimeoutErr      = "\nERROR CAUSE: Total timeout for all scripts within a single LifecycleEvent added up must not exceed 3600 seconds. :"
)
