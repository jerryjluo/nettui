package app

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines global keybindings.
type KeyMap struct {
	Quit      key.Binding
	NextTab   key.Binding
	PrevTab   key.Binding
	Tab1      key.Binding
	Tab2      key.Binding
	Tab3      key.Binding
	Tab4      key.Binding
	Tab5      key.Binding
	Tab6      key.Binding
	Up        key.Binding
	Down      key.Binding
	Filter    key.Binding
	Enter     key.Binding
	Escape    key.Binding
	GoTo      key.Binding
	Copy      key.Binding
	DNS       key.Binding
	Help      key.Binding
	Refresh   key.Binding
	PageUp      key.Binding
	PageDown    key.Binding
	ProtoFilter key.Binding
}

// DefaultKeyMap returns the default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("l", "tab"),
			key.WithHelp("l/tab", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("h", "shift+tab"),
			key.WithHelp("h/shift+tab", "prev tab"),
		),
		Tab1: key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "Interfaces")),
		Tab2: key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "Routes")),
		Tab3: key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "Sockets")),
		Tab4: key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "Unix")),
		Tab5: key.NewBinding(key.WithKeys("5"), key.WithHelp("5", "Processes")),
		Tab6: key.NewBinding(key.WithKeys("6"), key.WithHelp("6", "Firewall")),
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/up", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/down", "down"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "detail panel"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "close/clear"),
		),
		GoTo: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "go to ref"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy"),
		),
		DNS: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "toggle DNS"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "page down"),
		),
		ProtoFilter: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "protocol filter"),
		),
	}
}
