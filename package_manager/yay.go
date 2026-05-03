package package_manager

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	. "os/exec"
	. "pm-tui/utils"
)

type Yay struct{}

func (y Yay) ListInstalled() ([]string, error) {
	return GetListStr(Command("yay", "-Qq"))
}

func (y Yay) SearchForPackage(input string) ([]SearchResult, error) {
	lines, err := GetListStr(Command("pacman", "-Ss", input))
	if err != nil {
		if err.Error() == "exit status 1: " {
			return []SearchResult{}, nil
		}
		return nil, err
	}

	results := []SearchResult{}

	for i := 0; i < len(lines); i += 2 {
		if strings.Contains(lines[i], "[installed]") || i+1 >= len(lines) {
			continue
		}

		segments := strings.SplitN(lines[i], "/", 2)
		if len(segments) < 2 {
			continue
		}
		title := strings.SplitN(segments[1], " ", 2)
		if len(title) < 2 {
			continue
		}

		results = append(results, SearchResult{
			Title:       title[0],
			Description: lines[i+1],
		})
	}

	Sort(results, input)

	return results, nil
}

func (y Yay) Info(pkg string) (string, error) {
	isInstalled, err := y.IsInstalled(pkg)
	if err != nil {
		return "", err
	}

	flag := "-Qi"
	if !isInstalled {
		flag = "-Si"
	}

	return GetStr(Command("yay", flag, pkg))
}

func (y Yay) Install(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("yay", "-S", pkg)
}

func (y Yay) UpdatePackage(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("yay", "-S", pkg)
}

func (y Yay) Remove(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("yay", "-R", pkg)
}

func (y Yay) UpdateSystem() *Cmd {
	return Command("yay", "-Sua")
}

func (y Yay) IsInstalled(pkg string) (bool, error) {
	cmd := exec.Command("yay", "-Qi", pkg)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
		msg := stderr.String()
		if strings.Contains(msg, "was not found") {
			return false, nil
		}
		return false, fmt.Errorf("%w: %s", err, msg)
	}

	return false, err
}
