package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// ModalType represents different types of modals
type ModalType int

const (
	ModalInfo ModalType = iota
	ModalSuccess
	ModalWarning
	ModalError
	ModalConfirm
)

// Modal is a dialog box overlay component
type Modal struct {
	title        string
	message      string
	modalType    ModalType
	buttons      []string
	selectedBtn  int
	visible      bool
	width        int
	height       int
	onConfirm    func()
	onCancel     func()
}

// NewModal creates a new modal dialog
func NewModal(title, message string, modalType ModalType) Modal {
	buttons := []string{"OK"}
	if modalType == ModalConfirm {
		buttons = []string{"Confirm", "Cancel"}
	}

	return Modal{
		title:       title,
		message:     message,
		modalType:   modalType,
		buttons:     buttons,
		selectedBtn: 0,
		visible:     false,
		width:       50,
		height:      10,
	}
}

// NewConfirmModal creates a confirmation modal with custom buttons
func NewConfirmModal(title, message string, confirmText, cancelText string) Modal {
	return Modal{
		title:       title,
		message:     message,
		modalType:   ModalConfirm,
		buttons:     []string{confirmText, cancelText},
		selectedBtn: 0,
		visible:     false,
		width:       50,
		height:      10,
	}
}

// Update updates the modal
func (m Modal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			if m.selectedBtn > 0 {
				m.selectedBtn--
			}
		case "right", "l":
			if m.selectedBtn < len(m.buttons)-1 {
				m.selectedBtn++
			}
		case "tab":
			m.selectedBtn = (m.selectedBtn + 1) % len(m.buttons)
		case "enter", " ":
			// Execute callback based on selected button
			if m.selectedBtn == 0 && m.onConfirm != nil {
				m.onConfirm()
			} else if m.selectedBtn > 0 && m.onCancel != nil {
				m.onCancel()
			}
			m.visible = false
		case "esc":
			if m.onCancel != nil {
				m.onCancel()
			}
			m.visible = false
		}
	}

	return m, nil
}

// View renders the modal
func (m Modal) View() string {
	if !m.visible {
		return ""
	}

	theme := styles.GetTheme()

	// Determine border and title color based on modal type
	var borderColor, titleColor lipgloss.Color
	var icon string

	switch m.modalType {
	case ModalInfo:
		borderColor = lipgloss.Color(theme.Info)
		titleColor = lipgloss.Color(theme.Info)
		icon = "ℹ"
	case ModalSuccess:
		borderColor = lipgloss.Color(theme.Success)
		titleColor = lipgloss.Color(theme.Success)
		icon = "✓"
	case ModalWarning:
		borderColor = lipgloss.Color(theme.Warning)
		titleColor = lipgloss.Color(theme.Warning)
		icon = "⚠"
	case ModalError:
		borderColor = lipgloss.Color(theme.Error)
		titleColor = lipgloss.Color(theme.Error)
		icon = "✗"
	case ModalConfirm:
		borderColor = lipgloss.Color(theme.Primary)
		titleColor = lipgloss.Color(theme.Primary)
		icon = "?"
	}

	// Title style
	titleStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true).
		Padding(0, 1)

	// Message style
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground)).
		Padding(1, 2).
		Width(m.width - 4)

	// Buttons
	var buttonViews []string
	for i, btnText := range m.buttons {
		var btnStyle lipgloss.Style
		if i == m.selectedBtn {
			btnStyle = styles.ButtonStyle
		} else {
			btnStyle = styles.SecondaryButtonStyle
		}
		buttonViews = append(buttonViews, btnStyle.Render(btnText))
	}

	buttonsRow := lipgloss.JoinHorizontal(lipgloss.Left, buttonViews...)

	// Combine all parts
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(icon+" "+m.title),
		messageStyle.Render(m.message),
		"",
		lipgloss.NewStyle().Padding(0, 2).Render(buttonsRow),
	)

	// Box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(m.width)

	return boxStyle.Render(content)
}

// Show displays the modal
func (m *Modal) Show() {
	m.visible = true
}

// Hide hides the modal
func (m *Modal) Hide() {
	m.visible = false
}

// IsVisible returns true if the modal is visible
func (m Modal) IsVisible() bool {
	return m.visible
}

// SetOnConfirm sets the callback for the confirm button
func (m *Modal) SetOnConfirm(callback func()) {
	m.onConfirm = callback
}

// SetOnCancel sets the callback for the cancel button
func (m *Modal) SetOnCancel(callback func()) {
	m.onCancel = callback
}

// SetMessage updates the modal message
func (m *Modal) SetMessage(message string) {
	m.message = message
}

// SetTitle updates the modal title
func (m *Modal) SetTitle(title string) {
	m.title = title
}

// RenderOverlay renders the modal centered over background content
func (m Modal) RenderOverlay(background string, termWidth, termHeight int) string {
	if !m.visible {
		return background
	}

	modalView := m.View()

	// Calculate modal dimensions
	modalLines := strings.Split(modalView, "\n")
	modalHeight := len(modalLines)
	modalWidth := 0
	for _, line := range modalLines {
		if len(line) > modalWidth {
			modalWidth = len(line)
		}
	}

	// Center the modal
	verticalPadding := (termHeight - modalHeight) / 2
	horizontalPadding := (termWidth - modalWidth) / 2

	// Create overlay
	overlay := strings.Repeat("\n", verticalPadding)
	for _, line := range modalLines {
		overlay += strings.Repeat(" ", horizontalPadding) + line + "\n"
	}

	return overlay
}
