package output

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Format represents the output format type
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// Formatter defines the interface for output formatting
type Formatter interface {
	// Format formats any data structure
	Format(data interface{}) error
	// FormatList formats a list with headers
	FormatList(data interface{}, headers []string) error
	// FormatError formats an error message
	FormatError(err error) error
}

// NewFormatter creates a new formatter based on the format type
func NewFormatter(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatYAML:
		return &YAMLFormatter{}
	case FormatTable:
		fallthrough
	default:
		return &TableFormatter{}
	}
}

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// Format formats data as JSON
func (f *JSONFormatter) Format(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// FormatList formats a list as JSON
func (f *JSONFormatter) FormatList(data interface{}, headers []string) error {
	return f.Format(data)
}

// FormatError formats an error as JSON
func (f *JSONFormatter) FormatError(err error) error {
	errorData := map[string]string{
		"error": err.Error(),
	}
	return f.Format(errorData)
}

// YAMLFormatter formats output as YAML
type YAMLFormatter struct{}

// Format formats data as YAML
func (f *YAMLFormatter) Format(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	defer encoder.Close()
	return encoder.Encode(data)
}

// FormatList formats a list as YAML
func (f *YAMLFormatter) FormatList(data interface{}, headers []string) error {
	return f.Format(data)
}

// FormatError formats an error as YAML
func (f *YAMLFormatter) FormatError(err error) error {
	errorData := map[string]string{
		"error": err.Error(),
	}
	return f.Format(errorData)
}

// Print prints formatted output
func Print(format Format, data interface{}) error {
	formatter := NewFormatter(format)
	return formatter.Format(data)
}

// PrintList prints a formatted list
func PrintList(format Format, data interface{}, headers []string) error {
	formatter := NewFormatter(format)
	return formatter.FormatList(data, headers)
}

// PrintError prints a formatted error
func PrintError(format Format, err error) error {
	formatter := NewFormatter(format)
	return formatter.FormatError(err)
}

// Success prints a success message
func Success(message string) {
	fmt.Fprintf(os.Stdout, "✓ %s\n", message)
}

// Info prints an info message
func Info(message string) {
	fmt.Fprintf(os.Stdout, "%s\n", message)
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Fprintf(os.Stderr, "⚠ %s\n", message)
}

// Error prints an error message
func Error(message string) {
	fmt.Fprintf(os.Stderr, "✗ %s\n", message)
}
