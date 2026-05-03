package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"pm-tui/package_manager"
	. "pm-tui/utils"
	. "pm-tui/views"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	Log("state: %d, msg: %T %+v", m.state, msg, msg)

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:

		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)

		updated, cmd := handleKeyPress(m, msg)
		m = updated.(model)
		cmds = append(cmds, cmd)

		if m.state != StateShowInput {
			m.listView, cmd = m.listView.Update(msg)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.listView.SetSize(msg.Width, msg.Height)
		m.input.SetWidth(min(50, m.width))
		m.viewport.SetWidth(m.width)
		m.viewport.SetHeight(m.height - HeaderHeight)

	case PkgInfoMsg:
		if msg.Err != nil {
			m.isShowingInfo = false

			return m, showModal(msg.Err.Error())
		}
		m.viewport = NewViewport(msg.Content, m.width, m.height)
		return m, nil
	case ErrMsg:
		return m, showModal(msg.Err.Error())

	case OptionsListMsg:
		m.listView = NewListView(msg.Values, m.width, m.height)
		return m, nil
	case ModalMsg:
		m.modalContent = msg.Text
		return m, nil
	case ClearModalMsg:
		m.modalContent = ""
		return m, nil
	case SearchThrottleMsg:
		m.pendingTick = false
		items, err := m.pm.SearchForPackage(strings.TrimSpace(m.input.Value()))
		if err != nil {
			return m, showModal(err.Error())
		}

		return m, LoadSearchResults(items)
	}

	m.listView, cmd = m.listView.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func handleItemChoice(m model, choice string) (tea.Model, tea.Cmd) {
	choice = strings.TrimSpace(choice)
	if len(choice) == 0 {
		return m, nil
	}

	switch m.state {
	case StateListMenuActions:
		return handleMenuChoice(m, choice)

	case StateListProviders:
		m.provider = package_manager.GetProviderType(choice)
		m.pm = package_manager.NewPackageManager(m.provider)
		return LoadActions(m)

	case StateListInstalled:
		m.selectedPkg = choice
		return LoadInfo(m)

	case StateListSearchResults:
		m.selectedPkg = choice
		return LoadInfo(m)

	case StateListRecent:
		m.selectedPkg = choice
		return LoadInfo(m)
	}

	return m, nil
}

func handleMenuChoice(m model, choice string) (tea.Model, tea.Cmd) {
	switch choice {
	case "Explore installed":
		items, err := m.pm.ListInstalled()
		if err != nil {
			return m, showModal(err.Error())
		}

		m.state = StateListInstalled
		return m, LoadList(items)

	case "Install new":
		m.state = StateShowInput
		m.listView = NewListView([]list.Item{}, m.width, m.height)
		return m, nil

	case "Update system":
		return m, run(m.pm.UpdateSystem(), "System updated successfully")

	case "Recent packages":
		return LoadRecent(m)

	case "Select provider":
		return LoadProviders(m)

	case "Quit":
		return m, tea.Quit
	}

	return m, nil
}

func handleKeyPress(m model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.state == StateShowInput {
		switch msg.String() {
		case "enter":
			m.state = StateListSearchResults
			return m, nil
		case "esc":
			return LoadActions(m)
		}

		if len(m.input.Value()) == 0 {
			return m, LoadSearchResults([]package_manager.SearchResult{})
		}

		if !m.pendingTick {
			m.pendingTick = true
			return m, throttle()
		}

		return m, nil
	}

	switch {
	case key.Matches(msg, m.listView.KeyMap.CursorUp):
		if m.listView.Index() == 0 {
			m.listView.Select(len(m.listView.Items())) // wrap to bottom
		}
		return m, nil

	case key.Matches(msg, m.listView.KeyMap.CursorDown):
		if m.listView.Index() == len(m.listView.Items())-1 {
			m.listView.Select(-1) // wrap to top
		}
		return m, nil
	}

	if m.isShowingInfo {
		switch msg.String() {
		case "u":
			return m, run(m.pm.UpdatePackage(m.selectedPkg), fmt.Sprintf("Successfully updated %s", m.selectedPkg))
		case "d":
			return m, run(m.pm.Remove(m.selectedPkg), fmt.Sprintf("Successfully removed %s", m.selectedPkg))
		case "i":
			return m, run(m.pm.Install(m.selectedPkg), fmt.Sprintf("Successfully installed %s", m.selectedPkg))
		case "esc":
			if m.isShowingInfo {
				m.isShowingInfo = false
				return m, nil
			}

			return LoadActions(m)
		}
	}

	switch msg.String() {
	case "enter":
		if m.listView.FilterState() == list.Filtering {
			m.listView.SetFilterState(list.FilterApplied)
			return m, nil
		}

		if m.listView.SelectedItem() == nil {
			return m, nil
		}
		return handleItemChoice(m, m.listView.SelectedItem().(ListItem).Name)

	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		if m.state == StateListMenuActions {
			return m, tea.Quit
		}

		return LoadActions(m)
	case "q":
		if m.listView.FilterState() != list.Filtering {
			return m, tea.Quit
		}
	}

	return m, nil
}

func showModal(msg string) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg { return ModalMsg{Text: msg} },
		clearUI(),
	)
}

func run(cmd *exec.Cmd, msg string) tea.Cmd {
	var stderr bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	return tea.Sequence(
		tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return ErrMsg{Err: FormatErr(stderr.Bytes(), err)}
			}

			return ModalMsg{Text: msg}
		}),
		clearUI(),
	)
}

func clearUI() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return ClearModalMsg{}
	})
}

func throttle() tea.Cmd {
	return tea.Tick(SearchThrottle, func(t time.Time) tea.Msg {
		return SearchThrottleMsg{}
	})
}
