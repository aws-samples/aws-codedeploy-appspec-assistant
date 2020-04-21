package assistant

import (
	"fmt"
	"testing"

	"aws-codedeploy-appspec-assistant/models"
)

// Test getServerAppSpecObjFromString
func TestGetServerAppSpecObjFromString_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		fileStrInput     string
		objectStrOutput  string
		fileExtensionVal string
	}{
		{"Valid JSON",
			serverJsonString, serverOutputStr, "json"},

		{"Valid YAML",
			serverYamlString, serverOutputStr, "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		appSpecModel, err := getServerAppSpecObjFromString([]byte(test.fileStrInput))
		if err != nil || fmt.Sprintf("%v", appSpecModel) != test.objectStrOutput {
			t.Errorf(appSpecStrConversionError)
		}
	}
}

// Test validateServerAppSpec
func TestValidateServerAppSpec_ValidInput(t *testing.T) {
	var tests = []struct {
		name             string
		fileStrInput     string
		fileExtensionVal string
	}{
		{"Valid YAML",
			serverYamlString, "yml"},
	}

	for _, test := range tests {
		fileExtension = test.fileExtensionVal
		appSpecModel, modelErr := getServerAppSpecObjFromString([]byte(test.fileStrInput))
		if modelErr != nil {
			t.Errorf("getServerAppSpecObjFromString FAILED")
		}
		err := validateServerAppSpec(appSpecModel)
		if err != nil {
			t.Errorf(appSpecObjValidationError)
		}
	}
}

// Test checkOS
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

// Test validateServerHooks
func TestValidateServerHooks_ValidInput(t *testing.T) {
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
		output := validateServerHooks(test.hooksInput)
		if output != true {
			t.Errorf("The validateServerHooks function failed for: %v", test)
		}
	}
}

func TestValidateServerHooks_InvalidInput(t *testing.T) {
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
		output := validateServerHooks(test.hooksInput)
		if output == true {
			t.Errorf("The validateServerHooks function succeeded but should have failed for: %v", test)
		}
	}
}

// Test validateServerFiles
func TestValidateServerFiles_ValidInput(t *testing.T) {
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
		output := validateServerFiles(test.filesInput)
		if output != true {
			t.Errorf("The validateServerFiles failed validation for: %v", test)
		}
	}
}

func TestValidateServerFiles_InvalidInput(t *testing.T) {
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
		output := validateServerFiles(test.filesInput)
		if output == true {
			t.Errorf("The validateServerFiles succeeded but should have failed validation for: %v", test)
		}
	}
}

// Test validateServerPermissions
func TestValidateServerPermissions_ValidInput(t *testing.T) {
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
		output := validateServerPermissions(test.permissionsInput)
		if output != true {
			t.Errorf("The validateServerPermissions failed validation for: %v", test)
		}
	}
}

func TestValidateServerPermissions_InvalidInput(t *testing.T) {
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
		output := validateServerPermissions(test.permissionsInput)
		if output == true {
			t.Errorf("The validateServerPermissions succeeded but should have failed validation for: %v", test)
		}
	}
}
