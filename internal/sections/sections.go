package sections

const (
	StatsSection = iota
	PolicySection
	RulesSection
)

// All section models implement the following common methods:
//   - Focus()         - activates the section and resets cursor position
//   - Blur()          - deactivates the section and clears overlays/menus
//   - HasOpenMenu()   - returns true if a menu overlay is displayed
//   - GetMenu()       - returns the current menu or nil
//
// Note: Update() and View() have different signatures per section
// because they depend on different data (stats, policy, or rules).
