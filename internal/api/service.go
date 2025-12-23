package api

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/models"
)

// ListServices retrieves all services for the current user/organization
func (c *Client) ListServices(ctx context.Context) ([]models.Service, error) {
	var result []models.Service
	if err := c.Get(ctx, "/api/v1/services", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetService retrieves a specific service by ID
func (c *Client) GetService(ctx context.Context, serviceID string) (*models.Service, error) {
	var result models.Service
	path := fmt.Sprintf("/api/v1/services/%s", serviceID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateServiceRequest represents the request to create a new service
type CreateServiceRequest struct {
	Name        string              `json:"name"`
	Type        models.ServiceType  `json:"type"`
	Tier        models.ServiceTier  `json:"tier"`
	Region      string              `json:"region,omitempty"`
	Description string              `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// CreateService creates a new service
func (c *Client) CreateService(ctx context.Context, req CreateServiceRequest) (*models.Service, error) {
	var result models.Service
	if err := c.Post(ctx, "/api/v1/services", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateServiceRequest represents the request to update a service
type UpdateServiceRequest struct {
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// UpdateService updates an existing service
func (c *Client) UpdateService(ctx context.Context, serviceID string, req UpdateServiceRequest) (*models.Service, error) {
	var result models.Service
	path := fmt.Sprintf("/api/v1/services/%s", serviceID)
	if err := c.Patch(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteService deletes a service by ID
func (c *Client) DeleteService(ctx context.Context, serviceID string) error {
	path := fmt.Sprintf("/api/v1/services/%s", serviceID)
	return c.Delete(ctx, path, nil)
}

// ScaleService scales a service to a different tier
func (c *Client) ScaleService(ctx context.Context, serviceID string, tier models.ServiceTier) (*models.Service, error) {
	var result models.Service
	req := struct {
		Tier models.ServiceTier `json:"tier"`
	}{
		Tier: tier,
	}
	path := fmt.Sprintf("/api/v1/services/%s/scale", serviceID)
	if err := c.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetServiceLogs retrieves logs for a service
func (c *Client) GetServiceLogs(ctx context.Context, serviceID string, lines int) ([]models.ServiceLogEntry, error) {
	var result []models.ServiceLogEntry
	path := fmt.Sprintf("/api/v1/services/%s/logs?lines=%d", serviceID, lines)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetServiceMetrics retrieves metrics for a service
func (c *Client) GetServiceMetrics(ctx context.Context, serviceID string) (*models.ServiceMetrics, error) {
	var result models.ServiceMetrics
	path := fmt.Sprintf("/api/v1/services/%s/metrics", serviceID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
