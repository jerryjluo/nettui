package interfaces

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/data"
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
		HighlightStyle(model.SelectedRowStyle)
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
	if m.navKey != "" {
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

// IsFiltering implements Tab.
func (m *Model) IsFiltering() bool {
	return m.table.GetIsFilterInputFocused()
}
