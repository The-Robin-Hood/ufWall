package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit        key.Binding
	Refresh     key.Binding
	
	NextSection key.Binding
	PrevSection key.Binding
	CursorUp    key.Binding
	CursorDown  key.Binding
}

var Keys = KeyMap{
	CursorUp: key.NewBinding(
		key.WithKeys("up"),
	),

	CursorDown: key.NewBinding(
		key.WithKeys("down"),
	),

	PrevSection: key.NewBinding(
		key.WithKeys("shift+tab", "left"),
	),

	NextSection: key.NewBinding(
		key.WithKeys("tab", "right"),
	),

	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),

	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
}
