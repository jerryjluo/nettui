package model

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#7C3AED")
	SecondaryColor = lipgloss.Color("#06B6D4")
	AccentColor    = lipgloss.Color("#F59E0B")
	ErrorColor     = lipgloss.Color("#EF4444")
	SuccessColor   = lipgloss.Color("#10B981")
	MutedColor     = lipgloss.Color("#6B7280")
	BgColor        = lipgloss.Color("#1F2937")
	FgColor        = lipgloss.Color("#F9FAFB")

	// Tab styles
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(PrimaryColor).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(MutedColor).
				Padding(0, 2)

	TabWarningStyle = lipgloss.NewStyle().
			Foreground(AccentColor)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Padding(0, 1)

	StatusBadgeStyle = lipgloss.NewStyle().
				Foreground(ErrorColor).
				Bold(true)

	// Side panel
	PanelBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(PrimaryColor).
				Padding(0, 1)

	PanelHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(SecondaryColor)

	PanelLabelStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	PanelValueStyle = lipgloss.NewStyle().
			Foreground(FgColor)

	// Help
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	// Table
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(PrimaryColor)

	SelectedRowStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#374151")).
				Foreground(FgColor)

	// Misc
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	NeedsRootStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			Align(lipgloss.Center)
)
