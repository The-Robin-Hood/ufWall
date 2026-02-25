package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit              key.Binding
	Refresh           key.Binding
	MoveToNextSection key.Binding
	MoveToPrevSection key.Binding
}

var Keys = KeyMap{
	MoveToPrevSection: key.NewBinding(
		key.WithKeys("shift+tab", "left"),
	),
	MoveToNextSection: key.NewBinding(
		key.WithKeys("tab", "right"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
}
