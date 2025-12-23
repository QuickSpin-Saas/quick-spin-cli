package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	theme = GetTheme()

	// Border styles
	RoundedBorder = lipgloss.RoundedBorder()
	ThickBorder   = lipgloss.ThickBorder()
	DoubleBorder  = lipgloss.DoubleBorder()

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(theme.Secondary).
			Padding(0, 1)

	H1Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary).
		MarginBottom(1)

	H2Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Secondary).
		MarginBottom(1)

	H3Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Accent)

	// Text styles
	BoldStyle = lipgloss.NewStyle().
			Bold(true)

	ItalicStyle = lipgloss.NewStyle().
			Italic(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(theme.Muted)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(theme.Subtle)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(theme.Info).
			Bold(true)

	// Box/Container styles
	BoxStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(theme.Border).
			Padding(1, 2)

	FocusedBoxStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(theme.Primary).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(theme.Border).
			Padding(0, 1).
			MarginBottom(1)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(theme.Background).
			Background(theme.Primary).
			Padding(0, 3).
			MarginRight(2)

	SecondaryButtonStyle = lipgloss.NewStyle().
				Foreground(theme.Foreground).
				Background(theme.Secondary).
				Padding(0, 3).
				MarginRight(2)

	DisabledButtonStyle = lipgloss.NewStyle().
				Foreground(theme.Muted).
				Background(theme.Highlight).
				Padding(0, 3).
				MarginRight(2)

	// List styles
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	SelectedListItemStyle = lipgloss.NewStyle().
				Foreground(theme.Primary).
				Background(theme.Highlight).
				Padding(0, 2).
				Bold(true)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Border(RoundedBorder).
			BorderForeground(theme.Border).
			Padding(0, 1)

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(theme.Foreground).
				Border(RoundedBorder).
				BorderForeground(theme.Primary).
				Padding(0, 1)

	// Badge styles
	BadgeStyle = lipgloss.NewStyle().
			Foreground(theme.Background).
			Background(theme.Accent).
			Padding(0, 1).
			MarginRight(1)

	SuccessBadgeStyle = lipgloss.NewStyle().
				Foreground(theme.Background).
				Background(theme.Success).
				Padding(0, 1).
				MarginRight(1)

	ErrorBadgeStyle = lipgloss.NewStyle().
			Foreground(theme.Background).
			Background(theme.Error).
			Padding(0, 1).
			MarginRight(1)

	InfoBadgeStyle = lipgloss.NewStyle().
			Foreground(theme.Background).
			Background(theme.Info).
			Padding(0, 1).
			MarginRight(1)

	// Status bar style
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Background(theme.Highlight).
			Padding(0, 1)

	// Help/keyboard shortcut styles
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(theme.Accent).
			Bold(true)

	HelpTextStyle = lipgloss.NewStyle().
			Foreground(theme.Subtle)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.Primary).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(theme.Border)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(theme.Foreground)

	TableSelectedRowStyle = lipgloss.NewStyle().
				Foreground(theme.Primary).
				Background(theme.Highlight).
				Bold(true)

	// Divider/separator
	DividerStyle = lipgloss.NewStyle().
			Foreground(theme.Border)

	// Error message box
	ErrorBoxStyle = lipgloss.NewStyle().
			Border(RoundedBorder).
			BorderForeground(theme.Error).
			Foreground(theme.Error).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Loading/spinner styles
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(theme.Primary)

	// Code/monospace styles
	CodeStyle = lipgloss.NewStyle().
			Foreground(theme.Accent).
			Background(theme.Highlight).
			Padding(0, 1)
)

// RenderKeyValue renders a key-value pair with consistent styling
func RenderKeyValue(key, value string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true).
		Width(20)

	valueStyle := lipgloss.NewStyle().
		Foreground(theme.Foreground)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		keyStyle.Render(key+":"),
		valueStyle.Render(value),
	)
}

// RenderBadge renders a status badge with the given text and style
func RenderBadge(text string, badgeStyle lipgloss.Style) string {
	return badgeStyle.Render(" " + text + " ")
}

// RenderHelp renders a keyboard shortcut help item
func RenderHelp(key, description string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		HelpKeyStyle.Render(key),
		HelpTextStyle.Render(": "+description+"  "),
	)
}

// RenderDivider renders a horizontal divider line
func RenderDivider(width int) string {
	divider := ""
	for i := 0; i < width; i++ {
		divider += "─"
	}
	return DividerStyle.Render(divider)
}

// CenterText centers text within a given width
func CenterText(text string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(text)
}

// RefreshTheme updates all styles with the current theme
// Call this after changing themes
func RefreshTheme() {
	theme = GetTheme()

	// Update styles that depend on theme colors
	TitleStyle = TitleStyle.Foreground(theme.Primary)
	SubtitleStyle = SubtitleStyle.Foreground(theme.Secondary)
	H1Style = H1Style.Foreground(theme.Primary)
	H2Style = H2Style.Foreground(theme.Secondary)
	H3Style = H3Style.Foreground(theme.Accent)

	SuccessStyle = SuccessStyle.Foreground(theme.Success)
	ErrorStyle = ErrorStyle.Foreground(theme.Error)
	WarningStyle = WarningStyle.Foreground(theme.Warning)
	InfoStyle = InfoStyle.Foreground(theme.Info)

	BoxStyle = BoxStyle.BorderForeground(theme.Border)
	FocusedBoxStyle = FocusedBoxStyle.BorderForeground(theme.Primary)
	PanelStyle = PanelStyle.BorderForeground(theme.Border)

	ButtonStyle = ButtonStyle.Foreground(theme.Background).Background(theme.Primary)
	SecondaryButtonStyle = SecondaryButtonStyle.Foreground(theme.Foreground).Background(theme.Secondary)

	SelectedListItemStyle = SelectedListItemStyle.Foreground(theme.Primary).Background(theme.Highlight)

	InputStyle = InputStyle.BorderForeground(theme.Border)
	FocusedInputStyle = FocusedInputStyle.BorderForeground(theme.Primary)

	BadgeStyle = BadgeStyle.Foreground(theme.Background).Background(theme.Accent)
	SuccessBadgeStyle = SuccessBadgeStyle.Foreground(theme.Background).Background(theme.Success)
	ErrorBadgeStyle = ErrorBadgeStyle.Foreground(theme.Background).Background(theme.Error)
	InfoBadgeStyle = InfoBadgeStyle.Foreground(theme.Background).Background(theme.Info)

	StatusBarStyle = StatusBarStyle.Foreground(theme.Foreground).Background(theme.Highlight)

	HelpKeyStyle = HelpKeyStyle.Foreground(theme.Accent)

	TableHeaderStyle = TableHeaderStyle.Foreground(theme.Primary).BorderForeground(theme.Border)
	TableSelectedRowStyle = TableSelectedRowStyle.Foreground(theme.Primary).Background(theme.Highlight)

	SpinnerStyle = SpinnerStyle.Foreground(theme.Primary)
}
