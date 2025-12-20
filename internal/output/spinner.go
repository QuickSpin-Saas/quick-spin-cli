package output

import (
	"time"

	"github.com/briandowns/spinner"
)

// Spinner wraps the spinner library
type Spinner struct {
	s *spinner.Spinner
}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " " + message
	return &Spinner{s: s}
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.s.Start()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.s.Stop()
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.s.Suffix = " " + message
}

// Success stops the spinner and shows a success message
func (s *Spinner) Success(message string) {
	s.s.Stop()
	Success(message)
}

// Fail stops the spinner and shows an error message
func (s *Spinner) Fail(message string) {
	s.s.Stop()
	Error(message)
}
