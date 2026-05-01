package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	. "pm-tui/package_manager"
	. "pm-tui/utils"
	. "pm-tui/views"
)

func (m model) View() tea.View {
	var v tea.View
	v.AltScreen = true                    // use the full size of the terminal in its "alternate screen buffer"
	v.MouseMode = tea.MouseModeCellMotion // turn on mouse support so we can track the mouse wheel

	var b strings.Builder
	if m.modalContent != "" {
		v.SetContent(modal(m.modalContent, m.width, m.height))
		return v
	}

	if m.isShowingInfo {
		isInstalled, err := m.pm.IsInstalled(m.selectedPkg)
		if err != nil {
			v.SetContent(modal(err.Error(), m.width, m.height))
			return v
		} else {
			b.WriteString(InfoView(m.selectedPkg, isInstalled, &m.viewport))
		}
		b.WriteString(footer(m))
	} else {
		switch m.state {
		case StateShowInput:
			b.WriteString(SearchView(m.input))

		case StateListSearchResults, StateListInstalled, StateListProviders, StateListMenuActions, StateListRecent:
			b.WriteString(lipgloss.JoinVertical(lipgloss.Left, m.listView.View(), footer(m)))
		}
	}

	v.SetContent(b.String())
	return v
}

func modal(str string, width, height int) string {
	box := lipgloss.NewStyle().
		Padding(1, 4).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Render(str)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

func footer(m model) string {

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(Secondary).
		Bold(true).
		Padding(0, 1)

	sep := " • "

	manager := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"manager: ",
		lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true).
			Render(GetProviderStr(m.provider)),
	)

	controls := lipgloss.JoinHorizontal(
		lipgloss.Center,
		sep,
		keyStyle.Render("↑/↓"), " move",
		sep,
		keyStyle.Render("enter"), " select",
		sep,
		keyStyle.Render("q"), " quit",
		sep,
		keyStyle.Render("esc"), " back",
		sep,
		keyStyle.Render("/"), " filter",
	)

	footer := lipgloss.JoinHorizontal(
		lipgloss.Left,
		manager,
		controls,
	)

	footer = lipgloss.NewStyle().
		Width(m.listView.Width()).
		Padding(0, 1).
		Render(footer)

	return footer
}
