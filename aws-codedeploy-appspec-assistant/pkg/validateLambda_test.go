package assistant

import (
	"fmt"
	"testing"

	"aws-codedeploy-appspec-assistant/models"
)

// Test getLambdaAppSpecObjFromString
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
		appSpecModel, err := getLambdaAppSpecObjFromString([]byte(test.fileStrInput))
		if err != nil || fmt.Sprintf("%v", appSpecModel) != test.objectStrOutput {
			t.Errorf(appSpecStrConversionError)
		}
	}
}

// Test validateLambdaAppSpec
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
		appSpecModel, modelErr := getLambdaAppSpecObjFromString([]byte(test.fileStrInput))
		if modelErr != nil {
			t.Errorf("getLambdaAppSpecObjFromString FAILED")
		}
		err := validateLambdaAppSpec(appSpecModel)
		if err != nil {
			t.Errorf(appSpecObjValidationError)
		}
	}
}

// Test validateLambdaHooks
func TestValidateLambdaHooks_ValidInput(t *testing.T) {
	var tests = []struct {
		name       string
		hooksInput []map[string]string
	}{
		{"One hook",
			[]map[string]string{{"BeforeAllowTraffic": "<SanityTestHookLambdaFunctionName>"}}},
		{"Multiple hooks",
			[]map[string]string{{"BeforeAllowTraffic": "<SanityTestHookLambdaFunctionName>"}, {"AfterAllowTraffic": "<ValidationTestHookLambdaFunctionName>"}}},
	}

	for _, test := range tests {
		output := validateLambdaHooks(test.hooksInput)
		if output != true {
			t.Errorf("The validateLambdaHooks function failed for: %v", test)
		}
	}
}

func TestValidateLambdaHooks_InvalidInput(t *testing.T) {
	var tests = []struct {
		name       string
		hooksInput []map[string]string
	}{
		{"One hook",
			[]map[string]string{{"NotHook": "<SanityTestHookLambdaFunctionName>"}}},
		{"One hook, no value",
			[]map[string]string{{"BeforeInstall": ""}}},
		{"Multiple hooks",
			[]map[string]string{{"BeforeAllowTraffic": "<SanityTestHookLambdaFunctionName>"}, {"NotHook": "<ValidationTestHookLambdaFunctionName>"}}},
	}

	for _, test := range tests {
		output := validateLambdaHooks(test.hooksInput)
		if output == true {
			t.Errorf("The validateLambdaHooks function succeeded but should have failed for: %v", test)
		}
	}
}

// Test ValidateLambdaResources
func TestValidateLambdaResources_ValidInput(t *testing.T) {
	var tests = []struct {
		name           string
		resourcesInput []map[string]models.Function
	}{
		{"One resource",
			[]map[string]models.Function{
				{
					"myLambdaFunction": {
						"AWS::Lambda::Function",
						models.LambdaProperties{
							"myLambdaFunction",
							"myLambdaFunctionAlias",
							"1",
							"2",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		output := validateLambdaResources(test.resourcesInput)
		if output != true {
			t.Errorf("The validateLambdaResources function failed for: %v", test)
		}
	}
}

func TestValidateLambdaResources_InvalidInput(t *testing.T) {
	var tests = []struct {
		name           string
		resourcesInput []map[string]models.Function
	}{
		{"One resource, missing Function name",
			[]map[string]models.Function{
				{
					"": {
						"AWS::Lambda::Function",
						models.LambdaProperties{
							"myLambdaFunction",
							"myLambdaFunctionAlias",
							"1",
							"2",
						},
					},
				},
			},
		},
		{"One resource, wrong Type",
			[]map[string]models.Function{
				{
					"myLambdaFunction": {
						"WRONG",
						models.LambdaProperties{
							"myLambdaFunction",
							"myLambdaFunctionAlias",
							"1",
							"2",
						},
					},
				},
			},
		},
		{"One resource, missing Name property",
			[]map[string]models.Function{
				{
					"myLambdaFunction": {
						"AWS::Lambda::Function",
						models.LambdaProperties{
							"",
							"myLambdaFunctionAlias",
							"1",
							"2",
						},
					},
				},
			},
		},
		{"One resource, missing Alias property",
			[]map[string]models.Function{
				{
					"myLambdaFunction": {
						"AWS::Lambda::Function",
						models.LambdaProperties{
							"myLambdaFunction",
							"",
							"1",
							"2",
						},
					},
				},
			},
		},
		{"One resource, missing CurrentVersion property",
			[]map[string]models.Function{
				{
					"myLambdaFunction": {
						"AWS::Lambda::Function",
						models.LambdaProperties{
							"myLambdaFunction",
							"myLambdaFunctionAlias",
							"",
							"2",
						},
					},
				},
			},
		},
		{"One resource, missing TargetVersion property",
			[]map[string]models.Function{
				{
					"myLambdaFunction": {
						"AWS::Lambda::Function",
						models.LambdaProperties{
							"myLambdaFunction",
							"myLambdaFunctionAlias",
							"1",
							"",
						},
					},
				},
			},
		},
		{"Multiple resources",
			[]map[string]models.Function{
				{
					"myLambdaFunction": {
						"AWS::Lambda::Function",
						models.LambdaProperties{
							"myLambdaFunction",
							"myLambdaFunctionAlias",
							"1",
							"2",
						},
					},
				},
				{
					"myLambdaFunction2": {
						"AWS::Lambda::Function",
						models.LambdaProperties{
							"myLambdaFunction2",
							"myLambdaFunctionAlias2",
							"1",
							"2",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		output := validateLambdaResources(test.resourcesInput)
		if output == true {
			t.Errorf("The validateLambdaResources function succeeded but should have failed for: %v", test)
		}
	}
}
