package config

import (
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewViewCmd creates the config view command
func NewViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View configuration",
		Long:  "Display all configuration settings",
		RunE:  runView,
	}

	return cmd
}

func runView(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Configuration file:", cfg.GetConfigFile())
	fmt.Println()

	// Get all settings
	settings := viper.AllSettings()

	// Display configuration
	formatType := output.Format(viper.GetString("defaults.output"))
	return output.Print(formatType, settings)
}
