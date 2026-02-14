package processes

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

	b.WriteString(model.PanelHeaderStyle.Render("Process Details"))
	b.WriteString("\n\n")

	fields := []struct {
		label string
		key   string
	}{
		{"PID", "pid"},
		{"Name", "name"},
		{"User", "user"},
		{"Connections", "conns"},
	}

	for _, f := range fields {
		val := fmt.Sprintf("%v", rowData[f.key])
		b.WriteString(model.PanelLabelStyle.Render(fmt.Sprintf("%-14s", f.label)))
		b.WriteString(model.PanelValueStyle.Render(val))
		b.WriteString("\n")
	}

	return b.String()
}
