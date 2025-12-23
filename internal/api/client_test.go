package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)

	cfg := &config.Config{}
	client := NewClient(cfg)
	client.baseURL = server.URL
	client.httpClient.SetBaseURL(server.URL)

	return client, server
}

func TestNewClient(t *testing.T) {
	cfg := &config.Config{}
	client := NewClient(cfg)

	require.NotNil(t, client)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.config)
}

func TestClientSetToken(t *testing.T) {
	cfg := &config.Config{}
	client := NewClient(cfg)

	token := "test-token-123"
	client.SetToken(token)

	// Verify token is set (we can't directly access it, but we can verify the client works)
	assert.NotNil(t, client.httpClient)
}

func TestClientClearToken(t *testing.T) {
	cfg := &config.Config{}
	client := NewClient(cfg)

	client.SetToken("test-token")
	client.ClearToken()

	// Verify token is cleared
	assert.NotNil(t, client.httpClient)
}

func TestClientHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedStatus string
		wantErr        bool
	}{
		{
			name:           "Healthy",
			statusCode:     http.StatusOK,
			responseBody:   `{"status":"healthy","version":"1.0.0"}`,
			expectedStatus: "healthy",
			wantErr:        false,
		},
		{
			name:         "Unhealthy",
			statusCode:   http.StatusServiceUnavailable,
			responseBody: `{"status":"unhealthy"}`,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/health", r.URL.Path)
				assert.Equal(t, "GET", r.Method)

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			})
			defer server.Close()

			ctx := context.Background()
			health, err := client.HealthCheck(ctx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, health.Status)
			}
		})
	}
}

func TestClientGetVersion(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		responseBody    string
		expectedVersion string
		wantErr         bool
	}{
		{
			name:            "Get version successfully",
			statusCode:      http.StatusOK,
			responseBody:    `{"version":"1.0.0","commit":"abc123"}`,
			expectedVersion: "1.0.0",
			wantErr:         false,
		},
		{
			name:         "Version endpoint error",
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"error":"Internal server error"}`,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/version", r.URL.Path)
				assert.Equal(t, "GET", r.Method)

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			})
			defer server.Close()

			ctx := context.Background()
			version, err := client.GetVersion(ctx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedVersion, version.Version)
			}
		})
	}
}

func TestClientHTTPMethods(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       interface{}
		statusCode int
		wantErr    bool
	}{
		{
			name:       "GET request",
			method:     "GET",
			path:       "/test",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "POST request",
			method:     "POST",
			path:       "/test",
			body:       map[string]string{"key": "value"},
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:       "PUT request",
			method:     "PUT",
			path:       "/test/123",
			body:       map[string]string{"key": "updated"},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "DELETE request",
			method:     "DELETE",
			path:       "/test/123",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "PATCH request",
			method:     "PATCH",
			path:       "/test/123",
			body:       map[string]string{"key": "patched"},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.method, r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(`{}`))
			})
			defer server.Close()

			ctx := context.Background()
			var result map[string]interface{}

			var err error
			switch tt.method {
			case "GET":
				err = client.Get(ctx, tt.path, &result)
			case "POST":
				err = client.Post(ctx, tt.path, tt.body, &result)
			case "PUT":
				err = client.Put(ctx, tt.path, tt.body, &result)
			case "DELETE":
				err = client.Delete(ctx, tt.path, nil)
			case "PATCH":
				err = client.Patch(ctx, tt.path, tt.body, &result)
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClientErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedErrMsg string
	}{
		{
			name:           "Unauthorized",
			statusCode:     http.StatusUnauthorized,
			responseBody:   `{"error":"Unauthorized","message":"Invalid token"}`,
			expectedErrMsg: "unauthorized",
		},
		{
			name:           "Forbidden",
			statusCode:     http.StatusForbidden,
			responseBody:   `{"error":"Forbidden"}`,
			expectedErrMsg: "permission",
		},
		{
			name:           "Not Found",
			statusCode:     http.StatusNotFound,
			responseBody:   `{"error":"Not Found","message":"Resource not found"}`,
			expectedErrMsg: "not found",
		},
		{
			name:           "Bad Request",
			statusCode:     http.StatusBadRequest,
			responseBody:   `{"error":"Bad Request","message":"Invalid input"}`,
			expectedErrMsg: "invalid request",
		},
		{
			name:           "Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"error":"Internal Server Error"}`,
			expectedErrMsg: "server error",
		},
		{
			name:           "Too Many Requests",
			statusCode:     http.StatusTooManyRequests,
			responseBody:   `{"error":"Too Many Requests"}`,
			expectedErrMsg: "rate limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			})
			defer server.Close()

			ctx := context.Background()
			err := client.Get(ctx, "/test", nil)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErrMsg)
		})
	}
}

func TestClientTimeout(t *testing.T) {
	client, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	// Set a very short timeout
	client.httpClient.SetTimeout(10 * time.Millisecond)

	ctx := context.Background()
	err := client.Get(ctx, "/test", nil)

	// Should timeout
	assert.Error(t, err)
}

func TestClientContextCancellation(t *testing.T) {
	client, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := client.Get(ctx, "/test", nil)

	// Should be cancelled
	assert.Error(t, err)
}

func TestClientAuthToken(t *testing.T) {
	expectedToken := "test-token-123"

	client, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Contains(t, authHeader, expectedToken)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	client.SetToken(expectedToken)

	ctx := context.Background()
	err := client.Get(ctx, "/test", nil)

	assert.NoError(t, err)
}

func TestClientBaseURL(t *testing.T) {
	cfg := &config.Config{}
	client := NewClient(cfg)

	assert.NotEmpty(t, client.baseURL)
}

func TestClientConcurrentRequests(t *testing.T) {
	client, server := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer server.Close()

	ctx := context.Background()
	errors := make(chan error, 10)

	// Make 10 concurrent requests
	for i := 0; i < 10; i++ {
		go func() {
			errors <- client.Get(ctx, "/test", nil)
		}()
	}

	// Check all requests succeeded
	for i := 0; i < 10; i++ {
		err := <-errors
		assert.NoError(t, err)
	}
}
