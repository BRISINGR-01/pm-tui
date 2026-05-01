package utils

import "charm.land/bubbles/v2/list"

type AppState int

const (
	StateListMenuActions AppState = iota // show Explore installed
	//								  Install new
	//								  Update system
	//								  Select provider
	//								  Quit
	StateListInstalled     // (fzf) locally installed
	StateListSearchResults // (fzf) publically available
	StateListProviders     // (fzf) choose pacman/apt/npm/...
	StateListRecent        // last updated/installed/removed pkgs
	StateShowInput         // open search input
)

// When the items for a list are ready
type OptionsListMsg struct {
	Values []list.Item
}

type ChoiceMsg struct {
	Value string
}

type ChooseProviderMsg struct {
	Provider string
}

type PkgInfoMsg struct {
	Content string
	Err     error
}

type ClearModalMsg struct{}
type ModalMsg struct{ Text string }

type ErrMsg struct {
	Err error
}

type ListItem struct {
	Name string
	Desc string
}

func (i ListItem) FilterValue() string { return i.Name }
func (i ListItem) Title() string       { return i.Name }
func (i ListItem) Description() string { return i.Desc }
