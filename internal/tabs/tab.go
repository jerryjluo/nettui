package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jerryluo/nettui/internal/model"
	"github.com/jerryluo/nettui/internal/data"
)

// Tab is the interface that all tab models must implement.
type Tab interface {
	tea.Model

	// SetData updates the tab with a new store snapshot.
	SetData(store *data.Store)

	// SetSize sets the available dimensions for this tab.
	SetSize(width, height int)

	// TabID returns this tab's identifier.
	TabID() model.TabID

	// SelectedRow returns a summary string of the currently selected row (for copy).
	SelectedRow() string

	// DetailContent returns the detail view content for the currently selected row.
	DetailContent() string

	// CrossRef returns a CrossRefMsg for the currently selected row, or nil.
	CrossRef() *model.CrossRefMsg

	// NavigateTo applies a cross-reference filter/highlight.
	NavigateTo(key, val string)

	// IsFiltering returns true if the tab is in filter mode.
	IsFiltering() bool

	// HasActiveFilter returns true if a search filter is currently applied.
	HasActiveFilter() bool

	// ClearFilter clears the current search filter text.
	ClearFilter()

	// SortHint returns the chord hint for sort keys on this tab.
	SortHint() string

	// ApplySort applies or toggles sort by the given chord key.
	ApplySort(key string)

	// SortLabel returns the current sort indicator (e.g. "[↑Proto]"), or "".
	SortLabel() string

	// SetPanelWidth stores the detail panel width for content wrapping.
	SetPanelWidth(width int)

	// YankHint returns the chord hint for yank keys on this tab.
	YankHint() string

	// YankField returns the value for a specific yank target key, or "".
	YankField(key string) string
}
