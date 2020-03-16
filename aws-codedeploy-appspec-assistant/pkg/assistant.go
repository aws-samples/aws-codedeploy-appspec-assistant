package assistant

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"encoding/json"
	"gopkg.in/yaml.v3"

	//"aws-codedeploy-appspec-assistant/models"
)

var fileExtension string

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

	var appSpecMap map[string]interface{}
	var err error

	if fileExtension == "yml" || fileExtension == "yaml" {
		appSpecMap, err = getMapFromYaml(appSpec)
	} else {
		appSpecMap, err = getMapFromJson(appSpec)
	}

	if err != nil {
		handleError(err)
	}

	//check the appspec

	fmt.Println("AppSpec file is valid")
	fmt.Println(appSpec)
	fmt.Println(appSpecMap)

}

func getMapFromYaml(appSpecYamlString string) (map[string]interface{}, error) {
	var appSpecMapFromYaml map[string]interface{}

	err := yaml.Unmarshal([]byte(appSpecYamlString), &appSpecMapFromYaml)

	return appSpecMapFromYaml, err
}

func getMapFromJson(appSpecJaonString string) (map[string]interface{}, error) {
	var appSpecMapFromJson map[string]interface{}
	err := json.Unmarshal([]byte(appSpecJaonString), &appSpecMapFromJson)

	return appSpecMapFromJson, err
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
