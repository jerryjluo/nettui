package model

import "github.com/jerryluo/nettui/internal/data"

// DataRefreshMsg is sent when new data is available.
type DataRefreshMsg struct {
	Store *data.Store
}

// CrossRefMsg requests navigation to a different tab with a filter/highlight.
type CrossRefMsg struct {
	TargetTab TabID
	FilterKey string
	FilterVal string
}

// ToggleSidePanelMsg toggles the detail side panel.
type ToggleSidePanelMsg struct{}

// CloseSidePanelMsg closes the detail side panel.
type CloseSidePanelMsg struct{}

// ToggleDNSMsg toggles DNS resolution.
type ToggleDNSMsg struct{}

// CopyMsg requests copying selected data to clipboard.
type CopyMsg struct {
	Content string
}

// CopyResultMsg reports the result of a clipboard copy.
type CopyResultMsg struct {
	Success bool
	Error   string
}

// HelpToggleMsg toggles the help overlay.
type HelpToggleMsg struct{}

// WindowSizeMsg carries terminal dimensions.
type WindowSizeMsg struct {
	Width  int
	Height int
}

// ErrMsg carries a non-fatal error for display.
type ErrMsg struct {
	Error string
}
