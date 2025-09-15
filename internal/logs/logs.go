package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func LogPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
<title>Logs</title>
<style>
body { font-family: monospace; background: #111; color: #0f0; padding: 10px; }
pre { white-space: pre-wrap; word-wrap: break-word; }
</style>
</head>
<body>
<h2>Freedom Core Logs</h2>
<pre id="log"></pre>
<script>
const logBox = document.getElementById("log");
const evtSource = new EventSource("/logs/stream");
evtSource.onmessage = function(e) {
	logBox.innerText += e.data + "\n";
	logBox.scrollTop = logBox.scrollHeight;
};
</script>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	}
}

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
