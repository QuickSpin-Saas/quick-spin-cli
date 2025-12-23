package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// Spinner is a component that displays an animated loading spinner
type Spinner struct {
	spinner spinner.Model
	message string
	style   lipgloss.Style
}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	return Spinner{
		spinner: s,
		message: message,
		style:   styles.SpinnerStyle,
	}
}

// NewSpinnerWithStyle creates a new spinner with a custom style
func NewSpinnerWithStyle(message string, spinnerType spinner.Spinner) Spinner {
	s := spinner.New()
	s.Spinner = spinnerType
	s.Style = styles.SpinnerStyle

	return Spinner{
		spinner: s,
		message: message,
		style:   styles.SpinnerStyle,
	}
}

// Init initializes the spinner
func (s Spinner) Init() tea.Cmd {
	return s.spinner.Tick
}

// Update updates the spinner
func (s Spinner) Update(msg tea.Msg) (Spinner, tea.Cmd) {
	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

// View renders the spinner
func (s Spinner) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		s.spinner.View(),
		" ",
		s.message,
	)
}

// SetMessage updates the spinner message
func (s *Spinner) SetMessage(message string) {
	s.message = message
}

// SetStyle updates the spinner style
func (s *Spinner) SetStyle(style lipgloss.Style) {
	s.style = style
	s.spinner.Style = style
}

// Available spinner types for easy access
var (
	SpinnerDot       = spinner.Dot
	SpinnerLine      = spinner.Line
	SpinnerMiniDot   = spinner.MiniDot
	SpinnerJump      = spinner.Jump
	SpinnerPulse     = spinner.Pulse
	SpinnerPoints    = spinner.Points
	SpinnerGlobe     = spinner.Globe
	SpinnerMoon      = spinner.Moon
	SpinnerMonkey    = spinner.Monkey
	SpinnerMeter     = spinner.Meter
	SpinnerHamburger = spinner.Hamburger
)
