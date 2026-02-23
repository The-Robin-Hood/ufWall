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

	uptimeLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("UpTime:"),
		m.styles.Value.Render(m.status.UpTime),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		statusLine,
		loggingLine,
		uptimeLine,
	)

	return RenderBoxWithTitle("Firewall Status", content, m.styles, -1)
}

func (m model) renderPoliciesSection() string {

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
		incomingLine,
		outgoingLine,
		routedLine,
	)
	return RenderBoxWithTitle("Default Policies", content, m.styles, -1)
}

func (m model) renderActiveRulesCountSection() string {

	ruleCountText := fmt.Sprintf("%d rule(s) configured", len(m.status.Rules))
	ruleCountLine := m.styles.Value.Render(ruleCountText)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		ruleCountLine,
	)

	return RenderBoxWithTitle("Active Rules", content, m.styles,-1)
}

func (m model) renderAllRulesSection(width int) string {
	title := m.styles.SectionTitle.Render("All Rules")
	var rows []string

	header := fmt.Sprintf(
		"%-1s %-3s %-5s %-7s",
		"NUM", "TO", "ACTION", "FROM",
	)
	rows = append(rows, m.styles.Label.Render(header))
	// Rule rows
	for i, r := range m.status.Rules {
		if i > 1 {
			continue
		}
		row := fmt.Sprintf(
			"%-5d %-15s %-10s %-20s",
			r.Num, r.To, r.Action, r.From,
		)
		rows = append(rows, m.styles.Value.Render(row))
	}
	table := strings.Join(rows, "\n")
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		table,
	)
	return m.styles.SectionBorder.Width(width).Render(content)
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

func (m model) renderFooter(width int) string {
	keys := []string{
		m.styles.FooterKey.Render("r") + " refresh",
		m.styles.FooterKey.Render("?") + " help",
		m.styles.FooterKey.Render("q") + " quit",
	}

	footerText := strings.Join(keys, "  •  ")
	return m.styles.Footer.Width(width).Render(footerText)
}
