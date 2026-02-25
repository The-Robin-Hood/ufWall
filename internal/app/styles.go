package app

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles holds all the styling for the TUI
type Styles struct {
	// Title and header
	Title    lipgloss.Style
	Subtitle lipgloss.Style

	// Status indicators
	ActiveStatus   lipgloss.Style
	InactiveStatus lipgloss.Style
	StatusLabel    lipgloss.Style

	// Sections
	SectionBorder       lipgloss.Style
	SectionTitle        lipgloss.Style
	SectionBorderActive lipgloss.Style

	// Policies
	AllowPolicy  lipgloss.Style
	DenyPolicy   lipgloss.Style
	RejectPolicy lipgloss.Style

	// Labels and values
	Label lipgloss.Style
	Value lipgloss.Style

	// Footer
	Footer    lipgloss.Style
	FooterKey lipgloss.Style

	// Error
	Error lipgloss.Style
}

// NewStyles creates and returns styled components with Catppuccin Mocha theme
func NewStyles() Styles {
	// Catppuccin Mocha Color Palette
	var (
		// rosewater = lipgloss.Color("#f5e0dc")
		// flamingo  = lipgloss.Color("#f2cdcd")
		// pink      = lipgloss.Color("#f5c2e7")
		mauve  = lipgloss.Color("#cba6f7")
		red    = lipgloss.Color("#f38ba8")
		maroon = lipgloss.Color("#eba0ac")
		peach  = lipgloss.Color("#fab387")
		// yellow    = lipgloss.Color("#f9e2af")
		green = lipgloss.Color("#a6e3a1")
		// teal      = lipgloss.Color("#94e2d5")
		sky = lipgloss.Color("#89dceb")
		// sapphire  = lipgloss.Color("#74c7ec")
		// blue      = lipgloss.Color("#89b4fa")
		lavender = lipgloss.Color("#b4befe")

		text = lipgloss.Color("#cdd6f4")
		// subtext1  = lipgloss.Color("#bac2de")
		subtext0 = lipgloss.Color("#a6adc8")
		overlay2 = lipgloss.Color("#9399b2")
		overlay1 = lipgloss.Color("#7f849c")
		// overlay0  = lipgloss.Color("#6c7086")
		surface2  = lipgloss.Color("#585b70")
		// surface1  = lipgloss.Color("#45475a")
		// surface0  = lipgloss.Color("#313244")
		// base      = lipgloss.Color("#1e1e2e")
		// mantle    = lipgloss.Color("#181825")
		// crust     = lipgloss.Color("#11111b")
	)

	return Styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(mauve).
			Align(lipgloss.Center),

		Subtitle: lipgloss.NewStyle().
			Foreground(overlay2).
			Italic(true).
			Align(lipgloss.Center),

		ActiveStatus: lipgloss.NewStyle().
			Bold(true).
			Foreground(green),

		InactiveStatus: lipgloss.NewStyle().
			Bold(true).
			Foreground(red),

		StatusLabel: lipgloss.NewStyle().
			Foreground(text).
			Bold(true),

		SectionBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface2),

		SectionBorderActive: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mauve),

		SectionTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lavender).
			Underline(true),

		AllowPolicy: lipgloss.NewStyle().
			Bold(true).
			Foreground(green),

		DenyPolicy: lipgloss.NewStyle().
			Bold(true).
			Foreground(red),

		RejectPolicy: lipgloss.NewStyle().
			Bold(true).
			Foreground(peach),

		Label: lipgloss.NewStyle().
			Foreground(subtext0).
			Width(15),

		Value: lipgloss.NewStyle().
			Foreground(text).
			Bold(true),

		Footer: lipgloss.NewStyle().
			Foreground(overlay1).
			PaddingBottom(1).
			Align(lipgloss.Center),

		FooterKey: lipgloss.NewStyle().
			Foreground(sky).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(red).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(maroon).
			Padding(1, 2),
	}
}
