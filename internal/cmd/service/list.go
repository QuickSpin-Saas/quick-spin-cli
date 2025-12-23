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

// NewListCmd creates the service list command
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all services",
		Long:    "List all managed services for your account",
		RunE:    runList,
	}

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	// Check if we should use TUI mode
	outputFormat := viper.GetString("defaults.output")
	if outputpkg.ShouldUseTUI(outputFormat) {
		// Launch TUI service list
		return tui.LaunchView(tui.ViewServiceList)
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
	spinner := outputpkg.NewSpinner("Loading services...")
	spinner.Start()

	// List services
	services, err := client.ListServices(ctx)
	spinner.Stop()

	if err != nil {
		outputpkg.Error(fmt.Sprintf("Failed to list services: %s", err))
		return err
	}

	if len(services) == 0 {
		outputpkg.Info("No services found")
		return nil
	}

	// Display services
	outputpkg.Success(fmt.Sprintf("Found %d service(s)", len(services)))
	fmt.Println()

	formatType := outputpkg.Format(viper.GetString("defaults.output"))
	return outputpkg.Print(formatType, services)
}
