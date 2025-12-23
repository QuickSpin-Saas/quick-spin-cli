package models

import (
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			t.Time = parsed
			return nil
		}
	}

	return nil
}

type APIError struct {
	StatusCode int                    `json:"status_code"`
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

// HealthCheck represents a health check response
type HealthCheck struct {
	Status  string                 `json:"status"`
	Version string                 `json:"version,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// VersionInfo represents version information
type VersionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit,omitempty"`
	Date    string `json:"date,omitempty"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Total       int `json:"total"`
	Page        int `json:"page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}
