package config

import (
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewInitCmd creates the config init command
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration",
		Long:  "Create a new QuickSpin CLI configuration file with default values",
		RunE:  runInit,
	}

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	if err := config.InitConfig(); err != nil {
		output.Error(fmt.Sprintf("Failed to initialize config: %s", err))
		return err
	}

	cfg := config.New()
	output.Success(fmt.Sprintf("Configuration initialized at %s", cfg.GetConfigFile()))
	return nil
}
