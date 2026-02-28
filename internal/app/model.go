package app

import (
	"ufWall/internal/sections/policy"
	"ufWall/internal/sections/rules"
	"ufWall/internal/sections/stats"
	"ufWall/internal/ufw"
	"ufWall/internal/ui"
)

type model struct {
	stats ufw.Stats
	policy ufw.Policy
	rules []ufw.Rule
  err error 

	activeSection int

	width  int
	height int
	styles ui.Styles

	statsSection stats.Model	
	policySection policy.Model
	rulesSection  rules.Model

	openMenu bool
}

func InitialModel() model {
	data := ufw.GetUFWData()
	styles := ui.NewStyles()
	return model{
		stats:  data.Stats,
		rules:  data.Rules,
		policy: data.Policy,

		styles: styles,
		width:  87,
		height: 30,

		statsSection: stats.New(styles),
		policySection: policy.New(styles),
		rulesSection: rules.New(styles),

		openMenu : false,
	}
}
