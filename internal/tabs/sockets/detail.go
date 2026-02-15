package sockets

import (
	"fmt"
	"strings"

	"github.com/jerryluo/nettui/internal/model"
)

const labelWidth = 14

func detailContent(rowData map[string]interface{}, panelWidth int) string {
	if rowData == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(model.PanelHeaderStyle.Render("Socket Details"))
	b.WriteString("\n\n")

	fields := []struct {
		label string
		key   string
		wrap  bool
	}{
		{"Protocol", "proto", false},
		{"Local", "local", false},
		{"Remote", "remote", false},
		{"State", "state", false},
		{"PID", "pid", false},
		{"Process", "process", false},
		{"Command", "command", true},
	}

	// Available width for values: panel minus border/padding (4) minus label.
	valWidth := panelWidth - 4 - labelWidth
	if valWidth < 20 {
		valWidth = 20
	}

	for _, f := range fields {
		val := fmt.Sprintf("%v", rowData[f.key])
		label := model.PanelLabelStyle.Render(fmt.Sprintf("%-*s", labelWidth, f.label))

		if f.wrap && len(val) > valWidth {
			lines := wrapText(val, valWidth)
			b.WriteString(label)
			b.WriteString(model.PanelValueStyle.Render(lines[0]))
			b.WriteString("\n")
			indent := strings.Repeat(" ", labelWidth)
			for _, line := range lines[1:] {
				b.WriteString(indent)
				b.WriteString(model.PanelValueStyle.Render(line))
				b.WriteString("\n")
			}
		} else {
			b.WriteString(label)
			b.WriteString(model.PanelValueStyle.Render(val))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// wrapText breaks s into lines of at most width characters, splitting on spaces.
func wrapText(s string, width int) []string {
	if width <= 0 || len(s) <= width {
		return []string{s}
	}
	var lines []string
	for len(s) > 0 {
		if len(s) <= width {
			lines = append(lines, s)
			break
		}
		// Find last space within width.
		cut := strings.LastIndex(s[:width], " ")
		if cut <= 0 {
			// No space found; hard-break at width.
			cut = width
		}
		lines = append(lines, s[:cut])
		s = strings.TrimLeft(s[cut:], " ")
	}
	return lines
}
