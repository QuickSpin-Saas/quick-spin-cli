package types

// ViewType represents different views/screens in the TUI
type ViewType int

const (
	ViewDashboard ViewType = iota
	ViewAuthMenu
	ViewAuthLogin
	ViewAuthLogout
	ViewAuthWhoami
	ViewServiceList
	ViewServiceCreate
	ViewServiceDetail
	ViewServiceLogs
	ViewConfigMenu
	ViewConfigEditor
	ViewConfigView
	ViewHelp
	ViewExit
)

// String returns the string representation of a ViewType
func (v ViewType) String() string {
	switch v {
	case ViewDashboard:
		return "Dashboard"
	case ViewAuthMenu:
		return "Authentication Menu"
	case ViewAuthLogin:
		return "Login"
	case ViewAuthLogout:
		return "Logout"
	case ViewAuthWhoami:
		return "Current User"
	case ViewServiceList:
		return "Services"
	case ViewServiceCreate:
		return "Create Service"
	case ViewServiceDetail:
		return "Service Details"
	case ViewServiceLogs:
		return "Service Logs"
	case ViewConfigMenu:
		return "Configuration Menu"
	case ViewConfigEditor:
		return "Config Editor"
	case ViewConfigView:
		return "View Configuration"
	case ViewHelp:
		return "Help"
	case ViewExit:
		return "Exit"
	default:
		return "Unknown"
	}
}

// NavigationMsg is sent when navigating to a new view
type NavigationMsg struct {
	View ViewType
	Data interface{} // Optional data to pass to the new view
}

// BackMsg is sent when navigating back
type BackMsg struct{}

// ExitMsg is sent when exiting the application
type ExitMsg struct{}
