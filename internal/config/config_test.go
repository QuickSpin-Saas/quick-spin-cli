package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := New()
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.v)
}

func TestGetAPIURL(t *testing.T) {
	cfg := New()

	// Test default value
	url := cfg.GetAPIURL()
	assert.Equal(t, "https://api.quickspin.cloud", url)

	// Test environment variable override
	os.Setenv("QUICKSPIN_API_URL", "http://localhost:8000")
	defer os.Unsetenv("QUICKSPIN_API_URL")

	url = cfg.GetAPIURL()
	assert.Equal(t, "http://localhost:8000", url)
}

func TestSetAndGet(t *testing.T) {
	cfg := New()

	// Set a value
	cfg.Set("test.key", "test-value")

	// Get the value
	value := cfg.GetString("test.key")
	assert.Equal(t, "test-value", value)
}

func TestInitConfig(t *testing.T) {
	// Create a temporary directory for test config
	tempDir, err := os.MkdirTemp("", "qspin-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Initialize config
	err = InitConfig()
	assert.NoError(t, err)

	// Verify config file was created
	configFile := filepath.Join(tempDir, ".quickspin", "config.yaml")
	_, err = os.Stat(configFile)
	assert.NoError(t, err)

	// Try to initialize again (should fail)
	err = InitConfig()
	assert.Error(t, err)
}

func TestGetConfigDir(t *testing.T) {
	cfg := New()
	configDir := cfg.GetConfigDir()
	assert.Contains(t, configDir, ".quickspin")
}

func TestGetConfigFile(t *testing.T) {
	cfg := New()
	configFile := cfg.GetConfigFile()
	assert.Contains(t, configFile, ".quickspin")
	assert.Contains(t, configFile, "config.yaml")
}

func TestGetDefaultValues(t *testing.T) {
	cfg := New()

	tests := []struct {
		name     string
		getter   func() string
		expected string
	}{
		{"DefaultRegion", cfg.GetDefaultRegion, "us-east-1"},
		{"DefaultOutput", cfg.GetDefaultOutput, "table"},
		{"DefaultServiceType", cfg.GetDefaultServiceType, "redis"},
		{"DefaultTier", cfg.GetDefaultTier, "developer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := tt.getter()
			assert.Equal(t, tt.expected, value)
		})
	}
}
