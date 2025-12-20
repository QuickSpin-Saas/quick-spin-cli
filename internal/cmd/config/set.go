package config

import (
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewSetCmd creates the config set command
func NewSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long:  "Set a configuration key-value pair",
		Args:  cobra.ExactArgs(2),
		RunE:  runSet,
	}

	return cmd
}

func runSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Set value
	cfg.Set(key, value)

	// Save config
	if err := cfg.Save(); err != nil {
		output.Error(fmt.Sprintf("Failed to save config: %s", err))
		return err
	}

	output.Success(fmt.Sprintf("Set %s = %s", key, value))
	return nil
}
