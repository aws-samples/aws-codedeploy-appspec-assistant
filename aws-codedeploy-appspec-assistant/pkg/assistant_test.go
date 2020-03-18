package assistant

import (
	"fmt"
	"strings"
	"testing"

	"aws-codedeploy-appspec-assistant/models"
)

var appSpecStrConversionError string = "The appSpecModel does not match the expected output model"
var appSpecObjValidationError string = "The valid appSpec object threw errors during validation"

// Test ValidateAppSpec
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
		defer func() {
			if r := recover(); r == nil {
				t.Log(r)
				t.Errorf("The code did not panic for: %v", test)
			}
		}()
		ValidateAppSpec(test.filePathInput, test.computeTypeInput)
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
		defer func() {
			if r := recover(); r == nil {
				t.Log(r)
				t.Errorf("The code did not panic for: %v", test)
			} else if !strings.Contains(fmt.Sprintf("%v", r), "no such file or directory") {
				t.Errorf("The test failed before checking if the path exists for: %v", test)
			}
		}()
		ValidateAppSpec(test.filePathInput, test.computeTypeInput)
	}
}

// Test GetEcsAppSpecObjFromString
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

// Test GetLambdaAppSpecObjFromString
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

// Test GetEc2OnPremAppSpecObjFromString
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

// Test ValidateEcsAppSpec
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

// Test ValidateLambdaAppSpec
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

// Test ValidateEc2OnPremAppSpec
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

// Test ValidateVersionString
func TestValidateVersionString_ValidInput(t *testing.T) {
	var tests = []struct {
		name               string
		appSpecStringInput string
		fileExtensionVal   string
	}{
		{"Valid JSON version",
			"{\"version\": 0.0}", "json"},
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

// Test CheckOS
func TestCheckOS_ValidInput(t *testing.T) {
	var tests = []struct {
		name          string
		osStringInput string
	}{
		{"Valid OS linux",
			"linux"},
		{"Valid OS windows",
			"windows"},
	}

	for _, test := range tests {
		output := checkOS(test.osStringInput)
		if output != true {
			t.Errorf("The checkOS function failed for: %v", test)
		}
	}
}

func TestCheckOS_InvalidInput(t *testing.T) {
	var tests = []struct {
		name          string
		osStringInput string
	}{
		{"Invalid OS",
			"ubuntu"},
		{"Invalid OS list",
			"linux,windows"},
	}

	for _, test := range tests {
		output := checkOS(test.osStringInput)
		if output == true {
			t.Errorf("The checkOS function succeeded but should have failed for: %v", test)
		}
	}
}

// Test ValidateEcsHooks
func TestValidateEcsHooks_ValidInput(t *testing.T) {
	var tests = []struct {
		name       string
		hooksInput []map[string]string
	}{
		{"One hook",
			[]map[string]string{{"BeforeInstall": "BeforeInstallHookLambdaFunctionName"}}},
		{"Multiple hooks",
			[]map[string]string{{"BeforeInstall": "BeforeInstallHookLambdaFunctionName"}, {"AfterInstall": "AfterInstallHookLambdaFunctionName"}}},
	}

	for _, test := range tests {
		output := validateEcsHooks(test.hooksInput)
		if output != true {
			t.Errorf("The validateEcsHooks function failed for: %v", test)
		}
	}
}

func TestValidateEcsHooks_InvalidInput(t *testing.T) {
	var tests = []struct {
		name       string
		hooksInput []map[string]string
	}{
		{"One hook",
			[]map[string]string{{"NotHook": "BeforeInstallHookLambdaFunctionName"}}},
		{"One hook, no value",
			[]map[string]string{{"BeforeInstall": ""}}},
		{"Multiple hooks",
			[]map[string]string{{"BeforeInstall": "BeforeInstallHookLambdaFunctionName"}, {"NotHook": "AfterInstallHookLambdaFunctionName"}}},
	}

	for _, test := range tests {
		output := validateEcsHooks(test.hooksInput)
		if output == true {
			t.Errorf("The validateEcsHooks function succeeded but should have failed for: %v", test)
		}
	}
}

// Test ValidateLambdaHooks
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

// Test ValidateEc2OnPremHooks
func TestValidateEc2OnPremHooks_ValidInput(t *testing.T) {
	var tests = []struct {
		name       string
		hooksInput map[string][]models.Hook
	}{
		{"One hook, one script",
			map[string][]models.Hook{"ApplicationStop": {
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
			}},
		},
		{"One hook, multiple scripts",
			map[string][]models.Hook{"ApplicationStop": {
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
			}},
		},
		{"Multiple hooks, multiple scripts",
			map[string][]models.Hook{"ApplicationStop": {
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
			},
				"BeforeInstall": {
					{
						Location: "script-location",
						Timeout:  "10",
						Runas:    "user-name",
					},
				}},
		},
	}

	for _, test := range tests {
		output := validateEc2OnPremHooks(test.hooksInput)
		if output != true {
			t.Errorf("The validateEc2OnPremHooks function failed for: %v", test)
		}
	}
}

func TestValidateEc2OnPremHooks_InvalidInput(t *testing.T) {
	var tests = []struct {
		name       string
		hooksInput map[string][]models.Hook
	}{
		{"One bad hook, one script",
			map[string][]models.Hook{"NotHook": {
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
			}},
		},
		{"One hook, multiple scripts, bad timeout",
			map[string][]models.Hook{"ApplicationStop": {
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
				{
					Location: "script-location",
					Timeout:  "3600",
					Runas:    "user-name",
				},
			}},
		},
		{"One hook, multiple scripts, missing location value",
			map[string][]models.Hook{"ApplicationStop": {
				{
					Location: "",
					Timeout:  "10",
					Runas:    "user-name",
				},
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
			}},
		},
		{"One hook, multiple scripts, missing location key",
			map[string][]models.Hook{"ApplicationStop": {
				{
					Timeout: "10",
					Runas:   "user-name",
				},
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
			}},
		},
		{"Multiple hooks, multiple scripts, one bad hook",
			map[string][]models.Hook{"ApplicationStop": {
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
				{
					Location: "script-location",
					Timeout:  "10",
					Runas:    "user-name",
				},
			},
				"NotHook": {
					{
						Location: "script-location",
						Timeout:  "10",
						Runas:    "user-name",
					},
				}},
		},
	}

	for _, test := range tests {
		output := validateEc2OnPremHooks(test.hooksInput)
		if output == true {
			t.Errorf("The validateEc2OnPremHooks function succeeded but should have failed for: %v", test)
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

// Test ValidateEcsResources
func TestValidateEcsResources_ValidInput(t *testing.T) {
	var tests = []struct {
		name           string
		resourcesInput []models.Resource
	}{
		{"One resource, ENABLED",
			[]models.Resource{{
				models.TargetService{
					"AWS::ECS::Service",
					models.EcsProperties{
						"[Your task definition arn]",
						models.LoadBalancerInfo{
							"[Your container Name]",
							8000,
						},
						"[Version number, ex: 1.3.0]",
						models.NetworkConfiguration{
							models.AwsvpcConfiguration{
								[]string{
									"SubnetId1",
									"SubnetId2",
								},
								[]string{
									"ecs-security-group-1",
								},
								"ENABLED",
							},
						},
					},
				},
			}},
		},
		{"One resource, DISABLED",
			[]models.Resource{{
				models.TargetService{
					"AWS::ECS::Service",
					models.EcsProperties{
						"[Your task definition arn]",
						models.LoadBalancerInfo{
							"[Your container Name]",
							8000,
						},
						"[Version number, ex: 1.3.0]",
						models.NetworkConfiguration{
							models.AwsvpcConfiguration{
								[]string{
									"SubnetId1",
									"SubnetId2",
								},
								[]string{
									"ecs-security-group-1",
									"ecs-security-group-2",
								},
								"DISABLED",
							},
						},
					},
				},
			}},
		},
	}

	for _, test := range tests {
		output := validateEcsResources(test.resourcesInput)
		if output != true {
			t.Errorf("The validateEcsResources function failed for: %v", test)
		}
	}
}

func TestValidateEcsResources_InvalidInput(t *testing.T) {
	var tests = []struct {
		name           string
		resourcesInput []models.Resource
	}{
		{"One resource, invalid Type",
			[]models.Resource{
				{
					models.TargetService{
						"WRONG",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"SubnetId1",
										"SubnetId2",
									},
									[]string{
										"ecs-security-group-1",
										"ecs-security-group-2",
									},
									"DISABLED",
								},
							},
						},
					},
				},
			},
		},
		{"One resource, missing TaskDefinition ARN",
			[]models.Resource{
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"SubnetId1",
										"SubnetId2",
									},
									[]string{
										"ecs-security-group-1",
										"ecs-security-group-2",
									},
									"DISABLED",
								},
							},
						},
					},
				},
			},
		},
		{"One resource, missing LoadBalancer Container Name",
			[]models.Resource{
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"SubnetId1",
										"SubnetId2",
									},
									[]string{
										"ecs-security-group-1",
										"ecs-security-group-2",
									},
									"DISABLED",
								},
							},
						},
					},
				},
			},
		},
		{"One resource, missing Subnets",
			[]models.Resource{
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{},
									[]string{
										"ecs-security-group-1",
										"ecs-security-group-2",
									},
									"DISABLED",
								},
							},
						},
					},
				},
			},
		},
		{"One resource, empty Subnet strings",
			[]models.Resource{
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"",
									},
									[]string{
										"ecs-security-group-1",
										"ecs-security-group-2",
									},
									"DISABLED",
								},
							},
						},
					},
				},
			},
		},
		{"One resource, missing SecurityGroups",
			[]models.Resource{
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"SubnetId1",
										"SubnetId2",
									},
									[]string{},
									"DISABLED",
								},
							},
						},
					},
				},
			},
		},
		{"One resource, empty SecurityGroup strings",
			[]models.Resource{
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"SubnetId1",
										"SubnetId2",
									},
									[]string{
										"",
									},
									"DISABLED",
								},
							},
						},
					},
				},
			},
		},
		{"One resource, invalid AssignPublicIp value",
			[]models.Resource{
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"SubnetId1",
										"SubnetId2",
									},
									[]string{
										"ecs-security-group-1",
										"ecs-security-group-2",
									},
									"INVALID",
								},
							},
						},
					},
				},
			},
		},
		{"Multiple resources",
			[]models.Resource{
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"SubnetId1",
										"SubnetId2",
									},
									[]string{
										"ecs-security-group-1",
										"ecs-security-group-2",
									},
									"DISABLED",
								},
							},
						},
					},
				},
				{
					models.TargetService{
						"AWS::ECS::Service",
						models.EcsProperties{
							"[Your task definition arn]",
							models.LoadBalancerInfo{
								"[Your container Name]",
								8000,
							},
							"[Version number, ex: 1.3.0]",
							models.NetworkConfiguration{
								models.AwsvpcConfiguration{
									[]string{
										"SubnetId1",
										"SubnetId2",
									},
									[]string{
										"ecs-security-group-1",
										"ecs-security-group-2",
									},
									"DISABLED",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		output := validateEcsResources(test.resourcesInput)
		if output == true {
			t.Errorf("The validateEcsResources function succeeded but should have failed for: %v", test)
		}
	}
}

// Test IsEcsNetworkConfigurationFilledOut
func TestIsEcsNetworkConfigurationFilledOut_FilledOutInput(t *testing.T) {
	var tests = []struct {
		name               string
		networkConfigInput models.NetworkConfiguration
	}{
		{"A filled out EcsNetworkConfiguration",
			models.NetworkConfiguration{
				models.AwsvpcConfiguration{
					[]string{
						"SubnetId1",
						"SubnetId2",
					},
					[]string{
						"ecs-security-group-1",
						"ecs-security-group-2",
					},
					"DISABLED",
				},
			},
		},
		{"A partially filled out EcsNetworkConfiguration",
			models.NetworkConfiguration{
				models.AwsvpcConfiguration{
					[]string{
						"SubnetId1",
					},
					[]string{},
					"",
				},
			},
		},
	}

	for _, test := range tests {
		output := isEcsNetworkConfigurationFilledOut(test.networkConfigInput)
		if output != true {
			t.Errorf("The isEcsNetworkConfigurationFilledOut function did not validate that the Object was not empty for: %v", test)
		}
	}
}

func TestIsEcsNetworkConfigurationFilledOut_NotFilledOutInput(t *testing.T) {
	var tests = []struct {
		name               string
		networkConfigInput models.NetworkConfiguration
	}{
		{"A not filled out EcsNetworkConfiguration",
			models.NetworkConfiguration{
				models.AwsvpcConfiguration{
					[]string{},
					[]string{},
					"",
				},
			},
		},
	}

	for _, test := range tests {
		output := isEcsNetworkConfigurationFilledOut(test.networkConfigInput)
		if output == true {
			t.Errorf("The isEcsNetworkConfigurationFilledOut function did not validate that the Object was empty for: %v", test)
		}
	}
}

// Test validateEc2OnPremFiles
func TestValidateEc2OnPremFiles_ValidInput(t *testing.T) {
	var tests = []struct {
		name       string
		filesInput []models.File
	}{
		{"Valid File info",
			[]models.File{
				models.File{
					"sourceString",
					"destinationString",
				},
			},
		},
	}

	for _, test := range tests {
		output := validateEc2OnPremFiles(test.filesInput)
		if output != true {
			t.Errorf("The validateEc2OnPremFiles failed validation for: %v", test)
		}
	}
}

func TestValidateEc2OnPremFiles_InvalidInput(t *testing.T) {
	var tests = []struct {
		name       string
		filesInput []models.File
	}{
		{"Invalid File info, missing both",
			[]models.File{
				models.File{
					"",
					"",
				},
			},
		},
		{"Invalid File info, missing source",
			[]models.File{
				models.File{
					"",
					"destinationString",
				},
			},
		},
		{"Invalid File info, missing destination",
			[]models.File{
				models.File{
					"sourceString",
					"",
				},
			},
		},
	}

	for _, test := range tests {
		output := validateEc2OnPremFiles(test.filesInput)
		if output == true {
			t.Errorf("The validateEc2OnPremFiles succeeded but should have failed validation for: %v", test)
		}
	}
}

// Test validateEc2OnPremPermissions
func TestValidateEc2OnPremPermissions_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		permissionsInput []models.Permission
	}{
		{"Valid Permission info",
			[]models.Permission{
				models.Permission{
					"Object",
					"",
					"",
					"",
					"",
					"",
					[]string{},
					models.Context{},
					[]string{},
				},
			},
		},
		{"Valid Permission info, 2 permissions",
			[]models.Permission{
				models.Permission{
					"Object",
					"",
					"",
					"",
					"",
					"",
					[]string{},
					models.Context{},
					[]string{},
				},
				models.Permission{
					"Object2",
					"",
					"",
					"",
					"",
					"",
					[]string{},
					models.Context{},
					[]string{},
				},
			},
		},
		{"Valid Permission info, with 1 Type value",
			[]models.Permission{
				models.Permission{
					"Object",
					"",
					"",
					"",
					"",
					"",
					[]string{},
					models.Context{},
					[]string{
						"file",
					},
				},
			},
		},
		{"Valid Permission info, with 2 Type values",
			[]models.Permission{
				models.Permission{
					"Object",
					"",
					"",
					"",
					"",
					"",
					[]string{},
					models.Context{},
					[]string{
						"file", "directory",
					},
				},
			},
		},
	}

	for _, test := range tests {
		output := validateEc2OnPremPermissions(test.permissionsInput)
		if output != true {
			t.Errorf("The validateEc2OnPremPermissions failed validation for: %v", test)
		}
	}
}

func TestValidateEc2OnPremPermissions_InvalidInput(t *testing.T) {
	var tests = []struct {
		name             string
		permissionsInput []models.Permission
	}{
		{"Invalid Permission info",
			[]models.Permission{
				models.Permission{
					"",
					"",
					"",
					"",
					"",
					"",
					[]string{},
					models.Context{},
					[]string{},
				},
			},
		},
		{"Invalid Permission info, with Type",
			[]models.Permission{
				models.Permission{
					"Object",
					"",
					"",
					"",
					"",
					"",
					[]string{},
					models.Context{},
					[]string{
						"wrong",
					},
				},
			},
		},
		{"Invalid Permission info, with 2 Type values",
			[]models.Permission{
				models.Permission{
					"Object",
					"",
					"",
					"",
					"",
					"",
					[]string{},
					models.Context{},
					[]string{
						"file", "wrong",
					},
				},
			},
		},
	}

	for _, test := range tests {
		output := validateEc2OnPremPermissions(test.permissionsInput)
		if output == true {
			t.Errorf("The validateEc2OnPremPermissions succeeded but should have failed validation for: %v", test)
		}
	}
}
