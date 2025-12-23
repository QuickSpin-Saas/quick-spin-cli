package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// ToastType represents different types of toast notifications
type ToastType int

const (
	ToastInfo ToastType = iota
	ToastSuccess
	ToastWarning
	ToastError
)

// Toast is a temporary notification component
type Toast struct {
	message   string
	toastType ToastType
	visible   bool
	duration  time.Duration
}

// NewToast creates a new toast notification
func NewToast(message string, toastType ToastType, duration time.Duration) Toast {
	return Toast{
		message:   message,
		toastType: toastType,
		visible:   false,
		duration:  duration,
	}
}

// Show displays the toast and starts the timer
func (t *Toast) Show() tea.Cmd {
	t.visible = true
	return tea.Tick(t.duration, func(time.Time) tea.Msg {
		return toastTimeoutMsg{}
	})
}

// Hide hides the toast
func (t *Toast) Hide() {
	t.visible = false
}

// IsVisible returns true if the toast is visible
func (t Toast) IsVisible() bool {
	return t.visible
}

// SetMessage updates the toast message
func (t *Toast) SetMessage(message string) {
	t.message = message
}

// Update updates the toast
func (t Toast) Update(msg tea.Msg) (Toast, tea.Cmd) {
	switch msg.(type) {
	case toastTimeoutMsg:
		t.visible = false
	}
	return t, nil
}

// View renders the toast
func (t Toast) View() string {
	if !t.visible {
		return ""
	}

	theme := styles.GetTheme()

	var style lipgloss.Style
	var icon string

	switch t.toastType {
	case ToastInfo:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Background)).
			Background(lipgloss.Color(theme.Info)).
			Padding(0, 2).
			Bold(true)
		icon = "ℹ"
	case ToastSuccess:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Background)).
			Background(lipgloss.Color(theme.Success)).
			Padding(0, 2).
			Bold(true)
		icon = "✓"
	case ToastWarning:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Background)).
			Background(lipgloss.Color(theme.Warning)).
			Padding(0, 2).
			Bold(true)
		icon = "⚠"
	case ToastError:
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Background)).
			Background(lipgloss.Color(theme.Error)).
			Padding(0, 2).
			Bold(true)
		icon = "✗"
	}

	return style.Render(icon + " " + t.message)
}

// toastTimeoutMsg is sent when the toast should be hidden
type toastTimeoutMsg struct{}

// ToastManager manages multiple toast notifications
type ToastManager struct {
	toasts []Toast
	maxToasts int
}

// NewToastManager creates a new toast manager
func NewToastManager(maxToasts int) ToastManager {
	return ToastManager{
		toasts:    []Toast{},
		maxToasts: maxToasts,
	}
}

// Add adds a new toast notification
func (tm *ToastManager) Add(message string, toastType ToastType, duration time.Duration) tea.Cmd {
	toast := NewToast(message, toastType, duration)

	// Remove oldest toast if at capacity
	if len(tm.toasts) >= tm.maxToasts {
		tm.toasts = tm.toasts[1:]
	}

	tm.toasts = append(tm.toasts, toast)

	return toast.Show()
}

// Update updates all toasts
func (tm ToastManager) Update(msg tea.Msg) (ToastManager, tea.Cmd) {
	var cmds []tea.Cmd

	// Update all toasts
	for i := range tm.toasts {
		var cmd tea.Cmd
		tm.toasts[i], cmd = tm.toasts[i].Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Remove hidden toasts
	var activeToasts []Toast
	for _, toast := range tm.toasts {
		if toast.IsVisible() {
			activeToasts = append(activeToasts, toast)
		}
	}
	tm.toasts = activeToasts

	return tm, tea.Batch(cmds...)
}

// View renders all visible toasts
func (tm ToastManager) View() string {
	var views []string
	for _, toast := range tm.toasts {
		if toast.IsVisible() {
			views = append(views, toast.View())
		}
	}

	if len(views) == 0 {
		return ""
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

// Clear clears all toasts
func (tm *ToastManager) Clear() {
	tm.toasts = []Toast{}
}

// Helper functions for common toast notifications
func ShowSuccessToast(message string) Toast {
	return NewToast(message, ToastSuccess, 3*time.Second)
}

func ShowErrorToast(message string) Toast {
	return NewToast(message, ToastError, 5*time.Second)
}

func ShowInfoToast(message string) Toast {
	return NewToast(message, ToastInfo, 3*time.Second)
}

func ShowWarningToast(message string) Toast {
	return NewToast(message, ToastWarning, 4*time.Second)
}
