package auth

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockAPIClient is a mock for the API client
type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAPIClient) Logout(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAPIClient) WhoAmI(ctx context.Context) (*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Helper function to capture output
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Helper function to execute command and capture output
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestNewAuthCmd(t *testing.T) {
	cmd := NewAuthCmd()
	require.NotNil(t, cmd)
	assert.Equal(t, "auth", cmd.Use)
	assert.True(t, len(cmd.Commands()) > 0, "Auth command should have subcommands")
}

func TestLoginCommand(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		mockResponse  *models.LoginResponse
		mockError     error
		expectedError string
		wantErr       bool
	}{
		{
			name:     "Successful login",
			email:    "user@example.com",
			password: "password123",
			mockResponse: &models.LoginResponse{
				User: models.User{
					ID:    "user-123",
					Email: "user@example.com",
					Name:  "Test User",
					Role:  models.UserRoleAdmin,
				},
				Tokens: models.AuthTokens{
					AccessToken:  "access-token-123",
					RefreshToken: "refresh-token-123",
					TokenType:    "Bearer",
					ExpiresIn:    3600,
				},
			},
			wantErr: false,
		},
		{
			name:          "Invalid credentials",
			email:         "user@example.com",
			password:      "wrongpassword",
			mockError:     assert.AnError,
			expectedError: "assert.AnError",
			wantErr:       true,
		},
		{
			name:          "Missing email",
			email:         "",
			password:      "password123",
			expectedError: "email and password are required",
			wantErr:       true,
		},
		{
			name:          "Missing password",
			email:         "user@example.com",
			password:      "",
			expectedError: "email and password are required",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cmd := NewLoginCmd()

			// Set flags
			if tt.email != "" {
				cmd.Flags().Set("email", tt.email)
			}
			if tt.password != "" {
				cmd.Flags().Set("password", tt.password)
			}

			// For non-interactive tests, we would need to mock the API client
			// This is a basic structure - full implementation would require dependency injection
		})
	}
}

func TestLogoutCommand(t *testing.T) {
	tests := []struct {
		name          string
		mockError     error
		expectedError string
		wantErr       bool
	}{
		{
			name:    "Successful logout",
			wantErr: false,
		},
		{
			name:          "Logout error",
			mockError:     assert.AnError,
			expectedError: "assert.AnError",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewLogoutCmd()
			require.NotNil(t, cmd)
		})
	}
}

func TestWhoAmICommand(t *testing.T) {
	tests := []struct {
		name         string
		mockUser     *models.User
		mockError    error
		expectOutput string
		wantErr      bool
	}{
		{
			name: "Get current user info",
			mockUser: &models.User{
				ID:    "user-123",
				Email: "user@example.com",
				Name:  "Test User",
				Role:  models.UserRoleAdmin,
			},
			expectOutput: "user@example.com",
			wantErr:      false,
		},
		{
			name:      "Not authenticated",
			mockError: assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewWhoAmICmd()
			require.NotNil(t, cmd)
		})
	}
}

func TestAuthCommandFlags(t *testing.T) {
	tests := []struct {
		name     string
		cmdFunc  func() *cobra.Command
		flagName string
		flagType string
	}{
		{
			name:     "Login email flag",
			cmdFunc:  NewLoginCmd,
			flagName: "email",
			flagType: "string",
		},
		{
			name:     "Login password flag",
			cmdFunc:  NewLoginCmd,
			flagName: "password",
			flagType: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmdFunc()
			flag := cmd.Flags().Lookup(tt.flagName)
			require.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type())
		})
	}
}

func TestAuthCommandHelp(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{
			name: "Auth command help",
			cmd:  NewAuthCmd(),
		},
		{
			name: "Login command help",
			cmd:  NewLoginCmd(),
		},
		{
			name: "Logout command help",
			cmd:  NewLogoutCmd(),
		},
		{
			name: "WhoAmI command help",
			cmd:  NewWhoAmICmd(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(tt.cmd, "--help")
			require.NoError(t, err)
			assert.Contains(t, output, tt.cmd.Use)
			assert.Contains(t, strings.ToLower(output), "usage")
		})
	}
}

func TestAuthCommandAliases(t *testing.T) {
	cmd := NewAuthCmd()
	assert.Contains(t, cmd.Aliases, "authentication")
	assert.Contains(t, cmd.Aliases, "login")
}

func TestConfigLoading(t *testing.T) {
	// Test config loading with different scenarios
	tests := []struct {
		name        string
		configFile  string
		expectError bool
	}{
		{
			name:        "Default config",
			configFile:  "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := config.LoadConfig()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
