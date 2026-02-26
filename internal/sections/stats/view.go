package stats

import (
	"strconv"
	"strings"
	"ufWall/internal/ufw"
	"ufWall/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View(stats ufw.Stats, active bool) string {
	var statusText string
	var statusStyle lipgloss.Style

	if stats.Active {
		statusText = "ACTIVE"
		statusStyle = m.styles.ActiveStatus
	} else {
		statusText = "INACTIVE"
		statusStyle = m.styles.InactiveStatus
	}

	statusLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Status:"),
		statusStyle.Render(statusText),
	)

	loggingLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Logging:"),
		m.styles.Value.Render(strings.ToUpper(stats.Logging)),
	)

	totalRulesLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Total Rules:"),
		m.styles.Value.Render(strconv.Itoa(stats.TotalRules)),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		statusLine,
		loggingLine,
		totalRulesLine,
	)

	return ui.TitledBox("Firewall Stats", content, m.styles, -1, active)
}
