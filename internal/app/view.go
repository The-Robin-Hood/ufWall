package app

import (
	"fmt"
	"github.com/The-Robin-Hood/ufWall/internal/sections/rules"
	"github.com/The-Robin-Hood/ufWall/internal/ui"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderCenteredOverlay(overlay, background string) string {
	dimmed := "\x1b[2m" + lipgloss.NewStyle().Faint(true).Render(background) + "\x1b[0m"
	overlayW := lipgloss.Width(overlay)
	overlayH := lipgloss.Height(overlay)
	x := (m.width - overlayW) / 2
	y := (m.height - overlayH) / 2
	return PlaceOverlay(x, y, overlay, dimmed)
}

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

	if !m.stats.Active {
		prompt := lipgloss.JoinVertical(
			lipgloss.Center,
			m.styles.Error.Render("Firewall is disabled"),
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press Space to activate"),
			lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Press Esc or Q to quit"),
		)
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			prompt,
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
		m.rulesSection.View(rules.RulesData{IPv4: m.ipv4Rules, IPv6: m.ipv6Rules}),
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

	if m.rulesSection.ShowingDeleteConfirm() {
		return m.renderCenteredOverlay(m.rulesSection.DeleteConfirmView(), layout)
	}

	if m.rulesSection.ShowingDetails() {
		return m.renderCenteredOverlay(m.rulesSection.DetailView(), layout)
	}

	if m.rulesSection.ShowingAddWizard() {
		return m.renderCenteredOverlay(m.rulesSection.AddWizardView(), layout)
	}

	if menu := m.statsSection.GetMenu(); menu != nil {
		return m.renderCenteredOverlay(menu.View(m.styles), layout)
	}

	if menu := m.rulesSection.GetMenu(); menu != nil {
		return m.renderCenteredOverlay(menu.View(m.styles), layout)
	}

	return layout
}
