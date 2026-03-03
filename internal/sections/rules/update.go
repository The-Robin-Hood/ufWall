package rules

import (
	"log"
	"github.com/The-Robin-Hood/ufWall/internal/keys"
	"github.com/The-Robin-Hood/ufWall/internal/ufw"
	"github.com/The-Robin-Hood/ufWall/internal/ui"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	ActionToggle   = "toggle"
	ActionMoveUp   = "move_up"
	ActionMoveDown = "move_down"
	ActionAdd      = "add"
)

func (m Model) Update(msg tea.Msg, data RulesData) (Model, tea.Cmd) {
	rules := m.getActiveRules(data)

	if m.addWizard != nil {
		return m.updateAddWizard(msg)
	}

	if m.showDeleteConfirm {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Bindings.Execute):
				if m.deleteRule != nil {
					log.Printf("Deleting rule #%d (confirmed)", m.deleteRule.Num)
					ufw.DeleteRule(m.deleteRule.Num)
				}
				m.showDeleteConfirm = false
				m.deleteRule = nil
				return m, keys.Refresh()
			case key.Matches(msg, keys.Bindings.Quit):
				m.showDeleteConfirm = false
				m.deleteRule = nil
				return m, nil
			}
		}
		return m, nil
	}

	if m.showDetails {
		switch msg.(type) {
		case tea.KeyMsg:
			m.showDetails = false
			m.detailRule = nil
			return m, nil
		}
		return m, nil
	}

	if m.menu != nil {
		if quit := m.menu.Update(msg); quit {
			if m.menuContext != nil && m.menuContext.PendingSubmenu {
				return m.handlePendingSubmenu()
			}
			m.menu = nil
			m.menuContext = nil
			return m, keys.Refresh()
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Bindings.SwitchTable):
			if m.activeTable == IPv4Table {
				m.activeTable = IPv6Table
			} else {
				m.activeTable = IPv4Table
			}
			return m, nil

		case key.Matches(msg, keys.Bindings.CursorUp):
			m.moveCursorUp()
			return m, nil

		case key.Matches(msg, keys.Bindings.CursorDown):
			m.moveCursorDown(rules)
			return m, nil

		case key.Matches(msg, keys.Bindings.Execute):
			if len(rules) > 0 {
				return m.openMainMenu(data)
			}
			m.addWizard = NewAddWizard()
			return m, nil

		case key.Matches(msg, keys.Bindings.AddRule):
			m.addWizard = NewAddWizard()
			return m, nil

		case key.Matches(msg, keys.Bindings.Info):
			cursorLine := m.getCurrentCursor()
			if len(rules) > 0 && cursorLine < len(rules) {
				ruleCopy := rules[cursorLine]
				m.showDetails = true
				m.detailRule = &ruleCopy
				return m, nil
			}

		case key.Matches(msg, keys.Bindings.Delete):
			cursorLine := m.getCurrentCursor()
			if len(rules) > 0 && cursorLine < len(rules) {
				ruleCopy := rules[cursorLine]
				m.showDeleteConfirm = true
				m.deleteRule = &ruleCopy
				return m, nil
			}
		}
	}
	return m, nil
}

func (m Model) updateAddWizard(msg tea.Msg) (Model, tea.Cmd) {
	w := m.addWizard

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Bindings.Quit):
			if w.InputMode {
				w.InputMode = false
				w.Input = ""
				return m, nil
			}
			m.addWizard = nil
			return m, nil
		}

		if w.InputMode {
			switch {
			case key.Matches(msg, keys.Bindings.Execute):
				return m.confirmWizardInput()
			case key.Matches(msg, keys.Bindings.Back):
				if len(w.Input) > 0 {
					w.Input = w.Input[:len(w.Input)-1]
				}
			default:
				if len(msg.String()) == 1 {
					w.Input += msg.String()
				}
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, keys.Bindings.CursorUp):
			if w.Cursor > 0 {
				w.Cursor--
			}
		case key.Matches(msg, keys.Bindings.CursorDown):
			if w.Cursor < len(w.Options)-1 {
				w.Cursor++
			}
		case key.Matches(msg, keys.Bindings.Execute):
			return m.selectWizardOption()
		case key.Matches(msg, keys.Bindings.CustomInput):
			if w.Step == StepPort || w.Step == StepSource || w.Step == StepDestination {
				w.InputMode = true
				w.Input = ""
			}
		case key.Matches(msg, keys.Bindings.Back):
			return m.wizardPrevStep()
		}
	}

	return m, nil
}

func (m Model) selectWizardOption() (Model, tea.Cmd) {
	w := m.addWizard
	if w.Cursor >= len(w.Options) {
		return m, nil
	}

	selected := w.Options[w.Cursor]

	switch w.Step {
	case StepAction:
		w.Params.Action = selected
		w.Step = StepDirection
		w.Options = append([]string{"Both (in & out)"}, ufw.Directions...)
		w.Cursor = 0

	case StepDirection:
		if selected == "Both (in & out)" {
			w.Params.Direction = ""
		} else {
			w.Params.Direction = selected
		}
		w.Step = StepProtocol
		w.Options = ufw.Protocols
		w.Cursor = 0

	case StepProtocol:
		w.Params.Protocol = selected
		w.Step = StepPort
		w.Options = []string{"Any (no port filter)"}
		for _, p := range ufw.CommonPorts {
			w.Options = append(w.Options, p.Port+" ("+p.Name+")")
		}
		w.Cursor = 0

	case StepPort:
		if selected == "Any (no port filter)" {
			w.Params.Port = ""
		} else {
			for i, c := range selected {
				if c == ' ' {
					w.Params.Port = selected[:i]
					break
				}
			}
		}
		w.Step = StepSource
		w.Options = []string{"Any", "Custom..."}
		w.Cursor = 0

	case StepSource:
		switch selected {
		case "Any":
			w.Params.FromAddr = ""
		case "Custom...":
			w.InputMode = true
			w.Input = ""
			return m, nil
		}
		w.Step = StepDestination
		w.Options = []string{"Any", "Custom..."}
		w.Cursor = 0

	case StepDestination:
		switch selected {
		case "Any":
			w.Params.ToAddr = ""
		case "Custom...":
			w.InputMode = true
			w.Input = ""
			return m, nil
		}
		w.Step = StepInterface
		w.Options = append([]string{"All interfaces"}, w.Interfaces...)
		w.Cursor = 0

	case StepInterface:
		if selected == "All interfaces" {
			w.Params.Interface = ""
		} else {
			w.Params.Interface = selected
		}
		w.Step = StepConfirm
		w.Options = []string{"Confirm and Add Rule", "Cancel"}
		w.Cursor = 0

	case StepConfirm:
		if selected == "Confirm and Add Rule" {
			// Execute the command
			log.Printf("Adding rule: %+v", w.Params)
			_, stderr, err := ufw.AddNewRule(w.Params)
			if err != nil {
				log.Printf("Error adding rule: %s", stderr)
				w.Error = stderr
				return m, nil
			}
			m.addWizard = nil
			return m, keys.Refresh()
		} else {
			// Cancel
			m.addWizard = nil
			return m, nil
		}
	}

	return m, nil
}

func (m Model) confirmWizardInput() (Model, tea.Cmd) {
	w := m.addWizard
	input := w.Input
	w.InputMode = false
	w.Input = ""

	switch w.Step {
	case StepPort:
		w.Params.Port = input
		w.Step = StepSource
		w.Options = []string{"Any", "Custom..."}
		w.Cursor = 0

	case StepSource:
		w.Params.FromAddr = input
		w.Step = StepDestination
		w.Options = []string{"Any", "Custom..."}
		w.Cursor = 0

	case StepDestination:
		w.Params.ToAddr = input
		w.Step = StepInterface
		w.Options = append([]string{"All interfaces"}, w.Interfaces...)
		w.Cursor = 0
	}

	return m, nil
}

func (m Model) wizardPrevStep() (Model, tea.Cmd) {
	w := m.addWizard

	switch w.Step {
	case StepAction:
		m.addWizard = nil
		return m, nil

	case StepDirection:
		w.Step = StepAction
		w.Options = ufw.Actions
		w.Cursor = 0

	case StepProtocol:
		w.Step = StepDirection
		w.Options = append([]string{"Both (in & out)"}, ufw.Directions...)
		w.Cursor = 0

	case StepPort:
		w.Step = StepProtocol
		w.Options = ufw.Protocols
		w.Cursor = 0

	case StepSource:
		w.Step = StepPort
		w.Options = []string{"Any (no port filter)"}
		for _, p := range ufw.CommonPorts {
			w.Options = append(w.Options, p.Port+" ("+p.Name+")")
		}
		w.Cursor = 0

	case StepDestination:
		w.Step = StepSource
		w.Options = []string{"Any", "Custom..."}
		w.Cursor = 0

	case StepInterface:
		w.Step = StepDestination
		w.Options = []string{"Any", "Custom..."}
		w.Cursor = 0

	case StepConfirm:
		w.Step = StepInterface
		w.Options = append([]string{"All interfaces"}, w.Interfaces...)
		w.Cursor = 0
	}

	return m, nil
}

func (m Model) getActiveRules(data RulesData) []ufw.Rule {
	if m.activeTable == IPv6Table {
		return data.IPv6
	}
	return data.IPv4
}

func (m Model) getCurrentCursor() int {
	if m.activeTable == IPv6Table {
		return m.ipv6CursorLine
	}
	return m.ipv4CursorLine
}

func (m *Model) moveCursorUp() {
	if m.activeTable == IPv6Table {
		if m.ipv6CursorLine > 0 {
			m.ipv6CursorLine--
			// Scroll up if cursor goes above visible area
			if m.ipv6CursorLine < m.ipv6ScrollOffset {
				m.ipv6ScrollOffset = m.ipv6CursorLine
			}
		}
	} else {
		if m.ipv4CursorLine > 0 {
			m.ipv4CursorLine--
			// Scroll up if cursor goes above visible area
			if m.ipv4CursorLine < m.ipv4ScrollOffset {
				m.ipv4ScrollOffset = m.ipv4CursorLine
			}
		}
	}
}

func (m *Model) moveCursorDown(rules []ufw.Rule) {
	if m.activeTable == IPv6Table {
		if m.ipv6CursorLine < len(rules)-1 {
			m.ipv6CursorLine++
			// Scroll down if cursor goes below visible area
			if m.ipv6CursorLine >= m.ipv6ScrollOffset+MaxVisibleRules {
				m.ipv6ScrollOffset = m.ipv6CursorLine - MaxVisibleRules + 1
			}
		}
	} else {
		if m.ipv4CursorLine < len(rules)-1 {
			m.ipv4CursorLine++
			// Scroll down if cursor goes below visible area
			if m.ipv4CursorLine >= m.ipv4ScrollOffset+MaxVisibleRules {
				m.ipv4ScrollOffset = m.ipv4CursorLine - MaxVisibleRules + 1
			}
		}
	}
}

func (m Model) handlePendingSubmenu() (Model, tea.Cmd) {
	if m.menuContext == nil {
		m.menu = nil
		return m, keys.Refresh()
	}

	log.Printf("handlePendingSubmenu: Action=%s, PendingSubmenu=%v", m.menuContext.Action, m.menuContext.PendingSubmenu)
	m.menuContext.PendingSubmenu = false

	switch m.menuContext.Action {
	case ActionToggle:
		return m.openToggleActionMenu()

	case ActionMoveUp:
		rule := m.menuContext.SelectedRule
		if rule != nil {
			log.Printf("Moving rule #%d up (IPv6=%v)", rule.Num, rule.IPv6)
			err := ufw.MoveRule(*rule, -1, rule.Action)
			if err != nil {
				log.Printf("Error moving rule up: %v", err)
			}
		}
		m.menu = nil
		m.menuContext = nil
		return m, keys.Refresh()

	case ActionMoveDown:
		rule := m.menuContext.SelectedRule
		if rule != nil {
			log.Printf("Moving rule #%d down (IPv6=%v)", rule.Num, rule.IPv6)
			err := ufw.MoveRule(*rule, 1, rule.Action)
			if err != nil {
				log.Printf("Error moving rule down: %v", err)
			}
		}
		m.menu = nil
		m.menuContext = nil
		return m, keys.Refresh()

	case ActionAdd:
		m.menu = nil
		m.menuContext = nil
		m.addWizard = NewAddWizard()
		return m, nil
	}

	m.menu = nil
	m.menuContext = nil
	return m, keys.Refresh()
}

func (m Model) openMainMenu(data RulesData) (Model, tea.Cmd) {
	rules := m.getActiveRules(data)
	cursorLine := m.getCurrentCursor()

	if len(rules) == 0 || cursorLine >= len(rules) {
		return m, nil
	}

	selectedRule := rules[cursorLine]
	ruleCopy := selectedRule // Make a copy to store in context

	menuLabels := []string{
		"Toggle Action",
		"Move Up",
		"Move Down",
		"Add New Rule",
	}

	m.menuContext = &MenuContext{
		SelectedRule: &ruleCopy,
		TotalRules:   len(rules),
	}

	options := ui.MakeMenuItems(menuLabels, func(label string) tea.Cmd {
		switch label {
		case "Toggle Action":
			m.menuContext.Action = ActionToggle
		case "Move Up":
			m.menuContext.Action = ActionMoveUp
		case "Move Down":
			m.menuContext.Action = ActionMoveDown
		case "Add New Rule":
			m.menuContext.Action = ActionAdd
		}
		m.menuContext.PendingSubmenu = true
		return nil
	})

	menu := ui.NewMenu(options, m.styles)
	m.menu = &menu
	return m, nil
}

func (m Model) openToggleActionMenu() (Model, tea.Cmd) {
	rule := m.menuContext.SelectedRule

	options := ui.MakeMenuItems(ufw.Actions, func(newAction string) tea.Cmd {
		if rule != nil && newAction != rule.Action {
			log.Printf("Toggling rule #%d from %s to %s (IPv6=%v)", rule.Num, rule.Action, newAction, rule.IPv6)

			err := ufw.MoveRule(*rule, 0, newAction)
			if err != nil {
				log.Printf("Error toggling rule: %v", err)
			}
		}
		m.menuContext = nil
		return nil
	})

	menu := ui.NewMenu(options, m.styles)
	m.menu = &menu
	return m, nil
}
