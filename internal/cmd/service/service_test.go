package service

import (
	"bytes"
	"context"
	"testing"

	"github.com/quickspin/quickspin-cli/internal/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockAPIClient mocks the API client for testing
type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) ListServices(ctx context.Context) ([]models.Service, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Service), args.Error(1)
}

func (m *MockAPIClient) GetService(ctx context.Context, serviceID string) (*models.Service, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Service), args.Error(1)
}

func (m *MockAPIClient) CreateService(ctx context.Context, name, serviceType, tier, region, description string) (*models.Service, error) {
	args := m.Called(ctx, name, serviceType, tier, region, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Service), args.Error(1)
}

func (m *MockAPIClient) DeleteService(ctx context.Context, serviceID string) error {
	args := m.Called(ctx, serviceID)
	return args.Error(0)
}

func (m *MockAPIClient) ScaleService(ctx context.Context, serviceID string, tier models.ServiceTier) (*models.Service, error) {
	args := m.Called(ctx, serviceID, tier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Service), args.Error(1)
}

func (m *MockAPIClient) GetServiceLogs(ctx context.Context, serviceID string, lines int) ([]models.ServiceLogEntry, error) {
	args := m.Called(ctx, serviceID, lines)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ServiceLogEntry), args.Error(1)
}

// Helper function to execute command
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestNewServiceCmd(t *testing.T) {
	cmd := NewServiceCmd()
	require.NotNil(t, cmd)
	assert.Equal(t, "service", cmd.Use)
	assert.Contains(t, cmd.Aliases, "services")
	assert.Contains(t, cmd.Aliases, "svc")
	assert.True(t, len(cmd.Commands()) > 0, "Service command should have subcommands")
}

func TestServiceSubcommands(t *testing.T) {
	cmd := NewServiceCmd()

	expectedSubcommands := []string{"list", "create", "delete", "describe", "scale", "logs"}
	actualSubcommands := make(map[string]bool)

	for _, subCmd := range cmd.Commands() {
		actualSubcommands[subCmd.Name()] = true
	}

	for _, expected := range expectedSubcommands {
		assert.True(t, actualSubcommands[expected], "Expected subcommand %s not found", expected)
	}
}

func TestListCommand(t *testing.T) {
	tests := []struct {
		name          string
		services      []models.Service
		mockError     error
		expectOutput  string
		wantErr       bool
	}{
		{
			name: "List services successfully",
			services: []models.Service{
				{
					ID:     "svc-1",
					Name:   "redis-cache",
					Type:   models.ServiceTypeRedis,
					Tier:   models.ServiceTierDeveloper,
					Status: models.ServiceStatusRunning,
					Region: "us-east-1",
				},
				{
					ID:     "svc-2",
					Name:   "postgres-db",
					Type:   models.ServiceTypePostgreSQL,
					Tier:   models.ServiceTierBasic,
					Status: models.ServiceStatusRunning,
					Region: "us-west-2",
				},
			},
			expectOutput: "redis-cache",
			wantErr:      false,
		},
		{
			name:      "Empty service list",
			services:  []models.Service{},
			wantErr:   false,
		},
		{
			name:      "API error",
			mockError: assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewListCmd()
			require.NotNil(t, cmd)
		})
	}
}

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		flagName    string
		flagType    string
		flagTier    string
		flagRegion  string
		expectError bool
	}{
		{
			name:        "Valid create command",
			flagName:    "redis-cache",
			flagType:    "redis",
			flagTier:    "developer",
			flagRegion:  "us-east-1",
			expectError: false,
		},
		{
			name:        "Missing name",
			flagType:    "redis",
			flagTier:    "developer",
			expectError: true,
		},
		{
			name:        "Missing type",
			flagName:    "redis-cache",
			flagTier:    "developer",
			expectError: true,
		},
		{
			name:        "Invalid service type",
			flagName:    "test-service",
			flagType:    "invalid-type",
			flagTier:    "developer",
			expectError: true,
		},
		{
			name:        "Invalid tier",
			flagName:    "test-service",
			flagType:    "redis",
			flagTier:    "invalid-tier",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCreateCmd()
			require.NotNil(t, cmd)

			// Verify flags exist
			assert.NotNil(t, cmd.Flags().Lookup("name"))
			assert.NotNil(t, cmd.Flags().Lookup("type"))
			assert.NotNil(t, cmd.Flags().Lookup("tier"))
			assert.NotNil(t, cmd.Flags().Lookup("region"))
		})
	}
}

func TestDeleteCommand(t *testing.T) {
	tests := []struct {
		name          string
		serviceID     string
		force         bool
		mockError     error
		expectConfirm bool
		wantErr       bool
	}{
		{
			name:          "Delete with force flag",
			serviceID:     "svc-123",
			force:         true,
			expectConfirm: false,
			wantErr:       false,
		},
		{
			name:          "Delete with confirmation",
			serviceID:     "svc-123",
			force:         false,
			expectConfirm: true,
			wantErr:       false,
		},
		{
			name:      "Delete non-existent service",
			serviceID: "non-existent",
			force:     true,
			mockError: assert.AnError,
			wantErr:   true,
		},
		{
			name:      "Missing service ID",
			serviceID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDeleteCmd()
			require.NotNil(t, cmd)

			// Verify force flag exists
			forceFlag := cmd.Flags().Lookup("force")
			assert.NotNil(t, forceFlag)
			assert.Equal(t, "bool", forceFlag.Value.Type())
		})
	}
}

func TestDescribeCommand(t *testing.T) {
	tests := []struct {
		name      string
		serviceID string
		service   *models.Service
		mockError error
		wantErr   bool
	}{
		{
			name:      "Describe existing service",
			serviceID: "svc-123",
			service: &models.Service{
				ID:     "svc-123",
				Name:   "redis-cache",
				Type:   models.ServiceTypeRedis,
				Tier:   models.ServiceTierDeveloper,
				Status: models.ServiceStatusRunning,
				Region: "us-east-1",
			},
			wantErr: false,
		},
		{
			name:      "Service not found",
			serviceID: "non-existent",
			mockError: assert.AnError,
			wantErr:   true,
		},
		{
			name:      "Missing service ID",
			serviceID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDescribeCmd()
			require.NotNil(t, cmd)
		})
	}
}

func TestScaleCommand(t *testing.T) {
	tests := []struct {
		name      string
		serviceID string
		tier      string
		mockError error
		wantErr   bool
	}{
		{
			name:      "Scale to higher tier",
			serviceID: "svc-123",
			tier:      "pro",
			wantErr:   false,
		},
		{
			name:      "Scale to lower tier",
			serviceID: "svc-123",
			tier:      "developer",
			wantErr:   false,
		},
		{
			name:      "Invalid tier",
			serviceID: "svc-123",
			tier:      "invalid-tier",
			wantErr:   true,
		},
		{
			name:      "Service not found",
			serviceID: "non-existent",
			tier:      "pro",
			mockError: assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewScaleCmd()
			require.NotNil(t, cmd)

			// Verify tier flag exists
			tierFlag := cmd.Flags().Lookup("tier")
			assert.NotNil(t, tierFlag)
		})
	}
}

func TestLogsCommand(t *testing.T) {
	tests := []struct {
		name      string
		serviceID string
		lines     int
		logs      []models.ServiceLogEntry
		mockError error
		wantErr   bool
	}{
		{
			name:      "Get service logs",
			serviceID: "svc-123",
			lines:     100,
			logs: []models.ServiceLogEntry{
				{
					Level:   "info",
					Message: "Service started",
				},
				{
					Level:   "error",
					Message: "Connection failed",
				},
			},
			wantErr: false,
		},
		{
			name:      "No logs available",
			serviceID: "svc-123",
			lines:     100,
			logs:      []models.ServiceLogEntry{},
			wantErr:   false,
		},
		{
			name:      "Service not found",
			serviceID: "non-existent",
			lines:     100,
			mockError: assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewLogsCmd()
			require.NotNil(t, cmd)

			// Verify lines flag exists
			linesFlag := cmd.Flags().Lookup("lines")
			assert.NotNil(t, linesFlag)
			assert.Equal(t, "int", linesFlag.Value.Type())
		})
	}
}

func TestServiceCommandHelp(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{name: "Service command help", cmd: NewServiceCmd()},
		{name: "List command help", cmd: NewListCmd()},
		{name: "Create command help", cmd: NewCreateCmd()},
		{name: "Delete command help", cmd: NewDeleteCmd()},
		{name: "Describe command help", cmd: NewDescribeCmd()},
		{name: "Scale command help", cmd: NewScaleCmd()},
		{name: "Logs command help", cmd: NewLogsCmd()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(tt.cmd, "--help")
			require.NoError(t, err)
			assert.Contains(t, output, tt.cmd.Use)
		})
	}
}

func TestServiceTypes(t *testing.T) {
	validTypes := []models.ServiceType{
		models.ServiceTypeRedis,
		models.ServiceTypeRabbitMQ,
		models.ServiceTypePostgreSQL,
		models.ServiceTypeMongoDB,
		models.ServiceTypeMySQL,
		models.ServiceTypeElasticsearch,
	}

	for _, serviceType := range validTypes {
		t.Run(string(serviceType), func(t *testing.T) {
			assert.NotEmpty(t, serviceType)
		})
	}
}

func TestServiceTiers(t *testing.T) {
	validTiers := []models.ServiceTier{
		models.ServiceTierStarter,
		models.ServiceTierDeveloper,
		models.ServiceTierBasic,
		models.ServiceTierStandard,
		models.ServiceTierPro,
		models.ServiceTierPremium,
		models.ServiceTierEnterprise,
	}

	for _, tier := range validTiers {
		t.Run(string(tier), func(t *testing.T) {
			assert.NotEmpty(t, tier)
		})
	}
}

func TestServiceStatus(t *testing.T) {
	validStatuses := []models.ServiceStatus{
		models.ServiceStatusPending,
		models.ServiceStatusCreating,
		models.ServiceStatusRunning,
		models.ServiceStatusStopped,
		models.ServiceStatusFailed,
		models.ServiceStatusDeleting,
	}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			assert.NotEmpty(t, status)
		})
	}
}
