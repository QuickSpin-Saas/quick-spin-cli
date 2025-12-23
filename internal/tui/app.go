package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quickspin/quickspin-cli/internal/config"
	"github.com/quickspin/quickspin-cli/internal/output"
	"github.com/quickspin/quickspin-cli/internal/tui/models"
)

// Model is the main Bubble Tea model for the application
type Model struct {
	state      *AppState
	activeView tea.Model
	err        error
	quitting   bool
}

// NewModel creates a new TUI model with the given configuration
func NewModel(cfg *config.Config) Model {
	state := NewAppState(cfg)

	// Initialize with dashboard view
	dashboard := models.NewDashboardModel(state)

	return Model{
		state:      state,
		activeView: dashboard,
		err:        nil,
		quitting:   false,
	}
}

// NewModelWithView creates a new TUI model starting at a specific view
func NewModelWithView(cfg *config.Config, view ViewType) Model {
	state := NewAppStateWithView(cfg, view)

	// Create a temporary model to use createView
	m := Model{
		state:      state,
		activeView: nil,
		err:        nil,
		quitting:   false,
	}

	// Initialize with specified view
	m.activeView = m.createView(view)

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Get terminal size
	cmds = append(cmds, getTerminalSize)

	// Initialize active view
	if m.activeView != nil {
		cmds = append(cmds, m.activeView.Init())
	}

	return tea.Batch(cmds...)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global keyboard shortcuts
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "?":
			m.state.ToggleHelp()
			return m, nil

		case "esc":
			// Go back if we can
			if m.state.Router.CanGoBack() {
				m.state.Router.Pop()
				// TODO: Switch to the previous view
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.state.SetTerminalSize(msg.Width, msg.Height)
		// Forward to active view
		if m.activeView != nil {
			var cmd tea.Cmd
			m.activeView, cmd = m.activeView.Update(msg)
			return m, cmd
		}

	case NavigationMsg:
		m.state.Router.Push(msg.View)
		// Switch to the new view
		m.activeView = m.createView(msg.View)
		if m.activeView != nil {
			return m, m.activeView.Init()
		}
		return m, nil

	case BackMsg:
		if m.state.Router.CanGoBack() {
			prevView := m.state.Router.Pop()
			// Switch to the previous view
			m.activeView = m.createView(prevView)
			if m.activeView != nil {
				return m, m.activeView.Init()
			}
		}
		return m, nil

	case ExitMsg:
		m.quitting = true
		return m, tea.Quit

	case terminalSizeMsg:
		m.state.SetTerminalSize(msg.width, msg.height)
	}

	// Forward messages to the active view
	if m.activeView != nil {
		var cmd tea.Cmd
		m.activeView, cmd = m.activeView.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the model
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if m.activeView != nil {
		return m.activeView.View()
	}

	// Fallback view if no active view is set
	return "Initializing QuickSpin TUI...\n"
}

// createView creates a view model based on the ViewType
func (m Model) createView(viewType ViewType) tea.Model {
	switch viewType {
	case ViewDashboard:
		return models.NewDashboardModel(m.state)
	case ViewAuthMenu:
		return models.NewAuthMenuModel(m.state)
	case ViewAuthLogin:
		return models.NewAuthLoginModel(m.state)
	case ViewAuthLogout:
		return models.NewAuthLogoutModel(m.state)
	case ViewAuthWhoami:
		return models.NewAuthWhoamiModel(m.state)
	case ViewServiceList:
		return models.NewServiceListModel(m.state)
	case ViewServiceCreate:
		return models.NewServiceCreateModel(m.state)
	case ViewHelp:
		return models.NewHelpModel(m.state)
	// Views that still need implementation
	// case ViewServiceDetail, ViewServiceLogs, ViewConfigEditor, ViewConfigView:
	//     return models.NewDashboardModel(m.state) // Fallback for now
	default:
		// Return dashboard as fallback
		return models.NewDashboardModel(m.state)
	}
}

// terminalSizeMsg is sent when terminal size is detected
type terminalSizeMsg struct {
	width  int
	height int
}

// getTerminalSize gets the current terminal size
func getTerminalSize() tea.Msg {
	width, height, err := output.GetTerminalSize()
	if err != nil {
		// Return default size
		return terminalSizeMsg{width: 80, height: 24}
	}
	return terminalSizeMsg{width: width, height: height}
}

// LaunchApp starts the TUI application at the dashboard
func LaunchApp() error {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create model
	m := NewModel(cfg)

	// Create program with alt screen
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}

// LaunchView starts the TUI application at a specific view
func LaunchView(view ViewType) error {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create model with specific view
	m := NewModelWithView(cfg, view)

	// Create program with alt screen
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
