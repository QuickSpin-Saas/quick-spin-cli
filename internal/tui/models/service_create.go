package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/api"
	"github.com/quickspin/quickspin-cli/internal/models"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
	"github.com/quickspin/quickspin-cli/internal/tui/components"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// ServiceCreateModel is a wizard for creating a new service
type ServiceCreateModel struct {
	state           *types.AppState
	currentStep     int
	nameInput       components.TextInput
	typeSelect      components.SelectInput
	tierSelect      components.SelectInput
	descriptionInput components.TextInput
	focusedField    int
	loading         bool
	error           error
	spinner         components.Spinner
	progress        components.MultiStepProgress
	width           int
	height          int
	statusBar       components.StatusBar
}

// NewServiceCreateModel creates a new service create model
func NewServiceCreateModel(state *types.AppState) ServiceCreateModel {
	serviceTypes := []string{"redis", "rabbitmq", "postgresql", "mongodb", "mysql", "elasticsearch"}
	serviceTiers := []string{"developer", "basic", "standard", "pro", "premium"}

	steps := []string{"Service Name", "Service Type", "Service Tier", "Description", "Confirm"}
	progress := components.NewMultiStepProgress(steps, 60)

	return ServiceCreateModel{
		state:            state,
		currentStep:      0,
		nameInput:        components.NewTextInput("Service Name", "my-service", 40),
		typeSelect:       components.NewSelectInput("Service Type", serviceTypes, 40),
		tierSelect:       components.NewSelectInput("Service Tier", serviceTiers, 40),
		descriptionInput: components.NewTextInput("Description (optional)", "", 40),
		focusedField:     0,
		loading:          false,
		error:            nil,
		spinner:          components.NewSpinner("Creating service..."),
		progress:         progress,
		width:            state.Width,
		height:           state.Height,
		statusBar:        components.NewStatusBar(state.Width),
	}
}

// Init initializes the service create model
func (m ServiceCreateModel) Init() tea.Cmd {
	return m.nameInput.Focus()
}

// Update handles messages
func (m ServiceCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}

		// Handle special navigation keys first
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			if m.currentStep == 0 {
				// Go back to service list
				return m, func() tea.Msg {
					return types.NavigationMsg{View: types.ViewServiceList}
				}
			} else {
				// Go to previous step
				m.currentStep--
				m.progress.PrevStep()
				return m, nil
			}

		case "ctrl+n", "ctrl+j":
			// Next step
			if m.currentStep < 4 {
				m.currentStep++
				m.progress.NextStep()
				return m, nil
			}

		case "ctrl+p", "ctrl+k":
			// Previous step
			if m.currentStep > 0 {
				m.currentStep--
				m.progress.PrevStep()
				return m, nil
			}

		default:
			// Let input fields handle all other keys including enter
			// Update active field based on current step
			switch m.currentStep {
			case 0:
				m.nameInput, cmd = m.nameInput.Update(msg)
				cmds = append(cmds, cmd)
			case 1:
				m.typeSelect, cmd = m.typeSelect.Update(msg)
				cmds = append(cmds, cmd)
			case 2:
				m.tierSelect, cmd = m.tierSelect.Update(msg)
				cmds = append(cmds, cmd)
			case 3:
				m.descriptionInput, cmd = m.descriptionInput.Update(msg)
				cmds = append(cmds, cmd)
			case 4:
				// Confirmation step - enter creates service
				if msg.String() == "enter" {
					return m, m.createService()
				}
			}
			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.state.SetTerminalSize(msg.Width, msg.Height)
		m.statusBar.SetWidth(msg.Width)

	case serviceCreatedMsg:
		m.loading = false
		// Navigate to service list
		return m, func() tea.Msg {
			return types.NavigationMsg{View: types.ViewServiceList}
		}

	case serviceCreateErrorMsg:
		m.loading = false
		m.error = msg.err
	}

	// Update spinner if loading
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the service create wizard
func (m ServiceCreateModel) View() string {
	// Header
	header := m.renderHeader()

	// Content
	var content string
	if m.loading {
		content = lipgloss.NewStyle().Padding(2).Render(m.spinner.View())
	} else {
		content = m.renderWizard()
	}

	// Status bar
	m.statusBar.SetLeft(fmt.Sprintf("Create Service - Step %d of 5", m.currentStep+1))
	m.statusBar.SetCenter(m.state.Router.PrintBreadcrumb())
	m.statusBar.SetRight("Enter: Next • Esc: Back • q: Quit")
	statusBar := m.statusBar.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		statusBar,
	)
}

// renderHeader renders the header
func (m ServiceCreateModel) renderHeader() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 2)

	header := titleStyle.Render("⚡ Create New Service")

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

// renderWizard renders the wizard steps
func (m ServiceCreateModel) renderWizard() string {
	var stepContent string

	switch m.currentStep {
	case 0:
		stepContent = m.nameInput.View()
	case 1:
		stepContent = m.typeSelect.View()
	case 2:
		stepContent = m.tierSelect.View()
	case 3:
		stepContent = m.descriptionInput.View()
	case 4:
		// Confirmation step
		stepContent = m.renderConfirmation()
	}

	// Error message
	var errorView string
	if m.error != nil {
		errorView = styles.ErrorBoxStyle.Render("✗ " + m.error.Error())
	}

	wizardContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.progress.View(),
		"",
		stepContent,
	)

	if errorView != "" {
		wizardContent = lipgloss.JoinVertical(
			lipgloss.Left,
			wizardContent,
			"",
			errorView,
		)
	}

	return styles.BoxStyle.Width(70).Render(wizardContent)
}

// renderConfirmation renders the confirmation step
func (m ServiceCreateModel) renderConfirmation() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		MarginBottom(1)

	var details []string
	details = append(details, styles.RenderKeyValue("Name", m.nameInput.Value()))
	details = append(details, styles.RenderKeyValue("Type", m.typeSelect.Value()))
	details = append(details, styles.RenderKeyValue("Tier", m.tierSelect.Value()))
	if m.descriptionInput.Value() != "" {
		details = append(details, styles.RenderKeyValue("Description", m.descriptionInput.Value()))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, details...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Confirm Service Creation"),
		content,
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted)).Render("Press Enter to create, Esc to go back"),
	)
}

// createService creates the service
func (m *ServiceCreateModel) createService() tea.Cmd {
	m.loading = true
	m.error = nil

	name := m.nameInput.Value()
	serviceType := models.ServiceType(m.typeSelect.Value())
	tier := models.ServiceTier(m.tierSelect.Value())
	description := m.descriptionInput.Value()

	return func() tea.Msg {
		ctx := context.Background()

		req := api.CreateServiceRequest{
			Name:        name,
			Type:        serviceType,
			Tier:        tier,
			Description: description,
		}

		_, err := m.state.Client.CreateService(ctx, req)
		if err != nil {
			return serviceCreateErrorMsg{err: err}
		}

		return serviceCreatedMsg{}
	}
}

// serviceCreatedMsg is sent when service is created
type serviceCreatedMsg struct{}

// serviceCreateErrorMsg is sent when service creation fails
type serviceCreateErrorMsg struct {
	err error
}
