package assistant

import (
	"fmt"
	"strings"
	"testing"
)

var testJsonLambdaAppSpecFileString string = `{
  "version": 0,
  "Resources": [
    {
      "myLambdaFunction": {
        "Type": "AWS::Lambda::Function",
        "Properties": {
          "Name": "",
          "Alias": "",
          "CurrentVersion": "",
          "TargetVersion": ""
        }
      }
    }
  ]
}`

var testYamlLambdaAppSpecFileString string = `# https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file.html#appspec-reference-lambda
version: 0.0
# https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file-structure-resources.html#reference-appspec-file-structure-resources-lambda
Resources:
  - myLambdaFunction:
      Type: AWS::Lambda::Function
      Properties:
        Name: ""
        Alias: ""
        CurrentVersion: ""
        TargetVersion: ""
# https://docs.aws.amazon.com/codedeploy/latest/userguide/reference-appspec-file-structure-hooks.html#appspec-hooks-lambda
#Hooks:
#  - BeforeAllowTraffic: ""
#  - AfterAllowTraffic: ""`

var testJsonLambdaAppSpecMap string = `map[Resources:[map[myLambdaFunction:map[Properties:map[Alias: CurrentVersion: Name: TargetVersion:] Type:AWS::Lambda::Function]]] version:0]`

var testYamlLambdaAppSpecMap string = `map[Resources:[map[myLambdaFunction:map[Properties:map[Alias: CurrentVersion: Name: TargetVersion:] Type:AWS::Lambda::Function]]] version:0]`

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

func TestGetMapFromYaml(t *testing.T) {
	var tests = []struct {
		name            string
		inputFileString string
		outputString    string
	}{
		{"YAML file",
			testYamlLambdaAppSpecFileString, testYamlLambdaAppSpecMap},
	}

	for _, test := range tests {
		if output, err := getMapFromYaml(test.inputFileString); err != nil || fmt.Sprintf("%v", output) != test.outputString {
			t.Error("Test Failed: {} inputted, {} expected, recieved: {}", test.inputFileString, test.outputString, output)
		}
	}
}

func TestGetMapFromJson(t *testing.T) {
	var tests = []struct {
		name            string
		inputFileString string
		outputString    string
	}{
		{"JSON file",
			testJsonLambdaAppSpecFileString, testJsonLambdaAppSpecMap},
	}

	for _, test := range tests {
		if output, err := getMapFromJson(test.inputFileString); err != nil || fmt.Sprintf("%v", output) != test.outputString {
			t.Error("Test Failed: {} inputted, {} expected, recieved: {}", test.inputFileString, test.outputString, output)
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
