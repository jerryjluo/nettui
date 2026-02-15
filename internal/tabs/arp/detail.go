package arp

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

	b.WriteString(model.PanelHeaderStyle.Render("ARP Entry Details"))
	b.WriteString("\n\n")

	fields := []struct {
		label string
		key   string
	}{
		{"IP", "ip"},
		{"Hostname", "hostname"},
		{"MAC", "mac"},
		{"Interface", "iface"},
		{"Flags", "flags"},
		{"Type", "type"},
	}

	for _, f := range fields {
		val := fmt.Sprintf("%v", rowData[f.key])
		if val == "" || val == "<nil>" {
			continue
		}
		b.WriteString(model.PanelLabelStyle.Render(fmt.Sprintf("%-14s", f.label)))
		b.WriteString(model.PanelValueStyle.Render(val))
		b.WriteString("\n")
	}

	return b.String()
}
