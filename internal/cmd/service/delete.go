package service

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	outputpkg "github.com/quickspin/quickspin-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

// NewDeleteCmd creates the service delete command
func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete SERVICE_ID",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete a service",
		Long:    "Delete a managed microservice",
		Args:    cobra.ExactArgs(1),
		RunE:    runDelete,
	}

	cmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func runDelete(cmd *cobra.Command, args []string) error {
	serviceID := args[0]

	ctx := context.Background()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client := api.NewClient(cfg)

	// Get service details first
	service, err := client.GetService(ctx, serviceID)
	if err != nil {
		outputpkg.Error(fmt.Sprintf("Failed to get service: %s", err))
		return err
	}

	// Confirm deletion unless --force flag is used
	if !deleteForce {
		fmt.Printf("Are you sure you want to delete service '%s' (%s)? This action cannot be undone.\n", service.Name, serviceID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirmation string
		fmt.Scanln(&confirmation)

		if confirmation != "yes" {
			outputpkg.Info("Deletion cancelled")
			return nil
		}
	}

	// Show spinner
	spinner := outputpkg.NewSpinner(fmt.Sprintf("Deleting service '%s'...", service.Name))
	spinner.Start()

	// Delete service
	err = client.DeleteService(ctx, serviceID)
	spinner.Stop()

	if err != nil {
		outputpkg.Error(fmt.Sprintf("Failed to delete service: %s", err))
		return err
	}

	// Success message
	outputpkg.Success(fmt.Sprintf("Successfully deleted service '%s'", service.Name))

	return nil
}
