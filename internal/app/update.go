package app

import (
	"ufWall/internal/keys"
	"ufWall/internal/sections"
	"ufWall/internal/ufw"

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
			}
		}

		switch m.activeSection {
		case sections.StatsSection:
			newStats, sectionCmd := m.statsSection.Update(msg)
			m.statsSection = newStats
			return m, sectionCmd
		}

	case keys.RefreshMsg:
		data := ufw.GetUFWData()
		m.rules = data.Rules
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
	return m.statsSection.GetMenu() != nil
}
