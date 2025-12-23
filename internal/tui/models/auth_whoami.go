package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
	"github.com/quickspin/quickspin-cli/internal/tui/components"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// AuthWhoamiModel displays current user information
type AuthWhoamiModel struct {
	state     *types.AppState
	width     int
	height    int
	statusBar components.StatusBar
}

// NewAuthWhoamiModel creates a new whoami model
func NewAuthWhoamiModel(state *types.AppState) AuthWhoamiModel {
	return AuthWhoamiModel{
		state:     state,
		width:     state.Width,
		height:    state.Height,
		statusBar: components.NewStatusBar(state.Width),
	}
}

// Init initializes the whoami model
func (m AuthWhoamiModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the whoami model
func (m AuthWhoamiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc", "enter", " ":
			// Go back to auth menu
			return m, func() tea.Msg {
				return types.NavigationMsg{View: types.ViewAuthMenu}
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

// View renders the current user information
func (m AuthWhoamiModel) View() string {
	// Header
	header := m.renderHeader()

	// User info
	userInfo := m.renderUserInfo()

	// Status bar
	m.statusBar.SetLeft("Current User")
	m.statusBar.SetCenter(m.state.Router.PrintBreadcrumb())
	m.statusBar.SetRight("Enter/Esc: Back • q: Quit")
	statusBar := m.statusBar.View()

	// Calculate available height
	headerHeight := lipgloss.Height(header)
	statusBarHeight := 1
	contentHeight := m.height - headerHeight - statusBarHeight - 2

	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		Width(m.width).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	// Build final view
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		contentStyle.Render(userInfo),
		statusBar,
	)

	return view
}

// renderHeader renders the whoami header
func (m AuthWhoamiModel) renderHeader() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 2)

	header := titleStyle.Render("👤 Current User")

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

// renderUserInfo renders the user information
func (m AuthWhoamiModel) renderUserInfo() string {
	theme := styles.GetTheme()

	if !m.state.IsAuthenticated || m.state.CurrentUser == nil {
		// Not authenticated
		notAuthStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Warning)).
			Bold(true).
			Padding(2)

		return styles.BoxStyle.
			Width(60).
			Render(notAuthStyle.Render("⚠ Not authenticated\n\nPlease login first."))
	}

	user := m.state.CurrentUser

	// User details
	var details []string

	details = append(details, styles.RenderKeyValue("Email", user.Email))
	details = append(details, styles.RenderKeyValue("Name", fmt.Sprintf("%s %s", user.FirstName, user.LastName)))
	details = append(details, styles.RenderKeyValue("User ID", user.ID))

	if user.OrganizationID != "" {
		details = append(details, styles.RenderKeyValue("Organization", user.OrganizationID))
	}

	if user.Role != "" {
		roleBadge := styles.RenderBadge(string(user.Role), styles.InfoBadgeStyle)
		details = append(details, styles.RenderKeyValue("Role", roleBadge))
	}

	// Account status
	statusBadge := styles.RenderBadge("Active", styles.SuccessBadgeStyle)
	details = append(details, styles.RenderKeyValue("Status", statusBadge))

	// Created/Updated dates
	if !user.CreatedAt.IsZero() {
		details = append(details, styles.RenderKeyValue("Created", user.CreatedAt.Format("2006-01-02 15:04:05")))
	}
	if !user.UpdatedAt.IsZero() {
		details = append(details, styles.RenderKeyValue("Updated", user.UpdatedAt.Format("2006-01-02 15:04:05")))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, details...)

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		MarginBottom(1)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Account Information"),
		content,
	)

	return styles.BoxStyle.
		Width(70).
		Render(fullContent)
}
