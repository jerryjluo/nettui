package interfaces

import (
	"fmt"
	"strings"

	"github.com/jerryluo/nettui/internal/model"
)

func detailContent(rowData map[string]interface{}) string {
	if rowData == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(model.PanelHeaderStyle.Render("Interface Details"))
	b.WriteString("\n\n")

	fields := []struct {
		label string
		key   string
	}{
		{"Name", "name"},
		{"IPs", "addrs"},
		{"MAC", "mac"},
		{"MTU", "mtu"},
		{"Status", "status"},
		{"TX Bytes", "tx_bytes"},
		{"RX Bytes", "rx_bytes"},
		{"TX Rate", "tx_rate"},
		{"RX Rate", "rx_rate"},
		{"TX Packets", "tx_pkts"},
		{"RX Packets", "rx_pkts"},
		{"Flags", "flags"},
	}

	for _, f := range fields {
		val := fmt.Sprintf("%v", rowData[f.key])
		b.WriteString(model.PanelLabelStyle.Render(fmt.Sprintf("%-12s", f.label)))
		b.WriteString(model.PanelValueStyle.Render(val))
		b.WriteString("\n")
	}

	return b.String()
}
