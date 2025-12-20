package models

import "time"

// ServiceType represents the type of service
type ServiceType string

const (
	ServiceTypeRedis         ServiceType = "redis"
	ServiceTypeRabbitMQ      ServiceType = "rabbitmq"
	ServiceTypeElasticsearch ServiceType = "elasticsearch"
	ServiceTypePostgreSQL    ServiceType = "postgresql"
	ServiceTypeMongoDB       ServiceType = "mongodb"
	ServiceTypeMySQL         ServiceType = "mysql"
)

// ServiceTier represents the pricing tier
type ServiceTier string

const (
	ServiceTierStarter    ServiceTier = "starter"
	ServiceTierDeveloper  ServiceTier = "developer"
	ServiceTierBasic      ServiceTier = "basic"
	ServiceTierStandard   ServiceTier = "standard"
	ServiceTierPro        ServiceTier = "pro"
	ServiceTierPremium    ServiceTier = "premium"
	ServiceTierEnterprise ServiceTier = "enterprise"
)

// ServiceStatus represents the current status of a service
type ServiceStatus string

const (
	ServiceStatusPending  ServiceStatus = "pending"
	ServiceStatusCreating ServiceStatus = "creating"
	ServiceStatusRunning  ServiceStatus = "running"
	ServiceStatusStopped  ServiceStatus = "stopped"
	ServiceStatusFailed   ServiceStatus = "failed"
	ServiceStatusDeleting ServiceStatus = "deleting"
)

// Service represents a managed service instance
type Service struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Type           ServiceType            `json:"type"`
	Tier           ServiceTier            `json:"tier"`
	Status         ServiceStatus          `json:"status"`
	Region         string                 `json:"region"`
	OrganizationID string                 `json:"organization_id"`
	Config         map[string]interface{} `json:"config,omitempty"`
	Labels         map[string]string      `json:"labels,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Credentials    *ServiceCredentials    `json:"credentials,omitempty"`
	Resources      *ServiceResources      `json:"resources,omitempty"`
	Namespace      string                 `json:"namespace,omitempty"`
}

// ServiceCredentials holds connection information for a service
type ServiceCredentials struct {
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Database string            `json:"database,omitempty"`
	URI      string            `json:"uri,omitempty"`
	Extra    map[string]string `json:"extra,omitempty"`
}

// ServiceResources holds resource allocation information
type ServiceResources struct {
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
}

// ServiceCreateRequest represents a request to create a new service
type ServiceCreateRequest struct {
	Name   string                 `json:"name"`
	Type   ServiceType            `json:"type"`
	Tier   ServiceTier            `json:"tier"`
	Region string                 `json:"region"`
	Config map[string]interface{} `json:"config,omitempty"`
	Labels map[string]string      `json:"labels,omitempty"`
}

// ServiceUpdateRequest represents a request to update a service
type ServiceUpdateRequest struct {
	Tier   *ServiceTier           `json:"tier,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
	Labels map[string]string      `json:"labels,omitempty"`
}

// ServiceScaleRequest represents a request to scale a service
type ServiceScaleRequest struct {
	Tier     *ServiceTier `json:"tier,omitempty"`
	Replicas *int         `json:"replicas,omitempty"`
}

// ServiceListResponse represents a list of services
type ServiceListResponse struct {
	Services []Service `json:"services"`
	Total    int       `json:"total"`
}

// ServiceLogEntry represents a single log entry
type ServiceLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source,omitempty"`
}

// ServiceMetrics represents service metrics
type ServiceMetrics struct {
	ServiceID string                 `json:"service_id"`
	Timestamp time.Time              `json:"timestamp"`
	CPU       float64                `json:"cpu_percent"`
	Memory    float64                `json:"memory_percent"`
	Storage   float64                `json:"storage_percent"`
	Network   *NetworkMetrics        `json:"network,omitempty"`
	Custom    map[string]interface{} `json:"custom,omitempty"`
}

// NetworkMetrics represents network metrics
type NetworkMetrics struct {
	BytesIn  int64 `json:"bytes_in"`
	BytesOut int64 `json:"bytes_out"`
}

// ServiceTypeInfo represents information about a service type
type ServiceTypeInfo struct {
	Type        ServiceType `json:"type"`
	DisplayName string      `json:"display_name"`
	Description string      `json:"description"`
	Icon        string      `json:"icon,omitempty"`
	Tiers       []TierInfo  `json:"tiers"`
}

// TierInfo represents information about a pricing tier
type TierInfo struct {
	Tier        ServiceTier `json:"tier"`
	DisplayName string      `json:"display_name"`
	CPU         string      `json:"cpu"`
	Memory      string      `json:"memory"`
	Storage     string      `json:"storage"`
	Price       float64     `json:"price_monthly"`
	Features    []string    `json:"features,omitempty"`
}
