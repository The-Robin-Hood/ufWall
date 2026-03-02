package ui

import (
	"strings"
	"ufWall/internal/sections"
)

func Footer(styles Styles, activeSection int, width int) string {
	var keys []string

	switch activeSection {
	case sections.StatsSection:
		keys = []string{"↑↓: navigate", "space: toggle menu"}
	case sections.PolicySection:
		keys = []string{"↑↓: navigate", "space: toggle allow/deny"}

	case sections.RulesSection:
		keys = []string{"enter: actions", "i: info", "a: add rule", "d: delete rule", "s: switch IPv4/IPv6"}
	default:
		keys = []string{"tab: next section", "r: refresh"}
	}

	if len(keys) < 3 {
		keys = append(keys, "q: quit")
	}
	footerText := strings.Join(keys, " | ")
	return styles.Footer.Width(width).Render(footerText)
}
