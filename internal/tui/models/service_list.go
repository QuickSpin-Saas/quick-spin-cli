package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/models"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
	"github.com/quickspin/quickspin-cli/internal/tui/components"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// ServiceListModel displays the list of services
type ServiceListModel struct {
	state     *types.AppState
	table     components.Table
	services  []models.Service
	loading   bool
	spinner   components.Spinner
	error     error
	width     int
	height    int
	statusBar components.StatusBar
}

// NewServiceListModel creates a new service list model
func NewServiceListModel(state *types.AppState) ServiceListModel {
	columns := []table.Column{
		components.NewColumn("Name", 25),
		components.NewColumn("Type", 15),
		components.NewColumn("Tier", 12),
		components.NewColumn("Status", 12),
		components.NewColumn("Region", 12),
	}

	serviceTable := components.NewTableWithTitle("Services", columns, []table.Row{}, state.Width-10, 15)

	return ServiceListModel{
		state:     state,
		table:     serviceTable,
		services:  []models.Service{},
		loading:   true,
		spinner:   components.NewSpinner("Loading services..."),
		error:     nil,
		width:     state.Width,
		height:    state.Height,
		statusBar: components.NewStatusBar(state.Width),
	}
}

// Init initializes the service list model
func (m ServiceListModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Init(),
		m.fetchServices(),
	)
}

// Update handles messages and updates the service list
func (m ServiceListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// Go back to dashboard
			return m, func() tea.Msg {
				return types.NavigationMsg{View: types.ViewDashboard}
			}

		case "r":
			// Refresh services
			m.loading = true
			m.error = nil
			return m, m.fetchServices()

		case "c":
			// Create new service
			return m, func() tea.Msg {
				return types.NavigationMsg{View: types.ViewServiceCreate}
			}

		case "enter":
			// View service details (if we have selected service)
			if len(m.services) > 0 {
				return m, func() tea.Msg {
					return types.NavigationMsg{View: types.ViewServiceDetail}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.state.SetTerminalSize(msg.Width, msg.Height)
		m.statusBar.SetWidth(msg.Width)
		m.table.SetWidth(msg.Width - 10)

	case servicesLoadedMsg:
		m.loading = false
		m.services = msg.services
		m.state.UpdateServices(msg.services)
		m.updateTable()

	case servicesErrorMsg:
		m.loading = false
		m.error = msg.err
		// If error is authentication-related, clear the user state
		if msg.err != nil && isAuthError(msg.err) {
			m.state.ClearUser()
		}
	}

	// Update table if not loading
	if !m.loading {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update spinner if loading
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the service list
func (m ServiceListModel) View() string {
	theme := styles.GetTheme()

	// Header
	header := m.renderHeader()

	// Content
	var content string
	if m.loading {
		// Show loading spinner
		content = lipgloss.NewStyle().
			Padding(2).
			Render(m.spinner.View())
	} else if m.error != nil {
		// Show error
		content = styles.ErrorBoxStyle.Render("✗ " + m.error.Error())
	} else if len(m.services) == 0 {
		// Show empty state
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Padding(2).
			Italic(true)
		content = emptyStyle.Render("No services found\n\nPress 'c' to create your first service")
	} else {
		// Show table
		content = m.table.View()
	}

	// Status bar
	m.statusBar.SetLeft(fmt.Sprintf("Services (%d)", len(m.services)))
	m.statusBar.SetCenter(m.state.Router.PrintBreadcrumb())
	m.statusBar.SetRight("↑/↓: Navigate • Enter: Details • c: Create • r: Refresh • Esc: Back • q: Quit")
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
		contentStyle.Render(content),
		statusBar,
	)

	return view
}

// renderHeader renders the service list header
func (m ServiceListModel) renderHeader() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 2)

	header := titleStyle.Render("⚡ Services")

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

// updateTable updates the table with current services
func (m *ServiceListModel) updateTable() {
	var rows []table.Row

	for _, service := range m.services {
		// Status badge
		status := string(service.Status)
		if service.Status == models.ServiceStatusRunning {
			status = "✓ Running"
		} else if service.Status == models.ServiceStatusStopped {
			status = "● Stopped"
		} else if service.Status == models.ServiceStatusFailed {
			status = "✗ Failed"
		}

		row := components.NewRow(
			service.Name,
			string(service.Type),
			string(service.Tier),
			status,
			service.Region,
		)
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
}

// fetchServices fetches services from the API
func (m *ServiceListModel) fetchServices() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		services, err := m.state.Client.ListServices(ctx)
		if err != nil {
			return servicesErrorMsg{err: err}
		}

		return servicesLoadedMsg{services: services}
	}
}

// servicesLoadedMsg is sent when services are loaded
type servicesLoadedMsg struct {
	services []models.Service
}

// servicesErrorMsg is sent when loading services fails
type servicesErrorMsg struct {
	err error
}

// isAuthError checks if an error is authentication-related
func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "unauthorized") ||
	       strings.Contains(errStr, "not authenticated") ||
	       strings.Contains(errStr, "authentication required")
}
