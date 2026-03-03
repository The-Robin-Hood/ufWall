package stats

import (
	"log"
	"github.com/The-Robin-Hood/ufWall/internal/keys"
	"github.com/The-Robin-Hood/ufWall/internal/ufw"
	"github.com/The-Robin-Hood/ufWall/internal/ui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
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
				if m.cursorLine < m.totalOpts - 1 {
					m.cursorLine++
				}
			case key.Matches(msg, keys.Bindings.Execute):
				return m.openMenu()
			}
		}
	}
	return m, nil
}

func (m Model) openMenu() (Model, tea.Cmd) {
	var options []ui.MenuItem

	switch m.cursorLine {
	case 0:
		options = ui.MakeMenuItems(
			[]string{"Enable Firewall", "Disable Firewall"},
			func(label string) tea.Cmd {
				switch label {
				case "Enable Firewall":
					log.Println("Enabling Firewall")
					ufw.Enable()
				case "Disable Firewall":
					log.Println("Disabling Firewall")
					ufw.Disable()
					return nil
				}
				return nil
			},
		)

	case 1:
		options = ui.MakeMenuItems(
			[]string{"off", "low", "medium", "high", "full"},
			func(level string) tea.Cmd {
				log.Println("Setting Log Level :", level)
				ufw.SetLogging(level)
				return nil
			},
		)
	}

	menu := ui.NewMenu(options, m.styles)
	m.menu = &menu
	return m, nil
}
