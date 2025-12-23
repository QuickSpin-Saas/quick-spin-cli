package auth

import (
	"context"
	"fmt"
	"syscall"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	outputpkg "github.com/quickspin/quickspin-cli/internal/output"
	"github.com/quickspin/quickspin-cli/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	email    string
	password string
)

// NewLoginCmd creates the login command
func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to QuickSpin",
		Long:  "Authenticate with QuickSpin using email and password",
		RunE:  runLogin,
	}

	cmd.Flags().StringVar(&email, "email", "", "Email address")
	cmd.Flags().StringVar(&password, "password", "", "Password (not recommended, will prompt if not provided)")

	return cmd
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Check if we should use TUI mode
	outputFormat := viper.GetString("defaults.output")
	if outputpkg.ShouldUseTUI(outputFormat) && email == "" && password == "" {
		// Launch TUI login form
		return tui.LaunchView(tui.ViewAuthLogin)
	}

	// Traditional CLI mode
	ctx := context.Background()

	// Get email if not provided
	if email == "" {
		fmt.Print("Email: ")
		if _, err := fmt.Scanln(&email); err != nil {
			return fmt.Errorf("failed to read email: %w", err)
		}
	}

	// Get password if not provided (securely)
	if password == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // Add newline after password input
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password = string(passwordBytes)
	}

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client := api.NewClient(cfg)

	// Show spinner
	spinner := outputpkg.NewSpinner("Logging in...")
	spinner.Start()

	// Perform login
	result, err := client.Login(ctx, email, password)
	spinner.Stop()

	if err != nil {
		outputpkg.Error(fmt.Sprintf("Login failed: %s", err))
		return err
	}

	// Success message
	outputpkg.Success(fmt.Sprintf("Successfully logged in as %s", result.User.Email))

	// Display user info
	fmt.Println()
	formatType := outputpkg.Format(viper.GetString("defaults.output"))
	return outputpkg.Print(formatType, result.User)
}
