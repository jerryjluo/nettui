package unixsockets

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

	b.WriteString(model.PanelHeaderStyle.Render("Unix Socket Details"))
	b.WriteString("\n\n")

	fields := []struct {
		label string
		key   string
	}{
		{"Path", "path"},
		{"Type", "type"},
		{"PID", "pid"},
		{"Process", "process"},
		{"FD", "fd"},
	}

	for _, f := range fields {
		val := fmt.Sprintf("%v", rowData[f.key])
		b.WriteString(model.PanelLabelStyle.Render(fmt.Sprintf("%-12s", f.label)))
		b.WriteString(model.PanelValueStyle.Render(val))
		b.WriteString("\n")
	}

	return b.String()
}
