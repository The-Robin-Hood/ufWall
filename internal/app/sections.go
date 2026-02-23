package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderStatusSection() string {
	var statusText string
	var statusStyle lipgloss.Style

	if m.status.Active {
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
		m.styles.Value.Render(strings.ToUpper(m.status.Logging)),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		statusLine,
		loggingLine,
	)

	return m.styles.SectionBorder.Render(content)
}

func (m model) renderPoliciesSection() string {
	title := m.styles.SectionTitle.Render("Default Policies")

	incomingStyle := m.getPolicyStyle(m.status.DefaultIn)
	outgoingStyle := m.getPolicyStyle(m.status.DefaultOut)
	routedStyle := m.getPolicyStyle(m.status.DefaultRouted)

	incomingLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Incoming:"),
		incomingStyle.Render(m.status.DefaultIn),
	)

	outgoingLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Outgoing:"),
		outgoingStyle.Render(m.status.DefaultOut),
	)

	routedLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Routed:"),
		routedStyle.Render(m.status.DefaultRouted),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		incomingLine,
		outgoingLine,
		routedLine,
	)

	return m.styles.SectionBorder.Render(content)
}

func (m model) renderRulesSection() string {
	title := m.styles.SectionTitle.Render("Active Rules")

	ruleCountText := fmt.Sprintf("%d rule(s) configured", len(m.status.Rules))
	ruleCountLine := m.styles.Value.Render(ruleCountText)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		ruleCountLine,
	)

	return m.styles.SectionBorder.Render(content)
}

func (m model) getPolicyStyle(policy string) lipgloss.Style {
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

func (m model) renderFooter() string {
	keys := []string{
		m.styles.FooterKey.Render("r") + " refresh",
		m.styles.FooterKey.Render("?") + " help",
		m.styles.FooterKey.Render("q") + " quit",
	}

	footerText := strings.Join(keys, "  •  ")
	return m.styles.Footer.Width(m.width).Render(footerText)
}
