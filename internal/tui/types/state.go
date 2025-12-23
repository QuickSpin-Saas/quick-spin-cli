package types

import (
	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/models"
)

// Router interface to avoid import cycles
type Router interface {
	Push(view ViewType)
	Pop() ViewType
	Current() ViewType
	CanGoBack() bool
	Reset()
	GetBreadcrumb() string
	PrintBreadcrumb() string
}

// AppState holds the global application state
type AppState struct {
	// Configuration
	Config *config.Config

	// API client
	Client *api.Client

	// Authentication
	IsAuthenticated bool
	CurrentUser     *models.User

	// Navigation
	Router Router

	// Terminal dimensions
	Width  int
	Height int

	// UI state
	ShowHelp bool
	Error    error

	// Service state (cached)
	Services []models.Service

	// Recent activity (for dashboard)
	RecentServices []models.Service
}

// SetUser sets the current authenticated user
func (s *AppState) SetUser(user *models.User) {
	s.CurrentUser = user
	s.IsAuthenticated = user != nil
}

// ClearUser clears the current user (logout)
func (s *AppState) ClearUser() {
	s.CurrentUser = nil
	s.IsAuthenticated = false
}

// SetTerminalSize updates the terminal dimensions
func (s *AppState) SetTerminalSize(width, height int) {
	s.Width = width
	s.Height = height
}

// SetError sets the global error state
func (s *AppState) SetError(err error) {
	s.Error = err
}

// ClearError clears the global error state
func (s *AppState) ClearError() {
	s.Error = nil
}

// ToggleHelp toggles the help screen visibility
func (s *AppState) ToggleHelp() {
	s.ShowHelp = !s.ShowHelp
}

// UpdateServices updates the cached services list
func (s *AppState) UpdateServices(services []models.Service) {
	s.Services = services

	// Update recent services (keep last 5)
	maxRecent := 5
	if len(services) > 0 {
		if len(services) > maxRecent {
			s.RecentServices = services[:maxRecent]
		} else {
			s.RecentServices = services
		}
	}
}
