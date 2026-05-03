package main

import (
	"fmt"
	"os"
	"strings"

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
	pendingTick   bool

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
	m.pendingTick = false
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
	providers := listAvailableProviders()
	if len(providers) == 0 {
		fmt.Println("No supported proviers were found. Supported providers are:\n", strings.Join(supportedProviders(), "\n"))
		os.Exit(1)
	}

	provider := GetProviderType(providers[0])

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
