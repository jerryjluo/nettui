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
	sort   tabs.SortState
}

var sortEntries = []tabs.SortEntry{
	{Key: "i", ColKey: "pid", SortKey: "raw_pid", Label: "PID"},
	{Key: "n", ColKey: "name", SortKey: "name", Label: "Name"},
	{Key: "u", ColKey: "user", SortKey: "user", Label: "User"},
	{Key: "c", ColKey: "conns", SortKey: "conns", Label: "#Conns"},
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
		HighlightStyle(model.SelectedRowStyle)
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
	m.table, cmd = m.table.Update(msg)
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
	if m.sort.Active() {
		m.sort.SortRows(rows)
	} else if m.navKey != "" {
		rows = m.reorderRows(rows, m.navKey, m.navVal)
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
	return detailContent(row.Data)
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
	m.sort.Clear()
	m.navKey = key
	m.navVal = val
	rows := m.reorderRows(m.buildRows(), key, val)
	m.table = m.table.WithRows(rows).WithHighlightedRow(0)
}

func (m *Model) reorderRows(rows []table.Row, key, val string) []table.Row {
	reordered := make([]table.Row, 0, len(rows))
	var rest []table.Row
	for _, r := range rows {
		if fmt.Sprintf("%v", r.Data[key]) == val {
			reordered = append(reordered, r)
		} else {
			rest = append(rest, r)
		}
	}
	return append(reordered, rest...)
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
	m.navKey = ""
	m.navVal = ""
	rows := m.buildRows()
	m.sort.SortRows(rows)
	m.table = m.table.WithRows(rows)
}

// SortLabel implements Tab.
func (m *Model) SortLabel() string {
	return m.sort.Label()
}

// IsFiltering implements Tab.
func (m *Model) IsFiltering() bool {
	return m.table.GetIsFilterInputFocused()
}
