package processes

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/tabs"
	"github.com/jerryluo/nettui/internal/util"
)

// Model is the Processes tab model.
type Model struct {
	table  table.Model
	store  *data.Store
	width  int
	height int
	tabID  model.TabID
	navKey string
	navVal string
	sort       tabs.SortState
	panelWidth int
}

var sortEntries = []tabs.SortEntry{
	{Key: "i", ColKey: "pid", SortKey: "raw_pid", Label: "PID"},
	{Key: "n", ColKey: "name", SortKey: "name", Label: "Name"},
	{Key: "m", ColKey: "command", SortKey: "command", Label: "Command"},
	{Key: "u", ColKey: "user", SortKey: "user", Label: "User"},
	{Key: "c", ColKey: "conns", SortKey: "conns", Label: "#Sockets"},
	{Key: "x", ColKey: "unix_socks", SortKey: "unix_socks", Label: "#Unix"},
}

// New creates a new Processes tab model.
func New() *Model {
	m := &Model{
		tabID: model.TabProcesses,
	}
	m.table = table.New(columns()).
		WithBaseStyle(lipgloss.NewStyle()).
		Focused(true).
		WithPageSize(20).
		Filtered(true).
		HeaderStyle(model.TableHeaderStyle).
		HighlightStyle(model.SelectedRowStyle).
		WithPaginationWrapping(false)
	return m
}

func (m *Model) buildRows() []table.Row {
	if m.store == nil {
		return nil
	}
	rows := make([]table.Row, 0, len(m.store.Processes))
	for _, p := range m.store.Processes {
		rows = append(rows, table.NewRow(table.RowData{
			"pid":        util.FormatPID(p.PID),
			"name":       util.FormatProcess(p.Name),
			"command":    p.Command,
			"user":       p.User,
			"conns":      fmt.Sprintf("%d", p.NumConns),
			"unix_socks": fmt.Sprintf("%d", p.NumUnixSocks),
			"raw_pid":    p.PID,
		}))
	}
	return rows
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = tabs.ClampedUpdate(m.table, msg)
	return m, cmd
}

// View implements tea.Model.
func (m *Model) View() string {
	return m.table.View()
}

// SetData implements Tab.
func (m *Model) SetData(store *data.Store) {
	m.store = store
	rows := m.buildRows()
	if m.navKey != "" {
		rows = tabs.FilterNavRows(rows, m.navKey, m.navVal)
	}
	if m.sort.Active() {
		m.sort.SortRows(rows)
	}
	m.table = m.table.WithRows(rows)
}

// SetSize implements Tab.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.table = m.table.WithPageSize(height - 6).WithTargetWidth(width)
}

// TabID implements Tab.
func (m *Model) TabID() model.TabID {
	return m.tabID
}

// SelectedRow implements Tab.
func (m *Model) SelectedRow() string {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	pid, _ := row.Data["pid"].(string)
	name, _ := row.Data["name"].(string)
	return fmt.Sprintf("PID %s: %s", pid, name)
}

// DetailContent implements Tab.
func (m *Model) DetailContent() string {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	return detailContent(row.Data, m.panelWidth)
}

// CrossRef implements Tab.
func (m *Model) CrossRef() *model.CrossRefMsg {
	return nil // chord mode handles goto now
}

// CrossRefTo returns a CrossRefMsg targeting the given tab, filtered by the selected PID.
func (m *Model) CrossRefTo(target model.TabID) *model.CrossRefMsg {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return nil
	}
	pid, _ := row.Data["pid"].(string)
	if pid == "" || pid == "--" {
		return nil
	}
	return &model.CrossRefMsg{
		TargetTab: target,
		FilterKey: "pid",
		FilterVal: pid,
	}
}

// NavigateTo implements Tab.
func (m *Model) NavigateTo(key, val string) {
	if key != "pid" {
		return
	}
	m.navKey = key
	m.navVal = val
	rows := m.buildRows()
	rows = tabs.FilterNavRows(rows, key, val)
	if m.sort.Active() {
		m.sort.SortRows(rows)
	}
	m.table = m.table.WithRows(rows).WithHighlightedRow(0)
}

// SortHint implements Tab.
func (m *Model) SortHint() string {
	return tabs.Hint(sortEntries)
}

// ApplySort implements Tab.
func (m *Model) ApplySort(key string) {
	if !m.sort.Apply(sortEntries, key) {
		return
	}
	rows := m.buildRows()
	if m.navKey != "" {
		rows = tabs.FilterNavRows(rows, m.navKey, m.navVal)
	}
	m.sort.SortRows(rows)
	m.table = m.table.WithRows(rows)
}

// SortLabel implements Tab.
func (m *Model) SortLabel() string {
	return m.sort.Label()
}

// SetPanelWidth implements Tab.
func (m *Model) SetPanelWidth(width int) {
	m.panelWidth = width
}

// YankHint implements Tab.
func (m *Model) YankHint() string {
	return "y→  p:PID  n:Name  m:Command  y:All"
}

// YankField implements Tab.
func (m *Model) YankField(key string) string {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	switch key {
	case "p":
		v, _ := row.Data["pid"].(string)
		return v
	case "n":
		v, _ := row.Data["name"].(string)
		return v
	case "m":
		v, _ := row.Data["command"].(string)
		return v
	case "y":
		return m.SelectedRow()
	}
	return ""
}

// IsFiltering implements Tab.
func (m *Model) IsFiltering() bool {
	return m.table.GetIsFilterInputFocused()
}

// HasActiveFilter implements Tab.
func (m *Model) HasActiveFilter() bool {
	return m.table.GetCurrentFilter() != "" || m.navKey != ""
}

// ClearFilter implements Tab.
func (m *Model) ClearFilter() {
	if m.table.GetCurrentFilter() != "" {
		m.table = m.table.WithFilterInputValue("")
		return
	}
	if m.navKey != "" {
		m.navKey = ""
		m.navVal = ""
		rows := m.buildRows()
		if m.sort.Active() {
			m.sort.SortRows(rows)
		}
		m.table = m.table.WithRows(rows)
	}
}

// NavFilterLabel implements Tab.
func (m *Model) NavFilterLabel() string {
	if m.navKey == "" {
		return ""
	}
	return fmt.Sprintf("[→%s: %s]", m.navKey, m.navVal)
}
