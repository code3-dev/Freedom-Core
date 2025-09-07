package logger

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
)

var logLevelStrings = map[LogLevel]string{
	INFO:  "INFO 🚀",
	WARN:  "WARN ⚠️",
	ERROR: "ERROR ❌",
}

func Log(level LogLevel, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] %s\n", timestamp, logLevelStrings[level], message)
}
