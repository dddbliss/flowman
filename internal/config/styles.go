package config

import "github.com/charmbracelet/lipgloss"

var (
	ScreenWidth  int
	ScreenHeight int
	Divisor      int = 3
)

type styles struct {
	ColumnStyle       lipgloss.Style
	FocusedStyle      lipgloss.Style
	TitleStyle        lipgloss.Style
	DialogBoxStyle    lipgloss.Style
	FormFieldStyle    lipgloss.Style
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
	ItemTitleStyle    lipgloss.Style
	ItemDescStyle     lipgloss.Style
	ActiveColor       lipgloss.Color
	InActiveColor     lipgloss.Color
}

func (styles styles) GetStatefulColor(condition bool) lipgloss.Color {
	if condition {
		return styles.ActiveColor
	} else {
		return styles.InActiveColor
	}
}

var Styles styles = styles{
	ColumnStyle:  lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.HiddenBorder()),
	FocusedStyle: lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")),
	TitleStyle: lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2),

	DialogBoxStyle: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1).
		Margin(1),

	FormFieldStyle: lipgloss.NewStyle().Padding(0).Margin(1, 1, 0, 1),

	ItemStyle:         lipgloss.NewStyle().Padding(0).Margin(0).Border(lipgloss.HiddenBorder(), false, false, false, true),
	SelectedItemStyle: lipgloss.NewStyle().Padding(0).Margin(0).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("62")),
	ItemTitleStyle:    lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("62")).Bold(true),
	ItemDescStyle:     lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("244")),
	ActiveColor:       lipgloss.Color("62"),
	InActiveColor:     lipgloss.Color("249"),
}
