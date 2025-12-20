package auth

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewWhoAmICmd creates the whoami command
func NewWhoAmICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Display current user information",
		Long:  "Show details about the currently authenticated user",
		RunE:  runWhoAmI,
	}

	return cmd
}

func runWhoAmI(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client := api.NewClient(cfg)

	// Get current user
	user, err := client.WhoAmI(ctx)
	if err != nil {
		output.Error(fmt.Sprintf("Failed to get user info: %s", err))
		return err
	}

	// Display user info
	formatType := output.Format(viper.GetString("defaults.output"))
	return output.Print(formatType, user)
}
