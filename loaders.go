package main

import (
	. "pm-tui/package_manager"
	. "pm-tui/utils"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func LoadSearchResults(options []SearchResult) tea.Cmd {
	return func() tea.Msg {
		items := make([]list.Item, len(options))
		for i, res := range options {
			items[i] = ListItem{Name: res.Title, Desc: res.Description}
		}
		return OptionsListMsg{Values: items}
	}
}

func LoadList(options []string) tea.Cmd {
	return func() tea.Msg {
		items := make([]list.Item, len(options))
		for i, l := range options {
			items[i] = ListItem{Name: l}
		}
		return OptionsListMsg{Values: items}
	}
}

func LoadActions(m model) (model, tea.Cmd) {
	m = m.ClearState()
	m.state = StateListMenuActions
	return m, LoadList([]string{"Explore installed", "Install new", "Recent packages", "Update system", "Select provider", "Quit"})
}

func LoadRecent(m model) (model, tea.Cmd) {
	pkgs := GetRecentPkgs()
	if len(pkgs) == 0 {
		return m, showModal("No recent packages")
	}

	m.state = StateListRecent
	return m, LoadList(pkgs)
}

func LoadProviders(m model) (model, tea.Cmd) {
	m = m.ClearState()
	m.state = StateListProviders
	return m, LoadList(filterAvailableProviders([]string{"pacman", "yay", "npm", "npm (global)", "pip", "apt"}))
}

func LoadInfo(m model) (model, tea.Cmd) {
	m.isShowingInfo = true

	return m, func() tea.Msg {
		content, err := m.pm.Info((m.selectedPkg))
		if err != nil {
			return PkgInfoMsg{Content: "", Err: err}
		}

		return PkgInfoMsg{Content: content, Err: nil}
	}
}

func filterAvailableProviders(providers []string) []string {
	filtered := []string{}

	for _, provider := range providers {
		_, err := NewPackageManager(GetProviderType(provider)).IsInstalled("")

		if err == nil {
			filtered = append(filtered, provider)
		}
	}

	return filtered
}
