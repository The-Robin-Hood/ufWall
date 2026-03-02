package rules

import (
	"ufWall/internal/ufw"
	"ufWall/internal/ui"
)

type Table int

const (
	IPv4Table Table = iota
	IPv6Table
)

type WizardStep int

const (
	StepAction WizardStep = iota
	StepDirection
	StepProtocol
	StepPort
	StepSource
	StepDestination
	StepInterface
	StepConfirm
)

type AddWizard struct {
	Step       WizardStep
	Params     ufw.AddRuleParams
	Input      string
	InputMode  bool
	Options    []string
	Cursor     int
	Interfaces []string
	Error      string
}

type MenuContext struct {
	Action         string
	SelectedRule   *ufw.Rule
	TotalRules     int
	PendingSubmenu bool
}

type Model struct {
	styles     ui.Styles
	cursorLine int
	menu       *ui.Menu
	active     bool

	// Two-table navigation
	activeTable    Table // Which table (IPv4 or IPv6) is currently focused
	ipv4CursorLine int   
	ipv6CursorLine int  

	// Detail view overlay
	showDetails bool
	detailRule  *ufw.Rule

	// Delete confirmation
	showDeleteConfirm bool
	deleteRule        *ufw.Rule

	// Multi-step menu operations
	menuContext *MenuContext
	addWizard *AddWizard
}

func New(styles ui.Styles) Model {
	return Model{
		styles:            styles,
		cursorLine:        0,
		menu:              nil,
		active:            false,
		activeTable:       IPv4Table,
		ipv4CursorLine:    0,
		ipv6CursorLine:    0,
		showDetails:       false,
		detailRule:        nil,
		showDeleteConfirm: false,
		deleteRule:        nil,
		menuContext:       nil,
		addWizard:         nil,
	}
}

func (m *Model) Focus() {
	m.active = true
}

func (m *Model) Blur() {
	m.active = false
	m.menu = nil
	m.showDetails = false
	m.showDeleteConfirm = false
	m.menuContext = nil
	m.addWizard = nil
}

func (m Model) HasOpenMenu() bool {
	return m.menu != nil || m.showDetails || m.showDeleteConfirm || m.addWizard != nil
}

func (m Model) GetMenu() *ui.Menu {
	return m.menu
}

func (m Model) ShowingDetails() bool {
	return m.showDetails
}

func (m Model) GetDetailRule() *ufw.Rule {
	return m.detailRule
}

func (m Model) ShowingDeleteConfirm() bool {
	return m.showDeleteConfirm
}

func (m Model) GetDeleteRule() *ufw.Rule {
	return m.deleteRule
}

func (m Model) ActiveTable() Table {
	return m.activeTable
}

func (m Model) IPv4CursorLine() int {
	return m.ipv4CursorLine
}

func (m Model) IPv6CursorLine() int {
	return m.ipv6CursorLine
}

func (m Model) CurrentCursorLine() int {
	if m.activeTable == IPv6Table {
		return m.ipv6CursorLine
	}
	return m.ipv4CursorLine
}

func (m Model) ShowingAddWizard() bool {
	return m.addWizard != nil
}

func (m Model) GetAddWizard() *AddWizard {
	return m.addWizard
}

func NewAddWizard() *AddWizard {
	return &AddWizard{
		Step:       StepAction,
		Params:     ufw.AddRuleParams{},
		Input:      "",
		InputMode:  false,
		Options:    ufw.Actions,
		Cursor:     0,
		Interfaces: ufw.GetInterfaces(),
		Error:      "",
	}
}
