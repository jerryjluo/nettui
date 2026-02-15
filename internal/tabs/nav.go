package tabs

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
)

// FilterNavRows returns only rows where Data[key] matches val.
func FilterNavRows(rows []table.Row, key, val string) []table.Row {
	filtered := make([]table.Row, 0, len(rows))
	for _, r := range rows {
		if fmt.Sprintf("%v", r.Data[key]) == val {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
