package stats

import (
	"github.com/The-Robin-Hood/ufWall/internal/ui"
)

type Model struct {
	styles     ui.Styles
	totalOpts  int
	cursorLine int
	menu       *ui.Menu
	active     bool
}

func New(styles ui.Styles) Model {
	return Model{
		styles:     styles,
		cursorLine: 0,
		menu:       nil,
		active:     true,
		totalOpts:  2,
	}
}

func (m *Model) Focus() {
	m.cursorLine = 0
	m.active = true
}

func (m *Model) Blur() {
	m.active = false
	m.menu = nil
	m.cursorLine = 0
}

func (m Model) HasOpenMenu() bool {
	return m.menu != nil
}

func (m Model) GetMenu() *ui.Menu {
	return m.menu
}
