package models

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
	"github.com/quickspin/quickspin-cli/internal/tui/components"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// AuthLogoutModel is the logout confirmation screen
type AuthLogoutModel struct {
	state        *types.AppState
	modal        components.Modal
	confirmed    bool
	loading      bool
	spinner      components.Spinner
	width        int
	height       int
	statusBar    components.StatusBar
}

// NewAuthLogoutModel creates a new logout model
func NewAuthLogoutModel(state *types.AppState) AuthLogoutModel {
	modal := components.NewConfirmModal(
		"Confirm Logout",
		"Are you sure you want to logout from QuickSpin?",
		"Logout",
		"Cancel",
	)
	modal.Show()

	return AuthLogoutModel{
		state:     state,
		modal:     modal,
		confirmed: false,
		loading:   false,
		spinner:   components.NewSpinner("Logging out..."),
		width:     state.Width,
		height:    state.Height,
		statusBar: components.NewStatusBar(state.Width),
	}
}

// Init initializes the logout model
func (m AuthLogoutModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the logout model
func (m AuthLogoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			// Ignore input while loading
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			// Cancel and go back
			return m, func() tea.Msg {
				return types.NavigationMsg{View: types.ViewAuthMenu}
			}

		case "left", "h", "right", "l", "tab":
			// Update modal for button selection
			m.modal, cmd = m.modal.Update(msg)
			cmds = append(cmds, cmd)

		case "enter", " ":
			// Handle modal confirmation
			m.modal, cmd = m.modal.Update(msg)
			if !m.modal.IsVisible() {
				// Modal was closed, perform logout if confirmed
				if m.confirmed {
					return m, m.performLogout()
				} else {
					// Cancelled, go back
					return m, func() tea.Msg {
						return types.NavigationMsg{View: types.ViewAuthMenu}
					}
				}
			}
			cmds = append(cmds, cmd)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.state.SetTerminalSize(msg.Width, msg.Height)
		m.statusBar.SetWidth(msg.Width)

	case logoutSuccessMsg:
		m.loading = false
		m.state.ClearUser()
		// Navigate to dashboard
		return m, func() tea.Msg {
			return types.NavigationMsg{View: types.ViewDashboard}
		}
	}

	// Update spinner if loading
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the logout confirmation
func (m AuthLogoutModel) View() string {
	// Header
	header := m.renderHeader()

	// Content
	var content string
	if m.loading {
		// Show loading spinner
		content = lipgloss.NewStyle().
			Padding(2).
			AlignHorizontal(lipgloss.Center).
			Render(m.spinner.View())
	} else {
		// Show confirmation modal
		content = m.modal.RenderOverlay("", m.width, m.height-5)
	}

	// Status bar
	m.statusBar.SetLeft("Logout")
	m.statusBar.SetCenter(m.state.Router.PrintBreadcrumb())
	m.statusBar.SetRight("←/→: Select • Enter: Confirm • Esc: Cancel • q: Quit")
	statusBar := m.statusBar.View()

	// Build final view
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		statusBar,
	)

	return view
}

// renderHeader renders the logout header
func (m AuthLogoutModel) renderHeader() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 2)

	header := titleStyle.Render("🚪 Logout")

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground)).
		Background(lipgloss.Color(theme.Highlight)).
		Padding(0, 2).
		Width(m.width)

	divider := styles.RenderDivider(m.width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render(header),
		divider,
	)
}

// performLogout performs the logout operation
func (m *AuthLogoutModel) performLogout() tea.Cmd {
	m.confirmed = true
	m.loading = true

	return func() tea.Msg {
		ctx := context.Background()

		// Perform logout
		err := m.state.Client.Logout(ctx)
		if err != nil {
			// Even if logout fails on server, clear local state
			return logoutSuccessMsg{}
		}

		return logoutSuccessMsg{}
	}
}

// logoutSuccessMsg is sent when logout succeeds
type logoutSuccessMsg struct{}
