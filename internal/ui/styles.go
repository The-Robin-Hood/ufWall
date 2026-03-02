package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title    lipgloss.Style
	Subtitle lipgloss.Style

	ActiveStatus   lipgloss.Style
	InactiveStatus lipgloss.Style
	StatusLabel    lipgloss.Style

	SectionBorder       lipgloss.Style
	SectionTitle        lipgloss.Style
	SectionBorderActive lipgloss.Style

	ActiveCursorPointer lipgloss.Style

	Menu lipgloss.Style

	AllowPolicy  lipgloss.Style
	DenyPolicy   lipgloss.Style
	RejectPolicy lipgloss.Style

	Label lipgloss.Style
	Value lipgloss.Style

	Footer    lipgloss.Style
	FooterKey lipgloss.Style

	Error lipgloss.Style
}

func NewStyles() Styles {
	var (
		mauve    = lipgloss.Color("#cba6f7")
		red      = lipgloss.Color("#f38ba8")
		maroon   = lipgloss.Color("#eba0ac")
		peach    = lipgloss.Color("#fab387")
		yellow   = lipgloss.Color("#f9e2af")
		green    = lipgloss.Color("#a6e3a1")
		sky      = lipgloss.Color("#89dceb")
		lavender = lipgloss.Color("#b4befe")

		text     = lipgloss.Color("#cdd6f4")
		subtext0 = lipgloss.Color("#a6adc8")
		overlay2 = lipgloss.Color("#9399b2")
		overlay1 = lipgloss.Color("#7f849c")
		surface2 = lipgloss.Color("#585b70")
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

		ActiveCursorPointer: lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true),

		StatusLabel: lipgloss.NewStyle().
			Foreground(text).
			Bold(true),

		SectionBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface2),

		SectionBorderActive: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mauve),

		Menu: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(yellow)).
			Padding(1, 2),

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
			Foreground(subtext0),

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

func GetPolicyStyle(styles Styles, policy string) lipgloss.Style {
	switch strings.ToUpper(policy) {
	case "ALLOW":
		return styles.AllowPolicy
	case "DENY":
		return styles.DenyPolicy
	case "REJECT":
		return styles.RejectPolicy
	default:
		return styles.Value
	}
}
