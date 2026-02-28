package stats

import (
	"ufWall/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		if m.menu != nil {
			return m.handleMenuInput(key)
		}
		return m.handleNavigation(key)
	}
	return m, nil
}

func (m Model) handleNavigation(key tea.KeyMsg) (Model, tea.Cmd) {
	switch key.String() {
	case "up", "k":
		if m.cursorLine > 0 {
			m.cursorLine--
		}

	case "down", "j":
		if m.cursorLine < 1  {
			m.cursorLine++
		}

	case "enter":
		return m.openMenu()
	}

	return m, nil
}

func (m Model) openMenu() (Model, tea.Cmd) {
	var options []string

	switch m.cursorLine {
	case 0: 
		m.menuType = MenuFirewall
		options = []string{"Enable Firewall", "Disable Firewall"}

	case 1: 		
		m.menuType = MenuLogging
		options = []string{"off", "low", "medium", "high", "full"}
	
	}

	menu := ui.NewMenu(options, m.styles)
	m.menu = &menu
	return m, nil
}

func (m Model) handleMenuInput(key tea.KeyMsg) (Model, tea.Cmd) {
	switch key.String() {
	case "up", "k":
		m.menu.Up()

	case "down", "j":
		m.menu.Down()

	case "enter":
		return m.executeMenuAction()

	case "esc":
		m.menu = nil
		m.menuType = MenuNone
	}

	return m, nil
}

func (m Model) executeMenuAction() (Model, tea.Cmd) {
	if m.menu == nil {
		return m, nil
	}

	// selectedOption := m.menu.SelectedOption()

	switch m.menuType {
	case MenuFirewall:
		// Handle enable/disable
		// TODO: Call UFW command based on selectedOption

	case MenuLogging:
		// Handle logging level change
		// TODO: Call UFW command with selectedOption (off/on/low/etc)
	}

	// Close menu
	m.menu = nil
	m.menuType = MenuNone

	return m, nil
}

