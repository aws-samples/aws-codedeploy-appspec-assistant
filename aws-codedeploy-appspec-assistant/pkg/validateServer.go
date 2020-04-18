package assistant

import (
	"fmt"
	"strconv"

	"encoding/json"
	"gopkg.in/yaml.v3"

	"aws-codedeploy-appspec-assistant/globalVars"
	"aws-codedeploy-appspec-assistant/models"
)

// Convert Server (EC2/On-Prem) AppSpec string to Server (EC2/On-Prem) AppSpec Object
// Deals with JSON adn YAML
func getServerAppSpecObjFromString(appSpecBytes []byte) models.ServerAppSpecModel {
	var err error
	var serverAppSpecModel models.ServerAppSpecModel

	if fileExtension == "yml" {
		err = yaml.Unmarshal(appSpecBytes, &serverAppSpecModel)
	} else {
		err = json.Unmarshal(appSpecBytes, &serverAppSpecModel)
	}

	// Uncomment to print resulting Object for debug
	//fmt.Println(serverAppSpecModel)

	if err != nil {
		fmt.Print("\nERROR:")
		handleError(err)
	}

	return serverAppSpecModel
}

// Validate EC2/On-Prem (Server) AppSpec
// Calls validation on each section
func validateServerAppSpec(serverAppSpecModel models.ServerAppSpecModel) error {
	var err error

	// OS
	if serverAppSpecModel.OS == "" || !checkOS(serverAppSpecModel.OS) {
		numOfErrors++
		osError := globalVars.UnsupportedServerOSErr
		fmt.Println(osError)
		err = osError
	}

	// Files
	if serverAppSpecModel.Files == nil || len(serverAppSpecModel.Files) < 1 {
		numOfErrors++
		err = globalVars.MissingServerFileSpecErr
		fmt.Println(err)
	} else {
		if !validateServerFiles(serverAppSpecModel.Files) {
			err = globalVars.InvalidServerFileSpecsErr
			fmt.Println(err)
		}
	}

	// Permissions (Optional)
	if serverAppSpecModel.Permissions != nil && len(serverAppSpecModel.Permissions) > 0 {
		// All other values are optionsl
		fmt.Println("\nWARNING: All options besides Object are optional for permissions so there is very little to validate automatically.")
		if !validateServerPermissions(serverAppSpecModel.Permissions) {
			err = globalVars.InvalidServerPermissionsErr
			fmt.Println(err)
		}
	}

	// Hooks (Optional)
	if serverAppSpecModel.Hooks != nil && len(serverAppSpecModel.Hooks) > 0 {
		fmt.Println("\nWARNING: runas under Hook Scripts only applies to Amazon Linux and Ubuntu Server instances. The user also cannot require a password. Leave blank for agent default.")
		if !validateServerHooks(serverAppSpecModel.Hooks) {
			err = globalVars.InvalidServerHooksErr
			fmt.Println(err)
		}
	}

	return err
}

// EC2/OnPrem OS validation method
// Validate OS is Linux or Windows
func checkOS(appSpecOS string) bool {
	for _, supportedOS := range globalVars.AppSpecSupportedServerOSs {
		if supportedOS == appSpecOS {
			return true
		}
	}

	return false
}

// EC2/On-Prem (Server) Files Validation method
// Validate the files object values
func validateServerFiles(files []models.File) bool {
	filesValid := true
	for _, file := range files {
		if file.Source == "" {
			filesValid = false
			numOfErrors++
			fmt.Println(globalVars.MissingServerFileSourceErr)
		}

		if file.Destination == "" {
			filesValid = false
			numOfErrors++
			fmt.Println(globalVars.MissingServerFileDestinationErr)
		}
	}

	return filesValid
}

// EC2/On-Prem (Server) Permissions Validation method
// Validate the Permissions object values
func validateServerPermissions(permissions []models.Permission) bool {
	permissionsValid := true

	for _, permission := range permissions {
		if permission.Object == "" {
			permissionsValid = false
			numOfErrors++
			fmt.Println(globalVars.EmptyServerPermissionObjErr.Error())
			fmt.Println(permission)
		}

		if permission.Type != nil && len(permission.Type) > 0 {
			for _, typeStr := range permission.Type {
				if typeStr != "" && typeStr != "file" && typeStr != "directory" {
					permissionsValid = false
					numOfErrors++
					fmt.Println(globalVars.InvalidServerPermissionTypeErr.Error())
					fmt.Println(permission)
				}
			}
		}
	}

	return permissionsValid
}

// EC2/OnPrem Hooks validation methods
// Validate Hooks object
func validateServerHooks(serverHooks map[string][]models.Hook) bool {
	numValidHooks := 0
	hookScriptsValid := true
	for _, hook := range globalVars.AppSpecSupportedServerHooksWithoutLB {
		if val, ok := serverHooks[hook]; ok {
			hookScriptsValid = hookScriptsValid && validateServerHookScripts(val, hook)
			numValidHooks++
		}
	}

	withLBHooksUsed := false

	for _, hook := range globalVars.AppSpecSupportedServerHooksWithLB {
		if val, ok := serverHooks[hook]; ok {
			withLBHooksUsed = true
			hookScriptsValid = hookScriptsValid && validateServerHookScripts(val, hook)
			numValidHooks++
		}
	}

	if withLBHooksUsed {
		fmt.Println("\nWARNING: EC2/On-Prem (Server) hooks for LoadBalancers used, so the deployments should use a LoadBalancer for these scripts to be run.")
	}

	if numValidHooks == len(serverHooks) && hookScriptsValid {
		return true
	}

	if !(numValidHooks == len(serverHooks)) {
		numOfErrors++
		fmt.Println(globalVars.UnsupportedServerHooksErr)
	}

	return false
}

func validateServerHookScripts(hookScriptList []models.Hook, hook string) bool {
	scriptsValid := true
	totalTimeout := 0
	for _, hookScript := range hookScriptList {
		if hookScript.Location == "" {
			scriptsValid = false
			numOfErrors++
			fmt.Println(globalVars.MissingServerHookScriptLocationErr.Error() + hook)
		}

		if hookScript.Timeout != "" {
			i, err := strconv.Atoi(hookScript.Timeout)
			if err != nil {
				scriptsValid = false
				fmt.Println(err)
				continue
			}
			totalTimeout += i
			if totalTimeout > 3600 {
				numOfErrors++
				fmt.Println(globalVars.InvalidServerScriptTimeoutErr.Error() + hook)
				scriptsValid = false
			}
		}
	}

	return scriptsValid
}
