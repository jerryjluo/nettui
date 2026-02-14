package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/data/sources"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/tabs"
	socketsTab "github.com/jerryluo/nettui/internal/tabs/sockets"
	"github.com/jerryluo/nettui/internal/ui"
	"github.com/jerryluo/nettui/internal/util"
)

type tickMsg time.Time

// Model is the root application model.
type Model struct {
	tabs      []tabs.Tab
	activeTab model.TabID
	keys      KeyMap
	collector *sources.Collector
	store     *data.Store
	panel     ui.SidePanel
	layout    ui.Layout

	width    int
	height   int
	showHelp bool
	dnsOn    bool
	message  string // ephemeral status message
	msgTimer int    // ticks remaining for message display

	warnings map[model.TabID]bool // tabs with partial data
}

// New creates a new root Model with the given tabs.
func New(tabModels []tabs.Tab, collector *sources.Collector) Model {
	return Model{
		tabs:      tabModels,
		activeTab: model.TabInterfaces,
		keys:      DefaultKeyMap(),
		collector: collector,
		store:     data.NewStore(),
		panel:     ui.NewSidePanel(),
		warnings:  make(map[model.TabID]bool),
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
	)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcLayout()
		return m, nil

	case tickMsg:
		result := m.collector.Collect()
		m.store.Update(result)
		snap := m.store.Snapshot()

		// Update warnings
		m.warnings = make(map[model.TabID]bool)
		if !snap.IsRoot {
			m.warnings[model.TabSockets] = true
			m.warnings[model.TabUnixSockets] = true
			m.warnings[model.TabProcesses] = true
			m.warnings[model.TabFirewall] = true
		}

		for _, t := range m.tabs {
			t.SetData(snap)
		}

		// Update side panel content if open
		if m.panel.Visible() {
			content := m.tabs[m.activeTab].DetailContent()
			m.panel.SetContent(content)
		}

		// Decrement message timer
		if m.msgTimer > 0 {
			m.msgTimer--
			if m.msgTimer == 0 {
				m.message = ""
			}
		}

		return m, tick()

	case model.CopyResultMsg:
		if msg.Success {
			m.message = "Copied!"
		} else {
			m.message = "Copy failed: " + msg.Error
		}
		m.msgTimer = 3
		return m, nil

	case model.CrossRefMsg:
		m.activeTab = msg.TargetTab
		m.tabs[m.activeTab].NavigateTo(msg.FilterKey, msg.FilterVal)
		m.panel.Hide()
		m.recalcLayout()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Delegate to active tab
	if int(m.activeTab) < len(m.tabs) {
		var cmd tea.Cmd
		updated, cmd := m.tabs[m.activeTab].Update(msg)
		m.tabs[m.activeTab] = updated.(tabs.Tab)
		return m, cmd
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If current tab is filtering, let it handle all keys
	if m.tabs[m.activeTab].IsFiltering() {
		var cmd tea.Cmd
		updated, cmd := m.tabs[m.activeTab].Update(msg)
		m.tabs[m.activeTab] = updated.(tabs.Tab)
		return m, cmd
	}

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Help):
		m.showHelp = !m.showHelp
		return m, nil

	case key.Matches(msg, m.keys.NextTab):
		m.activeTab = model.TabID((int(m.activeTab) + 1) % model.TabCount)
		m.panel.Hide()
		m.recalcLayout()
		return m, nil

	case key.Matches(msg, m.keys.PrevTab):
		m.activeTab = model.TabID((int(m.activeTab) - 1 + model.TabCount) % model.TabCount)
		m.panel.Hide()
		m.recalcLayout()
		return m, nil

	case key.Matches(msg, m.keys.Tab1):
		m.activeTab = model.TabInterfaces
		m.panel.Hide()
		m.recalcLayout()
		return m, nil
	case key.Matches(msg, m.keys.Tab2):
		m.activeTab = model.TabRoutes
		m.panel.Hide()
		m.recalcLayout()
		return m, nil
	case key.Matches(msg, m.keys.Tab3):
		m.activeTab = model.TabSockets
		m.panel.Hide()
		m.recalcLayout()
		return m, nil
	case key.Matches(msg, m.keys.Tab4):
		m.activeTab = model.TabUnixSockets
		m.panel.Hide()
		m.recalcLayout()
		return m, nil
	case key.Matches(msg, m.keys.Tab5):
		m.activeTab = model.TabProcesses
		m.panel.Hide()
		m.recalcLayout()
		return m, nil
	case key.Matches(msg, m.keys.Tab6):
		m.activeTab = model.TabFirewall
		m.panel.Hide()
		m.recalcLayout()
		return m, nil

	case key.Matches(msg, m.keys.Enter):
		m.panel.Toggle()
		if m.panel.Visible() {
			content := m.tabs[m.activeTab].DetailContent()
			m.panel.SetContent(content)
		}
		m.recalcLayout()
		return m, nil

	case key.Matches(msg, m.keys.Escape):
		if m.panel.Visible() {
			m.panel.Hide()
			m.recalcLayout()
			return m, nil
		}
		// Delegate to tab (e.g. clear filter)
		var cmd tea.Cmd
		updated, cmd := m.tabs[m.activeTab].Update(msg)
		m.tabs[m.activeTab] = updated.(tabs.Tab)
		return m, cmd

	case key.Matches(msg, m.keys.GoTo):
		ref := m.tabs[m.activeTab].CrossRef()
		if ref != nil {
			return m.Update(*ref)
		}
		return m, nil

	case key.Matches(msg, m.keys.Copy):
		content := m.tabs[m.activeTab].SelectedRow()
		if content == "" {
			return m, nil
		}
		return m, func() tea.Msg {
			err := util.CopyToClipboard(content)
			if err != nil {
				return model.CopyResultMsg{Success: false, Error: err.Error()}
			}
			return model.CopyResultMsg{Success: true}
		}

	case key.Matches(msg, m.keys.DNS):
		m.dnsOn = !m.dnsOn
		// Propagate DNS state to the sockets tab.
		if sockTab, ok := m.tabs[model.TabSockets].(*socketsTab.Model); ok {
			sockTab.SetDNSEnabled(m.dnsOn)
		}
		return m, nil

	case key.Matches(msg, m.keys.PageDown):
		synth := tea.KeyMsg{Type: tea.KeyPgDown}
		var cmd tea.Cmd
		updated, cmd := m.tabs[m.activeTab].Update(synth)
		m.tabs[m.activeTab] = updated.(tabs.Tab)
		return m, cmd

	case key.Matches(msg, m.keys.PageUp):
		synth := tea.KeyMsg{Type: tea.KeyPgUp}
		var cmd tea.Cmd
		updated, cmd := m.tabs[m.activeTab].Update(synth)
		m.tabs[m.activeTab] = updated.(tabs.Tab)
		return m, cmd

	case key.Matches(msg, m.keys.Filter):
		// Delegate '/' to the active tab to start filtering
		var cmd tea.Cmd
		updated, cmd := m.tabs[m.activeTab].Update(msg)
		m.tabs[m.activeTab] = updated.(tabs.Tab)
		return m, cmd
	}

	// Default: delegate to active tab (for j/k/arrow navigation)
	var cmd tea.Cmd
	updated, cmd := m.tabs[m.activeTab].Update(msg)
	m.tabs[m.activeTab] = updated.(tabs.Tab)
	return m, cmd
}

func (m *Model) recalcLayout() {
	m.layout = ui.Calculate(m.width, m.height, m.panel.Visible())
	for _, t := range m.tabs {
		t.SetSize(m.layout.TableWidth, m.layout.ContentHeight)
	}
	if m.layout.PanelOpen {
		m.panel.SetSize(m.layout.PanelWidth, m.layout.ContentHeight)
	}
}

// View implements tea.Model.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	if m.showHelp {
		return m.helpView()
	}

	// Tab bar
	tabBar := ui.RenderTabBar(m.activeTab, m.width, m.warnings)

	// Content area
	var content string
	tabView := m.tabs[m.activeTab].View()

	if m.layout.PanelOpen {
		panelView := m.panel.View()
		content = lipgloss.JoinHorizontal(lipgloss.Top, tabView, panelView)
	} else {
		content = tabView
	}

	// Truncate/pad content to fit
	content = lipgloss.NewStyle().
		Width(m.width).
		Height(m.layout.ContentHeight).
		Render(content)

	// Status bar
	statusBar := ui.RenderStatusBar(ui.StatusBarState{
		IsRoot:     m.store.IsRoot,
		DNSEnabled: m.dnsOn,
		Message:    m.message,
	}, m.width)

	return lipgloss.JoinVertical(lipgloss.Left, tabBar, content, statusBar)
}

func (m Model) helpView() string {
	title := model.PanelHeaderStyle.Render("Keybindings")
	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	bindings := []struct {
		key  string
		desc string
	}{
		{"q / Ctrl+C", "Quit"},
		{"h/l / Tab/Shift+Tab", "Prev / next tab"},
		{"1-6", "Jump to tab"},
		{"j/k / arrows", "Navigate rows"},
		{"J/K", "Page down / up"},
		{"/", "Filter / search"},
		{"Enter", "Toggle side panel"},
		{"Esc", "Close panel / clear filter"},
		{"g", "Go to cross-referenced entity"},
		{"c", "Copy selection to clipboard"},
		{"d", "Toggle DNS resolution"},
		{"?", "Toggle this help"},
	}

	for _, b := range bindings {
		line := fmt.Sprintf("  %s  %s",
			model.HelpKeyStyle.Render(fmt.Sprintf("%-20s", b.key)),
			model.HelpDescStyle.Render(b.desc),
		)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, model.HelpDescStyle.Render("  Press ? to close"))

	helpText := strings.Join(lines, "\n")
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(model.PrimaryColor).
		Padding(1, 3).
		Render(helpText)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
