package ui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jerryluo/nettui/internal/model"
)

// SidePanel is a scrollable detail panel using viewport.
type SidePanel struct {
	viewport viewport.Model
	content  string
	width    int
	height   int
	visible  bool
}

// NewSidePanel creates a new side panel.
func NewSidePanel() SidePanel {
	vp := viewport.New(30, 10)
	return SidePanel{
		viewport: vp,
	}
}

// SetSize updates the panel dimensions.
func (s *SidePanel) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.viewport.Width = width - 4 // account for border padding
	s.viewport.Height = height - 2
}

// SetContent updates the panel content.
func (s *SidePanel) SetContent(content string) {
	s.content = content
	s.viewport.SetContent(content)
	s.viewport.GotoTop()
}

// Toggle toggles panel visibility.
func (s *SidePanel) Toggle() {
	s.visible = !s.visible
}

// Show opens the panel.
func (s *SidePanel) Show() {
	s.visible = true
}

// Hide closes the panel.
func (s *SidePanel) Hide() {
	s.visible = false
}

// Visible returns whether the panel is visible.
func (s *SidePanel) Visible() bool {
	return s.visible
}

// Update handles messages for scrolling.
func (s *SidePanel) Update(msg tea.Msg) {
	s.viewport, _ = s.viewport.Update(msg)
}

// ScrollDown scrolls the panel down.
func (s *SidePanel) ScrollDown() {
	s.viewport.LineDown(3)
}

// ScrollUp scrolls the panel up.
func (s *SidePanel) ScrollUp() {
	s.viewport.LineUp(3)
}

// View renders the panel.
func (s *SidePanel) View() string {
	if !s.visible {
		return ""
	}
	return model.PanelBorderStyle.
		Width(s.width - 2).
		Height(s.height - 2).
		Render(s.viewport.View())
}
