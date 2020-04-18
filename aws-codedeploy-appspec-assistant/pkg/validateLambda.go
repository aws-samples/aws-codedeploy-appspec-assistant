package assistant

import (
	"fmt"

	"encoding/json"
	"gopkg.in/yaml.v3"

	"aws-codedeploy-appspec-assistant/globalVars"
	"aws-codedeploy-appspec-assistant/models"
)

// Convert Lambda AppSpec string to Lambda AppSpec Object
// Deals with JSON adn YAML
func getLambdaAppSpecObjFromString(appSpecBytes []byte) models.LambdaAppSpecModel {
	var err error
	var lambdaAppSpecModel models.LambdaAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal(appSpecBytes, &lambdaAppSpecModel)
	} else {
		err = json.Unmarshal(appSpecBytes, &lambdaAppSpecModel)
	}

	// Uncomment to print resulting Object for debug
	//fmt.Println(lambdaAppSpecModel)

	if err != nil {
		fmt.Print("\nERROR:")
		handleError(err)
	}

	return lambdaAppSpecModel
}

// Validate Lambda AppSpec
// Calls validation on each section
func validateLambdaAppSpec(lambdaAppSpecModel models.LambdaAppSpecModel) error {
	var err error

	// Resources
	if lambdaAppSpecModel.Resources == nil || len(lambdaAppSpecModel.Resources) < 0 || !validateLambdaResources(lambdaAppSpecModel.Resources) {
		err = globalVars.InvalidLambdaResourcesErr
		fmt.Println(err)
	}

	// Hooks (Optional)
	if lambdaAppSpecModel.Hooks != nil && len(lambdaAppSpecModel.Hooks) > 0 {
		if !validateLambdaHooks(lambdaAppSpecModel.Hooks) {
			err = globalVars.InvalidLambdaHooksErr
			fmt.Println(err)
		}
	}

	return err
}

// Lambda Resource validation methods
// Validate Lambda Function information
// Currently we only support 1
func validateLambdaResources(lambdaResources []map[string]models.Function) bool {
	resourcesValid := true

	if len(lambdaResources) > 1 {
		fmt.Println(globalVars.UnsupportedNumberOfLambdaResourceErr)
		numOfErrors++
		return false
	}

	for _, lambdaResource := range lambdaResources {
		for functionResourceName, function := range lambdaResource {

			// Function Name
			if functionResourceName == "" {
				resourcesValid = false
				numOfErrors++
				fmt.Println(globalVars.EmptyLambdaResourceFunctionNameErr)
			}

			// Function Type
			if function.Type != "AWS::Lambda::Function" {
				resourcesValid = false
				numOfErrors++
				fmt.Println(globalVars.InvalidLambdaFunctionTypeErr)
			}

			// Function Properties
			if !validateLambdaResourceProperties(function.Properties, functionResourceName) {
				resourcesValid = false
				fmt.Println(globalVars.InvalidLambdaFunctionPropsErr)
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
		fmt.Println(globalVars.EmptyLambdaFunctionNameErr.Error() + functionResourceName)
	}

	if lambdaProperties.Alias == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(globalVars.EmptyLambdaFunctionAliasErr.Error() + functionResourceName)
	}

	if lambdaProperties.CurrentVersion == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(globalVars.EmptyLambdaFunctionCurrVersionErr.Error() + functionResourceName)
	}

	if lambdaProperties.TargetVersion == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(globalVars.EmptyLambdaFunctionTargetVersionErr.Error() + functionResourceName)
	}

	return propertiesValid
}

// Lambda Hooks validation method
// Validate Hooks object
func validateLambdaHooks(lambdaHooks []map[string]string) bool {
	numValidHooks := 0
	hooksValid := true

	for _, lambdaHook := range lambdaHooks {
		for _, hook := range globalVars.AppSpecSupportedLambdaHooks {
			if val, ok := lambdaHook[hook]; ok {
				if val == "" {
					fmt.Println(globalVars.EmptyLambdaHookValErr.Error() + hook)
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
