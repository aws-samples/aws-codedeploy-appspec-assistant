package assistant

import (
	"fmt"
	"strings"
	"testing"
)

var appSpecStrConversionError string = "The appSpecModel does not match the expected output model"
var appSpecObjValidationError string = "The valid appSpec object threw errors during validation"

func TestValidateAppSpec_InvalidInput(t *testing.T) {
	var tests = []struct {
		name             string
		filePathInput    string
		computeTypeInput string
	}{
		{"Empty filePath",
			"", "lambda"},

		{"Invalid filePath",
			"/appSpec_assistant_test/testPath.txt", "lambda"},

		{"Invalid yaml filePath",
			"/appSpec_assistant_test/testPath.yaml", "lambda"},

		{"Empty computeType",
			"/appSpec_assistant_test/testPath.yml", ""},

		{"Invalid computeType",
			"/appSpec_assistant_test/testPath.json", "invalidComputeType"},
	}

	for _, test := range tests {
		r := validateAppSpecPanics(ValidateAppSpec, test.filePathInput, test.computeTypeInput)
		t.Log(r)
		if r == nil {
			t.Errorf("The code did not panic")
		}
	}
}

func TestValidateAppSpec_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		filePathInput    string
		computeTypeInput string
	}{
		{"YAML and Lambda",
			"/appSpec_assistant_test/testPath.yml", "lambda"},

		{"YAML and EC2",
			"/appSpec_assistant_test/testPath.yml", "ec2/on-prem"},

		{"YAML and ECS",
			"/appSpec_assistant_test/testPath.yml", "ecs"},

		{"JSON and Lambda",
			"/appSpec_assistant_test/testPath.json", "lambda"},

		{"JSON and EC2",
			"/appSpec_assistant_test/testPath.json", "ec2/on-prem"},

		{"JSON and ECS",
			"/appSpec_assistant_test/testPath.json", "ecs"},
	}

	for _, test := range tests {
		r := validateAppSpecPanics(ValidateAppSpec, test.filePathInput, test.computeTypeInput)
		t.Log(r)
		if r == nil {
			t.Errorf("The code did not panic")
		} else if !strings.Contains(fmt.Sprintf("%v", r), "no such file or directory") {
			t.Errorf("The test failed before checking if the path exists")
		}
	}
}

func TestGetEcsAppSpecObjFromString_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		fileStrInput     string
		objectStrOutput  string
		fileExtensionVal string
	}{
		{"Valid JSON",
			ecsJsonString, ecsOutputStr, "json"},

		{"Valid YAML",
			ecsYamlString, ecsOutputStr, "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		appSpecModel := getEcsAppSpecObjFromString(test.fileStrInput)
		if fmt.Sprintf("%v", appSpecModel) != test.objectStrOutput {
			t.Errorf(appSpecStrConversionError)
		}
	}
}

func TestGetLambdaAppSpecObjFromString_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		fileStrInput     string
		objectStrOutput  string
		fileExtensionVal string
	}{
		{"Valid JSON",
			lambdaJsonString, lambdaOutputStr, "json"},

		{"Valid YAML",
			lambdaYamlString, lambdaOutputStr, "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		appSpecModel := getLambdaAppSpecObjFromString(test.fileStrInput)
		if fmt.Sprintf("%v", appSpecModel) != test.objectStrOutput {
			t.Errorf(appSpecStrConversionError)
		}
	}
}

func TestGetEc2OnPremAppSpecObjFromString_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		fileStrInput     string
		objectStrOutput  string
		fileExtensionVal string
	}{
		{"Valid JSON",
			ec2OnPremJsonString, ec2OnPremOutputStr, "json"},

		{"Valid YAML",
			ec2OnPremYamlString, ec2OnPremOutputStr, "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		appSpecModel := getEc2OnPremAppSpecObjFromString(test.fileStrInput)
		if fmt.Sprintf("%v", appSpecModel) != test.objectStrOutput {
			t.Errorf(appSpecStrConversionError)
		}
	}
}

func TestValidateEcsAppSpec_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		fileStrInput     string
		fileExtensionVal string
	}{
		{"Valid YAML",
			ecsYamlString, "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		appSpecModel := getEcsAppSpecObjFromString(test.fileStrInput)
		err := validateEcsAppSpec(appSpecModel)
		if err != nil {
			t.Errorf(appSpecObjValidationError)
		}
	}
}

func TestValidateLambdaAppSpec_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		fileStrInput     string
		fileExtensionVal string
	}{
		{"Valid YAML",
			lambdaYamlString, "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		appSpecModel := getLambdaAppSpecObjFromString(test.fileStrInput)
		err := validateLambdaAppSpec(appSpecModel)
		if err != nil {
			t.Errorf(appSpecObjValidationError)
		}
	}
}

func TestValidateEc2OnPremAppSpec_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		fileStrInput     string
		fileExtensionVal string
	}{
		{"Valid YAML",
			ec2OnPremYamlString, "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		appSpecModel := getEc2OnPremAppSpecObjFromString(test.fileStrInput)
		err := validateEc2OnPremAppSpec(appSpecModel)
		if err != nil {
			t.Errorf(appSpecObjValidationError)
		}
	}
}

func validateAppSpecPanics(validateAppSpec func(string, string), filePathInput string, computeTypeInput string) (r interface{}) {
	defer func() {
		r = recover()
	}()

	validateAppSpec(filePathInput, computeTypeInput)

	return
}
