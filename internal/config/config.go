package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the CLI configuration
type Config struct {
	v *viper.Viper
}

// New creates a new Config instance
func New() *Config {
	return &Config{
		v: viper.GetViper(),
	}
}

// LoadConfig loads the configuration from file
func LoadConfig() (*Config, error) {
	cfg := New()

	// Try to read config file
	if err := cfg.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found is ok, we'll use defaults
	}

	return cfg, nil
}

// GetAPIURL returns the API URL
func (c *Config) GetAPIURL() string {
	// Check environment variable first
	if url := os.Getenv("QUICKSPIN_API_URL"); url != "" {
		return url
	}

	// Check for environment-based URL
	env := c.GetEnvironment()
	envURL := c.v.GetString(fmt.Sprintf("api.environments.%s", env))
	if envURL != "" {
		return envURL
	}

	// Fall back to default URL
	url := c.v.GetString("api.url")
	if url == "" {
		// Default to production if not set
		return "https://api.quickspin.cloud"
	}
	return url
}

// GetEnvironment returns the current environment (dev, staging, prod)
func (c *Config) GetEnvironment() string {
	// Check environment variable first
	if env := os.Getenv("QUICKSPIN_ENV"); env != "" {
		return env
	}
	env := c.v.GetString("api.environment")
	if env == "" {
		return "prod" // Default to production
	}
	return env
}

// SetEnvironment sets the current environment
func (c *Config) SetEnvironment(env string) {
	c.v.Set("api.environment", env)
}

// GetAPITimeout returns the API timeout
func (c *Config) GetAPITimeout() string {
	return c.v.GetString("api.timeout")
}

// GetDefaultOrganization returns the default organization
func (c *Config) GetDefaultOrganization() string {
	if org := os.Getenv("QUICKSPIN_ORG"); org != "" {
		return org
	}
	return c.v.GetString("defaults.organization")
}

// GetDefaultRegion returns the default region
func (c *Config) GetDefaultRegion() string {
	return c.v.GetString("defaults.region")
}

// GetDefaultOutput returns the default output format
func (c *Config) GetDefaultOutput() string {
	return c.v.GetString("defaults.output")
}

// GetDefaultServiceType returns the default service type
func (c *Config) GetDefaultServiceType() string {
	return c.v.GetString("defaults.service_type")
}

// GetDefaultTier returns the default tier
func (c *Config) GetDefaultTier() string {
	return c.v.GetString("defaults.tier")
}

// GetToken returns the stored authentication token
func (c *Config) GetToken() (string, error) {
	// Check environment variable first
	if token := os.Getenv("QUICKSPIN_TOKEN"); token != "" {
		return token, nil
	}

	// Try to load from credential store
	creds, err := LoadCredentials()
	if err != nil {
		return "", err
	}

	return creds.AccessToken, nil
}

// GetRefreshToken returns the stored refresh token
func (c *Config) GetRefreshToken() (string, error) {
	// Check environment variable first
	if token := os.Getenv("QUICKSPIN_REFRESH_TOKEN"); token != "" {
		return token, nil
	}

	// Try to load from credential store
	creds, err := LoadCredentials()
	if err != nil {
		return "", err
	}

	return creds.RefreshToken, nil
}

// GetCredentials returns both access and refresh tokens
func (c *Config) GetCredentials() (accessToken, refreshToken string, err error) {
	// Try to load from credential store
	creds, err := LoadCredentials()
	if err != nil {
		return "", "", err
	}

	// Check environment variables (they override stored credentials)
	if token := os.Getenv("QUICKSPIN_TOKEN"); token != "" {
		creds.AccessToken = token
	}
	if token := os.Getenv("QUICKSPIN_REFRESH_TOKEN"); token != "" {
		creds.RefreshToken = token
	}

	return creds.AccessToken, creds.RefreshToken, nil
}

// SaveToken saves the authentication token
func (c *Config) SaveToken(accessToken, refreshToken string) error {
	creds := &Credentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return SaveCredentials(creds)
}

// ClearToken clears the stored authentication token
func (c *Config) ClearToken() error {
	return ClearCredentials()
}

// Set sets a configuration value
func (c *Config) Set(key string, value interface{}) {
	c.v.Set(key, value)
}

// Get gets a configuration value
func (c *Config) Get(key string) interface{} {
	return c.v.Get(key)
}

// GetString gets a string configuration value
func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

// GetBool gets a boolean configuration value
func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configDir := c.GetConfigDir()
	configFile := filepath.Join(configDir, "config.yaml")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	if err := c.v.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigDir returns the configuration directory
func (c *Config) GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".quickspin"
	}
	return filepath.Join(home, ".quickspin")
}

// GetConfigFile returns the configuration file path
func (c *Config) GetConfigFile() string {
	return filepath.Join(c.GetConfigDir(), "config.yaml")
}

// InitConfig initializes a new configuration file
func InitConfig() error {
	cfg := New()

	// Check if config already exists
	configFile := cfg.GetConfigFile()
	if _, err := os.Stat(configFile); err == nil {
		return fmt.Errorf("configuration already exists at %s", configFile)
	}

	// Create config directory
	if err := os.MkdirAll(cfg.GetConfigDir(), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set defaults
	cfg.Set("api.url", "https://api.quickspin.cloud")
	cfg.Set("api.environment", "prod")
	cfg.Set("api.environments.dev", "http://localhost:8000")
	cfg.Set("api.environments.staging", "https://staging-api.quickspin.cloud")
	cfg.Set("api.environments.prod", "https://api.quickspin.cloud")
	cfg.Set("api.timeout", "30s")
	cfg.Set("auth.method", "jwt")
	cfg.Set("defaults.region", "us-east-1")
	cfg.Set("defaults.output", "table")
	cfg.Set("defaults.service_type", "redis")
	cfg.Set("defaults.tier", "developer")
	cfg.Set("telemetry.enabled", true)
	cfg.Set("telemetry.anonymous", true)

	// Save configuration
	if err := cfg.Save(); err != nil {
		return err
	}

	return nil
}
