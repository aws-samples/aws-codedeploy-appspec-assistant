package assistant

import (
	"fmt"
	"strings"
	"testing"

	"aws-codedeploy-appspec-assistant/globalVars"
)

// Test validateUserInput
func TestValidateUserInput_InvalidInput(t *testing.T) {
	var tests = []struct {
		name                string
		filePathInput       string
		computeTypeInput    string
		expectedErrorOutput string
	}{
		{"Empty filePath",
			"", "lambda", globalVars.EmptyFilePathErr.Error()},

		{"Invalid filePath extension",
			"/appSpec_assistant_test/appspec.txt", "lambda", globalVars.InvalidFileNameOrExtensionErr.Error()},

		{"Invalid filePath file name",
			"/appSpec_assistant_test/incorrect.yml", "lambda", globalVars.InvalidFileNameOrExtensionErr.Error()},

		{"Empty computeType",
			"/appSpec_assistant_test/appspec.yml", "", globalVars.ComputePlatformErr.Error()},

		{"Invalid computeType",
			"/appSpec_assistant_test/appspec.json", "invalidComputeType", globalVars.ComputePlatformErr.Error()},

		{"YAML and Lambda",
			"/appSpec_assistant_test/appspec.yml", "lambda", "no such file or directory"},
	}

	for _, test := range tests {
		err := validateUserInput(test.filePathInput, test.computeTypeInput)
		if (err == nil) || (!strings.Contains(fmt.Sprintf("%v", err), test.expectedErrorOutput)) {
			t.Errorf("The code did not error correctly for: %v. Got this error instead of expected: %v", test, err)
		}
	}
}

func TestValidateUserInput_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		filePathInput    string
		computeTypeInput string
	}{
		{"YAML and Lambda",
			"/appSpec_assistant_test/appspec.yml", "lambda"},

		{"YAML and EC2",
			"/appSpec_assistant_test/appspec.yml", "server"},

		{"YAML and ECS",
			"/appSpec_assistant_test/appspec.yml", "ecs"},

		{"JSON and Lambda",
			"/appSpec_assistant_test/appspec.json", "lambda"},

		{"JSON and EC2",
			"/appSpec_assistant_test/appspec.json", "server"},

		{"JSON and ECS",
			"/appSpec_assistant_test/appspec.json", "ecs"},
	}

	for _, test := range tests {
		if err := validateUserInput(test.filePathInput, test.computeTypeInput); err != nil {
			if !strings.Contains(fmt.Sprintf("%v", err), "no such file or directory") {
				t.Errorf("The code should not error for: %v. Got this error: %v", test, err)
			}
		}
	}
}

// Test validateVersionString
func TestValidateVersionString_ValidInput(t *testing.T) {
	var tests = []struct {
		name               string
		appSpecStringInput string
		fileExtensionVal   string
	}{
		{"Valid JSON version",
			"{\"version\": 0.0}", "json"},
		{"Valid JSON version",
			"{ \"version\": 0.0 }", "json"},
		{"Valid JSON version",
			"{ \"version\": 0.0, \"os\": \"linux\" }", "json"},
		{"Valid YAML version",
			"version: 0.0", "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		err := validateVersionString(test.appSpecStringInput)
		if err != nil {
			t.Errorf("The validateVersionString function failed for: %v", test)
		}
	}
}

func TestValidateVersionString_InvalidInput(t *testing.T) {
	var tests = []struct {
		name               string
		appSpecStringInput string
		fileExtensionVal   string
	}{
		{"Invalid JSON version",
			"{\"version\": 0.02}", "json"},
		{"Invalid YAML version",
			"version: 0.02", "yml"},
		{"Invalid JSON version",
			"{\"version\": test}", "json"},
		{"Invalid YAML version",
			"version: test", "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		err := validateVersionString(test.appSpecStringInput)
		if err == nil {
			t.Errorf("The validateVersionString function did not fail for: %v", test)
		}
	}
}

func TestSaveFileExtension(t *testing.T) {
	var tests = []struct {
		name              string
		filePathInput     string
		expectedExtension string
	}{
		{"JSON",
			"test/file/path/appspec.json", "json"},
		{"YML",
			"test/file/path/appspec.yml", "yml"},
	}

	for _, test := range tests {
		saveFileExtension(test.filePathInput)
		if fileExtension != test.expectedExtension {
			t.Errorf("The validateVersionString function did not fail for: %v", test)
		}
	}
}
