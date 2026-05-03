package package_manager

import (
	"slices"
	"strings"

	. "os/exec"
	. "pm-tui/utils"
)

type Apt struct{}

func (a Apt) ListInstalled() ([]string, error) {
	input, err := GetListStr(Command("apt", "list", "--installed"))
	if err != nil {
		return nil, err
	}

	var names []string
	for _, pkgInfo := range input {
		name := getFirstSegment(pkgInfo)
		if len(name) < 1 {
			continue
		}

		names = append(names, name)
	}

	return names, err
}

func (a Apt) SearchForPackage(input string) ([]SearchResult, error) {
	searchResult, err := GetListStr(Command("apt", "search", input))
	if err != nil {
		return nil, err
	}

	results := []SearchResult{}

	i := slices.IndexFunc(searchResult, func(item string) bool {
		return strings.Contains(item, "/")
	})
	if i < 2 {
		return results, nil
	}
	for ; i < len(searchResult); i += 3 {
		results = append(results, SearchResult{
			Title:       getFirstSegment(searchResult[i]),
			Description: searchResult[i+1],
		})
	}

	Sort(results, input)

	return results, nil
}

func (a Apt) Info(pkg string) (string, error) {
	return GetStr(Command("apt", "show", pkg))
}

func (a Apt) Install(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", "apt", "install", pkg)
}

func (a Apt) UpdatePackage(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", "apt", "install", "--only-upgrade", pkg)
}

func (a Apt) Remove(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", "apt", "purge", pkg)
}

func (a Apt) UpdateSystem() *Cmd {
	return Command("bash", "-c", "sudo apt update && sudo apt upgrade")
}

func (a Apt) IsInstalled(pkg string) (bool, error) {
	pkgs, err := a.ListInstalled()
	if err != nil {
		return false, err
	}

	return slices.Contains(pkgs, pkg), nil
}
