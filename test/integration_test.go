package test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var binaryPath string

// TestMain builds the binary before running tests
func TestMain(m *testing.M) {
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "qspin-test", "../cmd/qspin")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to build binary: %v\n", err)
		os.Exit(1)
	}

	binaryPath, _ = filepath.Abs("qspin-test")
	defer os.Remove(binaryPath)

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// runCommand executes the CLI binary with given args
func runCommand(t *testing.T, args ...string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func TestCLIVersion(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Version command",
			args:     []string{"version"},
			expected: "qspin version",
		},
		{
			name:     "Version flag",
			args:     []string{"--version"},
			expected: "qspin version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, err := runCommand(t, tt.args...)
			assert.NoError(t, err)
			assert.Contains(t, stdout, tt.expected)
		})
	}
}

func TestCLIHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "Root help",
			args:     []string{"--help"},
			expected: []string{"QuickSpin CLI", "Usage:", "Available Commands:"},
		},
		{
			name:     "Auth help",
			args:     []string{"auth", "--help"},
			expected: []string{"auth", "Authentication", "login", "logout", "whoami"},
		},
		{
			name:     "Service help",
			args:     []string{"service", "--help"},
			expected: []string{"service", "list", "create", "delete"},
		},
		{
			name:     "Config help",
			args:     []string{"config", "--help"},
			expected: []string{"config", "Configuration", "init", "get", "set"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, err := runCommand(t, tt.args...)
			assert.NoError(t, err)

			for _, exp := range tt.expected {
				assert.Contains(t, stdout, exp, "Expected output to contain: %s", exp)
			}
		})
	}
}

func TestCLICommands(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedInHelp string
	}{
		{
			name:           "Auth command exists",
			args:           []string{"auth"},
			expectedInHelp: "Authentication commands",
		},
		{
			name:           "Service command exists",
			args:           []string{"service"},
			expectedInHelp: "Manage microservices",
		},
		{
			name:           "Config command exists",
			args:           []string{"config"},
			expectedInHelp: "Configuration management",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, _, _ := runCommand(t, append(tt.args, "--help")...)
			assert.Contains(t, stdout, tt.expectedInHelp)
		})
	}
}

func TestCLIGlobalFlags(t *testing.T) {
	flags := []string{
		"--api-url",
		"--config",
		"--debug",
		"--output",
		"--org",
		"--profile",
		"--verbose",
		"--no-color",
	}

	stdout, _, err := runCommand(t, "--help")
	require.NoError(t, err)

	for _, flag := range flags {
		assert.Contains(t, stdout, flag, "Expected global flag %s to exist", flag)
	}
}

func TestCLIOutputFormats(t *testing.T) {
	formats := []string{"table", "json", "yaml"}

	stdout, _, err := runCommand(t, "--help")
	require.NoError(t, err)

	for _, format := range formats {
		assert.Contains(t, stdout, format, "Expected output format %s to be available", format)
	}
}

func TestCLIAuthSubcommands(t *testing.T) {
	subcommands := []string{"login", "logout", "whoami"}

	stdout, _, err := runCommand(t, "auth", "--help")
	require.NoError(t, err)

	for _, cmd := range subcommands {
		assert.Contains(t, stdout, cmd, "Expected auth subcommand %s to exist", cmd)
	}
}

func TestCLIServiceSubcommands(t *testing.T) {
	subcommands := []string{"list", "create", "delete", "describe", "scale", "logs"}

	stdout, _, err := runCommand(t, "service", "--help")
	require.NoError(t, err)

	for _, cmd := range subcommands {
		assert.Contains(t, stdout, cmd, "Expected service subcommand %s to exist", cmd)
	}
}

func TestCLIConfigSubcommands(t *testing.T) {
	subcommands := []string{"init", "get", "set", "view"}

	stdout, _, err := runCommand(t, "config", "--help")
	require.NoError(t, err)

	for _, cmd := range subcommands {
		assert.Contains(t, stdout, cmd, "Expected config subcommand %s to exist", cmd)
	}
}

func TestCLIServiceTypes(t *testing.T) {
	serviceTypes := []string{"redis", "rabbitmq", "postgresql", "mongodb", "mysql", "elasticsearch"}

	stdout, _, err := runCommand(t, "service", "create", "--help")
	require.NoError(t, err)

	for _, serviceType := range serviceTypes {
		assert.Contains(t, strings.ToLower(stdout), serviceType,
			"Expected service type %s to be mentioned", serviceType)
	}
}

func TestCLIServiceTiers(t *testing.T) {
	tiers := []string{"developer", "basic", "standard", "pro", "premium", "enterprise"}

	stdout, _, err := runCommand(t, "service", "create", "--help")
	require.NoError(t, err)

	lowerStdout := strings.ToLower(stdout)
	for _, tier := range tiers {
		assert.Contains(t, lowerStdout, tier,
			"Expected service tier %s to be mentioned", tier)
	}
}

func TestCLIInvalidCommand(t *testing.T) {
	_, stderr, err := runCommand(t, "invalid-command")
	assert.Error(t, err)
	assert.Contains(t, stderr, "unknown command")
}

func TestCLIInvalidFlag(t *testing.T) {
	_, stderr, err := runCommand(t, "--invalid-flag")
	assert.Error(t, err)
	assert.Contains(t, stderr, "unknown flag")
}

func TestCLIMissingRequiredArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedErrMsg string
	}{
		{
			name:           "Service describe missing ID",
			args:           []string{"service", "describe"},
			expectedErrMsg: "requires",
		},
		{
			name:           "Service delete missing ID",
			args:           []string{"service", "delete"},
			expectedErrMsg: "requires",
		},
		{
			name:           "Service scale missing ID",
			args:           []string{"service", "scale"},
			expectedErrMsg: "requires",
		},
		{
			name:           "Service logs missing ID",
			args:           []string{"service", "logs"},
			expectedErrMsg: "requires",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, stderr, err := runCommand(t, tt.args...)
			assert.Error(t, err)
			assert.Contains(t, stderr, tt.expectedErrMsg)
		})
	}
}

func TestCLICommandAliases(t *testing.T) {
	tests := []struct {
		name         string
		command      []string
		alias        []string
		shouldMatch  bool
	}{
		{
			name:        "Service alias 'svc'",
			command:     []string{"service", "--help"},
			alias:       []string{"svc", "--help"},
			shouldMatch: true,
		},
		{
			name:        "Service alias 'services'",
			command:     []string{"service", "--help"},
			alias:       []string{"services", "--help"},
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdOut, _, err1 := runCommand(t, tt.command...)
			aliasOut, _, err2 := runCommand(t, tt.alias...)

			if tt.shouldMatch {
				assert.NoError(t, err1)
				assert.NoError(t, err2)
				assert.Equal(t, cmdOut, aliasOut, "Command and alias should produce same output")
			}
		})
	}
}

func TestCLIEnvironmentVariables(t *testing.T) {
	t.Run("API URL from environment", func(t *testing.T) {
		// This test would require more setup to actually verify the env var is used
		// For now, just verify the flag exists
		stdout, _, err := runCommand(t, "--help")
		assert.NoError(t, err)
		assert.Contains(t, stdout, "api-url")
	})

	t.Run("Organization from environment", func(t *testing.T) {
		stdout, _, err := runCommand(t, "--help")
		assert.NoError(t, err)
		assert.Contains(t, stdout, "--org")
	})
}

func TestCLICompletionCommands(t *testing.T) {
	shells := []string{"bash", "zsh", "fish", "powershell"}

	for _, shell := range shells {
		t.Run(fmt.Sprintf("Completion for %s", shell), func(t *testing.T) {
			stdout, _, err := runCommand(t, "completion", shell, "--help")
			assert.NoError(t, err)
			assert.Contains(t, stdout, shell)
		})
	}
}

func TestCLIErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "Invalid service create - missing flags",
			args:        []string{"service", "create"},
			expectError: true,
		},
		{
			name:        "Invalid service delete - missing ID",
			args:        []string{"service", "delete"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := runCommand(t, tt.args...)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCLIHelpFlag(t *testing.T) {
	commands := [][]string{
		{"--help"},
		{"-h"},
		{"auth", "--help"},
		{"service", "--help"},
		{"config", "--help"},
		{"version", "--help"},
	}

	for _, cmd := range commands {
		t.Run(strings.Join(cmd, " "), func(t *testing.T) {
			stdout, _, err := runCommand(t, cmd...)
			assert.NoError(t, err)
			assert.NotEmpty(t, stdout)
			assert.Contains(t, strings.ToLower(stdout), "usage")
		})
	}
}
