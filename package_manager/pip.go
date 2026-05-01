package package_manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	. "os/exec"
	. "pm-tui/utils"
)

type Pip struct{}

func (p Pip) ListInstalled() ([]string, error) {
	out, err := GetStr(Command("pip", "list", "--format", "json"))
	if err != nil {
		return nil, err
	}

	var pkgs []struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal([]byte(out), &pkgs); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(pkgs))
	for _, pkg := range pkgs {
		names = append(names, pkg.Name)
	}

	return names, nil
}

func (p Pip) SearchForPackage(_ string) ([]SearchResult, error) {
	return nil, fmt.Errorf("pip does not support search")
}

func (p Pip) Info(pkg string) (string, error) {
	return GetStr(Command("pip", "show", pkg))
}

func (p Pip) Install(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("pip", "install", pkg)
}

func (p Pip) UpdatePackage(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("pip", "install", "-U", pkg)
}

func (p Pip) Remove(pkg string) *Cmd {
	AddRecentPkg(pkg)
	return Command("pip", "uninstall", pkg)
}

func (p Pip) UpdateSystem() *Cmd {
	cmd := Command("pip", "list", "--outdated", "--format=json")

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
		command.WriteString(fmt.Sprintf("pip install --upgrade %s\n", v.Name))
	}

	// execute all upgrades in one shell
	return Command("sh", "-c", command.String())
}

// Check installed
func (p Pip) IsInstalled(pkg string) (bool, error) {
	cmd := exec.Command("pip", "show", pkg)

	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	if exit, ok := err.(*exec.ExitError); ok && exit.ExitCode() == 1 {
		return false, nil
	}

	return false, err
}
