package assistant

import (
	"fmt"
	"testing"

	"aws-codedeploy-appspec-assistant/models"
)

// Test getEcsAppSpecObjFromString
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
		appSpecModel, err := getEcsAppSpecObjFromString([]byte(test.fileStrInput))
		if err != nil || fmt.Sprintf("%v", appSpecModel) != test.objectStrOutput {
			t.Errorf(appSpecStrConversionError)
		}
	}
}

// Test validateEcsAppSpec
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
		appSpecModel, modelErr := getEcsAppSpecObjFromString([]byte(test.fileStrInput))
		if modelErr != nil {
			t.Errorf("getEcsAppSpecObjFromString FAILED")
		}
		err := validateEcsAppSpec(appSpecModel)
		if err != nil {
			t.Errorf(appSpecObjValidationError)
		}
	}
}

// Test validateEcsHooks
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

// Test validateEcsResources
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

// Test isEcsNetworkConfigurationFilledOut
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
