package rules

import (
	"ufWall/internal/ui"
)

type Model struct {
	styles ui.Styles
	active bool
}

func New(styles ui.Styles) Model {
	return Model{
		styles: styles,
		active: false,
	}
}

func (m *Model) Focus() {
	// m.cursorLine = 0
	m.active = true
}

func (m *Model) Blur() {
	m.active = false
	// m.showMenu = false
	// m.cursorLine = 0
}
