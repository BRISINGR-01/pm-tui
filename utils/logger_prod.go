//go:build !dev

package utils

func InitLog() {}

func DestroyLog() {}

func Log(format string, args ...any) {}
