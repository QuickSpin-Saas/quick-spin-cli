package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/models"
)

// Client represents the API client
type Client struct {
	httpClient *resty.Client
	config     *config.Config
	baseURL    string
}

// ClientOption represents an option for configuring the client
type ClientOption func(*Client)

// NewClient creates a new API client
func NewClient(cfg *config.Config, opts ...ClientOption) *Client {
	client := &Client{
		config:     cfg,
		baseURL:    cfg.GetAPIURL(),
		httpClient: resty.New(),
	}

	// Configure HTTP client
	client.httpClient.
		SetBaseURL(client.baseURL).
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on 5xx errors or network errors
			return r.StatusCode() >= 500 || err != nil
		})

	// Set auth token if available
	token, err := cfg.GetToken()
	if err == nil && token != "" {
		client.httpClient.SetAuthToken(token)
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	c.httpClient.SetAuthToken(token)
}

// ClearToken clears the authentication token
func (c *Client) ClearToken() {
	c.httpClient.SetAuthToken("")
}

// Do performs an HTTP request
func (c *Client) Do(ctx context.Context, method, path string, body, result interface{}) error {
	req := c.httpClient.R().SetContext(ctx)

	if body != nil {
		req.SetBody(body)
	}

	if result != nil {
		req.SetResult(result)
	}

	// Set error result
	apiErr := &models.APIError{}
	req.SetError(apiErr)

	// Perform request
	var resp *resty.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = req.Get(path)
	case http.MethodPost:
		resp, err = req.Post(path)
	case http.MethodPut:
		resp, err = req.Put(path)
	case http.MethodPatch:
		resp, err = req.Patch(path)
	case http.MethodDelete:
		resp, err = req.Delete(path)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// Handle error responses
	if resp.IsError() {
		return c.handleErrorResponse(resp, apiErr)
	}

	return nil
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.Do(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	return c.Do(ctx, http.MethodPost, path, body, result)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	return c.Do(ctx, http.MethodPut, path, body, result)
}

// Patch performs a PATCH request
func (c *Client) Patch(ctx context.Context, path string, body, result interface{}) error {
	return c.Do(ctx, http.MethodPatch, path, body, result)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.Do(ctx, http.MethodDelete, path, nil, result)
}

// handleErrorResponse handles API error responses
func (c *Client) handleErrorResponse(resp *resty.Response, apiErr *models.APIError) error {
	// Try to parse the error response
	if resp.Body() != nil {
		var errResp struct {
			Error   string                 `json:"error"`
			Message string                 `json:"message"`
			Detail  string                 `json:"detail"`
			Details map[string]interface{} `json:"details"`
		}
		if err := json.Unmarshal(resp.Body(), &errResp); err == nil {
			apiErr.StatusCode = resp.StatusCode()
			apiErr.Code = errResp.Error
			if errResp.Message != "" {
				apiErr.Message = errResp.Message
			} else if errResp.Detail != "" {
				apiErr.Message = errResp.Detail
			} else {
				apiErr.Message = http.StatusText(resp.StatusCode())
			}
			apiErr.Details = errResp.Details
		}
	}

	// If we couldn't parse the error, create a generic one
	if apiErr.Message == "" {
		apiErr.StatusCode = resp.StatusCode()
		apiErr.Message = http.StatusText(resp.StatusCode())
	}

	// Return user-friendly error messages
	return c.getUserFriendlyError(apiErr)
}

// getUserFriendlyError converts API errors to user-friendly messages
func (c *Client) getUserFriendlyError(apiErr *models.APIError) error {
	switch apiErr.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("session expired. Please run 'qspin auth login' to authenticate")
	case http.StatusForbidden:
		return fmt.Errorf("you don't have permission to perform this action")
	case http.StatusNotFound:
		return fmt.Errorf("resource not found: %s", apiErr.Message)
	case http.StatusTooManyRequests:
		return fmt.Errorf("rate limit exceeded. Please try again later")
	case http.StatusBadRequest:
		return fmt.Errorf("invalid request: %s", apiErr.Message)
	case http.StatusConflict:
		return fmt.Errorf("conflict: %s", apiErr.Message)
	case http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
		return fmt.Errorf("QuickSpin API is experiencing issues. Please try again later")
	default:
		if apiErr.StatusCode >= 500 {
			return fmt.Errorf("server error: %s", apiErr.Message)
		}
		return apiErr
	}
}

// HealthCheck performs a health check
func (c *Client) HealthCheck(ctx context.Context) (*models.HealthCheck, error) {
	var result models.HealthCheck
	if err := c.Get(ctx, "/health", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetVersion gets the API version
func (c *Client) GetVersion(ctx context.Context) (*models.VersionInfo, error) {
	var result models.VersionInfo
	if err := c.Get(ctx, "/version", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
