package package_manager

import (
	. "os/exec"
	"sort"
	"strings"

	"github.com/sahilm/fuzzy"

	. "pm-tui/utils"
)

func GetBytes(cmd *Cmd) ([]byte, error) {
	out, err := cmd.CombinedOutput()

	if err != nil {
		return nil, FormatErr(out, err)
	}

	return out, nil
}

func GetStr(cmd *Cmd) (string, error) {
	out, err := GetBytes(cmd)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func GetListStr(cmd *Cmd) ([]string, error) {
	pkgs, err := GetStr(cmd)
	if err != nil {
		return nil, err
	}

	return strings.Split(pkgs, "\n"), nil
}

func GetProviderType(provider string) ProviderType {
	switch provider {
	case "pacman":
		return ProviderPacman
	case "apt":
		return ProviderApt
	case "npm":
		return ProviderNpm
	case "pip":
		return ProviderPip
	case "yay":
		return ProviderYay
	case "yum":
		return ProviderYum
	case "rpm":
		return ProviderRpm
	}

	return ProviderPacman
}

func GetProviderStr(managerType ProviderType) string {
	switch managerType {
	case ProviderPacman:
		return "pacman"
	case ProviderYay:
		return "yay"
	case ProviderApt:
		return "apt"
	case ProviderNpm:
		return "npm"
	case ProviderNpmGlobal:
		return "npm (global)"
	case ProviderPip:
		return "pip"
	case ProviderRpm:
		return "rpm"
	case ProviderYum:
		return "yum"
	}

	panic("No package manager found")
}

func Sort(results []SearchResult, input string) {
	titles := make([]string, len(results))
	for k, r := range results {
		titles[k] = r.Title
	}

	matches := fuzzy.Find(input, titles)
	// higher score => better match
	scoreOf := make(map[string]int)
	for _, m := range matches {
		scoreOf[m.Str] = m.Score
	}

	sort.Slice(results, func(i, j int) bool {
		return scoreOf[results[i].Title] > scoreOf[results[j].Title]
	})
}

func getFirstSegment(input string) string {
	segment := strings.SplitN(input, "/", 2)
	if len(segment) < 2 {
		return ""
	}

	return segment[0]
}

func NewPackageManager(provider ProviderType) PackageManager {
	switch provider {
	case ProviderPacman:
		return &Pacman{}
	case ProviderYay:
		return &Yay{}
	case ProviderNpm:
		return &Npm{"npm"}
	case ProviderNpmGlobal:
		return &Npm{"npm -g"}
	case ProviderPip:
		return &Pip{}
	case ProviderApt:
		return &Apt{}
	default:
		panic("No package manager found")
	}
}
