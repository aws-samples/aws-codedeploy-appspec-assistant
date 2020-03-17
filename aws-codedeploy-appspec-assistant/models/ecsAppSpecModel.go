package models

type EcsAppSpecModel struct {
	Version   float32    `json:"version" yaml:"version"`
	Resources []Resource `json:"Resources" yaml:"Resources"`

	// Optional
	Hooks []map[string]string `json:"Hooks" yaml:"Hooks"`
}

type Resource struct {
	TargetService TargetService `json:"TargetService" yaml:"TargetService"`
}

type TargetService struct {
	Type       string        `json:"Type" yaml:"Type"`
	Properties EcsProperties `json:"Properties" yaml:"Properties"`
}

type EcsProperties struct {
	TaskDefinition   string           `json:"TaskDefinition" yaml:"TaskDefinition"`
	LoadBalancerInfo LoadBalancerInfo `json:"LoadBalancerInfo" yaml:"LoadBalancerInfo"`

	// Optional
	PlatformVersion      string               `json:"PlatformVersion" yaml:"PlatformVersion"`
	NetworkConfiguration NetworkConfiguration `json:"NetworkConfiguration" yaml:"NetworkConfiguration"`
}

type LoadBalancerInfo struct {
	ContainerName string `json:"ContainerName" yaml:"ContainerName"`
	ContainerPort int    `json:"ContainerPort" yaml:"ContainerPort"`
}

type NetworkConfiguration struct {
	AwsvpcConfiguration AwsvpcConfiguration `json:"AwsvpcConfiguration" yaml:"AwsvpcConfiguration"`
}

type AwsvpcConfiguration struct {
	Subnets        []string `json:"Subnets" yaml:"Subnets"`
	SecurityGroups []string `json:"SecurityGroups" yaml:"SecurityGroups"`
	AssignPublicIp string   `json:"AssignPublicIp" yaml:"AssignPublicIp"`
}
