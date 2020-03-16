package models

type EcsAppSpecModel struct {
	Version   string
	Resources []EcsResource
	Hooks     map[string]string
}

type EcsResource struct {
	Type          string
	PropertiesObj EcsProperties
}

type EcsProperties struct {
	TaskDefinition      string
	LoadBalancerInfoObj LoadBalancerInfo
}

type LoadBalancerInfo struct {
	ContainerName string
	ContainerPort int

	// Optional properties
	PlatformVersion         string
	NetworkConfigurationObj NetworkConfiguration
}

type NetworkConfiguration struct {
	AwsvpcConfigurationObj AwsvpcConfiguration
}

type AwsvpcConfiguration struct {
	Subnets        []string
	SecurityGroups []string
	AssignPublicIp string
}
