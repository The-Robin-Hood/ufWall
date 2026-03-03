package rules

import (
	"fmt"
	"regexp"
	"strings"
	"github.com/The-Robin-Hood/ufWall/internal/ufw"
	"github.com/The-Robin-Hood/ufWall/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

type RulesData struct {
	IPv4 []ufw.Rule
	IPv6 []ufw.Rule
}

func (m Model) View(data RulesData) string {
	sectionActiveNoMenu := m.menu == nil && m.active && !m.showDetails && !m.showDeleteConfirm

	ipTable := m.renderTable(data.IPv4, data.IPv6, sectionActiveNoMenu)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		ipTable,
	)

	return ui.TitledBox("Active Rules", content, m.styles, -1, m.active, 13)
}

func (m Model) renderTable(ipV4Rules []ufw.Rule, ipV6Rules []ufw.Rule, isActive bool) string {
	var rows []string

	titleStyle := m.styles.Label
	if isActive {
		titleStyle = titleStyle.Bold(true).Foreground(lipgloss.Color("12"))
	}

	activeTitle := "IPv4 Rules"
	inActiveTitle := "IPv6 Rules"

	rules := ipV4Rules
	scrollOffset := m.ipv4ScrollOffset
	cursorLine := m.ipv4CursorLine

	maxIPv4 := len(fmt.Sprintf("%s (%d)", "IPv4 Rules", 9999))
	maxIPv6 := len(fmt.Sprintf("%s (%d)", "IPv6 Rules", 9999))

	padTitle := func(title string, count int, width int) string {
		return fmt.Sprintf("%-*s", width, fmt.Sprintf("%s (%d)", title, count))
	}

	tabTitle := titleStyle.Render(padTitle(activeTitle, len(ipV4Rules), maxIPv4)) + "|  " +
		m.styles.Label.Foreground(lipgloss.Color("241")).Render(padTitle(inActiveTitle, len(ipV6Rules), maxIPv6))

	if m.activeTable == IPv6Table {
		scrollOffset = m.ipv6ScrollOffset
		cursorLine = m.ipv6CursorLine
		rules = ipV6Rules
		tabTitle = m.styles.Label.Foreground(lipgloss.Color("241")).Render(padTitle(activeTitle, len(ipV4Rules), maxIPv4)) + "|  " +
			titleStyle.Render(padTitle(inActiveTitle, len(ipV6Rules), maxIPv6))
	}

	if len(rules) > MaxVisibleRules {
		scrollInfo := m.renderScrollIndicator(len(rules), scrollOffset)
		lineWidth := lipgloss.Width(tabTitle)
		scrollInfoRight := lipgloss.PlaceHorizontal(lineWidth, lipgloss.Right, scrollInfo)
		tabTitle = lipgloss.JoinHorizontal(lipgloss.Left, tabTitle, scrollInfoRight)
	}

	rows = append(rows, tabTitle)
	rows = append(rows, "")

	if len(rules) == 0 {
		emptyMsg := m.styles.Label.Foreground(lipgloss.Color("241")).Render("  (no rules)")
		rows = append(rows, emptyMsg)
		return strings.Join(rows, "\n")
	}

	header := fmt.Sprintf(
		"%-3s │ %-6s │ %-5s │ %-16s │ %-5s │ %-16s │ %-5s",
		"#", "Action", "Proto", "Source", "sPort", "Destination", "dPort")

	headerContent := m.styles.Label.UnsetWidth().Render(header)
	line := strings.Repeat("─", lipgloss.Width(headerContent))
	rows = append(rows, "  "+headerContent, "  "+m.styles.Label.UnsetWidth().Render(line))

	startIdx := scrollOffset
	endIdx := min(scrollOffset+MaxVisibleRules, len(rules))

	for i := startIdx; i < endIdx; i++ {
		r := rules[i]
		action := fmt.Sprintf("%-6s", r.Action)
		row := fmt.Sprintf(
			"%-3d │ %6s │ %-5s │ %-16s │ %-5s │ %-16s │ %-5s",
			r.Num,
			action,
			r.ToProtocol,
			truncate(r.FromSource, 16),
			truncate(r.FromPort, 5),
			truncate(r.ToDest, 16),
			truncate(r.ToPort, 5),
		)
		rows = append(rows, ui.InsertCursorRulesSection(row, cursorLine == i && isActive, m.styles, r.Action))
	}
	return strings.Join(rows, "\n")
}

func (m Model) renderScrollIndicator(totalRules int, scrollOffset int) string {
	current := scrollOffset + 1
	end := min(scrollOffset+MaxVisibleRules, totalRules)
	return m.styles.Label.Foreground(lipgloss.Color("241")).Render(
		fmt.Sprintf("  Showing %d-%d of %d", current, end, totalRules))
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSuffix(s, " (v6)")
	if len(s) > maxLen {
		return s[:maxLen-1] + "…"
	}
	return s
}

func (m Model) DetailView() string {
	if m.detailRule == nil {
		return ""
	}

	r := m.detailRule

	var lines []string
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Rule #:"), m.styles.Value.Render(fmt.Sprintf("%d", r.Num))))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Action:"), ui.GetPolicyStyle(m.styles, r.Action).Render(r.Action)))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("Protocol:"), m.styles.Value.Render(r.ToProtocol)))

	ipVersion := "IPv4"
	if r.IPv6 {
		ipVersion = "IPv6"
	}
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Render("IP Version:"), m.styles.Value.Render(ipVersion)))

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
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Width(10).Render("Rule   :"), m.styles.Value.Render(fmt.Sprintf("%d", r.Num))))
	lines = append(lines, fmt.Sprintf("  %s  %s %s", m.styles.Label.Width(10).Render("Action :"), ui.GetPolicyStyle(m.styles, r.Action).Render(r.Action), m.styles.Value.Render(r.ToPort)))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Width(10).Render("From   :"), m.styles.Value.Render(r.FromSource)))
	lines = append(lines, fmt.Sprintf("  %s  %s", m.styles.Label.Width(10).Render("To     :"), m.styles.Value.Render(r.ToDest)))
	lines = append(lines, "")
	text := m.styles.Label.
		PaddingTop(1).
		UnsetWidth().
		Render("[Enter] Confirm  [Esc] Cancel")

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

func (m Model) AddWizardView() string {
	w := m.addWizard
	if w == nil {
		return ""
	}

	var lines []string
	lines = append(lines, "")

	if w.Params.Action != "" {
		lines = append(lines, fmt.Sprintf("  %s %s", m.styles.Label.Render("Action:"), ui.GetPolicyStyle(m.styles, w.Params.Action).Render(w.Params.Action)))
	}
	if w.Step > StepDirection && w.Params.Direction != "" {
		lines = append(lines, fmt.Sprintf("  %s %s", m.styles.Label.Render("Direction:"), m.styles.Value.Render(w.Params.Direction)))
	} else if w.Step > StepDirection {
		lines = append(lines, fmt.Sprintf("  %s %s", m.styles.Label.Render("Direction:"), m.styles.Value.Render("both")))
	}
	if w.Step > StepProtocol {
		lines = append(lines, fmt.Sprintf("  %s %s", m.styles.Label.Render("Protocol:"), m.styles.Value.Render(w.Params.Protocol)))
	}
	if w.Step > StepPort {
		port := w.Params.Port
		if port == "" {
			port = "any"
		}
		lines = append(lines, fmt.Sprintf("  %s %s", m.styles.Label.Render("Port:"), m.styles.Value.Render(port)))
	}
	if w.Step > StepSource {
		src := w.Params.FromAddr
		if src == "" {
			src = "any"
		}
		lines = append(lines, fmt.Sprintf("  %s %s", m.styles.Label.Render("Source:"), m.styles.Value.Render(src)))
	}
	if w.Step > StepDestination {
		dst := w.Params.ToAddr
		if dst == "" {
			dst = "any"
		}
		lines = append(lines, fmt.Sprintf("  %s %s", m.styles.Label.Render("Destination:"), m.styles.Value.Render(dst)))
	}
	if w.Step > StepInterface {
		iface := w.Params.Interface
		if iface == "" {
			iface = "all"
		}
		lines = append(lines, fmt.Sprintf("  %s %s", m.styles.Label.Render("Interface:"), m.styles.Value.Render(iface)))
	}

	if len(lines) > 1 {
		lines = append(lines, "")
	}

	stepTitle := m.getWizardStepTitle(w.Step)
	lines = append(lines, m.styles.Title.Render(fmt.Sprintf("  %s", stepTitle)))
	lines = append(lines, "")

	if w.Error != "" {
		lines = append(lines, m.styles.Error.Render(fmt.Sprintf("  Error: %s", w.Error)))
		lines = append(lines, "")
	}

	if w.InputMode {
		inputBox := fmt.Sprintf("  > %s_", w.Input)
		lines = append(lines, m.styles.Value.Render(inputBox))
		lines = append(lines, "")
		lines = append(lines, m.styles.Label.Foreground(lipgloss.Color("241")).Render("  [Enter] Confirm  [Esc] Cancel"))
	} else {
		for i, opt := range w.Options {
			prefix := "  "
			style := m.styles.Label
			if i == w.Cursor {
				prefix = "> "
				style = m.styles.Value.Bold(true)
			}
			lines = append(lines, style.Render(fmt.Sprintf("  %s%s", prefix, opt)))
		}
		lines = append(lines, "")

		hints := "[Enter] Select  [Esc] Cancel"
		if w.Step > StepAction {
			hints = "[b] Back  [Esc] Cancel"
		}
		if w.Step == StepPort || w.Step == StepSource || w.Step == StepDestination {
			hints = "[c] Custom  [b] Back  [Esc] Cancel"
		}
		lines = append(lines, m.styles.Label.UnsetWidth().Foreground(lipgloss.Color("241")).Render(hints))
	}

	lines = append(lines, "")

	// Show preview command at confirm step
	// if w.Step == StepConfirm {
	// 	args := ufw.BuildAddRuleCommand(w.Params)
	// 	cmdPreview := "sudo ufw " + strings.Join(args, " ")
	// 	lines = append(lines, m.styles.Label.Render("  Command:"))
	// 	lines = append(lines, m.styles.Value.Render(fmt.Sprintf("  %s", cmdPreview)))
	// 	lines = append(lines, "")
	// }

	content := strings.Join(lines, "\n")
	return ui.TitledBox("Add New Rule", content, m.styles, 40, true)
}

func (m Model) getWizardStepTitle(step WizardStep) string {
	switch step {
	case StepAction:
		return "Select Action"
	case StepDirection:
		return "Select Direction"
	case StepProtocol:
		return "Select Protocol"
	case StepPort:
		return "Select Port"
	case StepSource:
		return "Select Source Address"
	case StepDestination:
		return "Select Destination Address"
	case StepInterface:
		return "Select Interface"
	case StepConfirm:
		return "Confirm Rule"
	default:
		return "Add Rule"
	}
}
