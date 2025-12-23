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

var (
	logsLines int
)

// NewLogsCmd creates the service logs command
func NewLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs SERVICE_ID",
		Short: "View service logs",
		Long:  "View logs for a specific service",
		Args:  cobra.ExactArgs(1),
		RunE:  runLogs,
	}

	cmd.Flags().IntVarP(&logsLines, "lines", "n", 100, "Number of log lines to retrieve")

	return cmd
}

func runLogs(cmd *cobra.Command, args []string) error {
	serviceID := args[0]

	// Check if we should use TUI mode
	outputFormat := viper.GetString("defaults.output")
	if outputpkg.ShouldUseTUI(outputFormat) {
		// Launch TUI service logs view
		return tui.LaunchView(tui.ViewServiceLogs)
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
	spinner := outputpkg.NewSpinner(fmt.Sprintf("Loading last %d log lines...", logsLines))
	spinner.Start()

	// Get service logs
	logs, err := client.GetServiceLogs(ctx, serviceID, logsLines)
	spinner.Stop()

	if err != nil {
		outputpkg.Error(fmt.Sprintf("Failed to get service logs: %s", err))
		return err
	}

	if len(logs) == 0 {
		outputpkg.Info("No logs found")
		return nil
	}

	// Display logs
	outputpkg.Success(fmt.Sprintf("Retrieved %d log entries", len(logs)))
	fmt.Println()

	// Print logs in chronological order
	for _, log := range logs {
		timestamp := log.Timestamp.Format("2006-01-02 15:04:05")
		level := log.Level
		message := log.Message

		// Format log level prefix
		var prefix string
		switch level {
		case "error":
			prefix = "\033[31m[ERROR]\033[0m" // Red
		case "warning", "warn":
			prefix = "\033[33m[WARN]\033[0m" // Yellow
		case "info":
			prefix = "\033[36m[INFO]\033[0m" // Cyan
		default:
			prefix = fmt.Sprintf("[%s]", level)
		}

		fmt.Printf("%s %s %s\n", timestamp, prefix, message)
	}

	return nil
}
