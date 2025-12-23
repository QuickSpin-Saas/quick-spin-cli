package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// StatusBar is a bottom status bar component showing keyboard shortcuts and info
type StatusBar struct {
	leftText   string
	centerText string
	rightText  string
	width      int
}

// NewStatusBar creates a new status bar
func NewStatusBar(width int) StatusBar {
	return StatusBar{
		width: width,
	}
}

// View renders the status bar
func (s StatusBar) View() string {
	theme := styles.GetTheme()

	// Calculate spacing
	leftWidth := lipgloss.Width(s.leftText)
	centerWidth := lipgloss.Width(s.centerText)
	rightWidth := lipgloss.Width(s.rightText)

	// Calculate padding
	totalContent := leftWidth + centerWidth + rightWidth
	if totalContent >= s.width {
		// Not enough space, truncate
		return styles.StatusBarStyle.Width(s.width).Render(s.leftText)
	}

	// Calculate spacing to distribute content
	leftPadding := (s.width - totalContent) / 2
	rightPadding := s.width - totalContent - leftPadding

	// Style for different sections
	leftStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground)).
		Bold(true)

	centerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Muted))

	rightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent))

	// Build the status bar
	left := leftStyle.Render(s.leftText)
	center := centerStyle.Render(s.centerText)
	right := rightStyle.Render(s.rightText)

	// Join with spacing
	content := left +
		lipgloss.NewStyle().Width(leftPadding).Render("") +
		center +
		lipgloss.NewStyle().Width(rightPadding).Render("") +
		right

	return styles.StatusBarStyle.Width(s.width).Render(content)
}

// SetLeft sets the left text
func (s *StatusBar) SetLeft(text string) {
	s.leftText = text
}

// SetCenter sets the center text
func (s *StatusBar) SetCenter(text string) {
	s.centerText = text
}

// SetRight sets the right text
func (s *StatusBar) SetRight(text string) {
	s.rightText = text
}

// SetWidth sets the status bar width
func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

// RenderHelp renders a help status bar with keyboard shortcuts
func RenderHelp(shortcuts map[string]string, width int) string {
	theme := styles.GetTheme()

	var items []string
	for key, desc := range shortcuts {
		keyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Accent)).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Subtle))

		item := keyStyle.Render(key) + descStyle.Render(": "+desc)
		items = append(items, item)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Left, items...)

	return styles.StatusBarStyle.Width(width).Render(content)
}
