package auth

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	showToken    bool
	refreshToken bool
)

// NewTokenCmd creates the token command
func NewTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage authentication tokens",
		Long:  "Show or refresh authentication tokens",
		RunE:  runToken,
	}

	cmd.Flags().BoolVar(&showToken, "show", false, "Show current access token")
	cmd.Flags().BoolVar(&refreshToken, "refresh", false, "Refresh access token")

	return cmd
}

func runToken(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if refreshToken {
		// Refresh the token
		client := api.NewClient(cfg)

		spinner := output.NewSpinner("Refreshing token...")
		spinner.Start()

		tokens, err := client.RefreshToken(ctx)
		spinner.Stop()

		if err != nil {
			output.Error(fmt.Sprintf("Failed to refresh token: %s", err))
			return err
		}

		output.Success("Token refreshed successfully")
		if showToken {
			fmt.Println("\nAccess Token:", tokens.AccessToken)
		}
	} else if showToken {
		// Show current token
		token, err := cfg.GetToken()
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}

		if token == "" {
			output.Warning("No authentication token found. Please login first.")
			return nil
		}

		fmt.Println("Access Token:", token)
	} else {
		return fmt.Errorf("please specify --show or --refresh flag")
	}

	return nil
}
