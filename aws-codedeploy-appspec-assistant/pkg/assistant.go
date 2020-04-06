package assistant

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"encoding/json"
	"gopkg.in/yaml.v3"

	"aws-codedeploy-appspec-assistant/models"
)

var fileExtension string

var appSpecVersions = []string{"0.0"}
var appSpecVersionError = fmt.Errorf("\nERROR CAUSE: Version not supported. The only versions supported are: %v", appSpecVersions)

var appSpecSupportedOSs = []string{"linux", "windows"}

var appSpecSupportedEcsHooks = []string{"BeforeInstall", "AfterInstall", "AfterAllowTestTraffic", "BeforeAllowTraffic", "AfterAllowTraffic"}
var appSpecSupportedLambdaHooks = []string{"BeforeAllowTraffic", "AfterAllowTraffic"}
var appSpecSupportedEc2OnPremHooksWithLB = []string{"BeforeBlockTraffic", "AfterBlockTraffic", "BeforeAllowTraffic", "AfterAllowTraffic"}
var appSpecSupportedEc2OnPremHooksWithoutLB = []string{"ApplicationStop", "BeforeInstall", "AfterInstall", "ApplicationStart", "ValidateService"}
var appSpecEc2OnPremHooksError = fmt.Errorf("\nERROR CAUSE: The hooks must be one of the Ec2/OnPrem supported hooks.\n Deployments without a LoadBalancer: %v \n Deployments with a LoadBalancer: %v", appSpecSupportedEc2OnPremHooksWithoutLB, appSpecSupportedEc2OnPremHooksWithLB)

var appSpecEcsAssignPublicIpValues = []string{"ENABLED", "DISABLED"}

var numOfErrors int = 0

// Main function
func ValidateAppSpec(filePath string, computePlatform string) {
	fmt.Println("validateAppSpec called on:", filePath, ",", computePlatform)

	if err := validateUserInput(filePath, computePlatform); err != nil {
		handleError(err)
	}

	appSpec := loadAppSpec(filePath)

	if len(appSpec) < 1 {
		numOfErrors++
		handleError(fmt.Errorf("\nERROR CAUSE: AppSpec file is empty"))
	}

	runValidation(appSpec, computePlatform)
}

func validateUserInput(filePath string, computePlatform string) error {

	if len(filePath) < 1 {
		numOfErrors++
		return fmt.Errorf("\nERROR CAUSE: Empty filePath is not allowed")
	}

	if !isValidComputePlatform(computePlatform) {
		numOfErrors++
		return fmt.Errorf("\nERROR CAUSE: computePlatform must be server, lambda, or ecs")
	}

	if !isValidFileExtension(filePath) {
		numOfErrors++
		return fmt.Errorf("\nERROR CAUSE: File extension must be .json or .yml")
	}

	if _, err := os.Stat(filePath); err != nil { // Path does not exist
		numOfErrors++
		return err
	}

	return nil
}

func loadAppSpec(filePath string) string {
	raw_appSpec, err := ioutil.ReadFile(filePath)

	if err != nil {
		numOfErrors++
		handleError(err)
	}

	return string(raw_appSpec)
}

func isValidFileExtension(filePath string) bool {
	filePathSplit := strings.Split(filePath, ".")
	fileExtension = filePathSplit[len(filePathSplit)-1]

	if fileExtension == "json" || fileExtension == "yml" {
		return true
	}

	return false
}

func isValidComputePlatform(computePlatform string) bool {
	if computePlatform == "server" || computePlatform == "lambda" || computePlatform == "ecs" {
		return true
	}

	return false
}

// Starts validation fo the AppSpec file content
// Converts string into AppSpec Objects
// Runs validationon the AppSpec Objects
func runValidation(appSpec string, computePlatform string) {
	var err error

	// Validate version before converting AppSpec to objects
	err = validateVersionString(appSpec)

	if err != nil {
		numOfErrors++
		handleError(err)
	}

	if computePlatform == "ecs" {
		ecsAppSpecModel := getEcsAppSpecObjFromString(appSpec)
		err = validateEcsAppSpec(ecsAppSpecModel)
	} else if computePlatform == "lambda" {
		lambdaAppSpecModel := getLambdaAppSpecObjFromString(appSpec)
		err = validateLambdaAppSpec(lambdaAppSpecModel)
	} else {
		ec2OnPremAppSpecModel := getEc2OnPremAppSpecObjFromString(appSpec)
		err = validateEc2OnPremAppSpec(ec2OnPremAppSpecModel)
	}

	if err != nil {
		handleError(err)
	}

	fmt.Println("AppSpec file is valid")

}

// Convert ECS AppSpec string to ECS AppSpec Object
// Deals with JSON adn YAML
func getEcsAppSpecObjFromString(appSpecString string) models.EcsAppSpecModel {
	var err error
	var ecsAppSpecModel models.EcsAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal([]byte(appSpecString), &ecsAppSpecModel)
	} else {
		err = json.Unmarshal([]byte(appSpecString), &ecsAppSpecModel)
	}

	// Uncomment to print resulting Object for debug
	//fmt.Println(ecsAppSpecModel)

	if err != nil {
		handleError(err)
	}

	return ecsAppSpecModel
}

// Convert ECS AppSpec string to ECS AppSpec Object
// Deals with JSON adn YAML
func getLambdaAppSpecObjFromString(appSpecString string) models.LambdaAppSpecModel {
	var err error
	var lambdaAppSpecModel models.LambdaAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal([]byte(appSpecString), &lambdaAppSpecModel)
	} else {
		err = json.Unmarshal([]byte(appSpecString), &lambdaAppSpecModel)
	}

	// Uncomment to print resulting Object for debug
	//fmt.Println(lambdaAppSpecModel)

	if err != nil {
		handleError(err)
	}

	return lambdaAppSpecModel
}

// Convert ECS AppSpec string to ECS AppSpec Object
// Deals with JSON adn YAML
func getEc2OnPremAppSpecObjFromString(appSpecString string) models.Ec2OnPremAppSpecModel {
	var err error
	var ec2OnPremAppSpecModel models.Ec2OnPremAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal([]byte(appSpecString), &ec2OnPremAppSpecModel)
	} else {
		err = json.Unmarshal([]byte(appSpecString), &ec2OnPremAppSpecModel)
	}

	// Uncomment to print resulting Object for debug
	//fmt.Println(ec2OnPremAppSpecModel)

	if err != nil {
		handleError(err)
	}

	return ec2OnPremAppSpecModel
}

// Validate Version string in all types of AppSpec
// Called before validating the rest of the AppSPec content
//     since it is the same in all types of AppSpecs right now
func validateVersionString(appSpecString string) error {
	for _, version := range appSpecVersions {
		i := -1
		versionStrLen := 0
		if fileExtension == "yml" {
			i = strings.LastIndex(appSpecString, "version: "+version)
			versionStrLen = 12
		} else {
			i = strings.LastIndex(appSpecString, "\"version\": "+version)
			versionStrLen = 14
		}
		// If version is the only thing in the file
		if (i + versionStrLen) == len(appSpecString) {
			return nil
		}
		// Check there isn't anything like 0.02
		if i > -1 {
			// Check character after version
			followingChar := appSpecString[i+versionStrLen]
			followingRune, _ := utf8.DecodeRune([]byte{followingChar})
			if unicode.IsSpace(followingRune) || followingChar == ',' || followingChar == '}' {
				return nil
			}
		}
	}

	return appSpecVersionError
}

// Validate ECS AppSpec
// Calls validation on each section
func validateEcsAppSpec(ecsAppSpecModel models.EcsAppSpecModel) error {
	var err error

	// Resources
	if ecsAppSpecModel.Resources == nil || len(ecsAppSpecModel.Resources) < 0 || !validateEcsResources(ecsAppSpecModel.Resources) {
		err = fmt.Errorf("ERROR: ECS resources are required and need to be valid")
		fmt.Println(err)
	}

	// Hooks (Optional)
	if ecsAppSpecModel.Hooks != nil && len(ecsAppSpecModel.Hooks) > 0 {
		if !validateEcsHooks(ecsAppSpecModel.Hooks) {
			err = fmt.Errorf("ERROR: The hooks and their function values must be valid")
			fmt.Println(err)
		}
	}

	return err
}

// Validate Lambda AppSpec
// Calls validation on each section
func validateLambdaAppSpec(lambdaAppSpecModel models.LambdaAppSpecModel) error {
	var err error

	// Resources
	if lambdaAppSpecModel.Resources == nil || len(lambdaAppSpecModel.Resources) < 0 || !validateLambdaResources(lambdaAppSpecModel.Resources) {
		err = fmt.Errorf("ERROR: Lambda resources are required and need to be valid")
		fmt.Println(err)
	}

	// Hooks (Optional)
	if lambdaAppSpecModel.Hooks != nil && len(lambdaAppSpecModel.Hooks) > 0 {
		if !validateLambdaHooks(lambdaAppSpecModel.Hooks) {
			err = fmt.Errorf("ERROR: The hooks must be one of the Lambda supported hooks: %v", appSpecSupportedLambdaHooks)
			fmt.Println(err)
		}
	}

	return err
}

// Validate EC2/On-Prem (Server) AppSpec
// Calls validation on each section
func validateEc2OnPremAppSpec(ec2OnPremAppSpecModel models.Ec2OnPremAppSpecModel) error {
	var err error

	// OS
	if ec2OnPremAppSpecModel.OS == "" || !checkOS(ec2OnPremAppSpecModel.OS) {
		numOfErrors++
		osError := fmt.Errorf("\nERROR CAUSE: OS not supported. Only 1 OS supported at a time. The only OSs supported are: %v", appSpecSupportedOSs)
		fmt.Println(osError)
		err = osError
	}

	// Files
	if ec2OnPremAppSpecModel.Files == nil || len(ec2OnPremAppSpecModel.Files) < 1 {
		numOfErrors++
		err = fmt.Errorf("\nERROR CAUSE: There must be at least 1 File(source, destination) specification")
		fmt.Println(err)
	} else {
		if !validateEc2OnPremFiles(ec2OnPremAppSpecModel.Files) {
			err = fmt.Errorf("ERROR: The Files specifications are invalid")
			fmt.Println(err)
		}
	}

	// Permissions (Optional)
	if ec2OnPremAppSpecModel.Permissions != nil && len(ec2OnPremAppSpecModel.Permissions) > 0 {
		if !validateEc2OnPremPermissions(ec2OnPremAppSpecModel.Permissions) {
			err = fmt.Errorf("ERROR: The Permissions are invalid")
			fmt.Println(err)
		}
	}

	// Hooks (Optional)
	if ec2OnPremAppSpecModel.Hooks != nil && len(ec2OnPremAppSpecModel.Hooks) > 0 {
		if !validateEc2OnPremHooks(ec2OnPremAppSpecModel.Hooks) {
			err = fmt.Errorf("ERROR: The hooks are invalid")
			fmt.Println(err)
		}
	}

	return err
}

// ECS Resource validation methods
// Validate ECS TargetService information
// Currently we only support 1
func validateEcsResources(ecsResources []models.Resource) bool {
	resourcesValid := true

	if len(ecsResources) > 1 {
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Only 1 ECS resource (TargetService) supported in a deployment"))
		return false
	}

	for _, ecsResource := range ecsResources {
		// Resource Type
		if ecsResource.TargetService.Type != "AWS::ECS::Service" {
			resourcesValid = false
			numOfErrors++
			fmt.Println(fmt.Errorf("\nERROR CAUSE: TargetService Type must be AWS::ECS::Service"))
		}

		// Resource Properties
		if !validateEcsResourceProperties(ecsResource.TargetService.Properties) {
			resourcesValid = false
			fmt.Println(fmt.Errorf("ERROR: Invalid TargetService properties"))
		}
	}

	return resourcesValid
}

func validateEcsResourceProperties(ecsProperties models.EcsProperties) bool {
	propertiesValid := true

	// TaskDefinition
	if ecsProperties.TaskDefinition == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> TaskDefinition must not be empty (ECS TaskDefinition)"))
	}

	// LoadBalancerInfo
	if !validateEcsLoadBalancerInfo(ecsProperties.LoadBalancerInfo, ecsProperties.TaskDefinition) {
		propertiesValid = false
		fmt.Println(fmt.Errorf("ERROR: Resources -> TargetService -> Properties -> LoadBalancerInfo invalid"))
	}

	// PlatformVersion (Optional)

	// NetworkConfiguration (Optional)
	if isEcsNetworkConfigurationFilledOut(ecsProperties.NetworkConfiguration) {
		if !validateEcsAwsvpcConfiguration(ecsProperties.NetworkConfiguration.AwsvpcConfiguration, ecsProperties.TaskDefinition) {
			propertiesValid = false
			fmt.Println(fmt.Errorf("ERROR: Resources -> TargetService -> Properties -> NetworkConfiguration invalid"))
		}
	}

	return propertiesValid
}

func validateEcsLoadBalancerInfo(ecsLoadBalancerInfo models.LoadBalancerInfo, taskDefinition string) bool {
	infoValid := true

	if ecsLoadBalancerInfo.ContainerName == "" {
		infoValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> LoadBalancerInfo ... ContainerName missing for: %s", taskDefinition))
	}

	if ecsLoadBalancerInfo.ContainerPort == 0 {
		fmt.Println(fmt.Errorf("WARNING: Resources -> TargetService -> Properties -> LoadBalancerInfo ... ContainerPort is 0. Please check this was on purpose for: %s", taskDefinition))
	}

	return infoValid
}

func isEcsNetworkConfigurationFilledOut(ecsNetworkConfig models.NetworkConfiguration) bool {
	return isEcsAwsvpcConfigurationFilledOut(ecsNetworkConfig.AwsvpcConfiguration)
}

func isEcsAwsvpcConfigurationFilledOut(ecsAwsvpcConfig models.AwsvpcConfiguration) bool {
	if (ecsAwsvpcConfig.Subnets == nil || len(ecsAwsvpcConfig.Subnets) < 1) &&
		(ecsAwsvpcConfig.SecurityGroups == nil || len(ecsAwsvpcConfig.SecurityGroups) < 1) &&
		ecsAwsvpcConfig.AssignPublicIp == "" {
		return false
	}

	return true
}

func validateEcsAwsvpcConfiguration(ecsAwsvpcConfiguration models.AwsvpcConfiguration, taskDefinition string) bool {
	configValid := true

	if ecsAwsvpcConfiguration.Subnets == nil || len(ecsAwsvpcConfiguration.Subnets) < 1 {
		configValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... Subnets missing for: %s", taskDefinition))
	} else {
		for _, subnet := range ecsAwsvpcConfiguration.Subnets {
			if subnet == "" {
				configValid = false
				numOfErrors++
				fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... Subnets cannot be empty strings for: %s", taskDefinition))
			}
		}
	}

	if ecsAwsvpcConfiguration.SecurityGroups == nil || len(ecsAwsvpcConfiguration.SecurityGroups) < 1 {
		configValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... SecurityGroups missing for: %s", taskDefinition))
	} else {
		for _, securityGroup := range ecsAwsvpcConfiguration.SecurityGroups {
			if securityGroup == "" {
				configValid = false
				numOfErrors++
				fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... SecurityGroups cannot be empty strings for: %s", taskDefinition))
			}
		}
	}

	if ecsAwsvpcConfiguration.AssignPublicIp == "" {
		configValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... AssignPublicIp missing for: %s", taskDefinition))
	} else if !validateEcsAssignPublicIpValue(ecsAwsvpcConfiguration.AssignPublicIp) {
		configValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> TargetService -> Properties -> AwsvpcConfiguration ... AssignPublicIp invalid (should be ENABLED or DISABLED) for: %s", taskDefinition))
	}

	return configValid
}

func validateEcsAssignPublicIpValue(assignPublicIpValue string) bool {
	for _, supportedPublicIpValue := range appSpecEcsAssignPublicIpValues {
		if assignPublicIpValue == supportedPublicIpValue {
			return true
		}
	}

	return false
}

// Lambda Resource validation methods
// Validate Lambda Function information
// Currently we only support 1
func validateLambdaResources(lambdaResources []map[string]models.Function) bool {
	resourcesValid := true

	if len(lambdaResources) > 1 {
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Only 1 Lambda resource (Function) supported in a deployment"))
		numOfErrors++
		return false
	}

	for _, lambdaResource := range lambdaResources {
		for functionResourceName, function := range lambdaResource {

			// Function Name
			if functionResourceName == "" {
				resourcesValid = false
				numOfErrors++
				fmt.Println(fmt.Errorf("\nERROR CAUSE: Value should not be empty for FunctionName of Resource"))
			}

			// Function Type
			if function.Type != "AWS::Lambda::Function" {
				resourcesValid = false
				numOfErrors++
				fmt.Println(fmt.Errorf("\nERROR CAUSE: Function Type must be AWS::Lambda::Function"))
			}

			// Function Properties
			if !validateLambdaResourceProperties(function.Properties, functionResourceName) {
				resourcesValid = false
				fmt.Println(fmt.Errorf("ERROR: Invalid Function properties"))
			}
		}
	}

	return resourcesValid
}

func validateLambdaResourceProperties(lambdaProperties models.LambdaProperties, functionResourceName string) bool {
	propertiesValid := true

	if lambdaProperties.Name == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> <Function> -> Properties -> Name must not be empty (Lambda Function Name) : %s", functionResourceName))
	}

	if lambdaProperties.Alias == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> <Function> -> Properties -> Alias must not be empty (Lambda Function Alias) : %s", functionResourceName))
	}

	if lambdaProperties.CurrentVersion == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> <Function> -> Properties -> CurrentVersion must not be empty (Lambda Function current version, ex: 1) : %s", functionResourceName))
	}

	if lambdaProperties.TargetVersion == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: Resources -> <Function> -> Properties -> TargetVersion must not be empty (Lambda Function target version to flip to, ex: 2) : %s", functionResourceName))
	}

	return propertiesValid
}

// EC2/OnPrem OS validation method
// Validate OS is Linux or Windows
func checkOS(appSpecOS string) bool {
	for _, supportedOS := range appSpecSupportedOSs {
		if supportedOS == appSpecOS {
			return true
		}
	}

	return false
}

// EC2/On-Prem (Server) Files Validation method
// Validate the files object values
func validateEc2OnPremFiles(files []models.File) bool {
	filesValid := true
	for _, file := range files {
		if file.Source == "" {
			filesValid = false
			numOfErrors++
			fmt.Println(fmt.Errorf("\nERROR CAUSE: Missing File Source (Source-Destination pairs must be set together)"))
		}

		if file.Destination == "" {
			filesValid = false
			numOfErrors++
			fmt.Println(fmt.Errorf("\nERROR CAUSE: Missing File Destination (Source-Destination pairs must be set together)"))
		}
	}

	return filesValid
}

// EC2/On-Prem (Server) Permissions Validation method
// Validate the Permissions object values
func validateEc2OnPremPermissions(permissions []models.Permission) bool {
	permissionsValid := true

	for _, permission := range permissions {
		if permission.Object == "" {
			permissionsValid = false
			numOfErrors++
			fmt.Println(fmt.Errorf("\nERROR CAUSE: Object cannot be empty for permission: %v", permission))
		}

		if permission.Type != nil && len(permission.Type) > 0 {
			for _, typeStr := range permission.Type {
				if typeStr != "" && typeStr != "file" && typeStr != "directory" {
					permissionsValid = false
					numOfErrors++
					fmt.Println(fmt.Errorf("\nERROR CAUSE: If Permission Type is specifiec, it must be `file` or `direcory` for permission: %v", permission))
				}
			}
		}
	}

	// All other values are optionsl
	fmt.Println("WARNING: All options besides Object are optional for permissions so there is very little to validate automatically.")

	return permissionsValid
}

// ECS Hooks validation method
// Validate Hooks object
func validateEcsHooks(ecsHooks []map[string]string) bool {
	numValidHooks := 0
	hooksValid := true

	for _, ecsHook := range ecsHooks {
		for _, hook := range appSpecSupportedEcsHooks {
			if val, ok := ecsHook[hook]; ok {
				if val == "" {
					fmt.Println(fmt.Errorf("\nERROR CAUSE: Value cannot be empty for hook: %s", hook))
					numOfErrors++
					hooksValid = false
				}
				numValidHooks++
			}
		}
	}

	if numValidHooks != len(ecsHooks) {
		numOfErrors++
		fmt.Println(fmt.Errorf("\nERROR CAUSE: The hooks must be one of the ECS supported hooks: %v", appSpecSupportedEcsHooks))
		hooksValid = false
	}

	return hooksValid
}

// Lambda Hooks validation method
// Validate Hooks object
func validateLambdaHooks(lambdaHooks []map[string]string) bool {
	numValidHooks := 0
	hooksValid := true

	for _, lambdaHook := range lambdaHooks {
		for _, hook := range appSpecSupportedLambdaHooks {
			if val, ok := lambdaHook[hook]; ok {
				if val == "" {
					fmt.Println(fmt.Errorf("\nERROR CAUSE: Value cannot be empty for hook: %s", hook))
					numOfErrors++
					hooksValid = false
				}
				numValidHooks++
			}
		}
	}

	if numValidHooks != len(lambdaHooks) {
		numOfErrors++
		hooksValid = false
	}

	return hooksValid
}

// EC2/OnPrem Hooks validation methods
// Validate Hooks object
func validateEc2OnPremHooks(ec2OnPremHooks map[string][]models.Hook) bool {
	numValidHooks := 0
	hookScriptsValid := true
	for _, hook := range appSpecSupportedEc2OnPremHooksWithoutLB {
		if val, ok := ec2OnPremHooks[hook]; ok {
			hookScriptsValid = hookScriptsValid && validateEc2OnPremHookScripts(val, hook)
			numValidHooks++
		}
	}

	withLBHooksUsed := false

	for _, hook := range appSpecSupportedEc2OnPremHooksWithLB {
		if val, ok := ec2OnPremHooks[hook]; ok {
			withLBHooksUsed = true
			hookScriptsValid = hookScriptsValid && validateEc2OnPremHookScripts(val, hook)
			numValidHooks++
		}
	}

	fmt.Println("WARNING: runas under Hook Scripts only applies to Amazon Linux and Ubuntu Server instances. The user also cannot require a password. Leave blank for agent default.")

	if withLBHooksUsed {
		fmt.Println("WARNING: EC2/On-Prem (Server) hooks for LoadBalancers used, so the deployments should use a LoadBalancer for these scripts to be run.")
	}

	if numValidHooks == len(ec2OnPremHooks) && hookScriptsValid {
		return true
	}

	if !(numValidHooks == len(ec2OnPremHooks)) {
		numOfErrors++
		fmt.Println(appSpecEc2OnPremHooksError)
	}

	return false
}

func validateEc2OnPremHookScripts(hookScriptList []models.Hook, hook string) bool {
	scriptsValid := true
	totalTimeout := 0
	for _, hookScript := range hookScriptList {
		if hookScript.Location == "" {
			scriptsValid = false
			numOfErrors++
			fmt.Println(fmt.Errorf("\nERROR CAUSE: The hook must have a script location: %s", hook))
		}

		if hookScript.Timeout != "" {
			i, err := strconv.Atoi(hookScript.Timeout)
			if err != nil {
				scriptsValid = false
				fmt.Println(err)
				continue
			}
			totalTimeout += i
			if totalTimeout > 3600 {
				numOfErrors++
				fmt.Println(fmt.Errorf("\nERROR CAUSE: Total timeout for all scripts within a single LifecycleEvent added up must not exceed 3600 seconds. : %s", hook))
				scriptsValid = false
			}
		}
	}

	return scriptsValid
}

func handleError(err error) {
	if err != nil {
		defer func() {
			fmt.Println(err)
			fmt.Println(fmt.Errorf("\nThe AppSpec is not valid. There are %d errors.", numOfErrors))
			os.Exit(1)
		}()
		panic(err)
	}
}
