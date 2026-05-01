package package_manager

import (
	"encoding/json"
	"os/exec"

	. "os/exec"
	. "pm-tui/utils"
)

type Npm struct{ bin string }

type npmListResult struct {
	Dependencies map[string]struct {
		Version    string `json:"version"`
		Overridden bool   `json:"overridden"`
	} `json:"dependencies"`
}

type npmSearchResult struct {
	Array []struct {
		Title       string `json:"name"`
		Description string `json:"description"`
	}
}

func (n Npm) ListInstalled() ([]string, error) {
	input, err := GetStr(Command(n.bin, "ls", "--json"))
	if err != nil {
		return nil, err
	}

	var d npmListResult
	if err := json.Unmarshal([]byte(input), &d); err != nil {
		return nil, err
	}

	var names []string
	for name := range d.Dependencies {
		names = append(names, name)
	}

	return names, err
}

func (n Npm) SearchForPackage(input string) ([]SearchResult, error) {
	searchResult, err := GetStr(Command(n.bin, "search", "--json", input))
	if err != nil {
		return nil, err
	}

	var d npmSearchResult
	if err := json.Unmarshal([]byte(searchResult), &d); err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(d.Array))
	for i := 0; i < len(d.Array); i++ {
		results[i] = SearchResult(d.Array[i])
	}

	Sort(results, input)

	return results, nil
}

func (n Npm) Info(pkg string) (string, error) {
	return GetStr(Command(n.bin, "view", pkg))
}

func (n Npm) Install(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", n.bin, "install", pkg)
}

func (n Npm) UpdatePackage(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", n.bin, "update", pkg)
}

func (n Npm) Remove(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", n.bin, "uninstall", pkg)
}

func (n Npm) UpdateSystem() *Cmd {
	return Command("sudo", n.bin, "update")
}
func (n Npm) IsInstalled(pkg string) (bool, error) {
	cmd := exec.Command(n.bin, "ls", pkg, "--depth=0")
	err := cmd.Run()

	if err == nil {
		return true, nil
	}

	if _, ok := err.(*exec.ExitError); ok {
		// npm returns non-zero if not installed
		return false, nil
	}

	return false, err
}
