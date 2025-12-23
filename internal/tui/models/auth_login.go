package models

import (
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/models"
	"github.com/quickspin/quickspin-cli/internal/tui/types"
	"github.com/quickspin/quickspin-cli/internal/tui/components"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// AuthLoginModel is the login form screen
type AuthLoginModel struct {
	state         *types.AppState
	emailInput    components.TextInput
	passwordInput components.PasswordInput
	focusedField  int // 0 = email, 1 = password, 2 = submit button
	loading       bool
	error         error
	spinner       components.Spinner
	width         int
	height        int
	statusBar     components.StatusBar
}

// NewAuthLoginModel creates a new login form model
func NewAuthLoginModel(state *types.AppState) AuthLoginModel {
	emailInput := components.NewTextInput("Email", "your@email.com", 40)
	passwordInput := components.NewPasswordInput("Password", "Enter your password", 40)

	return AuthLoginModel{
		state:         state,
		emailInput:    emailInput,
		passwordInput: passwordInput,
		focusedField:  0,
		loading:       false,
		error:         nil,
		spinner:       components.NewSpinner("Logging in..."),
		width:         state.Width,
		height:        state.Height,
		statusBar:     components.NewStatusBar(state.Width),
	}
}

// Init initializes the login form
func (m AuthLoginModel) Init() tea.Cmd {
	return m.emailInput.Focus()
}

// Update handles messages and updates the login form
func (m AuthLoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// Go back to auth menu
			return m, func() tea.Msg {
				return types.NavigationMsg{View: types.ViewAuthMenu}
			}

		case "tab", "down":
			// Move focus to next field
			m.focusedField = (m.focusedField + 1) % 3
			m.updateFocus()

		case "shift+tab", "up":
			// Move focus to previous field
			m.focusedField = (m.focusedField - 1 + 3) % 3
			m.updateFocus()

		case "enter":
			if m.focusedField == 2 || (m.focusedField == 1 && len(m.passwordInput.Value()) > 0) {
				// Submit form
				return m, m.submitLogin()
			} else {
				// Move to next field
				m.focusedField = (m.focusedField + 1) % 3
				m.updateFocus()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.state.SetTerminalSize(msg.Width, msg.Height)
		m.statusBar.SetWidth(msg.Width)

	case loginSuccessMsg:
		m.loading = false
		m.state.SetUser(msg.user)
		// Navigate to dashboard
		return m, func() tea.Msg {
			return types.NavigationMsg{View: types.ViewDashboard}
		}

	case loginErrorMsg:
		m.loading = false
		m.error = msg.err
	}

	// Update active input
	if !m.loading {
		switch m.focusedField {
		case 0:
			m.emailInput, cmd = m.emailInput.Update(msg)
			cmds = append(cmds, cmd)
		case 1:
			m.passwordInput, cmd = m.passwordInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Update spinner if loading
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the login form
func (m AuthLoginModel) View() string {
	// Header
	header := m.renderHeader()

	// Form
	form := m.renderForm()

	// Status bar
	m.statusBar.SetLeft("Login")
	m.statusBar.SetCenter(m.state.Router.PrintBreadcrumb())
	m.statusBar.SetRight("Tab: Next Field • Enter: Submit • Esc: Back • q: Quit")
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
		contentStyle.Render(form),
		statusBar,
	)

	return view
}

// renderHeader renders the login form header
func (m AuthLoginModel) renderHeader() string {
	theme := styles.GetTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true).
		Padding(0, 2)

	header := titleStyle.Render("🔑 Login to QuickSpin")

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

// renderForm renders the login form fields
func (m AuthLoginModel) renderForm() string {
	if m.loading {
		// Show loading spinner
		return lipgloss.NewStyle().
			Padding(2).
			Render(m.spinner.View())
	}

	// Email input
	emailView := m.emailInput.View()

	// Password input
	passwordView := m.passwordInput.View()

	// Submit button
	var submitButtonStyle lipgloss.Style
	if m.focusedField == 2 {
		submitButtonStyle = styles.ButtonStyle
	} else {
		submitButtonStyle = styles.SecondaryButtonStyle
	}
	submitButton := submitButtonStyle.Render("Login")

	// Error message
	var errorView string
	if m.error != nil {
		errorView = styles.ErrorBoxStyle.Render("✗ " + m.error.Error())
	}

	// Combine all elements
	formContent := lipgloss.JoinVertical(
		lipgloss.Left,
		emailView,
		"",
		passwordView,
		"",
		submitButton,
	)

	if errorView != "" {
		formContent = lipgloss.JoinVertical(
			lipgloss.Left,
			formContent,
			"",
			errorView,
		)
	}

	// Add box around form
	formBox := styles.BoxStyle.
		Width(60).
		Render(formContent)

	return formBox
}

// updateFocus updates which field has focus
func (m *AuthLoginModel) updateFocus() {
	switch m.focusedField {
	case 0:
		m.emailInput.Focus()
		m.passwordInput.Blur()
	case 1:
		m.emailInput.Blur()
		m.passwordInput.Focus()
	case 2:
		m.emailInput.Blur()
		m.passwordInput.Blur()
	}
}

// submitLogin submits the login form
func (m *AuthLoginModel) submitLogin() tea.Cmd {
	email := strings.TrimSpace(m.emailInput.Value())
	password := m.passwordInput.Value()

	// Validate inputs
	if email == "" || password == "" {
		m.error = &models.APIError{
			Message: "Email and password are required",
		}
		return nil
	}

	// Clear previous error
	m.error = nil
	m.loading = true

	return func() tea.Msg {
		ctx := context.Background()

		// Perform login
		result, err := m.state.Client.Login(ctx, email, password)
		if err != nil {
			return loginErrorMsg{err: err}
		}

		return loginSuccessMsg{user: &result.User}
	}
}

// loginSuccessMsg is sent when login succeeds
type loginSuccessMsg struct {
	user *models.User
}

// loginErrorMsg is sent when login fails
type loginErrorMsg struct {
	err error
}
