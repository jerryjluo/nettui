package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jerryluo/nettui/internal/model"
)

// StatusBarState holds the current state for status bar rendering.
type StatusBarState struct {
	IsRoot      bool
	DNSEnabled  bool
	FilterText  string
	Message     string
	ChordHint   string
	ProtoFilter string
	SortLabel   string
}

// RenderStatusBar renders the bottom status bar.
func RenderStatusBar(state StatusBarState, width int) string {
	var left []string
	var right []string

	if state.FilterText != "" {
		left = append(left, fmt.Sprintf("/%s", state.FilterText))
	}

	if state.ChordHint != "" {
		left = append(left, lipgloss.NewStyle().Foreground(model.AccentColor).Bold(true).Render(state.ChordHint))
	} else {
		hints := []string{
			model.HelpKeyStyle.Render("c") + model.HelpDescStyle.Render(":copy"),
			model.HelpKeyStyle.Render("g") + model.HelpDescStyle.Render(":goto"),
			model.HelpKeyStyle.Render("D") + model.HelpDescStyle.Render(":DNS"),
			model.HelpKeyStyle.Render("?") + model.HelpDescStyle.Render(":help"),
		}
		left = append(left, strings.Join(hints, "  "))
	}

	if state.Message != "" {
		left = append(left, lipgloss.NewStyle().Foreground(model.SuccessColor).Render(state.Message))
	}

	if state.SortLabel != "" {
		right = append(right, lipgloss.NewStyle().Foreground(model.AccentColor).Bold(true).Render(state.SortLabel))
	}

	if state.ProtoFilter != "" {
		right = append(right, lipgloss.NewStyle().Foreground(model.AccentColor).Bold(true).Render(state.ProtoFilter))
	}

	if !state.IsRoot {
		right = append(right, model.StatusBadgeStyle.Render("[no root]"))
	}

	dnsStatus := "off"
	dnsColor := model.MutedColor
	if state.DNSEnabled {
		dnsStatus = "on"
		dnsColor = model.SuccessColor
	}
	right = append(right, lipgloss.NewStyle().Foreground(dnsColor).Render(fmt.Sprintf("DNS:%s", dnsStatus)))

	leftStr := strings.Join(left, "  ")
	rightStr := strings.Join(right, "  ")

	// StatusBarStyle has Padding(0, 1) = 2 chars of horizontal padding
	gap := width - 2 - lipgloss.Width(leftStr) - lipgloss.Width(rightStr)
	if gap < 1 {
		gap = 1
	}

	return model.StatusBarStyle.Render(leftStr + strings.Repeat(" ", gap) + rightStr)
}
