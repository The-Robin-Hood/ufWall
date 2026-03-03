package policy

import (
	"github.com/The-Robin-Hood/ufWall/internal/keys"
	"github.com/The-Robin-Hood/ufWall/internal/ufw"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg, policy ufw.Policy) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Bindings.CursorUp):
			if m.cursorLine > 0 {
				m.cursorLine--
			}
		case key.Matches(msg, keys.Bindings.CursorDown):
			if m.cursorLine < m.totalOpts {
				m.cursorLine++
			}
		case key.Matches(msg, keys.Bindings.Execute):
			switch m.cursorLine {
			case 0:
				ufw.DefaultIncoming(!(policy.DefaultIncoming == "ALLOW"))
			case 1:
				ufw.DefaultOutgoing(!(policy.DefaultOutgoing == "ALLOW"))
			case 2:
				ufw.DefaultRouted(!(policy.DefaultRouted == "ALLOW"))
			}
			return m, keys.Refresh()
		}
	}
	return m, nil
}
