package model

// TabID identifies a tab.
type TabID int

const (
	TabSockets TabID = iota
	TabUnixSockets
	TabProcesses
	TabInterfaces
	TabRoutes
	TabFirewall
)

// TabCount is the total number of tabs.
const TabCount = 6

// TabName returns the display name for a tab.
func TabName(id TabID) string {
	switch id {
	case TabInterfaces:
		return "Interfaces"
	case TabRoutes:
		return "Routes"
	case TabSockets:
		return "Sockets"
	case TabUnixSockets:
		return "Unix"
	case TabProcesses:
		return "Processes"
	case TabFirewall:
		return "Firewall"
	default:
		return "Unknown"
	}
}
