package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const ENABLE_LOGGING = true

const recentPkgFileName = "pm-tui_recent_pkgs"

var logFile *os.File

func InitLog() {
	var err error
	if ENABLE_LOGGING {
		logFile, err = os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}

		logFile.Write([]byte{})
	}
}

func DestroyLog() {
	if ENABLE_LOGGING {
		logFile.Close()
	}
}

func Log(format string, args ...any) {
	if !ENABLE_LOGGING || logFile == nil {
		return
	}
	fmt.Fprintf(logFile, "[%s] %s\n", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, args...))
}

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
