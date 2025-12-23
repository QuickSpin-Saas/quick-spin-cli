package styles

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the TUI
type Theme struct {
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Info       lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
	Border     lipgloss.Color
	Highlight  lipgloss.Color
	Muted      lipgloss.Color
	Subtle     lipgloss.Color
}

// DefaultTheme provides a modern, vibrant color scheme
var DefaultTheme = Theme{
	Primary:    lipgloss.Color("#00D9FF"), // Bright Cyan - main brand color
	Secondary:  lipgloss.Color("#8B5CF6"), // Vibrant Purple - secondary actions
	Accent:     lipgloss.Color("#F59E0B"), // Amber - highlights and accents
	Success:    lipgloss.Color("#10B981"), // Emerald Green - success states
	Warning:    lipgloss.Color("#F97316"), // Orange - warnings
	Error:      lipgloss.Color("#EF4444"), // Red - errors
	Info:       lipgloss.Color("#3B82F6"), // Blue - informational
	Background: lipgloss.Color("#0A0A0A"), // Very Dark - main background
	Foreground: lipgloss.Color("#E5E7EB"), // Light Gray - main text
	Border:     lipgloss.Color("#374151"), // Medium Gray - borders
	Highlight:  lipgloss.Color("#1F2937"), // Dark Gray - highlights/hover
	Muted:      lipgloss.Color("#6B7280"), // Gray - muted text
	Subtle:     lipgloss.Color("#9CA3AF"), // Light Gray - subtle text
}

// CurrentTheme is the active theme (can be switched dynamically)
var CurrentTheme = DefaultTheme

// GetTheme returns the current theme
func GetTheme() Theme {
	return CurrentTheme
}

// SetTheme changes the active theme
func SetTheme(theme Theme) {
	CurrentTheme = theme
}
