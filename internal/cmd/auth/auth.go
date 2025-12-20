package auth

import (
	"github.com/spf13/cobra"
)

// NewAuthCmd creates the auth command
func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
		Long:  "Manage authentication for QuickSpin CLI",
	}

	// Add subcommands
	cmd.AddCommand(NewLoginCmd())
	cmd.AddCommand(NewLogoutCmd())
	cmd.AddCommand(NewWhoAmICmd())
	cmd.AddCommand(NewTokenCmd())

	return cmd
}
