package globalVars

var AppSpecVersions = [...]string{"0.0"}

var AppSpecSupportedServerOSs = [...]string{"linux", "windows"}

var AppSpecSupportedEcsHooks = [...]string{"BeforeInstall", "AfterInstall", "AfterAllowTestTraffic", "BeforeAllowTraffic", "AfterAllowTraffic"}
var AppSpecSupportedLambdaHooks = [...]string{"BeforeAllowTraffic", "AfterAllowTraffic"}
var AppSpecSupportedServerHooksWithLB = [...]string{"BeforeBlockTraffic", "AfterBlockTraffic", "BeforeAllowTraffic", "AfterAllowTraffic"}
var AppSpecSupportedServerHooksWithoutLB = [...]string{"ApplicationStop", "BeforeInstall", "AfterInstall", "ApplicationStart", "ValidateService"}

var AppSpecEcsAssignPublicIpValues = [...]string{"ENABLED", "DISABLED"}
