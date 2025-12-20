package output

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		expected interface{}
	}{
		{"JSON", FormatJSON, &JSONFormatter{}},
		{"YAML", FormatYAML, &YAMLFormatter{}},
		{"Table", FormatTable, &TableFormatter{}},
		{"Default", "invalid", &TableFormatter{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(tt.format)
			assert.IsType(t, tt.expected, formatter)
		})
	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestJSONFormatter_Format(t *testing.T) {
	formatter := &JSONFormatter{}
	data := testStruct{Name: "test", Value: 42}

	output := captureOutput(func() {
		err := formatter.Format(data)
		assert.NoError(t, err)
	})

	assert.Contains(t, output, "\"name\"")
	assert.Contains(t, output, "\"test\"")
	assert.Contains(t, output, "\"value\"")
	assert.Contains(t, output, "42")
}

func TestYAMLFormatter_Format(t *testing.T) {
	formatter := &YAMLFormatter{}
	data := testStruct{Name: "test", Value: 42}

	output := captureOutput(func() {
		err := formatter.Format(data)
		assert.NoError(t, err)
	})

	assert.Contains(t, output, "name:")
	assert.Contains(t, output, "test")
	assert.Contains(t, output, "value:")
	assert.Contains(t, output, "42")
}

func TestTableFormatter_Format(t *testing.T) {
	formatter := &TableFormatter{}
	data := testStruct{Name: "test", Value: 42}

	output := captureOutput(func() {
		err := formatter.Format(data)
		assert.NoError(t, err)
	})

	assert.Contains(t, output, "name:")
	assert.Contains(t, output, "test")
	assert.Contains(t, output, "value:")
	assert.Contains(t, output, "42")
}

func TestPrint(t *testing.T) {
	data := testStruct{Name: "test", Value: 42}

	// Test JSON output
	output := captureOutput(func() {
		err := Print(FormatJSON, data)
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "\"name\"")

	// Test YAML output
	output = captureOutput(func() {
		err := Print(FormatYAML, data)
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "name:")

	// Test Table output
	output = captureOutput(func() {
		err := Print(FormatTable, data)
		assert.NoError(t, err)
	})
	assert.Contains(t, output, "name:")
}

func TestSuccess(t *testing.T) {
	output := captureOutput(func() {
		Success("Test success message")
	})
	assert.Contains(t, output, "Test success message")
}

func TestInfo(t *testing.T) {
	output := captureOutput(func() {
		Info("Test info message")
	})
	assert.Equal(t, "Test info message\n", output)
}
