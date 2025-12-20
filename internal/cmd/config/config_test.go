package config

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigCmd(t *testing.T) {
	cmd := NewConfigCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "config", cmd.Use)
	assert.Greater(t, len(cmd.Commands()), 0)
}

func TestInitCmd(t *testing.T) {
	// Create a temporary directory for test config
	tempDir, err := os.MkdirTemp("", "qspin-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	cmd := NewInitCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "init", cmd.Use)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Execute command
	err = cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Configuration initialized")
}

func TestSetCmd(t *testing.T) {
	// Create a temporary directory for test config
	tempDir, err := os.MkdirTemp("", "qspin-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Initialize config first
	initCmd := NewInitCmd()
	err = initCmd.Execute()
	require.NoError(t, err)

	// Test set command
	cmd := NewSetCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "set <key> <value>", cmd.Use)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Execute command with args
	cmd.SetArgs([]string{"test.key", "test-value"})
	err = cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Set test.key = test-value")
}

func TestGetCmd(t *testing.T) {
	// Create a temporary directory for test config
	tempDir, err := os.MkdirTemp("", "qspin-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Initialize config first
	initCmd := NewInitCmd()
	err = initCmd.Execute()
	require.NoError(t, err)

	// Set a value
	setCmd := NewSetCmd()
	setCmd.SetArgs([]string{"test.key", "test-value"})
	err = setCmd.Execute()
	require.NoError(t, err)

	// Test get command
	cmd := NewGetCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "get <key>", cmd.Use)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Execute command with args
	cmd.SetArgs([]string{"test.key"})
	err = cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "test-value")
}

func TestViewCmd(t *testing.T) {
	// Create a temporary directory for test config
	tempDir, err := os.MkdirTemp("", "qspin-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Initialize config first
	initCmd := NewInitCmd()
	err = initCmd.Execute()
	require.NoError(t, err)

	// Test view command
	cmd := NewViewCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "view", cmd.Use)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Execute command
	err = cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Configuration file:")
}
