package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func RenderBoxWithTitle(title, content string, styles Styles, width int, activeSession bool) string {
	
	borderColor := styles.SectionBorder.GetBorderTopForeground()
	if activeSession{
		borderColor = styles.SectionBorderActive.GetBorderTopForeground()
		title  = title + " ●"
	}

	boxWidth := 0
	if width != -1 {
		boxWidth = width
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2)

	if boxWidth > 0 {
		style = style.Width(boxWidth)
	}

	box := style.Render(content)

	lines := strings.Split(box, "\n")
	if len(lines) == 0 {
		return box
	}

	styledTitle := lipgloss.NewStyle().
		Foreground(styles.SectionTitle.GetForeground()).
		Bold(true).
		Render(" " + title + " ")

	topBorderWidth := lipgloss.Width(lines[0])
	titleWidth := lipgloss.Width(styledTitle)

	dashesAfterTitle := topBorderWidth - titleWidth - 3

	if dashesAfterTitle > 0 {
		cornerAndDash := lipgloss.NewStyle().Foreground(borderColor).Render("╭─")
		dashes := lipgloss.NewStyle().Foreground(borderColor).Render(strings.Repeat("─", dashesAfterTitle))
		corner := lipgloss.NewStyle().Foreground(borderColor).Render("╮")
		lines[0] = cornerAndDash + styledTitle + dashes + corner
	}

	return strings.Join(lines, "\n")
}
