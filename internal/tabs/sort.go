package tabs

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/evertras/bubble-table/table"
)

// SortEntry describes one sortable column for a tab.
type SortEntry struct {
	Key     string // chord second key ("p")
	ColKey  string // display column key ("proto")
	SortKey string // RowData key to sort on ("raw_pid" for numeric)
	Label   string // display name ("Proto")
}

// SortState tracks the current sort column and direction.
type SortState struct {
	Col string // current sort display column key
	Key string // RowData key being sorted on
	Asc bool
}

// Apply finds the matching entry for the given chord key and updates state.
// If the same column is already active, it toggles direction.
// Returns true if the key matched an entry.
func (s *SortState) Apply(entries []SortEntry, key string) bool {
	for _, e := range entries {
		if e.Key == key {
			if s.Col == e.ColKey {
				s.Asc = !s.Asc
			} else {
				s.Col = e.ColKey
				s.Key = e.SortKey
				s.Asc = true
			}
			return true
		}
	}
	return false
}

// Clear resets the sort state.
func (s *SortState) Clear() {
	s.Col = ""
	s.Key = ""
	s.Asc = false
}

// Active returns true if a sort is currently applied.
func (s *SortState) Active() bool {
	return s.Col != ""
}

// Hint builds a chord hint string like "s→  p:Proto l:Local ...".
func Hint(entries []SortEntry) string {
	var parts []string
	for _, e := range entries {
		parts = append(parts, e.Key+":"+e.Label)
	}
	return "s→  " + strings.Join(parts, "  ")
}

// Label returns a display indicator like "[↑Proto]" or "[↓PID]", or "" if inactive.
func (s *SortState) Label() string {
	if s.Col == "" {
		return ""
	}
	arrow := "↑"
	if !s.Asc {
		arrow = "↓"
	}
	return "[" + arrow + s.Col + "]"
}

// SortRows sorts rows in-place by the current sort key.
// It tries numeric comparison first, falling back to string comparison.
func (s *SortState) SortRows(rows []table.Row) {
	if s.Col == "" || len(rows) == 0 {
		return
	}
	sortKey := s.Key
	asc := s.Asc
	sort.SliceStable(rows, func(i, j int) bool {
		a := rows[i].Data[sortKey]
		b := rows[j].Data[sortKey]
		cmp := compareValues(a, b)
		if asc {
			return cmp < 0
		}
		return cmp > 0
	})
}

// compareValues compares two RowData values, trying numeric first then string.
func compareValues(a, b interface{}) int {
	na, aOk := toFloat(a)
	nb, bOk := toFloat(b)
	if aOk && bOk {
		if na < nb {
			return -1
		}
		if na > nb {
			return 1
		}
		return 0
	}
	sa := fmt.Sprintf("%v", a)
	sb := fmt.Sprintf("%v", b)
	return strings.Compare(sa, sb)
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float64:
		return n, true
	case string:
		f, err := strconv.ParseFloat(n, 64)
		if err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}
