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
	processesTab "github.com/jerryluo/nettui/internal/tabs/processes"
	socketsTab "github.com/jerryluo/nettui/internal/tabs/sockets"
	"github.com/jerryluo/nettui/internal/ui"
	"github.com/jerryluo/nettui/internal/util"
)

type refreshMsg struct{}
type clearMsgMsg struct{}
type clearChordMsg struct{}

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

	pendingChord rune   // first key of a chord sequence ('g' or 'f')
	chordHint    string // hint text shown in status bar during chord

	warnings map[model.TabID]bool // tabs with partial data
}

// New creates a new root Model with the given tabs.
func New(tabModels []tabs.Tab, collector *sources.Collector) Model {
	m := Model{
		tabs:      tabModels,
		activeTab: model.TabSockets,
		keys:      DefaultKeyMap(),
		collector: collector,
		store:     data.NewStore(),
		panel:     ui.NewSidePanel(),
		warnings:  make(map[model.TabID]bool),
	}
	m.panel.Show()
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return func() tea.Msg { return refreshMsg{} }
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcLayout()
		return m, nil

	case refreshMsg:
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

		return m, nil

	case clearMsgMsg:
		m.message = ""
		return m, nil

	case clearChordMsg:
		m.pendingChord = 0
		m.chordHint = ""
		return m, nil

	case model.CopyResultMsg:
		if msg.Success {
			m.message = "Copied!"
		} else {
			m.message = "Copy failed: " + msg.Error
		}
		return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg { return clearMsgMsg{} })

	case model.CrossRefMsg:
		m.activeTab = msg.TargetTab
		m.tabs[m.activeTab].NavigateTo(msg.FilterKey, msg.FilterVal)
		m.updatePanelContent()
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

	// If a chord is pending, dispatch to the second-key handler
	if m.pendingChord != 0 {
		return m.handleChordSecondKey(msg)
	}

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Help):
		m.showHelp = !m.showHelp
		return m, nil

	case key.Matches(msg, m.keys.NextTab):
		m.activeTab = model.TabID((int(m.activeTab) + 1) % model.TabCount)
		m.updatePanelContent()
		return m, nil

	case key.Matches(msg, m.keys.PrevTab):
		m.activeTab = model.TabID((int(m.activeTab) - 1 + model.TabCount) % model.TabCount)
		m.updatePanelContent()
		return m, nil

	case key.Matches(msg, m.keys.Tab1):
		m.activeTab = model.TabSockets
		m.updatePanelContent()
		return m, nil
	case key.Matches(msg, m.keys.Tab2):
		m.activeTab = model.TabUnixSockets
		m.updatePanelContent()
		return m, nil
	case key.Matches(msg, m.keys.Tab3):
		m.activeTab = model.TabProcesses
		m.updatePanelContent()
		return m, nil
	case key.Matches(msg, m.keys.Tab4):
		m.activeTab = model.TabInterfaces
		m.updatePanelContent()
		return m, nil
	case key.Matches(msg, m.keys.Tab5):
		m.activeTab = model.TabRoutes
		m.updatePanelContent()
		return m, nil
	case key.Matches(msg, m.keys.Tab6):
		m.activeTab = model.TabARP
		m.updatePanelContent()
		return m, nil
	case key.Matches(msg, m.keys.Tab7):
		m.activeTab = model.TabFirewall
		m.updatePanelContent()
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
		if m.tabs[m.activeTab].HasActiveFilter() {
			m.tabs[m.activeTab].ClearFilter()
			return m, nil
		}
		return m, nil

	case key.Matches(msg, m.keys.GoTo):
		// On Processes tab, enter chord mode for target selection
		if m.activeTab == model.TabProcesses {
			m.pendingChord = 'g'
			m.chordHint = "g→  s:Sockets  u:Unix"
			return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearChordMsg{} })
		}
		// On Sockets tab, enter chord mode for go-to selection
		if m.activeTab == model.TabSockets {
			m.pendingChord = 'g'
			m.chordHint = "g→  p:Process  r:Remote"
			return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearChordMsg{} })
		}
		// Other tabs: immediate cross-ref
		ref := m.tabs[m.activeTab].CrossRef()
		if ref != nil {
			return m.Update(*ref)
		}
		return m, nil

	case key.Matches(msg, m.keys.Sort):
		m.pendingChord = 's'
		m.chordHint = m.tabs[m.activeTab].SortHint()
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearChordMsg{} })

	case key.Matches(msg, m.keys.ProtoFilter):
		// On Sockets tab, enter chord mode for protocol filtering
		if m.activeTab == model.TabSockets {
			m.pendingChord = 'f'
			m.chordHint = "f→  t:TCP  u:UDP  4:IPv4  6:IPv6  c:clear"
			return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearChordMsg{} })
		}
		return m, nil

	case key.Matches(msg, m.keys.Copy):
		m.pendingChord = 'y'
		m.chordHint = m.tabs[m.activeTab].YankHint()
		return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg { return clearChordMsg{} })

	case key.Matches(msg, m.keys.Refresh):
		return m, func() tea.Msg { return refreshMsg{} }

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
		if m.panel.Visible() {
			m.panel.SetContent(m.tabs[m.activeTab].DetailContent())
		}
		return m, cmd

	case key.Matches(msg, m.keys.PageUp):
		synth := tea.KeyMsg{Type: tea.KeyPgUp}
		var cmd tea.Cmd
		updated, cmd := m.tabs[m.activeTab].Update(synth)
		m.tabs[m.activeTab] = updated.(tabs.Tab)
		if m.panel.Visible() {
			m.panel.SetContent(m.tabs[m.activeTab].DetailContent())
		}
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
	// Update panel content to reflect new selection
	if m.panel.Visible() {
		content := m.tabs[m.activeTab].DetailContent()
		m.panel.SetContent(content)
	}
	return m, cmd
}

func (m Model) handleChordSecondKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	chord := m.pendingChord
	m.pendingChord = 0
	m.chordHint = ""

	// Escape cancels the chord
	if key.Matches(msg, m.keys.Escape) {
		return m, nil
	}

	k := msg.String()
	switch chord {
	case 'g':
		return m.handleGotoChord(k)
	case 'f':
		return m.handleFilterChord(k)
	case 's':
		return m.handleSortChord(k)
	case 'y':
		return m.handleYankChord(k)
	}
	return m, nil
}

func (m Model) handleGotoChord(k string) (tea.Model, tea.Cmd) {
	switch m.activeTab {
	case model.TabProcesses:
		procTab, ok := m.tabs[model.TabProcesses].(*processesTab.Model)
		if !ok {
			return m, nil
		}
		var ref *model.CrossRefMsg
		switch k {
		case "s":
			ref = procTab.CrossRefTo(model.TabSockets)
		case "u":
			ref = procTab.CrossRefTo(model.TabUnixSockets)
		}
		if ref != nil {
			return m.Update(*ref)
		}

	case model.TabSockets:
		sockTab, ok := m.tabs[model.TabSockets].(*socketsTab.Model)
		if !ok {
			return m, nil
		}
		switch k {
		case "p":
			ref := sockTab.CrossRef()
			if ref != nil {
				return m.Update(*ref)
			}
		case "r":
			if sockTab.GoToRemotePeer() {
				m.updatePanelContent()
			}
		}
	}
	return m, nil
}

func (m Model) handleFilterChord(k string) (tea.Model, tea.Cmd) {
	sockTab, ok := m.tabs[model.TabSockets].(*socketsTab.Model)
	if !ok {
		return m, nil
	}

	switch k {
	case "t":
		sockTab.ToggleTransportFilter(socketsTab.TransportTCP)
	case "u":
		sockTab.ToggleTransportFilter(socketsTab.TransportUDP)
	case "4":
		sockTab.ToggleIPVersionFilter(socketsTab.IPVersion4)
	case "6":
		sockTab.ToggleIPVersionFilter(socketsTab.IPVersion6)
	case "c":
		sockTab.ClearProtoFilters()
	}
	return m, nil
}

func (m Model) handleYankChord(k string) (tea.Model, tea.Cmd) {
	content := m.tabs[m.activeTab].YankField(k)
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
}

func (m Model) handleSortChord(k string) (tea.Model, tea.Cmd) {
	m.tabs[m.activeTab].ApplySort(k)
	m.updatePanelContent()
	return m, nil
}

func (m *Model) updatePanelContent() {
	if m.panel.Visible() {
		content := m.tabs[m.activeTab].DetailContent()
		m.panel.SetContent(content)
	}
	m.recalcLayout()
}

func (m *Model) recalcLayout() {
	m.layout = ui.Calculate(m.width, m.height, m.panel.Visible())
	for _, t := range m.tabs {
		t.SetSize(m.layout.TableWidth, m.layout.ContentHeight)
	}
	if m.layout.PanelOpen {
		m.panel.SetSize(m.layout.PanelWidth, m.layout.ContentHeight)
	}
	for _, t := range m.tabs {
		t.SetPanelWidth(m.layout.PanelWidth)
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
		tabView = lipgloss.NewStyle().Width(m.layout.TableWidth).Render(tabView)
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

	// Extract proto filter label from sockets tab
	var protoFilter string
	if sockTab, ok := m.tabs[model.TabSockets].(*socketsTab.Model); ok {
		protoFilter = sockTab.ProtoFilterLabel()
	}

	// Extract sort label from active tab
	sortLabel := m.tabs[m.activeTab].SortLabel()

	// Extract nav filter label from active tab
	navFilter := m.tabs[m.activeTab].NavFilterLabel()

	// Status bar
	statusBar := ui.RenderStatusBar(ui.StatusBarState{
		IsRoot:      m.store.IsRoot,
		DNSEnabled:  m.dnsOn,
		Message:     m.message,
		ChordHint:   m.chordHint,
		ProtoFilter: protoFilter,
		SortLabel:   sortLabel,
		NavFilter:   navFilter,
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
		{"1-7", "Jump to tab"},
		{"j/k / arrows", "Navigate rows"},
		{"d/u", "Page down / up"},
		{"/", "Filter / search"},
		{"p", "Toggle side panel"},
		{"Esc", "Clear filter / close panel"},
		{"g", "Go to cross-referenced entity"},
		{"gs/gu", "Go to Sockets/Unix (Processes tab)"},
		{"gp/gr", "Go to Process/Remote (Sockets tab)"},
		{"f", "Protocol filter (Sockets tab)"},
		{"ft/fu/f4/f6/fc", "TCP/UDP/IPv4/IPv6/clear"},
		{"s", "Sort by column (chord)"},
		{"y", "Yank (copy) chord — field to clipboard"},
		{"yl/yr", "Yank local/remote addr (Sockets)"},
		{"yp/yn", "Yank PID/process name"},
		{"yy", "Yank full row summary"},
		{"r", "Refresh data"},
		{"D", "Toggle DNS resolution"},
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
