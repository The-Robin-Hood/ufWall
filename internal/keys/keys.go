package keys

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Quit    key.Binding
	Refresh key.Binding

	NextSection key.Binding
	PrevSection key.Binding
	CursorUp    key.Binding
	CursorDown  key.Binding
	Execute     key.Binding
	Info        key.Binding
	Delete      key.Binding
}

var Bindings = KeyMap{
	CursorUp: key.NewBinding(
		key.WithKeys("up", "k"),
	),

	CursorDown: key.NewBinding(
		key.WithKeys("down", "j"),
	),

	PrevSection: key.NewBinding(
		key.WithKeys("shift+tab", "left"),
	),

	NextSection: key.NewBinding(
		key.WithKeys("tab", "right"),
	),

	Refresh: key.NewBinding(
		key.WithKeys("r"),
	),

	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
	),

	Execute: key.NewBinding(
		key.WithKeys(" ", "enter"),
	),

	Info: key.NewBinding(
		key.WithKeys("i"),
	),

	Delete: key.NewBinding(
		key.WithKeys("d"),
	),
}

type RefreshMsg struct{}

func Refresh() tea.Cmd {
	return func() tea.Msg {
		return RefreshMsg{}
	}
}
