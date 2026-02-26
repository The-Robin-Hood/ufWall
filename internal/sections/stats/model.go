package stats

import (
	"ufWall/internal/ui"
)

type Model struct {
	styles ui.Styles
}

func New(styles ui.Styles) Model {
	return Model{
		styles,
	}
}
