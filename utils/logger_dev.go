//go:build dev

package utils

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

func InitLog() {
	var err error
	logFile, err = os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	logFile.Write([]byte{})
}

func DestroyLog() {
	logFile.Close()

}

func Log(format string, args ...any) {
	if logFile == nil {
		return
	}

	fmt.Fprintf(logFile, "[%s] %s\n", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, args...))
}
