package interfaces

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/tabs"
	"github.com/jerryluo/nettui/internal/util"
)

// Model is the Interfaces tab model.
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
	{Key: "n", ColKey: "name", SortKey: "name", Label: "Name"},
	{Key: "m", ColKey: "mtu", SortKey: "mtu", Label: "MTU"},
	{Key: "s", ColKey: "status", SortKey: "status", Label: "Status"},
	{Key: "t", ColKey: "tx_bytes", SortKey: "raw_tx", Label: "TX"},
	{Key: "r", ColKey: "rx_bytes", SortKey: "raw_rx", Label: "RX"},
}

// New creates a new Interfaces tab model.
func New() *Model {
	m := &Model{
		tabID: model.TabInterfaces,
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
	rows := make([]table.Row, 0, len(m.store.Interfaces))
	for _, iface := range m.store.Interfaces {
		status := "down"
		if iface.Up {
			status = "up"
		}
		rows = append(rows, table.NewRow(table.RowData{
			"name":    iface.Name,
			"addrs":   strings.Join(iface.Addrs, ", "),
			"mac":     iface.HWAddr,
			"mtu":     fmt.Sprintf("%d", iface.MTU),
			"status":  status,
			"tx_bytes": util.FormatBytes(iface.BytesSent),
			"rx_bytes": util.FormatBytes(iface.BytesRecv),
			"tx_rate": util.FormatRate(iface.TxRate),
			"rx_rate": util.FormatRate(iface.RxRate),
			"tx_pkts": fmt.Sprintf("%d", iface.PacketSent),
			"rx_pkts": fmt.Sprintf("%d", iface.PacketRecv),
			"flags":   iface.Flags.String(),
			"raw_tx":  iface.BytesSent,
			"raw_rx":  iface.BytesRecv,
		}))
	}
	return rows
}

// YankHint implements Tab.
func (m *Model) YankHint() string {
	return "y→  n:Name  a:Addrs  m:MAC  y:All"
}

// YankField implements Tab.
func (m *Model) YankField(key string) string {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	switch key {
	case "n":
		v, _ := row.Data["name"].(string)
		return v
	case "a":
		v, _ := row.Data["addrs"].(string)
		return v
	case "m":
		v, _ := row.Data["mac"].(string)
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
	name, _ := row.Data["name"].(string)
	addrs, _ := row.Data["addrs"].(string)
	return fmt.Sprintf("%s (%s)", name, addrs)
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
	name, _ := row.Data["name"].(string)
	if name == "" {
		return nil
	}
	return &model.CrossRefMsg{
		TargetTab: model.TabRoutes,
		FilterKey: "iface",
		FilterVal: name,
	}
}

// NavigateTo implements Tab.
func (m *Model) NavigateTo(key, val string) {
	if key != "name" {
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
func (m *Model) SetPanelWidth(width int) {}

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
