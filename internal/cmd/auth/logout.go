package auth

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewLogoutCmd creates the logout command
func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from QuickSpin",
		Long:  "Clear authentication credentials and logout",
		RunE:  runLogout,
	}

	return cmd
}

func runLogout(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client := api.NewClient(cfg)

	// Perform logout
	if err := client.Logout(ctx); err != nil {
		output.Warning(fmt.Sprintf("Logout request failed: %s", err))
		// Continue anyway to clear local credentials
	}

	output.Success("Successfully logged out")
	return nil
}
