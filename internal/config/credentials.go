package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Credentials represents stored authentication credentials
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GetCredentialsPath returns the path to the credentials file
func GetCredentialsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".quickspin/credentials.json"
	}
	return filepath.Join(home, ".quickspin", "credentials.json")
}

// LoadCredentials loads credentials from storage
func LoadCredentials() (*Credentials, error) {
	credPath := GetCredentialsPath()

	data, err := os.ReadFile(credPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Credentials{}, nil
		}
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

// SaveCredentials saves credentials to storage
func SaveCredentials(creds *Credentials) error {
	credPath := GetCredentialsPath()
	credDir := filepath.Dir(credPath)

	// Ensure credentials directory exists
	if err := os.MkdirAll(credDir, 0700); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}

	// Marshal credentials
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	// Write credentials file with restricted permissions
	if err := os.WriteFile(credPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

// ClearCredentials clears stored credentials
func ClearCredentials() error {
	credPath := GetCredentialsPath()

	if err := os.Remove(credPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}

	return nil
}

// CredentialsExist checks if credentials are stored
func CredentialsExist() bool {
	credPath := GetCredentialsPath()
	_, err := os.Stat(credPath)
	return err == nil
}
