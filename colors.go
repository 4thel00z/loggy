package loggy

import (
	"github.com/charmbracelet/lipgloss"
)

type ColorMap struct {
	Red600    lipgloss.Color
	Yellow500 lipgloss.Color
	Purple500 lipgloss.Color
	Purple400 lipgloss.Color
	Purple300 lipgloss.Color
	Purple200 lipgloss.Color
	Purple100 lipgloss.Color
}

var Colors = ColorMap{
	Red600:    lipgloss.Color("#dc2626"),
	Yellow500: lipgloss.Color("#eab308"),
	Purple500: lipgloss.Color("#8b5cf6"),
	Purple400: lipgloss.Color("#a78bfa"),
	Purple300: lipgloss.Color("#d2a8ff"),
	Purple200: lipgloss.Color("#d6bcfa"),
	Purple100: lipgloss.Color("#e9d8fd"),
}

type AdaptiveColorMap struct {
	DefaultText lipgloss.AdaptiveColor
	Subtle      lipgloss.AdaptiveColor
	Highlight   lipgloss.AdaptiveColor
	Special     lipgloss.AdaptiveColor
}

var AdaptiveColors = AdaptiveColorMap{
	DefaultText: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"},
	Subtle:      lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"},
	Highlight:   lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"},
	Special:     lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"},
}

var ListStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, true, false, false).
	BorderForeground(AdaptiveColors.Subtle).
	MarginRight(2)

var Divider = lipgloss.NewStyle().
	SetString("•").
	Padding(0, 1).
	Foreground(AdaptiveColors.Subtle).
	String()
