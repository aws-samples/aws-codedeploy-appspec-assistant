package assistant

import (
	"fmt"

	"encoding/json"
	"gopkg.in/yaml.v3"

	"aws-codedeploy-appspec-assistant/errorHandling"
	"aws-codedeploy-appspec-assistant/globalVars"
	"aws-codedeploy-appspec-assistant/models"
)

// Convert Lambda AppSpec string to Lambda AppSpec Object
// Deals with JSON adn YAML
func getLambdaAppSpecObjFromString(appSpecBytes []byte) (models.LambdaAppSpecModel, error) {
	var err error
	var lambdaAppSpecModel models.LambdaAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal(appSpecBytes, &lambdaAppSpecModel)
	} else {
		err = json.Unmarshal(appSpecBytes, &lambdaAppSpecModel)
	}

	// Uncomment to print resulting Object for debug
	//fmt.Println(lambdaAppSpecModel)

	return lambdaAppSpecModel, err
}

// Validate Lambda AppSpec
// Calls validation on each section
func validateLambdaAppSpec(lambdaAppSpecModel models.LambdaAppSpecModel) error {
	var err error

	// Resources
	if lambdaAppSpecModel.Resources == nil || len(lambdaAppSpecModel.Resources) < 0 || !validateLambdaResources(lambdaAppSpecModel.Resources) {
		err = fmt.Errorf(errorHandling.InvalidLambdaResourcesErr)
	}

	// Hooks (Optional)
	if lambdaAppSpecModel.Hooks != nil && len(lambdaAppSpecModel.Hooks) > 0 {
		// Print resource error (if there is one) since there could be Hooks errors that change the final error message
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		}

		if !validateLambdaHooks(lambdaAppSpecModel.Hooks) {
			err = fmt.Errorf(errorHandling.InvalidLambdaHooksErr, globalVars.AppSpecSupportedLambdaHooks)
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
		fmt.Println(errorHandling.UnsupportedNumberOfLambdaResourceErr)
		numOfErrors++
		return false
	}

	for _, lambdaResource := range lambdaResources {
		for functionResourceName, function := range lambdaResource {

			// Function Name
			if functionResourceName == "" {
				resourcesValid = false
				numOfErrors++
				fmt.Println(errorHandling.EmptyLambdaResourceFunctionNameErr)
			}

			// Function Type
			if function.Type != "AWS::Lambda::Function" {
				resourcesValid = false
				numOfErrors++
				fmt.Println(errorHandling.InvalidLambdaFunctionTypeErr)
			}

			// Function Properties
			if !validateLambdaResourceProperties(function.Properties, functionResourceName) {
				resourcesValid = false
				fmt.Println(errorHandling.InvalidLambdaFunctionPropsErr)
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
		fmt.Println(errorHandling.EmptyLambdaFunctionNameErr, functionResourceName)
	}

	if lambdaProperties.Alias == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(errorHandling.EmptyLambdaFunctionAliasErr, functionResourceName)
	}

	if lambdaProperties.CurrentVersion == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(errorHandling.EmptyLambdaFunctionCurrVersionErr, functionResourceName)
	}

	if lambdaProperties.TargetVersion == "" {
		propertiesValid = false
		numOfErrors++
		fmt.Println(errorHandling.EmptyLambdaFunctionTargetVersionErr, functionResourceName)
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
					fmt.Println(errorHandling.EmptyLambdaHookValErr, hook)
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
