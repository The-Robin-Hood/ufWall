package policy

import (
	"ufWall/internal/ufw"
	"ufWall/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View(policy ufw.Policy) string {

	incomingStyle := ui.GetPolicyStyle(m.styles, policy.DefaultIncoming)
	outgoingStyle := ui.GetPolicyStyle(m.styles, policy.DefaultOutgoing)
	routedStyle := ui.GetPolicyStyle(m.styles, policy.DefaultRouted)

	incomingLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Incoming:"),
		incomingStyle.Render(policy.DefaultIncoming),
	)

	outgoingLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Outgoing:"),
		outgoingStyle.Render(policy.DefaultOutgoing),
	)

	routedLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Routed:"),
		routedStyle.Render(policy.DefaultRouted),
	)

	sectionActiveNoMenu := m.menu == nil && m.active

	content := lipgloss.JoinVertical(
		lipgloss.Top,
		ui.InsertCursor(incomingLine, m.cursorLine == 0 && sectionActiveNoMenu, m.styles),
		ui.InsertCursor(outgoingLine, m.cursorLine == 1 && sectionActiveNoMenu, m.styles),
		ui.InsertCursor(routedLine, m.cursorLine == 2 && sectionActiveNoMenu, m.styles),
	)

	return ui.TitledBox("Default Policies", content, m.styles, 38, m.active)
}
