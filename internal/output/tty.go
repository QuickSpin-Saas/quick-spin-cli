package output

import (
	"os"

	"golang.org/x/term"
)

// IsInteractive checks if the CLI is running in an interactive terminal
func IsInteractive() bool {
	// Check if stdout is a terminal
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	// Check if stdin is a terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return false
	}

	// Check for environment variables that disable interactive mode
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	if os.Getenv("QUICKSPIN_NO_TUI") != "" {
		return false
	}

	// Check if running in CI/CD environment
	if os.Getenv("CI") != "" {
		return false
	}

	return true
}

// ShouldUseTUI determines whether to use the Terminal User Interface
// based on the output format and terminal capabilities
func ShouldUseTUI(outputFormat string) bool {
	// Force CLI mode for JSON/YAML output (for scripting/automation)
	if outputFormat == "json" || outputFormat == "yaml" {
		return false
	}

	// Use TUI only if running in an interactive terminal
	return IsInteractive()
}

// SupportsColor checks if the terminal supports color output
func SupportsColor() bool {
	// NO_COLOR environment variable takes precedence
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if stdout is a terminal
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	// Check TERM environment variable
	termEnv := os.Getenv("TERM")
	if termEnv == "dumb" {
		return false
	}

	return true
}

// GetTerminalSize returns the width and height of the terminal
func GetTerminalSize() (width, height int, err error) {
	width, height, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Return default size if we can't detect terminal size
		return 80, 24, err
	}
	return width, height, nil
}
