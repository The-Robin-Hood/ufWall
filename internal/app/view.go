package app

import (
	"fmt"

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

	if m.status.Error != nil {
		footer := m.renderFooter(containerWidth)
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			m.styles.Error.Render(fmt.Sprintf("%v", m.status.Error)),
			"",
			footer,
		)
	}

	title := m.styles.Title.
		Width(containerWidth).
		Render("Firewall Manager")

	infoSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderStatusSection(),
		m.renderPoliciesSection(),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		infoSection,
		m.renderRulesSection(),
	)

	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
		"",
		m.renderFooter(containerWidth),
	)

	layout = lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		layout,
	)

	return layout
}
