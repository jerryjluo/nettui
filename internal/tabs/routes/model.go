package routes

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/tabs"
)

// Model is the Routes tab model.
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
	{Key: "d", ColKey: "dest", SortKey: "dest", Label: "Dest"},
	{Key: "g", ColKey: "gateway", SortKey: "gateway", Label: "Gateway"},
	{Key: "n", ColKey: "netmask", SortKey: "netmask", Label: "Netmask"},
	{Key: "i", ColKey: "iface", SortKey: "iface", Label: "Interface"},
	{Key: "f", ColKey: "flags", SortKey: "flags", Label: "Flags"},
}

// New creates a new Routes tab model.
func New() *Model {
	m := &Model{
		tabID: model.TabRoutes,
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
	rows := make([]table.Row, 0, len(m.store.Routes))
	for _, r := range m.store.Routes {
		rows = append(rows, table.NewRow(table.RowData{
			"dest":    r.Destination,
			"gateway": r.Gateway,
			"netmask": r.Netmask,
			"iface":   r.Interface,
			"flags":   r.Flags,
		}))
	}
	return rows
}

// YankHint implements Tab.
func (m *Model) YankHint() string {
	return "y→  d:Dest  g:Gateway  i:Iface  y:All"
}

// YankField implements Tab.
func (m *Model) YankField(key string) string {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	switch key {
	case "d":
		v, _ := row.Data["dest"].(string)
		return v
	case "g":
		v, _ := row.Data["gateway"].(string)
		return v
	case "i":
		v, _ := row.Data["iface"].(string)
		return v
	case "y":
		return m.SelectedRow()
	}
	return ""
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
	dest, _ := row.Data["dest"].(string)
	gw, _ := row.Data["gateway"].(string)
	iface, _ := row.Data["iface"].(string)
	return fmt.Sprintf("%s via %s dev %s", dest, gw, iface)
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
	iface, _ := row.Data["iface"].(string)
	if iface == "" {
		return nil
	}
	return &model.CrossRefMsg{
		TargetTab: model.TabInterfaces,
		FilterKey: "name",
		FilterVal: iface,
	}
}

// NavigateTo implements Tab.
func (m *Model) NavigateTo(key, val string) {
	if key != "iface" {
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
