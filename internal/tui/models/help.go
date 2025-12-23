package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
	"github.com/quickspin/quickspin-cli/internal/tui/components"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// HelpModel displays help and keyboard shortcuts
type HelpModel struct {
	state     *types.AppState
	width     int
	height    int
	statusBar components.StatusBar
}

// NewHelpModel creates a new help model
func NewHelpModel(state *types.AppState) HelpModel {
	return HelpModel{
		state:     state,
		width:     state.Width,
		height:    state.Height,
		statusBar: components.NewStatusBar(state.Width),
	}
}

// Init initializes the help model
func (m HelpModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc", "enter", " ":
			// Go back to dashboard
			return m, func() tea.Msg {
				return types.NavigationMsg{View: types.ViewDashboard}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.state.SetTerminalSize(msg.Width, msg.Height)
		m.statusBar.SetWidth(msg.Width)
	}

	return m, nil
}

// View renders the help screen
func (m HelpModel) View() string {
	theme := styles.GetTheme()

	// Header
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 2)

	header := titleStyle.Render("❓ Help & Keyboard Shortcuts")

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground)).
		Background(lipgloss.Color(theme.Highlight)).
		Padding(0, 2).
		Width(m.width)

	divider := styles.RenderDivider(m.width)

	headerView := lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render(header),
		divider,
	)

	// Help content
	content := m.renderHelpContent()

	// Status bar
	m.statusBar.SetLeft("Help")
	m.statusBar.SetCenter(m.state.Router.PrintBreadcrumb())
	m.statusBar.SetRight("Enter/Esc: Back • q: Quit")
	statusBar := m.statusBar.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		content,
		statusBar,
	)
}

// renderHelpContent renders the help content
func (m HelpModel) renderHelpContent() string {
	theme := styles.GetTheme()

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	// Global shortcuts
	globalSection := sectionStyle.Render("Global Shortcuts")
	globalShortcuts := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.RenderHelp("Ctrl+C, q", "Quit application"),
		styles.RenderHelp("Esc", "Go back / Cancel"),
		styles.RenderHelp("?", "Toggle help"),
		styles.RenderHelp("Enter", "Select / Confirm"),
	)

	// Navigation shortcuts
	navSection := sectionStyle.Render("Navigation")
	navShortcuts := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.RenderHelp("↑/k, ↓/j", "Move up / down"),
		styles.RenderHelp("Tab", "Next field"),
		styles.RenderHelp("Shift+Tab", "Previous field"),
		styles.RenderHelp("PgUp/PgDn", "Scroll page"),
	)

	// Service shortcuts
	serviceSection := sectionStyle.Render("Service Management")
	serviceShortcuts := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.RenderHelp("c", "Create new service"),
		styles.RenderHelp("r", "Refresh list"),
		styles.RenderHelp("Enter", "View details"),
	)

	// Commands section
	commandsSection := sectionStyle.Render("CLI Commands")
	commands := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.RenderHelp("qspin", "Launch TUI dashboard"),
		styles.RenderHelp("qspin auth login", "Login to QuickSpin"),
		styles.RenderHelp("qspin service list", "List services"),
		styles.RenderHelp("qspin service create", "Create service"),
		styles.RenderHelp("qspin --output json", "Force CLI mode (no TUI)"),
	)

	// Combine all sections
	leftColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		globalSection,
		globalShortcuts,
		"",
		navSection,
		navShortcuts,
	)

	rightColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		serviceSection,
		serviceShortcuts,
		"",
		commandsSection,
		commands,
	)

	helpContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(m.width/2).Padding(2).Render(leftColumn),
		lipgloss.NewStyle().Width(m.width/2).Padding(2).Render(rightColumn),
	)

	return styles.BoxStyle.
		Width(m.width - 4).
		Render(helpContent)
}
