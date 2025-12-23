package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// TextInput is a single-line text input field
type TextInput struct {
	input       textinput.Model
	label       string
	placeholder string
	width       int
	focused     bool
	err         error
	validator   func(string) error
}

// NewTextInput creates a new text input field
func NewTextInput(label, placeholder string, width int) TextInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = width
	ti.CharLimit = 256

	theme := styles.GetTheme()
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Foreground))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted))

	return TextInput{
		input:       ti,
		label:       label,
		placeholder: placeholder,
		width:       width,
		focused:     false,
	}
}

// Init initializes the text input
func (t TextInput) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the text input
func (t TextInput) Update(msg tea.Msg) (TextInput, tea.Cmd) {
	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)

	// Run validator if provided
	if t.validator != nil {
		t.err = t.validator(t.input.Value())
	}

	return t, cmd
}

// View renders the text input
func (t TextInput) View() string {
	theme := styles.GetTheme()

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent)).
		Bold(true)

	var view string
	if t.label != "" {
		view = labelStyle.Render(t.label) + "\n"
	}

	// Add input field
	view += t.input.View()

	// Show error if validation failed
	if t.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Error)).
			Italic(true)
		view += "\n" + errorStyle.Render("✗ "+t.err.Error())
	}

	return view
}

// Focus gives focus to the input
func (t *TextInput) Focus() tea.Cmd {
	t.focused = true
	return t.input.Focus()
}

// Blur removes focus from the input
func (t *TextInput) Blur() {
	t.focused = false
	t.input.Blur()
}

// SetValue sets the input value
func (t *TextInput) SetValue(value string) {
	t.input.SetValue(value)
}

// Value returns the current input value
func (t TextInput) Value() string {
	return t.input.Value()
}

// SetValidator sets a validation function
func (t *TextInput) SetValidator(validator func(string) error) {
	t.validator = validator
}

// IsValid returns true if the input is valid
func (t TextInput) IsValid() bool {
	return t.err == nil
}

// PasswordInput is a password input field with masked characters
type PasswordInput struct {
	input       textinput.Model
	label       string
	placeholder string
	width       int
	focused     bool
}

// NewPasswordInput creates a new password input field
func NewPasswordInput(label, placeholder string, width int) PasswordInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = width
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	theme := styles.GetTheme()
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Primary))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Foreground))
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Muted))

	return PasswordInput{
		input:       ti,
		label:       label,
		placeholder: placeholder,
		width:       width,
		focused:     false,
	}
}

// Init initializes the password input
func (p PasswordInput) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the password input
func (p PasswordInput) Update(msg tea.Msg) (PasswordInput, tea.Cmd) {
	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	return p, cmd
}

// View renders the password input
func (p PasswordInput) View() string {
	theme := styles.GetTheme()

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent)).
		Bold(true)

	var view string
	if p.label != "" {
		view = labelStyle.Render(p.label) + "\n"
	}

	view += p.input.View()

	return view
}

// Focus gives focus to the input
func (p *PasswordInput) Focus() tea.Cmd {
	p.focused = true
	return p.input.Focus()
}

// Blur removes focus from the input
func (p *PasswordInput) Blur() {
	p.focused = false
	p.input.Blur()
}

// Value returns the current password value
func (p PasswordInput) Value() string {
	return p.input.Value()
}

// SetValue sets the password value
func (p *PasswordInput) SetValue(value string) {
	p.input.SetValue(value)
}

// SelectInput is a dropdown/selection field
type SelectInput struct {
	label    string
	options  []string
	selected int
	focused  bool
	width    int
	expanded bool
}

// NewSelectInput creates a new select input
func NewSelectInput(label string, options []string, width int) SelectInput {
	return SelectInput{
		label:    label,
		options:  options,
		selected: 0,
		focused:  false,
		width:    width,
		expanded: false,
	}
}

// Update updates the select input
func (s SelectInput) Update(msg tea.Msg) (SelectInput, tea.Cmd) {
	if !s.focused {
		return s, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.selected > 0 {
				s.selected--
			}
		case "down", "j":
			if s.selected < len(s.options)-1 {
				s.selected++
			}
		case "enter", " ":
			s.expanded = !s.expanded
		}
	}

	return s, nil
}

// View renders the select input
func (s SelectInput) View() string {
	theme := styles.GetTheme()

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent)).
		Bold(true)

	var view string
	if s.label != "" {
		view = labelStyle.Render(s.label) + "\n"
	}

	// Show selected option
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary)).
		Bold(true)

	if !s.expanded {
		view += selectedStyle.Render("▼ " + s.options[s.selected])
	} else {
		// Show all options
		for i, option := range s.options {
			if i == s.selected {
				view += selectedStyle.Render("▶ " + option) + "\n"
			} else {
				mutedStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(theme.Muted))
				view += mutedStyle.Render("  " + option) + "\n"
			}
		}
	}

	return view
}

// Focus gives focus to the select input
func (s *SelectInput) Focus() {
	s.focused = true
}

// Blur removes focus from the select input
func (s *SelectInput) Blur() {
	s.focused = false
	s.expanded = false
}

// Value returns the currently selected option
func (s SelectInput) Value() string {
	if s.selected >= 0 && s.selected < len(s.options) {
		return s.options[s.selected]
	}
	return ""
}

// SelectedIndex returns the currently selected index
func (s SelectInput) SelectedIndex() int {
	return s.selected
}

// SetSelected sets the selected option by index
func (s *SelectInput) SetSelected(index int) {
	if index >= 0 && index < len(s.options) {
		s.selected = index
	}
}
