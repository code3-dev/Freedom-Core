package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func LogStreamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		lastLen := 0
		for {
			entries := logger.GetLogs()
			if len(entries) > lastLen {
				for _, e := range entries[lastLen:] {
					fmt.Fprintf(w, "[%s] [%s] %s\n\n", e.Timestamp, levelString(e.Level), e.Message)
				}
				lastLen = len(entries)
				flusher.Flush()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func ClearLogsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.ClearLogs()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("logs cleared"))
	}
}


func levelString(level logger.LogLevel) string {
	switch level {
	case logger.INFO:
		return "INFO 🚀"
	case logger.WARN:
		return "WARN ⚠️"
	case logger.ERROR:
		return "ERROR ❌"
	case logger.DEBUG:
		return "DEBUG 🔍"
	default:
		return "UNKNOWN"
	}
}
