package firewall

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/util"
)

// Model is the Firewall tab model.
type Model struct {
	table  table.Model
	store  *data.Store
	width  int
	height int
	tabID  model.TabID
}

// New creates a new Firewall tab model.
func New() *Model {
	m := &Model{
		tabID: model.TabFirewall,
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
	rows := make([]table.Row, 0, len(m.store.Firewall))
	for _, r := range m.store.Firewall {
		rows = append(rows, table.NewRow(table.RowData{
			"rule":     fmt.Sprintf("%d", r.RuleNum),
			"action":   r.Action,
			"dir":      r.Direction,
			"proto":    r.Proto,
			"src":      r.Src,
			"dst":      r.Dst,
			"packets":  fmt.Sprintf("%d", r.Packets),
			"bytes":    util.FormatBytes(r.Bytes),
			"raw_rule": r.RawRule,
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
	if m.store == nil || !m.store.IsRoot {
		msg := "Requires root privileges"
		styled := model.NeedsRootStyle.Width(m.width).Render(msg)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, styled)
	}
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
	if m.store == nil || !m.store.IsRoot {
		return ""
	}
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	rule, _ := row.Data["rule"].(string)
	action, _ := row.Data["action"].(string)
	dir, _ := row.Data["dir"].(string)
	return fmt.Sprintf("Rule %s: %s %s", rule, action, dir)
}

// DetailContent implements Tab.
func (m *Model) DetailContent() string {
	if m.store == nil || !m.store.IsRoot {
		return ""
	}
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	return detailContent(row.Data)
}

// CrossRef implements Tab.
func (m *Model) CrossRef() *model.CrossRefMsg {
	return nil
}

// NavigateTo implements Tab.
func (m *Model) NavigateTo(key, val string) {
	// Firewall has no cross-reference navigation target.
}

// IsFiltering implements Tab.
func (m *Model) IsFiltering() bool {
	return m.table.GetIsFilterActive()
}
