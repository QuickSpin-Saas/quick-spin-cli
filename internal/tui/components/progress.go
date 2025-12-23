package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quickspin/quickspin-cli/internal/tui/styles"
)

// ProgressBar is a component that displays a progress bar
type ProgressBar struct {
	progress progress.Model
	percent  float64
	width    int
	label    string
}

// NewProgressBar creates a new progress bar with default styling
func NewProgressBar(width int) ProgressBar {
	theme := styles.GetTheme()

	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(width),
	)

	// Set custom colors using the theme
	prog.FullColor = string(theme.Primary)
	prog.EmptyColor = string(theme.Border)

	return ProgressBar{
		progress: prog,
		percent:  0.0,
		width:    width,
		label:    "",
	}
}

// NewProgressBarWithLabel creates a new progress bar with a label
func NewProgressBarWithLabel(width int, label string) ProgressBar {
	pb := NewProgressBar(width)
	pb.label = label
	return pb
}

// Update updates the progress bar
func (pb ProgressBar) Update(msg tea.Msg) (ProgressBar, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := pb.progress.Update(msg)
	pb.progress = model.(progress.Model)
	return pb, cmd
}

// View renders the progress bar
func (pb ProgressBar) View() string {
	if pb.label != "" {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			pb.label,
			pb.progress.ViewAs(pb.percent),
		)
	}
	return pb.progress.ViewAs(pb.percent)
}

// SetPercent updates the progress percentage (0.0 to 1.0)
func (pb *ProgressBar) SetPercent(percent float64) {
	if percent < 0.0 {
		percent = 0.0
	}
	if percent > 1.0 {
		percent = 1.0
	}
	pb.percent = percent
}

// Increment increments the progress by a given amount
func (pb *ProgressBar) Increment(amount float64) {
	pb.SetPercent(pb.percent + amount)
}

// SetLabel sets the label for the progress bar
func (pb *ProgressBar) SetLabel(label string) {
	pb.label = label
}

// SetWidth sets the width of the progress bar
func (pb *ProgressBar) SetWidth(width int) {
	pb.width = width
	pb.progress.Width = width
}

// GetPercent returns the current progress percentage
func (pb *ProgressBar) GetPercent() float64 {
	return pb.percent
}

// IsComplete returns true if progress is at 100%
func (pb *ProgressBar) IsComplete() bool {
	return pb.percent >= 1.0
}

// SimpleProgressBar renders a simple text-based progress bar without Bubble Tea
func SimpleProgressBar(percent float64, width int) string {
	if percent < 0.0 {
		percent = 0.0
	}
	if percent > 1.0 {
		percent = 1.0
	}

	theme := styles.GetTheme()

	filled := int(float64(width) * percent)
	empty := width - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	percentText := fmt.Sprintf(" %.0f%%", percent*100)

	style := lipgloss.NewStyle().
		Foreground(theme.Primary)

	return style.Render(bar) + percentText
}

// MultiStepProgress shows progress through multiple steps
type MultiStepProgress struct {
	currentStep int
	totalSteps  int
	stepLabels  []string
	width       int
}

// NewMultiStepProgress creates a multi-step progress indicator
func NewMultiStepProgress(steps []string, width int) MultiStepProgress {
	return MultiStepProgress{
		currentStep: 0,
		totalSteps:  len(steps),
		stepLabels:  steps,
		width:       width,
	}
}

// NextStep advances to the next step
func (msp *MultiStepProgress) NextStep() {
	if msp.currentStep < msp.totalSteps-1 {
		msp.currentStep++
	}
}

// PrevStep goes back to the previous step
func (msp *MultiStepProgress) PrevStep() {
	if msp.currentStep > 0 {
		msp.currentStep--
	}
}

// SetStep sets the current step by index
func (msp *MultiStepProgress) SetStep(step int) {
	if step >= 0 && step < msp.totalSteps {
		msp.currentStep = step
	}
}

// View renders the multi-step progress
func (msp MultiStepProgress) View() string {
	theme := styles.GetTheme()

	var steps []string
	for i, label := range msp.stepLabels {
		var stepStyle lipgloss.Style

		if i < msp.currentStep {
			// Completed step
			stepStyle = lipgloss.NewStyle().
				Foreground(theme.Success).
				Bold(true)
			steps = append(steps, stepStyle.Render("✓ "+label))
		} else if i == msp.currentStep {
			// Current step
			stepStyle = lipgloss.NewStyle().
				Foreground(theme.Primary).
				Bold(true)
			steps = append(steps, stepStyle.Render("▶ "+label))
		} else {
			// Future step
			stepStyle = lipgloss.NewStyle().
				Foreground(theme.Muted)
			steps = append(steps, stepStyle.Render("  "+label))
		}
	}

	// Add progress indicator
	percent := float64(msp.currentStep) / float64(msp.totalSteps-1)
	progressBar := SimpleProgressBar(percent, msp.width)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		strings.Join(steps, "\n"),
		"",
		progressBar,
	)
}

// GetCurrentStep returns the current step index
func (msp MultiStepProgress) GetCurrentStep() int {
	return msp.currentStep
}

// IsComplete returns true if all steps are complete
func (msp MultiStepProgress) IsComplete() bool {
	return msp.currentStep >= msp.totalSteps-1
}
