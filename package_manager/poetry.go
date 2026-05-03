package package_manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	. "os/exec"
	. "pm-tui/utils"
)

type Poetry struct{}

func (p Poetry) ListInstalled() ([]string, error) {
	out, err := Command("poetry", "show", "--format", "json").Output()

	if err != nil {
		return nil, FormatErr(out, err)
	}

	var pkgs []struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(out, &pkgs); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(pkgs))
	for _, pkg := range pkgs {
		names = append(names, pkg.Name)
	}

	return names, nil
}

func (p Poetry) SearchForPackage(input string) ([]SearchResult, error) {
	searchResult, err := GetListStr(Command("poetry", "search", input))
	// example output:
	//  Package Version Source Description
	//  flask   0.1     PyPI
	//  flask   0.2     PyPI

	if err != nil {
		return nil, err
	}

	var result []SearchResult
	spaces := regexp.MustCompile(`\s+`)
	for i, res := range searchResult {
		if i == 0 {
			continue
		}
		segments := spaces.Split(res, -1)
		if len(segments) >= 2 {
			result = append(result, SearchResult{
				Title: fmt.Sprintf("%s:%s", segments[1], segments[2]),
			})
		}
	}

	return result, nil
}

func (p Poetry) Info(pkg string) (string, error) {
	Log(pkg)
	if isInstalled, _ := p.IsInstalled(pkg); isInstalled {
		return GetStr(Command("poetry", "show", pkg))
	}

	return "", nil
}

func (p Poetry) Install(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("poetry", "add", pkg)
}

func (p Poetry) UpdatePackage(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("poetry", "update", pkg)
}

func (p Poetry) Remove(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("poetry", "remove", pkg)
}

func (p Poetry) UpdateSystem() *Cmd {
	cmd := Command("poetry", "list", "--outdated", "--format=json")

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return Command("echo", fmt.Sprintf("%v: %s", err, stderr.String()))
	}

	var pkgs []struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(out.Bytes(), &pkgs); err != nil {
		return Command("echo", err.Error())
	}

	if len(pkgs) == 0 {
		return Command("echo", "All packages are up to date")
	}

	var command strings.Builder
	for _, v := range pkgs {
		command.WriteString(fmt.Sprintf("poetry install --upgrade %s\n", v.Name))
	}

	// execute all upgrades in one shell
	return Command("sh", "-c", command.String())
}

func (p Poetry) IsInstalled(pkg string) (bool, error) {
	cmd := exec.Command("poetry", "show", pkg)

	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
		return false, nil
	}

	return false, err
}
