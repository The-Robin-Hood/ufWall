package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var b strings.Builder

	// // If there's an error, show it prominently
	if m.status.Error != nil {
		errorBox := m.styles.Error.Render(fmt.Sprintf("%v", m.status.Error))
		footer := m.renderFooter()
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			errorBox,
			"",
			footer,
		)
	}

	// Title
	title := m.styles.Title.Width(m.width).Render("Firewall Manager")

	b.WriteString(title)
	b.WriteString("\n\n")

	// Status Section
	statusSection := m.renderStatusSection()
	b.WriteString(statusSection)
	b.WriteString("\n")

	// Default Policies Section
	policiesSection := m.renderPoliciesSection()
	b.WriteString(policiesSection)
	b.WriteString("\n")

	// Rules Count Section
	rulesSection := m.renderRulesSection()
	b.WriteString(rulesSection)
	b.WriteString("\n")

	// Footer with keybindings
	footer := m.renderFooter()
	b.WriteString(footer)

	return b.String()
}
