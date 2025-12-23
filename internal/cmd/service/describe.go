package service

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	outputpkg "github.com/quickspin/quickspin-cli/internal/output"
	"github.com/quickspin/quickspin-cli/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewDescribeCmd creates the service describe command
func NewDescribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe SERVICE_ID",
		Aliases: []string{"get", "show"},
		Short:   "Get service details",
		Long:    "Get detailed information about a specific service",
		Args:    cobra.ExactArgs(1),
		RunE:    runDescribe,
	}

	return cmd
}

func runDescribe(cmd *cobra.Command, args []string) error {
	serviceID := args[0]

	// Check if we should use TUI mode
	outputFormat := viper.GetString("defaults.output")
	if outputpkg.ShouldUseTUI(outputFormat) {
		// Launch TUI service detail view
		// Note: We'd need to pass the serviceID to the view
		return tui.LaunchView(tui.ViewServiceDetail)
	}

	// Traditional CLI mode
	ctx := context.Background()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client := api.NewClient(cfg)

	// Show spinner
	spinner := outputpkg.NewSpinner("Loading service details...")
	spinner.Start()

	// Get service
	service, err := client.GetService(ctx, serviceID)
	spinner.Stop()

	if err != nil {
		outputpkg.Error(fmt.Sprintf("Failed to get service: %s", err))
		return err
	}

	// Display service details
	fmt.Println()
	formatType := outputpkg.Format(viper.GetString("defaults.output"))
	return outputpkg.Print(formatType, service)
}
