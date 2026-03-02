package rules

import (
	"fmt"
	"regexp"
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

	sectionActiveNoMenu := m.menu == nil && m.active && !m.showDetails && !m.showDeleteConfirm

	for i, r := range rules {
		action := fmt.Sprintf("%-6s", r.Action)
		row := fmt.Sprintf(
			"%-3d │ %6s │ %-5s │ %-16s │ %-5s │ %-16s │ %-5s",
			r.Num,
			action,
			r.ToProtocol,
			r.FromSource,
			r.FromPort,
			r.ToDest,
			r.ToPort,
		)
		rows = append(rows, ui.InsertCursorRulesSection(row, m.cursorLine == i && sectionActiveNoMenu, m.styles, r.Action))
	}
	table := strings.Join(rows, "\n")
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		table,
	)
	return ui.TitledBox("Active Rules", content, m.styles, -1, m.active)
}

// DetailView renders the rule detail overlay
func (m Model) DetailView() string {
	if m.detailRule == nil {
		return ""
	}

	r := m.detailRule

	// Build detail content
	var lines []string
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Rule #:"), m.styles.Value.Render(fmt.Sprintf("%d", r.Num))))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Action:"), ui.GetPolicyStyle(m.styles, r.Action).Render(r.Action)))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Protocol:"), m.styles.Value.Render(r.ToProtocol)))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Source:"), m.styles.Value.Render(r.FromSource)))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Source Port:"), m.styles.Value.Render(r.FromPort)))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Destination:"), m.styles.Value.Render(r.ToDest)))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Dest Port:"), m.styles.Value.Render(r.ToPort)))
	lines = append(lines, "")

	if r.Comment != "" {
		lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Comment:"), m.styles.Value.Render(r.Comment)))
		lines = append(lines, "")
	}

	re := regexp.MustCompile(`\s+`)
	rawText := re.ReplaceAllString(r.Raw, " ")
	lines = append(lines, fmt.Sprintf("  %s", m.styles.Label.Render("Raw:")))
	lines = append(lines, fmt.Sprintf("  %s", m.styles.Value.Render(rawText)))
	lines = append(lines, "")
	text := m.styles.Label.
		PaddingTop(1).
		UnsetWidth().
		Render("[Press any key to close]")

	centered := lipgloss.Place(
		50,
		1,
		lipgloss.Center,
		lipgloss.Center,
		text,
	)

	lines = append(lines, centered)

	lines = append(lines, "")

	content := strings.Join(lines, "\n")

	title := fmt.Sprintf("Rule #%d Details", r.Num)
	return ui.TitledBox(title, content, m.styles, -1, true)
}

func (m Model) DeleteConfirmView() string {
	if m.deleteRule == nil {
		return ""
	}

	r := m.deleteRule

	var lines []string
	lines = append(lines, "")
	lines = append(lines, m.styles.Error.Render("  Are you sure you want to delete this rule?"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Rule #:"), m.styles.Value.Render(fmt.Sprintf("%d", r.Num))))
	lines = append(lines, fmt.Sprintf("  %s  %s %s", m.styles.Label.Render("Action:"), ui.GetPolicyStyle(m.styles, r.Action).Render(r.Action), m.styles.Value.Render(r.ToPort)))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("From:"), m.styles.Value.Render(r.FromSource)))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("To:"), m.styles.Value.Render(r.ToDest)))
	lines = append(lines, "")
	text := m.styles.Label.
		PaddingTop(1).
		UnsetWidth().
		Render("[y] Yes, delete  [n/Esc] Cancel")

	centered := lipgloss.Place(
		50,
		1,
		lipgloss.Center,
		lipgloss.Center,
		text,
	)

	lines = append(lines, centered)
	lines = append(lines, "")

	content := strings.Join(lines, "\n")

	return ui.TitledBox("Confirm Delete", content, m.styles, -1, true)
}
