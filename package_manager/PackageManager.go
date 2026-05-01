package package_manager

import "os/exec"

type PackageManager interface {
	ListInstalled() ([]string, error)
	SearchForPackage(input string) ([]SearchResult, error)
	Install(pkg string) *exec.Cmd
	UpdatePackage(pkg string) *exec.Cmd
	UpdateSystem() *exec.Cmd
	Remove(pkg string) *exec.Cmd
	Info(pkg string) (string, error)
	IsInstalled(pkg string) (bool, error)
}

type ProviderType int

const (
	ProviderPacman ProviderType = iota
	ProviderYay
	ProviderApt
	ProviderNpm
	ProviderNpmGlobal
	ProviderPip
	ProviderRpm
	ProviderYum
)

type SearchResult struct {
	Title       string
	Description string
}
