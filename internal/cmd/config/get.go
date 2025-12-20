package config

import (
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/spf13/cobra"
)

// NewGetCmd creates the config get command
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Long:  "Retrieve a configuration value by key",
		Args:  cobra.ExactArgs(1),
		RunE:  runGet,
	}

	return cmd
}

func runGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get value
	value := cfg.Get(key)
	fmt.Println(value)

	return nil
}
