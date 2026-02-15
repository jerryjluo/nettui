package unixsockets

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

// Model is the Unix Sockets tab model.
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
	{Key: "a", ColKey: "path", SortKey: "path", Label: "Path"},
	{Key: "t", ColKey: "type", SortKey: "type", Label: "Type"},
	{Key: "s", ColKey: "state", SortKey: "state", Label: "State"},
	{Key: "i", ColKey: "pid", SortKey: "raw_pid", Label: "PID"},
	{Key: "n", ColKey: "process", SortKey: "process", Label: "Process"},
	{Key: "f", ColKey: "fd", SortKey: "fd", Label: "FD"},
}

// New creates a new Unix Sockets tab model.
func New() *Model {
	m := &Model{
		tabID: model.TabUnixSockets,
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
	rows := make([]table.Row, 0, len(m.store.UnixSockets))
	for _, s := range m.store.UnixSockets {
		state := s.State
		if state == "" {
			state = "-"
		}
		rows = append(rows, table.NewRow(table.RowData{
			"path":    s.Path,
			"type":    s.Type,
			"state":   state,
			"pid":     util.FormatPID(s.PID),
			"process": util.FormatProcess(s.Process),
			"fd":      s.FD,
			"raw_pid": s.PID,
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
	path, _ := row.Data["path"].(string)
	pid, _ := row.Data["pid"].(string)
	return fmt.Sprintf("%s (PID: %s)", path, pid)
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
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return nil
	}
	pid, _ := row.Data["raw_pid"].(int32)
	if pid <= 0 {
		return nil
	}
	return &model.CrossRefMsg{
		TargetTab: model.TabProcesses,
		FilterKey: "pid",
		FilterVal: fmt.Sprintf("%d", pid),
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

// SetPanelWidth implements Tab.
func (m *Model) SetPanelWidth(width int) {}

// IsFiltering implements Tab.
func (m *Model) IsFiltering() bool {
	return m.table.GetIsFilterInputFocused()
}

// HasActiveFilter implements Tab.
func (m *Model) HasActiveFilter() bool {
	return m.table.GetCurrentFilter() != ""
}

// ClearFilter implements Tab.
func (m *Model) ClearFilter() {
	m.table = m.table.WithFilterInputValue("")
}
