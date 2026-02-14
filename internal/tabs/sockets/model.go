package sockets

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/util"
)

// Model is the Sockets tab model.
type Model struct {
	table  table.Model
	store  *data.Store
	width  int
	height int
	tabID  model.TabID
}

// New creates a new Sockets tab model.
func New() *Model {
	m := &Model{
		tabID: model.TabSockets,
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
	rows := make([]table.Row, 0, len(m.store.Sockets))
	for _, s := range m.store.Sockets {
		rows = append(rows, table.NewRow(table.RowData{
			"proto":   s.Proto,
			"local":   util.FormatAddrPort(s.LocalAddr, s.LocalPort),
			"remote":  util.FormatAddrPort(s.RemoteAddr, s.RemotePort),
			"state":   s.State,
			"pid":     util.FormatPID(s.PID),
			"process": util.FormatProcess(s.Process),
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
	m.table = m.table.WithRows(m.buildRows())
}

// SetSize implements Tab.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.table = m.table.WithPageSize(height - 4).WithMaxTotalWidth(width)
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
	proto, _ := row.Data["proto"].(string)
	local, _ := row.Data["local"].(string)
	remote, _ := row.Data["remote"].(string)
	return fmt.Sprintf("%s %s -> %s", proto, local, remote)
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
	rows := m.buildRows()
	reordered := make([]table.Row, 0, len(rows))
	var rest []table.Row
	for _, r := range rows {
		if fmt.Sprintf("%v", r.Data["pid"]) == val {
			reordered = append(reordered, r)
		} else {
			rest = append(rest, r)
		}
	}
	reordered = append(reordered, rest...)
	m.table = m.table.WithRows(reordered)
}

// IsFiltering implements Tab.
func (m *Model) IsFiltering() bool {
	return m.table.GetIsFilterActive()
}
