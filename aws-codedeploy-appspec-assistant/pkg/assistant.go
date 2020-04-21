package assistant

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"aws-codedeploy-appspec-assistant/errorHandling"
	"aws-codedeploy-appspec-assistant/globalVars"
)

var fileExtension string

var numOfErrors int = 0

// Main function
func ValidateAppSpec(filePath string, computePlatform string) {
	fmt.Println("validateAppSpec called on:", filePath, ",", computePlatform)

	if err := validateUserInput(filePath, computePlatform); err != nil {
		errorHandling.HandleError(err)
	}

	// Load AppSpec
	raw_appSpec, err := ioutil.ReadFile(filePath)
	if err != nil {
		errorHandling.HandleError(err)
	}

	if len(string(raw_appSpec)) < 1 {
		numOfErrors++
		errorHandling.HandleError(fmt.Errorf(errorHandling.EmptyAppSpecFileErr))
	}

	if validationErr := runValidation(raw_appSpec, computePlatform); validationErr != nil {
		errorHandling.HandleError(validationErr)
	}

	fmt.Println("AppSpec file has passed available validation checks")
}

func validateUserInput(filePath string, computePlatform string) error {

	if len(filePath) < 1 {
		numOfErrors++
		return fmt.Errorf(errorHandling.EmptyFilePathErr)
	}

	if !isValidComputePlatform(computePlatform) {
		numOfErrors++
		return fmt.Errorf(errorHandling.ComputePlatformErr)
	}

	if !isValidFileNameAndExtension(filePath) {
		numOfErrors++
		return fmt.Errorf(errorHandling.InvalidFileNameOrExtensionErr)
	}

	// Very IMPORTANT. Do NOT delete. Need to set fileExtension variable
	saveFileExtension(filePath)

	if _, err := os.Stat(filePath); err != nil { // Path does not exist
		numOfErrors++
		return err
	}

	return nil
}

func isValidFileNameAndExtension(filePath string) bool {
	if strings.HasSuffix(filePath, "appspec.json") || strings.HasSuffix(filePath, "appspec.yml") {
		return true
	}

	return false
}

func saveFileExtension(filePath string) {
	filePathSplit := strings.Split(filePath, ".")
	fileExtension = filePathSplit[len(filePathSplit)-1]
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
func runValidation(appSpec []byte, computePlatform string) error {
	var err error

	// Validate version before converting AppSpec to objects
	err = validateVersionString(string(appSpec))

	if err != nil {
		numOfErrors++
		return err
	}

	if computePlatform == "ecs" {
		ecsAppSpecModel, modelErr := getEcsAppSpecObjFromString(appSpec)
		if modelErr != nil {
			return modelErr
		}
		err = validateEcsAppSpec(ecsAppSpecModel)
	} else if computePlatform == "lambda" {
		lambdaAppSpecModel, modelErr := getLambdaAppSpecObjFromString(appSpec)
		if modelErr != nil {
			return modelErr
		}
		err = validateLambdaAppSpec(lambdaAppSpecModel)
	} else {
		serverAppSpecModel, modelErr := getServerAppSpecObjFromString(appSpec)
		if modelErr != nil {
			return modelErr
		}
		err = validateServerAppSpec(serverAppSpecModel)
	}

	return err
}

// Validate Version string in all types of AppSpec
// Called before validating the rest of the AppSPec content
//     since it is the same in all types of AppSpecs right now
func validateVersionString(appSpecString string) error {
	appSpecStrSpaceSplit := strings.Fields(appSpecString)

	for _, version := range globalVars.AppSpecVersions {

		for i, appSpecSubStr := range appSpecStrSpaceSplit {

			// YAML
			if (len(appSpecStrSpaceSplit) > i+1) && (appSpecSubStr == "version:") {
				if appSpecStrSpaceSplit[i+1] == version {
					return nil
				}
				return fmt.Errorf(errorHandling.AppSpecVersionErr, globalVars.AppSpecVersions)
			}

			// JSON
			if (len(appSpecStrSpaceSplit) > i+1) && ((appSpecSubStr == "\"version\":") || (appSpecSubStr == "{\"version\":")) {
				if (appSpecStrSpaceSplit[i+1] == version+",") || (appSpecStrSpaceSplit[i+1] == version) ||
					(appSpecStrSpaceSplit[i+1] == version+"}") {
					return nil
				}
				return fmt.Errorf(errorHandling.AppSpecVersionErr, globalVars.AppSpecVersions)
			}
		}
	}

	return fmt.Errorf(errorHandling.AppSpecVersionErr, globalVars.AppSpecVersions)
}
