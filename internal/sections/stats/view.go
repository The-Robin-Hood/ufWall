package stats

import (
	"strconv"
	"strings"
	"github.com/The-Robin-Hood/ufWall/internal/ufw"
	"github.com/The-Robin-Hood/ufWall/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View(stats ufw.Stats) string {
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
		m.styles.Label.Width(15).Render("Status:"),
		statusStyle.Render(statusText),
	)

	loggingLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Width(15).Render("Logging:"),
		m.styles.Value.Render(strings.ToUpper(stats.Logging)),
	)

	rulesLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Width(15).Render("Total Rules:"),
		m.styles.Value.Render(strconv.Itoa(stats.TotalRules)),
	)

	sectionActiveNoMenu := m.menu == nil && m.active

	content := lipgloss.JoinVertical(
		lipgloss.Top,
		ui.InsertCursor(statusLine, m.cursorLine == 0 && sectionActiveNoMenu, m.styles),
		ui.InsertCursor(loggingLine, m.cursorLine == 1 && sectionActiveNoMenu, m.styles),
		ui.InsertCursor(rulesLine, false, m.styles),
	)

	return ui.TitledBox("Firewall Stats", content, m.styles, 38, m.active)
}
