// internal/ui/menu.go
package ui

import (
	"strings"
	"github.com/The-Robin-Hood/ufWall/internal/keys"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuItem struct {
	Label  string
	Action func() tea.Cmd
}

type Menu struct {
	Options  []MenuItem
	Selected int
	styles   Styles
}

func NewMenu(options []MenuItem, styles Styles) Menu {
	return Menu{
		Options:  options,
		Selected: 0,
		styles:   styles,
	}
}

func (m *Menu) Up() {
	if m.Selected > 0 {
		m.Selected--
	}
}

func (m *Menu) Down() {
	if m.Selected < len(m.Options)-1 {
		m.Selected++
	}
}

func (m Menu) ExecuteSelected() tea.Cmd {
	if m.Selected >= 0 && m.Selected < len(m.Options) {
		if m.Options[m.Selected].Action != nil {
			return m.Options[m.Selected].Action()
		}
	}
	return nil
}

func (m Menu) View(styles Styles) string {
	var lines []string

	for i, option := range m.Options {
		prefix := "  "
		style := m.styles.Value

		if i == m.Selected {
			prefix = "▶ "
			style = m.styles.ActiveStatus.Bold(true)
		}

		lines = append(lines, prefix+style.Render(option.Label))
	}

	content := strings.Join(lines, "\n")

	return styles.Menu.Render(content)
}

func (m *Menu) Update(msg tea.Msg) (bool) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Bindings.CursorUp):
			m.Up()
		case key.Matches(msg, keys.Bindings.CursorDown):
			m.Down()
		case key.Matches(msg, keys.Bindings.Execute):
			m.ExecuteSelected()
			return true
		case key.Matches(msg, keys.Bindings.Quit):
			return true
		}
	}
	return false
}

func MakeMenuItems(labels []string, handler func(string) tea.Cmd) []MenuItem {
	items := make([]MenuItem, len(labels))

	for i, label := range labels {
		l := label // capture loop variable!
		items[i] = MenuItem{
			Label: l,
			Action: func() tea.Cmd {
				return handler(l)
			},
		}
	}

	return items
}
