package models

type LambdaAppSpecModel struct {
	Version   float32               `json:"version" yaml:"version"`
	Resources []map[string]Function `json:"Resources" yaml:"Resources"`

	// Optional
	Hooks []map[string]string `json:"Hooks" yaml:"Hooks"`
}

type Function struct {
	Type       string           `json:"Type" yaml:"Type"`
	Properties LambdaProperties `json:"Properties" yaml:"Properties"`
}

type LambdaProperties struct {
	Name           string `json:"Name" yaml:"Name"`
	Alias          string `json:"Alias" yaml:"Alias"`
	CurrentVersion string `json:"CurrentVersion" yaml:"CurrentVersion"`
	TargetVersion  string `json:"TargetVersion" yaml:"TargetVersion"`
}
