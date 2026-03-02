package rules

import (
	"fmt"
	"strings"
	"ufWall/internal/ufw"
	"ufWall/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View(rules []ufw.Rule) string {
	var rows []string

	header := fmt.Sprintf(
		"%-3s │ %-6s │ %-5s │ %-16s │ %-5s │ %-16s │ %-5s",
		"#", "Action", "Proto", "Source", "sPort", "Destination", "dPort")

	headerContent := m.styles.Label.UnsetWidth().Render(header)
	line := strings.Repeat("─", lipgloss.Width(headerContent))
	rows = append(rows, "  "+headerContent, "  "+m.styles.Label.UnsetWidth().Render(line))

	sectionActiveNoMenu := m.menu == nil && m.active

	for i, r := range rules {
		row := fmt.Sprintf(
			"%-3d │ %-6s │ %-5s │ %-16s │ %-5s │ %-16s │ %-5s",
			r.Num,
			r.Action,
			r.ToProtocol,
			r.FromSource,
			r.FromPort,
			r.ToDest,
			r.ToPort,
		)
		rows = append(rows, ui.InsertCursor(row, m.cursorLine == i && sectionActiveNoMenu, m.styles))
	}
	table := strings.Join(rows, "\n")
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		table,
	)
	return ui.TitledBox("Active Rules", content, m.styles, -1, m.active)
}
