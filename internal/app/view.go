package app

import (
	"fmt"
	"ufWall/internal/sections"
	"ufWall/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	const minWidth, minHeight = 87, 30
	const containerWidth, containerHeight = minWidth - 5, minHeight - 3
	if m.width < minWidth || m.height < minHeight {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.styles.Error.Render(fmt.Sprintf("Terminal too small!\nMinimum: %dx%d\nCurrent: %dx%d", minWidth, minHeight, m.width, m.height)),
		)
	}

	if m.err != nil {
		footer := ui.Footer(m.styles, m.activeSection, containerWidth)
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			m.styles.Error.Render(fmt.Sprintf("%v", m.err)),
			"",
			footer,
		)
	}

	title := m.styles.Title.
		Width(containerWidth).
		Render("Firewall Manager")

	infoSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.statsSection.View(m.stats),
		m.policySection.View(m.policy),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		infoSection,
		m.rulesSection.View(m.rules),
	)

	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
		"",
		ui.Footer(m.styles, m.activeSection, containerWidth),
	)

	layout = lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		layout,
	)

	if m.activeSection == sections.StatsSection && m.statsSection.GetMenu() != nil {
		dimmed := "\x1b[2m" + lipgloss.NewStyle().Faint(true).Render(layout) + "\x1b[0m"
		menuView := m.statsSection.GetMenu().View(m.styles)

		menuW := lipgloss.Width(menuView)
		menuH := lipgloss.Height(menuView)
		x := (m.width - menuW) / 2
		y := (m.height - menuH) / 2

		return PlaceOverlay(x, y, menuView, dimmed)
	}

	return layout
}
