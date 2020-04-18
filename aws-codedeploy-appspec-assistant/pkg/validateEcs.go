package assistant

import (
	"fmt"

	"encoding/json"
	"gopkg.in/yaml.v3"

	"aws-codedeploy-appspec-assistant/globalVars"
	"aws-codedeploy-appspec-assistant/models"
)

// Convert ECS AppSpec string to ECS AppSpec Object
// Deals with JSON adn YAML
func getEcsAppSpecObjFromString(appSpecBytes []byte) models.EcsAppSpecModel {
	var err error
	var ecsAppSpecModel models.EcsAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal(appSpecBytes, &ecsAppSpecModel)
	} else {
		err = json.Unmarshal(appSpecBytes, &ecsAppSpecModel)
	}

	// Uncomment to print resulting Object for debug
	//fmt.Println(ecsAppSpecModel)

	if err != nil {
		fmt.Print("\nERROR:")
		handleError(err)
	}

	return ecsAppSpecModel
}

// Validate ECS AppSpec
// Calls validation on each section
func validateEcsAppSpec(ecsAppSpecModel models.EcsAppSpecModel) error {
	var err error

	// Resources
	if ecsAppSpecModel.Resources == nil || len(ecsAppSpecModel.Resources) < 0 || !validateEcsResources(ecsAppSpecModel.Resources) {
		err = globalVars.InvalidECSResourcesErr
		fmt.Println(err)
	}

	// Hooks (Optional)
	if ecsAppSpecModel.Hooks != nil && len(ecsAppSpecModel.Hooks) > 0 {
		if !validateEcsHooks(ecsAppSpecModel.Hooks) {
			err = globalVars.InvalidECSHooksAndFunctionsErr
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
		fmt.Println(globalVars.UnsupportedNumberOfECSResourcesErr)
		return false
	}

	for _, ecsResource := range ecsResources {
		// Resource Type
		if ecsResource.TargetService.Type != "AWS::ECS::Service" {
			resourcesValid = false
			numOfErrors++
			fmt.Println(globalVars.InvalidECSTargetServiceTypeErr)
		}

		// Resource Properties
		if !validateEcsResourceProperties(ecsResource.TargetService.Properties) {
			resourcesValid = false
			fmt.Println(globalVars.InvalidECSTargetServicePropsErr)
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
		fmt.Println(globalVars.EmptyECSTaskDefErr)
	}

	// LoadBalancerInfo
	if !validateEcsLoadBalancerInfo(ecsProperties.LoadBalancerInfo, ecsProperties.TaskDefinition) {
		propertiesValid = false
		fmt.Println(globalVars.InvalidECSLoadBalancerInfoErr)
	}

	// PlatformVersion (Optional)

	// NetworkConfiguration (Optional)
	if isEcsNetworkConfigurationFilledOut(ecsProperties.NetworkConfiguration) {
		if !validateEcsAwsvpcConfiguration(ecsProperties.NetworkConfiguration.AwsvpcConfiguration, ecsProperties.TaskDefinition) {
			propertiesValid = false
			fmt.Println(globalVars.InvalidECSNetworkConfigurationErr)
		}
	}

	return propertiesValid
}

func validateEcsLoadBalancerInfo(ecsLoadBalancerInfo models.LoadBalancerInfo, taskDefinition string) bool {
	infoValid := true

	if ecsLoadBalancerInfo.ContainerName == "" {
		infoValid = false
		numOfErrors++
		fmt.Println(globalVars.MissingECSContainerNameErr.Error() + taskDefinition)
	}

	if ecsLoadBalancerInfo.ContainerPort == 0 {
		fmt.Println(globalVars.ZeroECSContainerPortWarn.Error() + taskDefinition)
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
		fmt.Println(globalVars.MissingECSSubnetsErr.Error() + taskDefinition)
	} else {
		for _, subnet := range ecsAwsvpcConfiguration.Subnets {
			if subnet == "" {
				configValid = false
				numOfErrors++
				fmt.Println(globalVars.EmptyECSSubnetStrsErr.Error() + taskDefinition)
			}
		}
	}

	if ecsAwsvpcConfiguration.SecurityGroups == nil || len(ecsAwsvpcConfiguration.SecurityGroups) < 1 {
		configValid = false
		numOfErrors++
		fmt.Println(globalVars.MissingECSSecurityGroupsErr.Error() + taskDefinition)
	} else {
		for _, securityGroup := range ecsAwsvpcConfiguration.SecurityGroups {
			if securityGroup == "" {
				configValid = false
				numOfErrors++
				fmt.Println(globalVars.EmptyECSSecurityGroupStrsErr.Error() + taskDefinition)
			}
		}
	}

	if ecsAwsvpcConfiguration.AssignPublicIp == "" {
		configValid = false
		numOfErrors++
		fmt.Println(globalVars.MissingECSAssignPublicIpErr.Error() + taskDefinition)
	} else if !validateEcsAssignPublicIpValue(ecsAwsvpcConfiguration.AssignPublicIp) {
		configValid = false
		numOfErrors++
		fmt.Println(globalVars.InvalidECSAssignPublicIpErr.Error() + taskDefinition)
	}

	return configValid
}

func validateEcsAssignPublicIpValue(assignPublicIpValue string) bool {
	for _, supportedPublicIpValue := range globalVars.AppSpecEcsAssignPublicIpValues {
		if assignPublicIpValue == supportedPublicIpValue {
			return true
		}
	}

	return false
}

// ECS Hooks validation method
// Validate Hooks object
func validateEcsHooks(ecsHooks []map[string]string) bool {
	numValidHooks := 0
	hooksValid := true

	for _, ecsHook := range ecsHooks {
		for _, hook := range globalVars.AppSpecSupportedEcsHooks {
			if val, ok := ecsHook[hook]; ok {
				if val == "" {
					fmt.Println(globalVars.EmptyEcsHookValErr.Error() + hook)
					numOfErrors++
					hooksValid = false
				}
				numValidHooks++
			}
		}
	}

	if numValidHooks != len(ecsHooks) {
		numOfErrors++
		fmt.Println(globalVars.InvalidEcsHookStrErr)
		hooksValid = false
	}

	return hooksValid
}
