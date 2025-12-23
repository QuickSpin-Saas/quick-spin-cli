package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
	"github.com/quickspin/quickspin-cli/internal/tui/components"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// AuthMenuItem represents a menu item in the auth menu
type AuthMenuItem struct {
	Icon        string
	Title       string
	Description string
	View        types.ViewType
	RequiresAuth bool
}

// AuthMenuModel is the authentication menu screen
type AuthMenuModel struct {
	state        *types.AppState
	menuItems    []AuthMenuItem
	selectedItem int
	width        int
	height       int
	statusBar    components.StatusBar
}

// NewAuthMenuModel creates a new auth menu model
func NewAuthMenuModel(state *types.AppState) AuthMenuModel {
	menuItems := []AuthMenuItem{
		{
			Icon:         "🔑",
			Title:        "Login",
			Description:  "Authenticate with QuickSpin",
			View:         types.ViewAuthLogin,
			RequiresAuth: false,
		},
		{
			Icon:         "👤",
			Title:        "Current User",
			Description:  "View your account information",
			View:         types.ViewAuthWhoami,
			RequiresAuth: true,
		},
		{
			Icon:         "🚪",
			Title:        "Logout",
			Description:  "Sign out from QuickSpin",
			View:         types.ViewAuthLogout,
			RequiresAuth: true,
		},
		{
			Icon:         "⬅️",
			Title:        "Back to Dashboard",
			Description:  "Return to main menu",
			View:         types.ViewDashboard,
			RequiresAuth: false,
		},
	}

	return AuthMenuModel{
		state:        state,
		menuItems:    menuItems,
		selectedItem: 0,
		width:        state.Width,
		height:       state.Height,
		statusBar:    components.NewStatusBar(state.Width),
	}
}

// Init initializes the auth menu
func (m AuthMenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the auth menu
func (m AuthMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			// Go back to dashboard
			return m, func() tea.Msg {
				return types.NavigationMsg{View: types.ViewDashboard}
			}

		case "up", "k":
			if m.selectedItem > 0 {
				m.selectedItem--
			}

		case "down", "j":
			if m.selectedItem < len(m.menuItems)-1 {
				m.selectedItem++
			}

		case "enter", " ":
			// Check if item requires auth
			selectedItem := m.menuItems[m.selectedItem]
			if selectedItem.RequiresAuth && !m.state.IsAuthenticated {
				// Show error - not authenticated
				return m, nil
			}

			// Navigate to selected view
			return m, func() tea.Msg {
				return types.NavigationMsg{View: selectedItem.View}
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

// View renders the auth menu
func (m AuthMenuModel) View() string {
	// Header
	header := m.renderHeader()

	// Menu
	menu := m.renderMenu()

	// Status bar
	m.statusBar.SetLeft("Authentication")
	m.statusBar.SetCenter(m.state.Router.PrintBreadcrumb())
	m.statusBar.SetRight("↑/↓: Navigate • Enter: Select • Esc: Back • q: Quit")
	statusBar := m.statusBar.View()

	// Calculate available height
	headerHeight := lipgloss.Height(header)
	statusBarHeight := 1
	contentHeight := m.height - headerHeight - statusBarHeight - 2

	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		Width(m.width)

	// Build final view
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		contentStyle.Render(menu),
		statusBar,
	)

	return view
}

// renderHeader renders the auth menu header
func (m AuthMenuModel) renderHeader() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 2)

	var authStatus string
	if m.state.IsAuthenticated && m.state.CurrentUser != nil {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Success)).
			Bold(true)
		authStatus = statusStyle.Render("✓ Authenticated as " + m.state.CurrentUser.Email)
	} else {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Warning)).
			Bold(true)
		authStatus = statusStyle.Render("⚠ Not authenticated")
	}

	header := titleStyle.Render("🔐 Authentication")

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground)).
		Background(lipgloss.Color(theme.Highlight)).
		Padding(0, 2).
		Width(m.width)

	divider := styles.RenderDivider(m.width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render(header),
		lipgloss.NewStyle().Padding(0, 2).Render(authStatus),
		divider,
	)
}

// renderMenu renders the auth menu items
func (m AuthMenuModel) renderMenu() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		Padding(1, 2)

	var menuItems []string
	for i, item := range m.menuItems {
		var itemStyle lipgloss.Style
		disabled := item.RequiresAuth && !m.state.IsAuthenticated

		if i == m.selectedItem {
			// Selected item
			if disabled {
				itemStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.Muted)).
					Background(lipgloss.Color(theme.Highlight)).
					Padding(0, 2).
					Width(m.width - 8)
			} else {
				itemStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.Primary)).
					Background(lipgloss.Color(theme.Highlight)).
					Bold(true).
					Padding(0, 2).
					Width(m.width - 8)
			}
		} else {
			// Normal item
			if disabled {
				itemStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.Muted)).
					Padding(0, 2).
					Width(m.width - 8)
			} else {
				itemStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.Foreground)).
					Padding(0, 2).
					Width(m.width - 8)
			}
		}

		iconStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Accent)).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted))

		itemText := fmt.Sprintf("%s %s", iconStyle.Render(item.Icon), item.Title)
		if disabled {
			itemText += " " + lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Warning)).Render("(requires login)")
		}

		if i == m.selectedItem {
			itemText = fmt.Sprintf("▶ %s", itemText)
		} else {
			itemText = fmt.Sprintf("  %s", itemText)
		}

		itemView := lipgloss.JoinVertical(
			lipgloss.Left,
			itemStyle.Render(itemText),
			lipgloss.NewStyle().Padding(0, 4).Render(descStyle.Render(item.Description)),
		)

		menuItems = append(menuItems, itemView)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Options"),
		lipgloss.JoinVertical(lipgloss.Left, menuItems...),
	)
}
