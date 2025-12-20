package config

import (
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
		Long:  "Manage QuickSpin CLI configuration",
	}

	// Add subcommands
	cmd.AddCommand(NewInitCmd())
	cmd.AddCommand(NewSetCmd())
	cmd.AddCommand(NewGetCmd())
	cmd.AddCommand(NewViewCmd())

	return cmd
}
