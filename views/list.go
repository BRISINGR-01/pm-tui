package views

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"

	. "pm-tui/utils"
)

func NewListView(items []list.Item, w, h int) list.Model {
	delegate := list.NewDefaultDelegate()
	if len(items) > 0 && len(items[0].(ListItem).Desc) == 0 {
		delegate.SetHeight(1) // one line per item
		delegate.ShowDescription = false
	}
	delegate.SetSpacing(0)

	delegate.Styles.FilterMatch = lipgloss.NewStyle().
		Foreground(Primary).
		Underline(true).
		Bold(true)

	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(Primary).
		Foreground(Primary).
		Padding(0, 0, 0, 1)

	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(Primary).
		Foreground(Secondary).
		Padding(0, 0, 0, 1)

	l := list.New(items, delegate, w, h)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetHeight(l.Height() - 1)

	return l
}
