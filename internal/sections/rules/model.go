package rules

import (
	"ufWall/internal/ui"
)

type Model struct {
	styles     ui.Styles
	totalOpts  int
	cursorLine int 
	showMenu   bool
	menu       *ui.Menu
	active     bool
}

func New(styles ui.Styles) Model {
	return Model{
		styles:     styles,
		cursorLine: 0,
		showMenu:   false,
		menu:       nil,
		active:     false,
		totalOpts:  2,
	}
}

func (m *Model) Focus() {
	m.cursorLine = 0
	m.active = true
}

func (m *Model) Blur() {
	m.active = false
	m.showMenu = false
	m.cursorLine = 0
}

func (m Model) HasOpenMenu() bool {
	return m.menu != nil
}

func (m Model) GetMenu() *ui.Menu {
	return m.menu
}
