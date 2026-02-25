package app

import "ufWall/internal/ufw"

type model struct {
	status ufw.UFWStatus
	rules  []ufw.Rule

	activeSection int
	selectedRule int

	styles Styles
	width  int
	height int
}

func InitialModel() model {
	s,r := ufw.GetStatus()
	return model{
		status: s,
		rules: r,
		styles: NewStyles(),
		width:   87,
		height: 30,	
	}
}
