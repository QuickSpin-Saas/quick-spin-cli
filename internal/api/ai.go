package api

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/models"
)

// GetRecommendations retrieves AI-powered recommendations
func (c *Client) GetRecommendations(ctx context.Context, orgID string, req models.RecommendationRequest) (*models.RecommendationResponse, error) {
	var result models.RecommendationResponse
	path := fmt.Sprintf("/api/v1/organizations/%s/ai/recommendations", orgID)
	if err := c.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetServiceRecommendations retrieves recommendations for a specific service
func (c *Client) GetServiceRecommendations(ctx context.Context, serviceID string) (*models.RecommendationResponse, error) {
	var result models.RecommendationResponse
	path := fmt.Sprintf("/api/v1/services/%s/recommendations", serviceID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AnalyzeService performs AI analysis on a service
func (c *Client) AnalyzeService(ctx context.Context, serviceID string) (*models.AnalysisResult, error) {
	var result models.AnalysisResult
	path := fmt.Sprintf("/api/v1/services/%s/analyze", serviceID)
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AnalyzeOrganization performs AI analysis on an organization
func (c *Client) AnalyzeOrganization(ctx context.Context, orgID string) (*models.AnalysisResult, error) {
	var result models.AnalysisResult
	path := fmt.Sprintf("/api/v1/organizations/%s/ai/analyze", orgID)
	if err := c.Post(ctx, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetOptimizationSuggestions retrieves optimization suggestions
func (c *Client) GetOptimizationSuggestions(ctx context.Context, orgID, focus string) (*models.OptimizationResponse, error) {
	var result models.OptimizationResponse
	path := fmt.Sprintf("/api/v1/organizations/%s/ai/optimize?focus=%s", orgID, focus)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetServiceOptimization retrieves optimization suggestions for a specific service
func (c *Client) GetServiceOptimization(ctx context.Context, serviceID string) (*models.OptimizationResponse, error) {
	var result models.OptimizationResponse
	path := fmt.Sprintf("/api/v1/services/%s/optimize", serviceID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAnomalies retrieves detected anomalies
func (c *Client) ListAnomalies(ctx context.Context, orgID string) ([]models.Anomaly, error) {
	var result []models.Anomaly
	path := fmt.Sprintf("/api/v1/organizations/%s/ai/anomalies", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetServiceAnomalies retrieves anomalies for a specific service
func (c *Client) GetServiceAnomalies(ctx context.Context, serviceID string) ([]models.Anomaly, error) {
	var result []models.Anomaly
	path := fmt.Sprintf("/api/v1/services/%s/anomalies", serviceID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ResolveAnomaly marks an anomaly as resolved
func (c *Client) ResolveAnomaly(ctx context.Context, anomalyID string) error {
	path := fmt.Sprintf("/api/v1/ai/anomalies/%s/resolve", anomalyID)
	return c.Post(ctx, path, nil, nil)
}

// ChatRequest represents a request to the AI assistant
type ChatRequest struct {
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// ChatResponse represents a response from the AI assistant
type ChatResponse struct {
	Response    string   `json:"response"`
	Suggestions []string `json:"suggestions,omitempty"`
	Actions     []string `json:"actions,omitempty"`
}

// Chat sends a message to the AI assistant
func (c *Client) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	var result ChatResponse
	if err := c.Post(ctx, "/api/v1/ai/chat", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
