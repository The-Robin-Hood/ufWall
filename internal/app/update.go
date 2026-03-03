package app

import (
	"github.com/The-Robin-Hood/ufWall/internal/keys"
	"github.com/The-Robin-Hood/ufWall/internal/sections"
	"github.com/The-Robin-Hood/ufWall/internal/sections/rules"
	"github.com/The-Robin-Hood/ufWall/internal/ufw"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		if !m.isMenuOpen() {
			switch {
			case key.Matches(msg, keys.Bindings.NextSection):
				m.blurAllSections()
				m.activeSection = (m.activeSection + 1) % 3
				m.focusActiveSection()
				return m, nil

			case key.Matches(msg, keys.Bindings.PrevSection):
				m.blurAllSections()
				m.activeSection = (m.activeSection - 1 + 3) % 3
				m.focusActiveSection()
				return m, nil

			case key.Matches(msg, keys.Bindings.Refresh):
				return m, keys.Refresh()

			case key.Matches(msg, keys.Bindings.Quit):
				return m, tea.Quit

			case key.Matches(msg, keys.Bindings.Execute):
				if !m.stats.Active {
					ufw.Enable()
					return m, keys.Refresh()
				}
			}
		}

		switch m.activeSection {
		case sections.StatsSection:
			newStats, sectionCmd := m.statsSection.Update(msg)
			m.statsSection = newStats
			return m, sectionCmd

		case sections.PolicySection:
			newPolicy, sectionCmd := m.policySection.Update(msg, m.policy)
			m.policySection = newPolicy
			return m, sectionCmd

		case sections.RulesSection:
			newRules, sectionCmd := m.rulesSection.Update(msg, rules.RulesData{IPv4: m.ipv4Rules, IPv6: m.ipv6Rules})
			m.rulesSection = newRules
			return m, sectionCmd
		}

	case keys.RefreshMsg:
		data := ufw.GetUFWData()
		m.rules = data.Rules
		m.ipv4Rules = data.IPv4Rules
		m.ipv6Rules = data.IPv6Rules
		m.policy = data.Policy
		m.stats = data.Stats
		m.err = data.Error
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *model) blurAllSections() {
	m.statsSection.Blur()
	m.policySection.Blur()
	m.rulesSection.Blur()
}

func (m *model) focusActiveSection() {
	switch m.activeSection {
	case sections.StatsSection:
		m.statsSection.Focus()
	case sections.PolicySection:
		m.policySection.Focus()
	case sections.RulesSection:
		m.rulesSection.Focus()
	}
}

func (m *model) isMenuOpen() bool {
	return m.statsSection.HasOpenMenu() || m.rulesSection.HasOpenMenu()
}
