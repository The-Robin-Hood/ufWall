package app

import (
	"github.com/The-Robin-Hood/ufWall/internal/sections/policy"
	"github.com/The-Robin-Hood/ufWall/internal/sections/rules"
	"github.com/The-Robin-Hood/ufWall/internal/sections/stats"
	"github.com/The-Robin-Hood/ufWall/internal/ufw"
	"github.com/The-Robin-Hood/ufWall/internal/ui"
)

type model struct {
	stats     ufw.Stats
	policy    ufw.Policy
	rules     []ufw.Rule // (kept for backward compat)
	ipv4Rules []ufw.Rule
	ipv6Rules []ufw.Rule
	err       error

	activeSection int

	width  int
	height int
	styles ui.Styles

	statsSection  stats.Model
	policySection policy.Model
	rulesSection  rules.Model
}

func InitialModel() model {
	data := ufw.GetUFWData()
	styles := ui.NewStyles()
	return model{
		stats:     data.Stats,
		rules:     data.Rules,
		ipv4Rules: data.IPv4Rules,
		ipv6Rules: data.IPv6Rules,
		policy:    data.Policy,

		styles: styles,
		width:  87,
		height: 30,

		statsSection:  stats.New(styles),
		policySection: policy.New(styles),
		rulesSection:  rules.New(styles),
	}
}
