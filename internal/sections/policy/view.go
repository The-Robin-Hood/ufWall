package policy

import (
	"strings"
	"ufWall/internal/ufw"
	"ufWall/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) getPolicyStyle(policy string) lipgloss.Style {
	switch strings.ToUpper(policy) {
	case "ALLOW":
		return m.styles.AllowPolicy
	case "DENY":
		return m.styles.DenyPolicy
	case "REJECT":
		return m.styles.RejectPolicy
	default:
		return m.styles.Value
	}
}

func (m Model) View(policy ufw.Policy) string {

	incomingStyle := m.getPolicyStyle(policy.DefaultIncoming)
	outgoingStyle := m.getPolicyStyle(policy.DefaultOutgoing)
	routedStyle := m.getPolicyStyle(policy.DefaultRouted)

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

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		incomingLine,
		outgoingLine,
		routedLine,
	)
	return ui.TitledBox("Default Policies", content, m.styles, -1, m.active)
}
