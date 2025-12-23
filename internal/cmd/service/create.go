package service

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/models"
	outputpkg "github.com/quickspin/quickspin-cli/internal/output"
	"github.com/quickspin/quickspin-cli/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	createName        string
	createType        string
	createTier        string
	createRegion      string
	createDescription string
)

// NewCreateCmd creates the service create command
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new service",
		Long:  "Create a new managed microservice",
		RunE:  runCreate,
	}

	cmd.Flags().StringVar(&createName, "name", "", "Service name (required)")
	cmd.Flags().StringVar(&createType, "type", "", "Service type: redis, rabbitmq, postgresql, mongodb, mysql, elasticsearch (required)")
	cmd.Flags().StringVar(&createTier, "tier", "developer", "Service tier: starter, developer, basic, standard, pro, premium, enterprise")
	cmd.Flags().StringVar(&createRegion, "region", "", "Region (default: from config)")
	cmd.Flags().StringVar(&createDescription, "description", "", "Service description")

	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Check if we should use TUI mode
	outputFormat := viper.GetString("defaults.output")
	if outputpkg.ShouldUseTUI(outputFormat) && createName == "" && createType == "" {
		// Launch TUI service create wizard
		return tui.LaunchView(tui.ViewServiceCreate)
	}

	// Traditional CLI mode - require name and type
	if createName == "" {
		return fmt.Errorf("service name is required (use --name flag)")
	}
	if createType == "" {
		return fmt.Errorf("service type is required (use --type flag)")
	}

	ctx := context.Background()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client := api.NewClient(cfg)

	// Use region from config if not provided
	region := createRegion
	if region == "" {
		region = viper.GetString("defaults.region")
	}

	// Create service request
	req := api.CreateServiceRequest{
		Name:        createName,
		Type:        models.ServiceType(createType),
		Tier:        models.ServiceTier(createTier),
		Region:      region,
		Description: createDescription,
	}

	// Show spinner
	spinner := outputpkg.NewSpinner(fmt.Sprintf("Creating %s service '%s'...", createType, createName))
	spinner.Start()

	// Create service
	service, err := client.CreateService(ctx, req)
	spinner.Stop()

	if err != nil {
		outputpkg.Error(fmt.Sprintf("Failed to create service: %s", err))
		return err
	}

	// Success message
	outputpkg.Success(fmt.Sprintf("Successfully created service '%s'", service.Name))
	fmt.Println()

	formatType := outputpkg.Format(viper.GetString("defaults.output"))
	return outputpkg.Print(formatType, service)
}
