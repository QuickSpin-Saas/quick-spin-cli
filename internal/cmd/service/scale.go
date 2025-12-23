package service

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/models"
	outputpkg "github.com/quickspin/quickspin-cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewScaleCmd creates the service scale command
func NewScaleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scale SERVICE_ID TIER",
		Short: "Scale a service to a different tier",
		Long:  "Scale a service to a different tier (starter, developer, basic, standard, pro, premium, enterprise)",
		Args:  cobra.ExactArgs(2),
		RunE:  runScale,
	}

	return cmd
}

func runScale(cmd *cobra.Command, args []string) error {
	serviceID := args[0]
	tier := models.ServiceTier(args[1])

	ctx := context.Background()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client := api.NewClient(cfg)

	// Show spinner
	spinner := outputpkg.NewSpinner(fmt.Sprintf("Scaling service to %s tier...", tier))
	spinner.Start()

	// Scale service
	service, err := client.ScaleService(ctx, serviceID, tier)
	spinner.Stop()

	if err != nil {
		outputpkg.Error(fmt.Sprintf("Failed to scale service: %s", err))
		return err
	}

	// Success message
	outputpkg.Success(fmt.Sprintf("Successfully scaled service '%s' to %s tier", service.Name, tier))
	fmt.Println()

	formatType := outputpkg.Format(viper.GetString("defaults.output"))
	return outputpkg.Print(formatType, service)
}
