package models

type LambdaAppSpecModel struct {
	Version   string
	Resources []LambdaResource
	Hooks     map[string]string
}

type LambdaResource struct {
	Type      string
	Functions map[string]Function
}

type Function struct {
	Type          string
	PropertiesObj LambdaProperties
}

type LambdaProperties struct {
	Name           string
	Alias          string
	CurrentVersion string
	TargetVersion  string
}
