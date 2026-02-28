// internal/ui/menu.go
package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	Options  []string
	Selected int
	styles   Styles
}

func NewMenu(options []string, styles Styles) Menu {
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

func (m Menu) SelectedOption() string {
	if m.Selected >= 0 && m.Selected < len(m.Options) {
		return m.Options[m.Selected]
	}
	return ""
}

func (m Menu) View() string {
	var lines []string

	for i, option := range m.Options {
		prefix := "  "
		style := m.styles.Value

		if i == m.Selected {
			prefix = "▶ "
			style = m.styles.ActiveStatus.Bold(true)
		}

		lines = append(lines, prefix+style.Render(option))
	}

	content := strings.Join(lines, "\n")

	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#89b4fa")).
		Padding(1, 2).
		Background(lipgloss.Color("#1e1e2e"))

	return menuStyle.Render(content)
}
