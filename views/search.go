package views

import (
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"

	. "pm-tui/utils"
)

func NewSearchView(w int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search packages..."
	ti.SetWidth(min(50, w))

	ti.Focus()

	return ti
}

func SearchView(input textinput.Model, list list.Model) string {
	var b strings.Builder

	inputBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(0, 1).
		Render(input.View())

	b.WriteString(inputBox)
	b.WriteString(list.View())

	return b.String()
}
