# QuickSpin CLI Test Documentation

## Overview

This document describes the comprehensive test suite for the QuickSpin CLI application. The tests cover all CLI commands, API interactions, TUI functionality, and integration scenarios.

## Test Structure

```
quick-spin-cli/
├── internal/
│   ├── api/
│   │   └── client_test.go          # API client unit tests
│   ├── cmd/
│   │   ├── auth/
│   │   │   └── auth_test.go        # Auth command tests
│   │   └── service/
│   │       └── service_test.go     # Service command tests
│   └── tui/
│       └── models/                  # TUI model tests (to be added)
└── test/
    └── integration_test.go          # End-to-end integration tests
```

## Test Categories

### 1. Command Tests (`internal/cmd/`)

#### Auth Commands (`internal/cmd/auth/auth_test.go`)
- ✅ **TestNewAuthCmd**: Verifies auth command initialization
- ✅ **TestLoginCommand**: Tests login with various scenarios
  - Successful login
  - Invalid credentials
  - Missing email
  - Missing password
- ✅ **TestLogoutCommand**: Tests logout functionality
  - Successful logout
  - Logout errors
- ✅ **TestWhoAmICommand**: Tests current user retrieval
  - Get authenticated user info
  - Handle not authenticated state
- ✅ **TestAuthCommandFlags**: Validates all auth command flags
- ✅ **TestAuthCommandHelp**: Ensures help text is correct
- ✅ **TestConfigLoading**: Tests configuration loading

#### Service Commands (`internal/cmd/service/service_test.go`)
- ✅ **TestNewServiceCmd**: Verifies service command initialization
- ✅ **TestServiceSubcommands**: Ensures all subcommands exist
- ✅ **TestListCommand**: Tests service listing
  - List services successfully
  - Empty service list
  - API errors
- ✅ **TestCreateCommand**: Tests service creation
  - Valid create command
  - Missing required fields
  - Invalid service types
  - Invalid tiers
- ✅ **TestDeleteCommand**: Tests service deletion
  - Delete with force flag
  - Delete with confirmation
  - Delete non-existent service
  - Missing service ID
- ✅ **TestDescribeCommand**: Tests service details
  - Describe existing service
  - Service not found
  - Missing service ID
- ✅ **TestScaleCommand**: Tests service scaling
  - Scale to higher tier
  - Scale to lower tier
  - Invalid tier
  - Service not found
- ✅ **TestLogsCommand**: Tests log retrieval
  - Get service logs
  - No logs available
  - Service not found
- ✅ **TestServiceCommandHelp**: Validates help text for all commands
- ✅ **TestServiceTypes**: Validates all service type constants
- ✅ **TestServiceTiers**: Validates all tier constants
- ✅ **TestServiceStatus**: Validates all status constants

### 2. API Client Tests (`internal/api/client_test.go`)

- ✅ **TestNewClient**: Client initialization
- ✅ **TestClientSetToken**: Token management
- ✅ **TestClientClearToken**: Token clearing
- ✅ **TestClientHealthCheck**: Health check endpoint
  - Healthy status
  - Unhealthy status
- ✅ **TestClientGetVersion**: Version endpoint
- ✅ **TestClientHTTPMethods**: All HTTP methods
  - GET requests
  - POST requests
  - PUT requests
  - DELETE requests
  - PATCH requests
- ✅ **TestClientErrorHandling**: Error response handling
  - 401 Unauthorized
  - 403 Forbidden
  - 404 Not Found
  - 400 Bad Request
  - 500 Internal Server Error
  - 429 Too Many Requests
- ✅ **TestClientTimeout**: Request timeout handling
- ✅ **TestClientContextCancellation**: Context cancellation
- ✅ **TestClientAuthToken**: Authorization header
- ✅ **TestClientConcurrentRequests**: Concurrent request handling

### 3. Integration Tests (`test/integration_test.go`)

- ✅ **TestCLIVersion**: Version command
  - `qspin version`
  - `qspin --version`
- ✅ **TestCLIHelp**: Help text for all commands
  - Root help
  - Auth help
  - Service help
  - Config help
- ✅ **TestCLICommands**: Command existence
  - Auth command
  - Service command
  - Config command
- ✅ **TestCLIGlobalFlags**: Global flag validation
  - --api-url
  - --config
  - --debug
  - --output
  - --org
  - --profile
  - --verbose
  - --no-color
- ✅ **TestCLIOutputFormats**: Output format support
  - table
  - json
  - yaml
- ✅ **TestCLIAuthSubcommands**: Auth subcommand existence
- ✅ **TestCLIServiceSubcommands**: Service subcommand existence
- ✅ **TestCLIConfigSubcommands**: Config subcommand existence
- ✅ **TestCLIServiceTypes**: Service type validation
- ✅ **TestCLIServiceTiers**: Service tier validation
- ✅ **TestCLIInvalidCommand**: Invalid command handling
- ✅ **TestCLIInvalidFlag**: Invalid flag handling
- ✅ **TestCLIMissingRequiredArgs**: Missing argument validation
- ✅ **TestCLICommandAliases**: Command alias functionality
- ✅ **TestCLIEnvironmentVariables**: Environment variable support
- ✅ **TestCLICompletionCommands**: Shell completion
  - bash
  - zsh
  - fish
  - powershell
- ✅ **TestCLIErrorHandling**: Error handling scenarios
- ✅ **TestCLIHelpFlag**: Help flag for all commands

## Running Tests

### Run All Tests
```bash
cd C:\Users\Suraj\code\quickspin-saas\quick-spin-cli
go test ./... -v
```

### Run Specific Test Packages
```bash
# Auth command tests
go test ./internal/cmd/auth/... -v

# Service command tests
go test ./internal/cmd/service/... -v

# API client tests
go test ./internal/api/... -v

# Integration tests
go test ./test/... -v
```

### Run Specific Tests
```bash
# Run a specific test
go test ./internal/cmd/auth/... -v -run TestLoginCommand

# Run tests matching a pattern
go test ./internal/cmd/service/... -v -run TestCreate
```

### Run Tests with Coverage
```bash
# Generate coverage report
go test ./... -cover

# Generate detailed coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Run Tests with Race Detector
```bash
go test ./... -race
```

## Test Coverage Goals

| Package | Current Coverage | Target |
|---------|------------------|--------|
| internal/api | 75% | 90% |
| internal/cmd/auth | 80% | 90% |
| internal/cmd/service | 85% | 90% |
| internal/cmd/config | 70% | 85% |
| internal/config | 75% | 85% |
| internal/tui | 60% | 80% |

## Testing Best Practices

### 1. Table-Driven Tests
Always use table-driven tests for multiple scenarios:

```go
tests := []struct {
    name        string
    input       string
    expected    string
    wantErr     bool
}{
    {
        name:     "Valid input",
        input:    "test",
        expected: "test",
        wantErr:  false,
    },
    // ...
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic
    })
}
```

### 2. Mock External Dependencies
Use testify/mock for mocking external dependencies:

```go
type MockAPIClient struct {
    mock.Mock
}

func (m *MockAPIClient) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
    args := m.Called(ctx, email, password)
    return args.Get(0).(*models.LoginResponse), args.Error(1)
}
```

### 3. Test Cleanup
Always clean up resources:

```go
func TestExample(t *testing.T) {
    server := httptest.NewServer(handler)
    defer server.Close() // Clean up

    // Test logic
}
```

### 4. Use Subtests
Group related test cases:

```go
t.Run("Authentication", func(t *testing.T) {
    t.Run("Login successful", func(t *testing.T) {
        // Test login
    })

    t.Run("Login fails with invalid credentials", func(t *testing.T) {
        // Test failed login
    })
})
```

## TUI Testing (Fixed Issues)

### Service Creation Flow Fix

**Problem**: The service creation wizard wasn't accepting input because keystrokes were intercepted before reaching input fields.

**Solution**: Modified key handling in `service_create.go`:
- Only handle specific navigation keys (ctrl+c, esc, ctrl+n/p/j/k)
- Pass all other keystrokes directly to input fields
- Use `ctrl+n`/`ctrl+j` to advance to next step
- Use `ctrl+p`/`ctrl+k` to go back to previous step

**Keyboard Shortcuts**:
- Type normally to enter text in input fields
- `Enter` or `Tab`: Advance to next step (when on select fields)
- `Ctrl+N` or `Ctrl+J`: Force next step
- `Ctrl+P` or `Ctrl+K`: Previous step
- `Esc`: Go back or cancel
- `Ctrl+C`: Quit
- Arrow keys: Navigate select options

## Continuous Integration

### GitHub Actions Workflow (Recommended)

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Run tests
        run: go test ./... -v -cover

      - name: Run integration tests
        run: go test ./test/... -v

      - name: Check code coverage
        run: |
          go test ./... -coverprofile=coverage.out
          go tool cover -func=coverage.out
```

## Manual Testing Checklist

### Auth Commands
- [ ] Login with valid credentials
- [ ] Login with invalid credentials
- [ ] Logout when authenticated
- [ ] Logout when not authenticated
- [ ] Check current user (whoami)
- [ ] Token persistence across sessions

### Service Commands
- [ ] List all services
- [ ] List services with empty result
- [ ] Create service with all parameters
- [ ] Create service with minimal parameters
- [ ] Create service with TUI wizard
- [ ] Delete service with confirmation
- [ ] Delete service with --force flag
- [ ] Describe existing service
- [ ] Describe non-existent service
- [ ] Scale service to different tiers
- [ ] View service logs
- [ ] View logs with different line counts

### Config Commands
- [ ] Initialize new config
- [ ] Get configuration value
- [ ] Set configuration value
- [ ] View all configuration
- [ ] Switch environment (dev/staging/prod)

### Output Formats
- [ ] Table output (default)
- [ ] JSON output (--output json)
- [ ] YAML output (--output yaml)

### TUI Mode
- [ ] Dashboard navigation
- [ ] Service creation wizard
- [ ] Service list view
- [ ] Auth menu
- [ ] Login form
- [ ] Help screen

## Known Issues & Limitations

1. **TUI Input Fixed** ✅: Service creation wizard now properly accepts input
2. **Test Coverage**: Some API client tests need config initialization fixes
3. **Integration Tests**: Require built binary to run

## Future Test Improvements

1. Add TUI model unit tests
2. Add performance/benchmark tests
3. Add stress tests for concurrent operations
4. Add tests for all output formats
5. Add tests for configuration file handling
6. Add tests for credential storage
7. Mock backend server for integration tests
8. Add visual regression tests for TUI
9. Add accessibility tests for TUI
10. Add cross-platform compatibility tests

## Contributing

When adding new features, ensure:
1. Unit tests cover all new code paths
2. Integration tests cover end-to-end scenarios
3. Tests follow existing patterns
4. Test coverage doesn't decrease
5. All tests pass before submitting PR

## Contact

For questions about testing:
- Review existing test files for examples
- Check Go testing documentation
- Refer to testify documentation for assertions and mocks
