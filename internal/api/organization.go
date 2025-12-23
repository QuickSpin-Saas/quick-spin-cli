package api

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/models"
)

// ListOrganizations retrieves all organizations for the current user
func (c *Client) ListOrganizations(ctx context.Context) ([]models.Organization, error) {
	var result []models.Organization
	if err := c.Get(ctx, "/api/v1/organizations", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetOrganization retrieves a specific organization by ID
func (c *Client) GetOrganization(ctx context.Context, orgID string) (*models.Organization, error) {
	var result models.Organization
	path := fmt.Sprintf("/api/v1/organizations/%s", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateOrganizationRequest represents the request to create a new organization
type CreateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

// CreateOrganization creates a new organization
func (c *Client) CreateOrganization(ctx context.Context, req CreateOrganizationRequest) (*models.Organization, error) {
	var result models.Organization
	if err := c.Post(ctx, "/api/v1/organizations", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateOrganizationRequest represents the request to update an organization
type UpdateOrganizationRequest struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

// UpdateOrganization updates an existing organization
func (c *Client) UpdateOrganization(ctx context.Context, orgID string, req UpdateOrganizationRequest) (*models.Organization, error) {
	var result models.Organization
	path := fmt.Sprintf("/api/v1/organizations/%s", orgID)
	if err := c.Patch(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteOrganization deletes an organization by ID
func (c *Client) DeleteOrganization(ctx context.Context, orgID string) error {
	path := fmt.Sprintf("/api/v1/organizations/%s", orgID)
	return c.Delete(ctx, path, nil)
}

// ListOrganizationMembers retrieves all members of an organization
func (c *Client) ListOrganizationMembers(ctx context.Context, orgID string) ([]models.OrganizationMember, error) {
	var result []models.OrganizationMember
	path := fmt.Sprintf("/api/v1/organizations/%s/members", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// InviteMember invites a new member to an organization
func (c *Client) InviteMember(ctx context.Context, orgID string, req models.InviteMemberRequest) error {
	path := fmt.Sprintf("/api/v1/organizations/%s/members", orgID)
	return c.Post(ctx, path, req, nil)
}

// RemoveMember removes a member from an organization
func (c *Client) RemoveMember(ctx context.Context, orgID, userID string) error {
	path := fmt.Sprintf("/api/v1/organizations/%s/members/%s", orgID, userID)
	return c.Delete(ctx, path, nil)
}

// UpdateMemberRoleRequest represents the request to update a member's role
type UpdateMemberRoleRequest struct {
	Role models.UserRole `json:"role"`
}

// UpdateMemberRole updates a member's role in an organization
func (c *Client) UpdateMemberRole(ctx context.Context, orgID, userID string, role models.UserRole) error {
	req := UpdateMemberRoleRequest{Role: role}
	path := fmt.Sprintf("/api/v1/organizations/%s/members/%s", orgID, userID)
	return c.Patch(ctx, path, req, nil)
}

// SwitchOrganization switches the current active organization
func (c *Client) SwitchOrganization(ctx context.Context, orgID string) error {
	req := struct {
		OrganizationID string `json:"organization_id"`
	}{
		OrganizationID: orgID,
	}
	return c.Post(ctx, "/api/v1/auth/switch-org", req, nil)
}

// GetCurrentOrganization retrieves the currently active organization
func (c *Client) GetCurrentOrganization(ctx context.Context) (*models.Organization, error) {
	var result models.Organization
	if err := c.Get(ctx, "/api/v1/auth/current-org", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
