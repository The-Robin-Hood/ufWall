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
		keys = []string{"↑↓: navigate", "enter: actions", "i: info", "d: delete"}
	default:
		keys = []string{"tab: next section", "r: refresh"}
	}

	keys = append(keys, "q: quit")
	footerText := strings.Join(keys, "  |  ")
	return styles.Footer.Width(width).Render(footerText)
}
