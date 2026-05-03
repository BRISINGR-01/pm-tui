package package_manager

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	. "os/exec"
	. "pm-tui/utils"
)

type Pacman struct{}

func (p Pacman) ListInstalled() ([]string, error) {
	return GetListStr(Command("pacman", "-Qq"))
}

func (p Pacman) SearchForPackage(input string) ([]SearchResult, error) {
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

func (p Pacman) Info(pkg string) (string, error) {
	isInstalled, err := p.IsInstalled(pkg)
	if err != nil {
		return "", err
	}

	flag := "-Qi"
	if !isInstalled {
		flag = "-Si"
	}

	return GetStr(Command("pacman", flag, pkg))
}

func (p Pacman) Install(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", "pacman", "-Sy", pkg)
}

func (p Pacman) UpdatePackage(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", "pacman", "-Syu", pkg)
}

func (p Pacman) Remove(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("sudo", "pacman", "-R", pkg)
}

func (p Pacman) UpdateSystem() *Cmd {
	return Command("sudo", "pacman", "-Syu")
}

func (p Pacman) IsInstalled(pkg string) (bool, error) {
	cmd := exec.Command("pacman", "-Qi", pkg)
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
