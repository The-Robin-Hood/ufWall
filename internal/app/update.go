package app

import (
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
		switch {
		case key.Matches(msg, Keys.MoveToNextSection): 
			 m.activeSection = (m.activeSection + 1) % 3 
		case key.Matches(msg, Keys.MoveToPrevSection):
			 m.activeSection = (m.activeSection - 1 + 3) % 3 
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, Keys.Refresh):
			m.status,m.rules = ufw.GetStatus()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}
