package stats

import (
	"ufWall/internal/ui"
)

type MenuType int

const (
	MenuNone     MenuType = iota
	MenuFirewall          // Enable/Disable
	MenuLogging           // off/on/low/medium/high/full
)

type Model struct {
	styles     ui.Styles
	cursorLine int // 0=status, 1=logging
	showMenu   bool
	menu       *ui.Menu
	menuType   MenuType
	active     bool
}

func New(styles ui.Styles) Model {
	return Model{
		styles:     styles,
		cursorLine: 0,
		showMenu:   false,
		menu:       nil,
		menuType:   MenuNone,
		active:     true,
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

