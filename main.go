package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	. "pm-tui/package_manager"
	. "pm-tui/utils"
	. "pm-tui/views"
)

type model struct {
	pm       PackageManager
	provider ProviderType

	state         AppState
	selectedPkg   string
	isShowingInfo bool

	listView list.Model
	input    textinput.Model
	viewport viewport.Model

	width, height int
	modalContent  string
}

func (m model) ClearState() model {
	m.viewport = viewport.New()
	m.selectedPkg = ""
	m.modalContent = ""
	m.isShowingInfo = false
	m.listView = NewListView([]list.Item{}, 0, 0)
	m.input = NewSearchView(m.width)

	return m
}

func (m model) Init() tea.Cmd {
	_, cmd := LoadActions(m)
	return cmd
}

func initialModel() model {
	provider := ProviderPacman
	return model{
		pm:            NewPackageManager(provider),
		provider:      provider,
		state:         StateListMenuActions,
		isShowingInfo: false,
		listView:      NewListView([]list.Item{}, 0, 0),
		input:         NewSearchView(0),
	}
}

func main() {
	InitLog()
	defer DestroyLog()

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("Could not run program:", err)
		os.Exit(1)
	}
}
