package api

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/models"
)

// AdminUserListRequest represents filters for listing users
type AdminUserListRequest struct {
	Page     int    `json:"page,omitempty"`
	PerPage  int    `json:"per_page,omitempty"`
	Role     string `json:"role,omitempty"`
	Status   string `json:"status,omitempty"`
	Search   string `json:"search,omitempty"`
}

// AdminUserListResponse represents a paginated list of users
type AdminUserListResponse struct {
	Users      []models.User `json:"users"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	TotalPages int           `json:"total_pages"`
}

// ListAllUsers retrieves all users (admin only)
func (c *Client) ListAllUsers(ctx context.Context, req AdminUserListRequest) (*AdminUserListResponse, error) {
	var result AdminUserListResponse
	if err := c.Post(ctx, "/api/v1/admin/users", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserByID retrieves a user by ID (admin only)
func (c *Client) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var result models.User
	path := fmt.Sprintf("/api/v1/admin/users/%s", userID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateUserRequest represents fields that can be updated by admin
type UpdateUserRequest struct {
	Name   *string          `json:"name,omitempty"`
	Role   *models.UserRole `json:"role,omitempty"`
	Status *string          `json:"status,omitempty"`
}

// UpdateUser updates a user (admin only)
func (c *Client) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*models.User, error) {
	var result models.User
	path := fmt.Sprintf("/api/v1/admin/users/%s", userID)
	if err := c.Patch(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteUser deletes a user (admin only)
func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	path := fmt.Sprintf("/api/v1/admin/users/%s", userID)
	return c.Delete(ctx, path, nil)
}

// SuspendUser suspends a user account (admin only)
func (c *Client) SuspendUser(ctx context.Context, userID string) error {
	path := fmt.Sprintf("/api/v1/admin/users/%s/suspend", userID)
	return c.Post(ctx, path, nil, nil)
}

// ReactivateUser reactivates a suspended user account (admin only)
func (c *Client) ReactivateUser(ctx context.Context, userID string) error {
	path := fmt.Sprintf("/api/v1/admin/users/%s/reactivate", userID)
	return c.Post(ctx, path, nil, nil)
}

// OrganizationQuota represents resource quotas for an organization
type OrganizationQuota struct {
	OrganizationID string `json:"organization_id"`
	MaxServices    int    `json:"max_services"`
	MaxMembers     int    `json:"max_members"`
	MaxStorage     int64  `json:"max_storage_gb"`
	MaxBandwidth   int64  `json:"max_bandwidth_gb"`
	CustomLimits   map[string]interface{} `json:"custom_limits,omitempty"`
}

// GetOrganizationQuota retrieves quota for an organization (admin only)
func (c *Client) GetOrganizationQuota(ctx context.Context, orgID string) (*OrganizationQuota, error) {
	var result OrganizationQuota
	path := fmt.Sprintf("/api/v1/admin/organizations/%s/quota", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateOrganizationQuota updates quota for an organization (admin only)
func (c *Client) UpdateOrganizationQuota(ctx context.Context, orgID string, quota OrganizationQuota) (*OrganizationQuota, error) {
	var result OrganizationQuota
	path := fmt.Sprintf("/api/v1/admin/organizations/%s/quota", orgID)
	if err := c.Put(ctx, path, quota, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SystemHealth represents overall system health
type SystemHealth struct {
	Status       string                 `json:"status"`
	Version      string                 `json:"version"`
	Uptime       int64                  `json:"uptime_seconds"`
	Services     int                    `json:"total_services"`
	Users        int                    `json:"total_users"`
	Organizations int                   `json:"total_organizations"`
	Components   map[string]ComponentHealth `json:"components"`
}

// ComponentHealth represents health of a system component
type ComponentHealth struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency int64  `json:"latency_ms,omitempty"`
}

// GetSystemHealth retrieves overall system health (admin only)
func (c *Client) GetSystemHealth(ctx context.Context) (*SystemHealth, error) {
	var result SystemHealth
	if err := c.Get(ctx, "/api/v1/admin/health", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SystemMetrics represents system-wide metrics
type SystemMetrics struct {
	Timestamp         models.Time            `json:"timestamp"`
	CPUUsage          float64                `json:"cpu_usage_percent"`
	MemoryUsage       float64                `json:"memory_usage_percent"`
	DiskUsage         float64                `json:"disk_usage_percent"`
	ActiveConnections int                    `json:"active_connections"`
	RequestsPerMinute int                    `json:"requests_per_minute"`
	ErrorRate         float64                `json:"error_rate_percent"`
	Custom            map[string]interface{} `json:"custom,omitempty"`
}

// GetSystemMetrics retrieves system-wide metrics (admin only)
func (c *Client) GetSystemMetrics(ctx context.Context) (*SystemMetrics, error) {
	var result SystemMetrics
	if err := c.Get(ctx, "/api/v1/admin/metrics", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID             string                 `json:"id"`
	Timestamp      models.Time            `json:"timestamp"`
	UserID         string                 `json:"user_id"`
	UserEmail      string                 `json:"user_email,omitempty"`
	OrganizationID string                 `json:"organization_id,omitempty"`
	Action         string                 `json:"action"`
	Resource       string                 `json:"resource"`
	ResourceID     string                 `json:"resource_id,omitempty"`
	IPAddress      string                 `json:"ip_address,omitempty"`
	UserAgent      string                 `json:"user_agent,omitempty"`
	Details        map[string]interface{} `json:"details,omitempty"`
	Status         string                 `json:"status"`
}

// ListAuditLogs retrieves audit logs (admin only)
func (c *Client) ListAuditLogs(ctx context.Context, page, perPage int, filters map[string]string) ([]AuditLog, error) {
	var result []AuditLog
	path := fmt.Sprintf("/api/v1/admin/audit-logs?page=%d&per_page=%d", page, perPage)
	for key, value := range filters {
		path += fmt.Sprintf("&%s=%s", key, value)
	}
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// MaintenanceMode represents maintenance mode configuration
type MaintenanceMode struct {
	Enabled   bool   `json:"enabled"`
	Message   string `json:"message,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// GetMaintenanceMode retrieves maintenance mode status (admin only)
func (c *Client) GetMaintenanceMode(ctx context.Context) (*MaintenanceMode, error) {
	var result MaintenanceMode
	if err := c.Get(ctx, "/api/v1/admin/maintenance", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetMaintenanceMode enables or disables maintenance mode (admin only)
func (c *Client) SetMaintenanceMode(ctx context.Context, mode MaintenanceMode) error {
	return c.Post(ctx, "/api/v1/admin/maintenance", mode, nil)
}
