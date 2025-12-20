package api

import (
	"context"

	"github.com/quickspin/quickspin-cli/internal/models"
)

// Login authenticates a user with email and password
func (c *Client) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	req := models.LoginRequest{
		Email:    email,
		Password: password,
	}

	var result models.LoginResponse
	if err := c.Post(ctx, "/api/v1/auth/login", req, &result); err != nil {
		return nil, err
	}

	// Store the token
	c.SetToken(result.Tokens.AccessToken)
	if err := c.config.SaveToken(result.Tokens.AccessToken, result.Tokens.RefreshToken); err != nil {
		return nil, err
	}

	return &result, nil
}

// Logout invalidates the current session
func (c *Client) Logout(ctx context.Context) error {
	// Call logout endpoint (optional, depends on backend implementation)
	_ = c.Post(ctx, "/api/v1/auth/logout", nil, nil)

	// Clear local credentials
	c.ClearToken()
	return c.config.ClearToken()
}

// WhoAmI returns the current user information
func (c *Client) WhoAmI(ctx context.Context) (*models.User, error) {
	var result models.User
	if err := c.Get(ctx, "/api/v1/auth/me", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RefreshToken refreshes the access token using the refresh token
func (c *Client) RefreshToken(ctx context.Context) (*models.AuthTokens, error) {
	// Load credentials to get refresh token
	token, err := c.config.GetToken()
	if err != nil {
		return nil, err
	}

	req := models.RefreshTokenRequest{
		RefreshToken: token,
	}

	var result models.AuthTokens
	if err := c.Post(ctx, "/api/v1/auth/refresh", req, &result); err != nil {
		return nil, err
	}

	// Update stored tokens
	c.SetToken(result.AccessToken)
	if err := c.config.SaveToken(result.AccessToken, result.RefreshToken); err != nil {
		return nil, err
	}

	return &result, nil
}
