package api

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/models"
)

// DeployConfig deploys services from a configuration file
func (c *Client) DeployConfig(ctx context.Context, config models.DeploymentConfig) (*models.DeploymentResult, error) {
	var result models.DeploymentResult
	if err := c.Post(ctx, "/api/v1/deploy", config, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ValidateDeployConfig validates a deployment configuration
func (c *Client) ValidateDeployConfig(ctx context.Context, config models.DeploymentConfig) (*models.DeploymentResult, error) {
	var result models.DeploymentResult
	if err := c.Post(ctx, "/api/v1/deploy/validate", config, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDeploymentStatus retrieves the status of a deployment
func (c *Client) GetDeploymentStatus(ctx context.Context, deploymentID string) (*models.DeploymentResult, error) {
	var result models.DeploymentResult
	path := fmt.Sprintf("/api/v1/deploy/%s", deploymentID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDeployments retrieves all deployments for an organization
func (c *Client) ListDeployments(ctx context.Context, orgID string) ([]models.DeploymentResult, error) {
	var result []models.DeploymentResult
	path := fmt.Sprintf("/api/v1/organizations/%s/deployments", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// RollbackDeployment rolls back a deployment
func (c *Client) RollbackDeployment(ctx context.Context, deploymentID string) (*models.DeploymentResult, error) {
	var result models.DeploymentResult
	path := fmt.Sprintf("/api/v1/deploy/%s/rollback", deploymentID)
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExportConfig exports current services as deployment config
func (c *Client) ExportConfig(ctx context.Context, orgID string) (*models.DeploymentConfig, error) {
	var result models.DeploymentConfig
	path := fmt.Sprintf("/api/v1/organizations/%s/export", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
