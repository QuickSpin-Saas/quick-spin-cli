package models

// DeploymentConfig represents a deployment configuration file
type DeploymentConfig struct {
	Version      string            `yaml:"version"`
	Organization string            `yaml:"organization"`
	Services     []ServiceTemplate `yaml:"services"`
}

// ServiceTemplate represents a service definition in deployment config
type ServiceTemplate struct {
	Name   string                 `yaml:"name"`
	Type   ServiceType            `yaml:"type"`
	Tier   ServiceTier            `yaml:"tier"`
	Region string                 `yaml:"region"`
	Config map[string]interface{} `yaml:"config,omitempty"`
	Labels map[string]string      `yaml:"labels,omitempty"`
}

// DeploymentResult represents the result of a deployment operation
type DeploymentResult struct {
	Success         bool              `json:"success"`
	ServicesCreated []string          `json:"services_created,omitempty"`
	ServicesFailed  []DeploymentError `json:"services_failed,omitempty"`
	Message         string            `json:"message"`
}

// DeploymentError represents an error during deployment
type DeploymentError struct {
	ServiceName string `json:"service_name"`
	Error       string `json:"error"`
}
