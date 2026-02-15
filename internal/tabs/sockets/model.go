package sockets

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/data"
	"github.com/jerryluo/nettui/internal/data/sources"
	"github.com/jerryluo/nettui/internal/tabs"
	"github.com/jerryluo/nettui/internal/util"
)

// TransportFilter selects TCP or UDP sockets.
type TransportFilter int

const (
	TransportNone TransportFilter = iota
	TransportTCP
	TransportUDP
)

// IPVersionFilter selects IPv4 or IPv6 sockets.
type IPVersionFilter int

const (
	IPVersionNone IPVersionFilter = iota
	IPVersion4
	IPVersion6
)

// Model is the Sockets tab model.
type Model struct {
	table    table.Model
	store    *data.Store
	width    int
	height   int
	tabID    model.TabID
	navKey   string // cross-ref ordering key (Bug 3)
	navVal   string // cross-ref ordering value (Bug 3)
	dnsCache *sources.DNSCache // DNS cache (Bug 4)
	dnsOn    bool              // DNS resolution enabled (Bug 4)

	transportFilter TransportFilter
	ipVersionFilter IPVersionFilter

	sort       tabs.SortState
	panelWidth int
}

var sortEntries = []tabs.SortEntry{
	{Key: "p", ColKey: "proto", SortKey: "proto", Label: "Proto"},
	{Key: "l", ColKey: "local", SortKey: "local", Label: "Local"},
	{Key: "r", ColKey: "remote", SortKey: "remote", Label: "Remote"},
	{Key: "s", ColKey: "state", SortKey: "state", Label: "State"},
	{Key: "i", ColKey: "pid", SortKey: "raw_pid", Label: "PID"},
	{Key: "n", ColKey: "process", SortKey: "process", Label: "Process"},
	{Key: "m", ColKey: "command", SortKey: "command", Label: "Command"},
}

// New creates a new Sockets tab model.
func New(dns *sources.DNSCache) *Model {
	m := &Model{
		tabID:    model.TabSockets,
		dnsCache: dns,
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
	rows := make([]table.Row, 0, len(m.store.Sockets))
	for _, s := range m.store.Sockets {
		if !m.matchesProtoFilter(s.Proto) {
			continue
		}
		remoteAddr := s.RemoteAddr
		if m.dnsOn && m.dnsCache != nil && remoteAddr != "" {
			remoteAddr = m.dnsCache.Lookup(remoteAddr)
		}
		command := ""
		if m.store != nil && s.PID > 0 {
			if proc, ok := m.store.ProcessByPID[s.PID]; ok {
				command = proc.Command
			}
		}
		rows = append(rows, table.NewRow(table.RowData{
			"proto":           s.Proto,
			"local":           util.FormatAddrPort(s.LocalAddr, s.LocalPort),
			"remote":          util.FormatAddrPort(remoteAddr, s.RemotePort),
			"state":           s.State,
			"pid":             util.FormatPID(s.PID),
			"process":         util.FormatProcess(s.Process),
			"command":         command,
			"raw_pid":         s.PID,
			"raw_local_addr":  s.LocalAddr,
			"raw_local_port":  s.LocalPort,
			"raw_remote_addr": s.RemoteAddr,
			"raw_remote_port": s.RemotePort,
		}))
	}
	return rows
}

func (m *Model) matchesProtoFilter(proto string) bool {
	p := strings.ToLower(proto)
	if m.transportFilter != TransportNone {
		switch m.transportFilter {
		case TransportTCP:
			if !strings.HasPrefix(p, "tcp") {
				return false
			}
		case TransportUDP:
			if !strings.HasPrefix(p, "udp") {
				return false
			}
		}
	}
	if m.ipVersionFilter != IPVersionNone {
		switch m.ipVersionFilter {
		case IPVersion4:
			if strings.HasSuffix(p, "6") {
				return false
			}
		case IPVersion6:
			if !strings.HasSuffix(p, "6") {
				return false
			}
		}
	}
	return true
}

// ToggleTransportFilter toggles the given transport filter (or clears it if already active).
func (m *Model) ToggleTransportFilter(f TransportFilter) {
	if m.transportFilter == f {
		m.transportFilter = TransportNone
	} else {
		m.transportFilter = f
	}
	m.applyFilters()
}

// ToggleIPVersionFilter toggles the given IP version filter (or clears it if already active).
func (m *Model) ToggleIPVersionFilter(f IPVersionFilter) {
	if m.ipVersionFilter == f {
		m.ipVersionFilter = IPVersionNone
	} else {
		m.ipVersionFilter = f
	}
	m.applyFilters()
}

// ClearProtoFilters clears all protocol filters.
func (m *Model) ClearProtoFilters() {
	m.transportFilter = TransportNone
	m.ipVersionFilter = IPVersionNone
	m.applyFilters()
}

// ProtoFilterLabel returns a display label for the active protocol filters, or "" if none.
func (m *Model) ProtoFilterLabel() string {
	var parts []string
	switch m.transportFilter {
	case TransportTCP:
		parts = append(parts, "TCP")
	case TransportUDP:
		parts = append(parts, "UDP")
	}
	switch m.ipVersionFilter {
	case IPVersion4:
		parts = append(parts, "IPv4")
	case IPVersion6:
		parts = append(parts, "IPv6")
	}
	if len(parts) == 0 {
		return ""
	}
	return "[" + strings.Join(parts, "+") + "]"
}

func (m *Model) applyFilters() {
	rows := m.buildRows()
	if m.sort.Active() {
		m.sort.SortRows(rows)
	} else if m.navKey != "" {
		rows = m.reorderRows(rows, m.navKey, m.navVal)
	}
	m.table = m.table.WithRows(rows)
}

// SetDNSEnabled enables or disables DNS resolution for remote addresses.
func (m *Model) SetDNSEnabled(on bool) {
	m.dnsOn = on
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
	return detailContent(row.Data, m.panelWidth)
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
	m.applyFilters()
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
	return "y→  l:Local  r:Remote  p:PID  n:Process  m:Command  y:All"
}

// YankField implements Tab.
func (m *Model) YankField(key string) string {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return ""
	}
	switch key {
	case "l":
		v, _ := row.Data["local"].(string)
		return v
	case "r":
		v, _ := row.Data["remote"].(string)
		return v
	case "p":
		v, _ := row.Data["pid"].(string)
		return v
	case "n":
		v, _ := row.Data["process"].(string)
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
	return m.table.GetCurrentFilter() != ""
}

// ClearFilter implements Tab.
func (m *Model) ClearFilter() {
	m.table = m.table.WithFilterInputValue("")
}

// GoToRemotePeer navigates to the socket whose local address matches
// the current row's remote address, if the remote is localhost.
// Prefers a bidirectional match (peer's remote points back to us).
func (m *Model) GoToRemotePeer() bool {
	row := m.table.HighlightedRow()
	if row.Data == nil {
		return false
	}
	localAddr, _ := row.Data["raw_local_addr"].(string)
	localPort, _ := row.Data["raw_local_port"].(uint32)
	remoteAddr, _ := row.Data["raw_remote_addr"].(string)
	remotePort, _ := row.Data["raw_remote_port"].(uint32)
	if remoteAddr == "" || remotePort == 0 {
		return false
	}
	if !isLocalhost(remoteAddr) {
		return false
	}

	// Build rows in current display order
	rows := m.buildRows()
	if m.sort.Active() {
		m.sort.SortRows(rows)
	} else if m.navKey != "" {
		rows = m.reorderRows(rows, m.navKey, m.navVal)
	}

	// First pass: bidirectional match — peer's local matches our remote
	// AND peer's remote matches our local (the true connection peer).
	for i, r := range rows {
		rl, _ := r.Data["raw_local_addr"].(string)
		rlp, _ := r.Data["raw_local_port"].(uint32)
		rr, _ := r.Data["raw_remote_addr"].(string)
		rrp, _ := r.Data["raw_remote_port"].(uint32)
		if rl == remoteAddr && rlp == remotePort &&
			rr == localAddr && rrp == localPort {
			m.table = m.table.WithHighlightedRow(i)
			return true
		}
	}

	// Fallback: one-directional match on local address only
	target := util.FormatAddrPort(remoteAddr, remotePort)
	for i, r := range rows {
		local, _ := r.Data["local"].(string)
		if local == target {
			m.table = m.table.WithHighlightedRow(i)
			return true
		}
	}

	// Try wildcard match (*:port) for sockets bound to 0.0.0.0 or ::
	wildcard := util.FormatAddrPort("", remotePort)
	for i, r := range rows {
		local, _ := r.Data["local"].(string)
		if local == wildcard {
			m.table = m.table.WithHighlightedRow(i)
			return true
		}
	}

	return false
}

func isLocalhost(addr string) bool {
	return addr == "127.0.0.1" || addr == "::1" || addr == "localhost"
}
