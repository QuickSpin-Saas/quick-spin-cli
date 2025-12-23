package tui

import (
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/tui/types"
)

// Re-export types for backward compatibility
type ViewType = types.ViewType

const (
	ViewDashboard     = types.ViewDashboard
	ViewAuthMenu      = types.ViewAuthMenu
	ViewAuthLogin     = types.ViewAuthLogin
	ViewAuthLogout    = types.ViewAuthLogout
	ViewAuthWhoami    = types.ViewAuthWhoami
	ViewServiceList   = types.ViewServiceList
	ViewServiceCreate = types.ViewServiceCreate
	ViewServiceDetail = types.ViewServiceDetail
	ViewServiceLogs   = types.ViewServiceLogs
	ViewConfigMenu    = types.ViewConfigMenu
	ViewConfigEditor  = types.ViewConfigEditor
	ViewConfigView    = types.ViewConfigView
	ViewHelp          = types.ViewHelp
	ViewExit          = types.ViewExit
)

// Router manages navigation between views using a stack-based approach
type Router struct {
	viewStack []ViewType
	current   ViewType
}

// NewRouter creates a new router starting at the dashboard
func NewRouter() *Router {
	return &Router{
		viewStack: []ViewType{ViewDashboard},
		current:   ViewDashboard,
	}
}

// NewRouterWithView creates a new router starting at a specific view
func NewRouterWithView(view ViewType) *Router {
	return &Router{
		viewStack: []ViewType{view},
		current:   view,
	}
}

// Push navigates to a new view by pushing it onto the stack
func (r *Router) Push(view ViewType) {
	r.viewStack = append(r.viewStack, view)
	r.current = view
}

// Pop navigates back to the previous view by popping from the stack
// Returns the view we're going back to
func (r *Router) Pop() ViewType {
	if len(r.viewStack) <= 1 {
		// Don't pop the last view (dashboard)
		return r.current
	}

	// Remove current view
	r.viewStack = r.viewStack[:len(r.viewStack)-1]

	// Set current to the previous view
	r.current = r.viewStack[len(r.viewStack)-1]

	return r.current
}

// Current returns the current view
func (r *Router) Current() ViewType {
	return r.current
}

// CanGoBack returns true if there are views to go back to
func (r *Router) CanGoBack() bool {
	return len(r.viewStack) > 1
}

// Reset clears the navigation stack and returns to dashboard
func (r *Router) Reset() {
	r.viewStack = []ViewType{ViewDashboard}
	r.current = ViewDashboard
}

// GetBreadcrumb returns a breadcrumb string showing the navigation path
func (r *Router) GetBreadcrumb() string {
	if len(r.viewStack) == 0 {
		return ""
	}

	breadcrumb := r.viewStack[0].String()
	for i := 1; i < len(r.viewStack); i++ {
		breadcrumb += " > " + r.viewStack[i].String()
	}

	return breadcrumb
}

// Re-export message types for backward compatibility
type NavigationMsg = types.NavigationMsg
type BackMsg = types.BackMsg
type ExitMsg = types.ExitMsg

// PrintBreadcrumb returns a formatted breadcrumb for display
func (r *Router) PrintBreadcrumb() string {
	if len(r.viewStack) == 0 {
		return ""
	}

	result := ""
	for i, view := range r.viewStack {
		if i > 0 {
			result += " › "
		}
		if i == len(r.viewStack)-1 {
			// Current view in bold/highlighted
			result += fmt.Sprintf("[%s]", view.String())
		} else {
			result += view.String()
		}
	}

	return result
}
