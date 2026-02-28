package app

import (
	"log"
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
		switch m.activeSection {

		case sections.StatsSection:
			newStats, sectionCmd := m.statsSection.Update(msg)
			m.statsSection = newStats
			if m.statsSection.HasOpenMenu() {
				return m, sectionCmd
			}
		}

		switch {
		case key.Matches(msg, Keys.NextSection):
			m.blurAllSections()
			m.activeSection = (m.activeSection + 1) % 3
			m.focusActiveSection()
			return m, nil

		case key.Matches(msg, Keys.PrevSection):
			m.blurAllSections()
			m.activeSection = (m.activeSection - 1 + 3) % 3
			m.focusActiveSection()
			return m, nil

		case key.Matches(msg, Keys.Refresh):
			return m, m.refreshData()

		case key.Matches(msg, Keys.Quit):
			log.Println("EXITING")
			return m, tea.Quit
		}

	case RefreshMsg:
		m.rules = msg.Rules
		m.policy = msg.Policy
		m.stats = msg.Stats
		m.err = msg.Error
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

func (m model) refreshData() tea.Cmd {
	return func() tea.Msg {
		data := ufw.GetUFWData()
		return RefreshMsg{
			Rules:  data.Rules,
			Policy: data.Policy,
			Stats:  data.Stats,
			Error:  data.Error,
		}
	}
}

type RefreshMsg struct {
	Rules  []ufw.Rule
	Policy ufw.Policy
	Stats  ufw.Stats
	Error  error
}
