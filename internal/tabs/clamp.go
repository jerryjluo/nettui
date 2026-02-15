package tabs

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

var (
	rowDown = key.NewBinding(key.WithKeys("down", "j"))
	rowUp   = key.NewBinding(key.WithKeys("up", "k"))
)

// ClampedUpdate forwards msg to the table but suppresses row movement at
// boundaries so the cursor never wraps around.
func ClampedUpdate(t table.Model, msg tea.Msg) (table.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		total := t.TotalRows()
		idx := t.GetHighlightedRowIndex()
		if key.Matches(msg, rowDown) && idx >= total-1 {
			return t, nil
		}
		if key.Matches(msg, rowUp) && idx <= 0 {
			return t, nil
		}
	}
	return t.Update(msg)
}
