package app

import (
	"fmt"
	"strconv"
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

	totalRulesLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.styles.Label.Render("Total Rules:"),
		m.styles.Value.Render(strconv.Itoa(len(m.rules))),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		statusLine,
		loggingLine,
		totalRulesLine,
	)

	return RenderBoxWithTitle("Firewall Stats", content, m.styles, -1)
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

func (m model) renderRulesSection() string {
	var rows []string

	header := fmt.Sprintf(
	"%-3s │ %-6s │ %-5s │ %-16s │ %-5s │ %-16s │ %-5s",
		"#", "Action", "Proto", "Source", "sPort","Destination","dPort")

	headerContent := m.styles.Label.UnsetWidth().Render(header)
	line := strings.Repeat("─", lipgloss.Width(headerContent))
	rows = append(rows, headerContent, m.styles.Label.UnsetWidth().Render(line))

	for _, r := range m.rules {
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
		rows = append(rows, m.styles.Value.Render(row))
	}
	table := strings.Join(rows, "\n")
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		table,
	)
	return RenderBoxWithTitle("Active Rules", content, m.styles, -1)
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
