package main

import (
	"os"

	"github.com/quickspin/quickspin-cli/internal/cmd"
)

var (
	// Version is set during build
	Version = "dev"
	// Commit is set during build
	Commit = "none"
	// Date is set during build
	Date = "unknown"
)

func main() {
	// Set version info
	cmd.Version = Version
	cmd.Commit = Commit
	cmd.Date = Date

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
