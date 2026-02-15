package arp

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/tabs"
)

// Model is the ARP tab model.
type Model struct {
	table  table.Model
	store  *data.Store
	width  int
	height int
	tabID  model.TabID
	sort   tabs.SortState
}

var sortEntries = []tabs.SortEntry{
	{Key: "i", ColKey: "ip", SortKey: "ip", Label: "IP"},
	{Key: "m", ColKey: "mac", SortKey: "mac", Label: "MAC"},
	{Key: "n", ColKey: "iface", SortKey: "iface", Label: "Interface"},
	{Key: "f", ColKey: "flags", SortKey: "flags", Label: "Flags"},
}

// New creates a new ARP tab model.
func New() *Model {
	m := &Model{
		tabID: model.TabARP,
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
	rows := make([]table.Row, 0, len(m.store.ARPEntries))
	for _, e := range m.store.ARPEntries {
		rows = append(rows, table.NewRow(table.RowData{
			"ip":       e.IP,
			"mac":      e.MAC,
			"iface":    e.Interface,
			"hostname": e.Hostname,
			"flags":    e.Flags,
			"type":     e.Type,
		}))
	}
	return rows
}

// YankHint implements Tab.
func (m *Model) YankHint() string {
	return "y→  i:IP  m:MAC  n:Iface  y:All"
}

// YankField implements Tab.
func (m *Model) YankField(key string) string {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	switch key {
	case "i":
		v, _ := row.Data["ip"].(string)
		return v
	case "m":
		v, _ := row.Data["mac"].(string)
		return v
	case "n":
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
	ip, _ := row.Data["ip"].(string)
	mac, _ := row.Data["mac"].(string)
	iface, _ := row.Data["iface"].(string)
	return fmt.Sprintf("%s at %s on %s", ip, mac, iface)
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
	return nil
}

// NavigateTo implements Tab.
func (m *Model) NavigateTo(key, val string) {}

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
