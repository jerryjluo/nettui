package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jerryluo/nettui/internal/model"
)

// RenderTabBar renders the tab strip.
func RenderTabBar(activeTab model.TabID, width int, warnings map[model.TabID]bool) string {
	var tabs []string

	for i := range model.TabCount {
		id := model.TabID(i)
		name := model.TabName(id)

		if warnings[id] {
			name += " (!)"
		}

		if id == activeTab {
			tabs = append(tabs, model.ActiveTabStyle.Render(name))
		} else {
			tabs = append(tabs, model.InactiveTabStyle.Render(name))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	line := strings.Repeat("─", max(0, width-lipgloss.Width(row)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, lipgloss.NewStyle().Foreground(model.MutedColor).Render(line))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
