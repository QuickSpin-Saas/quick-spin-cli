package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
	"github.com/quickspin/quickspin-cli/internal/tui/components"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// MenuItem represents a menu item in the dashboard
type MenuItem struct {
	Icon        string
	Title       string
	Description string
	View        types.ViewType
}

// DashboardModel is the main dashboard/menu screen
type DashboardModel struct {
	state        *types.AppState
	menuItems    []MenuItem
	selectedItem int
	width        int
	height       int
	statusBar    components.StatusBar
	err          error
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel(state *types.AppState) DashboardModel {
	menuItems := []MenuItem{
		{
			Icon:        "⚡",
			Title:       "Services",
			Description: "Manage your microservices",
			View:        types.ViewServiceList,
		},
		{
			Icon:        "🔐",
			Title:       "Authentication",
			Description: "Login, logout, and manage tokens",
			View:        types.ViewAuthMenu,
		},
		{
			Icon:        "⚙️",
			Title:       "Configuration",
			Description: "Settings and preferences",
			View:        types.ViewConfigMenu,
		},
		{
			Icon:        "❓",
			Title:       "Help",
			Description: "Commands and keyboard shortcuts",
			View:        types.ViewHelp,
		},
	}

	return DashboardModel{
		state:        state,
		menuItems:    menuItems,
		selectedItem: 0,
		width:        state.Width,
		height:       state.Height,
		statusBar:    components.NewStatusBar(state.Width),
	}
}

// Init initializes the dashboard
func (m DashboardModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the dashboard
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.selectedItem > 0 {
				m.selectedItem--
			}

		case "down", "j":
			if m.selectedItem < len(m.menuItems)-1 {
				m.selectedItem++
			}

		case "enter", " ":
			// Navigate to selected view
			selectedView := m.menuItems[m.selectedItem].View
			return m, func() tea.Msg {
				return types.NavigationMsg{View: selectedView}
			}

		case "?":
			return m, func() tea.Msg {
				return types.NavigationMsg{View: types.ViewHelp}
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

// View renders the dashboard
func (m DashboardModel) View() string {
	// Header
	header := m.renderHeader()

	// Main menu
	menu := m.renderMenu()

	// User info section
	userInfo := m.renderUserInfo()

	// Recent services section
	recentServices := m.renderRecentServices()

	// Status bar
	m.statusBar.SetLeft("QuickSpin CLI")
	m.statusBar.SetCenter("")
	m.statusBar.SetRight("↑/↓: Navigate • Enter: Select • q: Quit • ?: Help")
	statusBar := m.statusBar.View()

	// Calculate available height
	headerHeight := lipgloss.Height(header)
	statusBarHeight := 1
	contentHeight := m.height - headerHeight - statusBarHeight - 2

	// Combine sections
	leftColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		menu,
	)

	rightColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		userInfo,
		"",
		recentServices,
	)

	// Create two-column layout
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(m.width/2).Render(leftColumn),
		lipgloss.NewStyle().Width(m.width/2).Render(rightColumn),
	)

	// Ensure content fits within available height
	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		Width(m.width)

	// Build final view
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		contentStyle.Render(content),
		statusBar,
	)

	return view
}

// renderHeader renders the dashboard header
func (m DashboardModel) renderHeader() string {
	theme := styles.GetTheme()

	// QuickSpin branding
	brandingStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 2)

	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Muted))

	// API status indicator
	var apiStatus string
	apiStatusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Success)).
		Bold(true)
	apiStatus = apiStatusStyle.Render("[●] Connected")

	// Build header
	left := brandingStyle.Render("⚡ QuickSpin CLI") + " " + versionStyle.Render("v1.0.0")
	right := apiStatus

	// Calculate spacing
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	spacing := m.width - leftWidth - rightWidth - 4

	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		left,
		strings.Repeat(" ", spacing),
		right,
	)

	// Add border
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

// renderMenu renders the main menu
func (m DashboardModel) renderMenu() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(1, 2)

	var menuItems []string
	for i, item := range m.menuItems {
		var itemStyle lipgloss.Style

		if i == m.selectedItem {
			// Selected item
			itemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Primary)).
				Background(lipgloss.Color(theme.Highlight)).
				Bold(true).
				Padding(0, 2).
				Width(m.width/2 - 4)
		} else {
			// Normal item
			itemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(theme.Foreground)).
				Padding(0, 2).
				Width(m.width/2 - 4)
		}

		iconStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Accent)).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted))

		itemText := fmt.Sprintf("%s %s", iconStyle.Render(item.Icon), item.Title)
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
		titleStyle.Render("Main Menu"),
		lipgloss.JoinVertical(lipgloss.Left, menuItems...),
	)
}

// renderUserInfo renders user/org information
func (m DashboardModel) renderUserInfo() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		Padding(1, 2)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Border)).
		Padding(1, 2).
		Width(m.width/2 - 6)

	var content string

	if m.state.IsAuthenticated && m.state.CurrentUser != nil {
		// Authenticated user
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			styles.RenderKeyValue("User", m.state.CurrentUser.Email),
			styles.RenderKeyValue("Name", m.state.CurrentUser.FirstName+" "+m.state.CurrentUser.LastName),
			styles.RenderKeyValue("Status", styles.RenderBadge("Active", styles.SuccessBadgeStyle)),
		)
	} else {
		// Not authenticated
		notAuthStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Warning)).
			Italic(true)

		content = notAuthStyle.Render("Not logged in\n\nSelect 'Authentication' to login")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Account"),
		boxStyle.Render(content),
	)
}

// renderRecentServices renders recent services section
func (m DashboardModel) renderRecentServices() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		Padding(1, 2)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Border)).
		Padding(1, 2).
		Width(m.width/2 - 6)

	var content string

	if len(m.state.RecentServices) > 0 {
		var items []string
		for _, service := range m.state.RecentServices {
			statusBadge := styles.RenderBadge("running", styles.SuccessBadgeStyle)
			item := fmt.Sprintf("• %s %s", service.Name, statusBadge)
			items = append(items, item)
		}
		content = strings.Join(items, "\n")
	} else {
		noServicesStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Italic(true)

		content = noServicesStyle.Render("No services yet\n\nSelect 'Services' to create your first service")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Recent Services"),
		boxStyle.Render(content),
	)
}
