package rules

import (
	"log"
	"ufWall/internal/keys"
	"ufWall/internal/ufw"
	"ufWall/internal/ui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg, rules []ufw.Rule) (Model, tea.Cmd) {
	if m.menu != nil {
		if quit := m.menu.Update(msg); quit {
			m.menu = nil
			return m, keys.Refresh()
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Bindings.CursorUp):
				if m.cursorLine > 0 {
					m.cursorLine--
				}
			case key.Matches(msg, keys.Bindings.CursorDown):
				if m.cursorLine < len(rules)-1 {
					m.cursorLine++
				}
			case key.Matches(msg, keys.Bindings.Execute):
				if len(rules) > 0 {
					return m.openMenu(rules)
				}
			}
		}
	}
	return m, nil
}

func (m Model) openMenu(rules []ufw.Rule) (Model, tea.Cmd) {
	if len(rules) == 0 || m.cursorLine >= len(rules) {
		return m, nil
	}

	selectedRule := rules[m.cursorLine]
	log.Println(selectedRule)

	options := ui.MakeMenuItems(
		[]string{"Delete Rule"},
		func(label string) tea.Cmd {
			switch label {
			case "Delete Rule":
				log.Printf("Deleting rule #%d", selectedRule.Num)
				ufw.DeleteRule(selectedRule.Num)
			}
			return nil
		},
	)

	menu := ui.NewMenu(options, m.styles)
	m.menu = &menu
	return m, nil
}
