package assistant

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"encoding/json"
	"gopkg.in/yaml.v3"

	"aws-codedeploy-appspec-assistant/models"
)

var fileExtension string

var appSpecVersions = []float32{0.0}

var appSpecVersionError = fmt.Errorf("ERROR: Version not supported. The only versions supported are: %v", appSpecVersions)

func ValidateAppSpec(filePath string, computePlatform string) {
	fmt.Println("validateAppSpec called on:", filePath, ",", computePlatform)

	if len(filePath) < 1 {
		handleError(fmt.Errorf("Empty filePath is not allowed"))
	}

	if !isValidComputePlatform(computePlatform) {
		handleError(fmt.Errorf("computePlatform must be ec2/on-prem, lambda, or ecs"))
	}

	if !isValidFileExtension(filePath) {
		handleError(fmt.Errorf("File extension must be .json or .yml"))
	}

	if _, err := os.Stat(filePath); err != nil { // Path does not exist
		handleError(err)
	}

	appSpec := loadAppSpec(filePath)

	if len(appSpec) < 1 {
		handleError(fmt.Errorf("AppSpec file is empty"))
	}

	runValidation(appSpec, computePlatform)
}

func loadAppSpec(filePath string) string {
	raw_appSpec, err := ioutil.ReadFile(filePath)

	if err != nil {
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
	if computePlatform == "ec2/on-prem" || computePlatform == "lambda" || computePlatform == "ecs" {
		return true
	}

	return false
}

func runValidation(appSpec string, computePlatform string) {
	var err error

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

func getEcsAppSpecObjFromString(appSpecString string) models.EcsAppSpecModel {
	var err error
	var ecsAppSpecModel models.EcsAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal([]byte(appSpecString), &ecsAppSpecModel)
	} else {
		err = json.Unmarshal([]byte(appSpecString), &ecsAppSpecModel)
	}

	fmt.Println(ecsAppSpecModel)

	if err != nil {
		handleError(err)
	}

	return ecsAppSpecModel
}

func getLambdaAppSpecObjFromString(appSpecString string) models.LambdaAppSpecModel {
	var err error
	var lambdaAppSpecModel models.LambdaAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal([]byte(appSpecString), &lambdaAppSpecModel)
	} else {
		err = json.Unmarshal([]byte(appSpecString), &lambdaAppSpecModel)
	}

	fmt.Println(lambdaAppSpecModel)

	if err != nil {
		handleError(err)
	}

	return lambdaAppSpecModel
}

func getEc2OnPremAppSpecObjFromString(appSpecString string) models.Ec2OnPremAppSpecModel {
	var err error
	var ec2OnPremAppSpecModel models.Ec2OnPremAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal([]byte(appSpecString), &ec2OnPremAppSpecModel)
	} else {
		err = json.Unmarshal([]byte(appSpecString), &ec2OnPremAppSpecModel)
	}

	fmt.Println(ec2OnPremAppSpecModel)

	if err != nil {
		handleError(err)
	}

	return ec2OnPremAppSpecModel
}

func validateEcsAppSpec(ecsAppSpecModel models.EcsAppSpecModel) error {
	var err error

	if !checkVersion(ecsAppSpecModel.Version) {
		fmt.Println(appSpecVersionError)
		err = appSpecVersionError
	}

	return err
}

func validateLambdaAppSpec(lambdaAppSpecModel models.LambdaAppSpecModel) error {
	var err error

	if !checkVersion(lambdaAppSpecModel.Version) {
		fmt.Println(appSpecVersionError)
		err = appSpecVersionError
	}

	return err
}

func validateEc2OnPremAppSpec(ec2OnPremAppSpecModel models.Ec2OnPremAppSpecModel) error {
	var err error

	if !checkVersion(ec2OnPremAppSpecModel.Version) {
		fmt.Println(appSpecVersionError)
		err = appSpecVersionError
	}

	return err
}

func checkVersion(version float32) bool {
	for _, versionNum := range appSpecVersions {
		if versionNum == version {
			return true
		}
	}

	return false
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
