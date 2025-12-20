package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoadCredentials(t *testing.T) {
	// Create a temporary directory for test credentials
	tempDir, err := os.MkdirTemp("", "qspin-creds-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create credentials
	creds := &Credentials{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
	}

	// Save credentials
	err = SaveCredentials(creds)
	assert.NoError(t, err)

	// Verify file was created
	credPath := filepath.Join(tempDir, ".quickspin", "credentials.json")
	_, err = os.Stat(credPath)
	assert.NoError(t, err)

	// Check file permissions (should be 0600)
	info, err := os.Stat(credPath)
	assert.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	// Load credentials
	loadedCreds, err := LoadCredentials()
	assert.NoError(t, err)
	assert.Equal(t, creds.AccessToken, loadedCreds.AccessToken)
	assert.Equal(t, creds.RefreshToken, loadedCreds.RefreshToken)
}

func TestLoadCredentials_NotFound(t *testing.T) {
	// Create a temporary directory with no credentials
	tempDir, err := os.MkdirTemp("", "qspin-creds-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Load credentials (should return empty credentials, no error)
	creds, err := LoadCredentials()
	assert.NoError(t, err)
	assert.Equal(t, "", creds.AccessToken)
	assert.Equal(t, "", creds.RefreshToken)
}

func TestClearCredentials(t *testing.T) {
	// Create a temporary directory for test credentials
	tempDir, err := os.MkdirTemp("", "qspin-creds-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create and save credentials
	creds := &Credentials{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
	}
	err = SaveCredentials(creds)
	require.NoError(t, err)

	// Verify credentials exist
	assert.True(t, CredentialsExist())

	// Clear credentials
	err = ClearCredentials()
	assert.NoError(t, err)

	// Verify credentials were removed
	assert.False(t, CredentialsExist())
}

func TestCredentialsExist(t *testing.T) {
	// Create a temporary directory for test credentials
	tempDir, err := os.MkdirTemp("", "qspin-creds-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Initially should not exist
	assert.False(t, CredentialsExist())

	// Create credentials
	creds := &Credentials{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
	}
	err = SaveCredentials(creds)
	require.NoError(t, err)

	// Now should exist
	assert.True(t, CredentialsExist())
}
