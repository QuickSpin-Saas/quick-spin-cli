package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	// Set version info for test
	Version = "1.0.0"
	Commit = "abc123"
	Date = "2025-01-01"

	cmd := NewVersionCmd()
	assert.NotNil(t, cmd)

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Execute command
	err := cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "qspin version 1.0.0")
	assert.Contains(t, output, "commit: abc123")
	assert.Contains(t, output, "built: 2025-01-01")
}

func TestVersionCmd_DevVersion(t *testing.T) {
	// Set dev version info
	Version = "dev"
	Commit = "none"
	Date = "unknown"

	cmd := NewVersionCmd()

	// Capture output
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Execute command
	err := cmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "qspin version dev")
	assert.Contains(t, output, "commit: none")
	assert.Contains(t, output, "built: unknown")
}
