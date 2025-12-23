package tui

import (
	"context"

	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/models"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
)

// Re-export AppState for backward compatibility
type AppState = types.AppState

// NewAppState creates a new application state
func NewAppState(cfg *config.Config) *AppState {
	client := api.NewClient(cfg)

	// Check if user has valid credentials
	isAuthenticated := false
	var currentUser *models.User = nil

	// Try to get token and validate authentication
	token, err := cfg.GetToken()
	if err == nil && token != "" {
		// We have a token, try to get current user to validate it
		ctx := context.Background()
		user, err := client.WhoAmI(ctx)
		if err == nil {
			// Token is valid and we got user info
			isAuthenticated = true
			currentUser = user
		} else {
			// Token is invalid or expired, clear it
			_ = cfg.ClearToken()
		}
	}

	return &AppState{
		Config:          cfg,
		Client:          client,
		IsAuthenticated: isAuthenticated,
		CurrentUser:     currentUser,
		Router:          NewRouter(),
		Width:           80,
		Height:          24,
		ShowHelp:        false,
		Services:        []models.Service{},
		RecentServices:  []models.Service{},
	}
}

// NewAppStateWithView creates a new application state starting at a specific view
func NewAppStateWithView(cfg *config.Config, view ViewType) *AppState {
	client := api.NewClient(cfg)

	// Check if user has valid credentials
	isAuthenticated := false
	var currentUser *models.User = nil

	// Try to get token and validate authentication
	token, err := cfg.GetToken()
	if err == nil && token != "" {
		// We have a token, try to get current user to validate it
		ctx := context.Background()
		user, err := client.WhoAmI(ctx)
		if err == nil {
			// Token is valid and we got user info
			isAuthenticated = true
			currentUser = user
		} else {
			// Token is invalid or expired, clear it
			_ = cfg.ClearToken()
		}
	}

	return &AppState{
		Config:          cfg,
		Client:          client,
		IsAuthenticated: isAuthenticated,
		CurrentUser:     currentUser,
		Router:          NewRouterWithView(view),
		Width:           80,
		Height:          24,
		ShowHelp:        false,
		Services:        []models.Service{},
		RecentServices:  []models.Service{},
	}
}

