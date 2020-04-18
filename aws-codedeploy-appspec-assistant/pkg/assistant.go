package assistant

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"aws-codedeploy-appspec-assistant/globalVars"
)

var fileExtension string

var numOfErrors int = 0

// Main function
func ValidateAppSpec(filePath string, computePlatform string) {
	fmt.Println("validateAppSpec called on:", filePath, ",", computePlatform)

	if err := validateUserInput(filePath, computePlatform); err != nil {
		handleError(err)
	}

	appSpec := loadAppSpec(filePath)

	if len(string(appSpec)) < 1 {
		numOfErrors++
		handleError(globalVars.EmptyAppSpecFileErr)
	}

	runValidation(appSpec, computePlatform)
}

func validateUserInput(filePath string, computePlatform string) error {

	if len(filePath) < 1 {
		numOfErrors++
		return globalVars.EmptyFilePathErr
	}

	if !isValidComputePlatform(computePlatform) {
		numOfErrors++
		return globalVars.ComputePlatformErr
	}

	if !isValidFileNameAndExtension(filePath) {
		numOfErrors++
		return globalVars.InvalidFileNameOrExtensionErr
	}

	// Very IMPORTANT. Do NOT delete
	saveFileExtension(filePath)

	if _, err := os.Stat(filePath); err != nil { // Path does not exist
		numOfErrors++
		if err != nil {
			fmt.Print("\nERROR:")
		}
		return err
	}

	return nil
}

func loadAppSpec(filePath string) []byte {
	raw_appSpec, err := ioutil.ReadFile(filePath)

	if err != nil {
		numOfErrors++
		fmt.Print("\nERROR:")
		handleError(err)
	}

	return raw_appSpec
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
func runValidation(appSpec []byte, computePlatform string) {
	var err error

	// Validate version before converting AppSpec to objects
	err = validateVersionString(string(appSpec))

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
		serverAppSpecModel := getServerAppSpecObjFromString(appSpec)
		err = validateServerAppSpec(serverAppSpecModel)
	}

	if err != nil {
		handleError(err)
	}

	fmt.Println("AppSpec file is valid")

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
				return globalVars.AppSpecVersionErr
			}

			// JSON
			if (len(appSpecStrSpaceSplit) > i+1) && ((appSpecSubStr == "\"version\":") || (appSpecSubStr == "{\"version\":")) {
				if (appSpecStrSpaceSplit[i+1] == version+",") || (appSpecStrSpaceSplit[i+1] == version) ||
					(appSpecStrSpaceSplit[i+1] == version+"}") {
					return nil
				}
				return globalVars.AppSpecVersionErr
			}
		}
	}

	return globalVars.AppSpecVersionErr
}

func handleError(err error) {
	if err != nil {
		defer func() {
			fmt.Println("Panic CAUSE: " + err.Error())
			fmt.Println(fmt.Errorf("The AppSpec is not valid. %d errors were found during validation.", numOfErrors))
			os.Exit(1)
		}()
		panic(err)
	}
}
