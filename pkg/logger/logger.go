package logger

import (
	"fmt"
	"sync"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
	DEBUG
)

var logLevelStrings = map[LogLevel]string{
	INFO:  "INFO 🚀",
	WARN:  "WARN ⚠️",
	ERROR: "ERROR ❌",
	DEBUG: "DEBUG 🔍",
}

type LogEntry struct {
	Timestamp string
	Level     LogLevel
	Message   string
}

var (
	mu   sync.Mutex
	logs []LogEntry
)

func Log(level LogLevel, message string) {
	mu.Lock()
	defer mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
	}
	logs = append(logs, entry)

	fmt.Printf("[%s] [%s] %s\n", timestamp, logLevelStrings[level], message)
}

func GetLogs() []LogEntry {
	mu.Lock()
	defer mu.Unlock()

	cp := make([]LogEntry, len(logs))
	copy(cp, logs)
	return cp
}

func ClearLogs() {
	mu.Lock()
	defer mu.Unlock()

	logs = []LogEntry{}
}
