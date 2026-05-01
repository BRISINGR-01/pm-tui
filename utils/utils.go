package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const recentPkgFileName = "pm-tui_recent_pkgs"

func AddRecentPkg(name string) {
	path := filepath.Join(os.TempDir(), recentPkgFileName)

	if data, err := os.ReadFile(path); err == nil {
		for line := range strings.SplitSeq(string(data), "\n") {
			if line == name {
				return
			}
		}
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Fprintln(f, name)
}

func GetRecentPkgs() []string {
	if f, err := os.ReadFile(filepath.Join(os.TempDir(), recentPkgFileName)); err == nil {
		return strings.Split(strings.TrimSpace(string(f)), "\n")
	}

	return []string{}
}

func FormatErr(out []byte, err error) error {
	return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
}
