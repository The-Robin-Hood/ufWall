package app

import "ufWall/internal/ufw"

type model struct {
	status ufw.UFWStatus
	styles Styles
	width  int
	height int
}

func InitialModel() model {
	return model{
		status: ufw.GetStatus(),
		styles: NewStyles(),
		width:  80,
		height: 24,
	}
}
