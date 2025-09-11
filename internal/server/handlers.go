package server

import (
	"fmt"
	"net/http"

	"github.com/Freedom-Guard/freedom-core/internal/hiddify"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func HiddifyStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	args := r.URL.Query()["args"]
	if len(args) == 0 {
		http.Error(w, "Missing arguments", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	logger.Log(logger.INFO, fmt.Sprintf("Hiddify streaming started with args: %v", args))

	ctx := r.Context()
	lines := make(chan string, 100)
	resultChan := make(chan bool, 1)

	go func() {
		found := hiddify.RunHiddifyStream(ctx, args, func(line string) {
			select {
			case lines <- line:
			case <-ctx.Done():
			}
		})
		resultChan <- found
		close(resultChan)
		close(lines)
	}()

	streamOpen := true

	for streamOpen {
		select {
		case line, ok := <-lines:
			if !ok {
				lines = nil
				continue
			}
			_, _ = w.Write([]byte(line + "\n"))
			flusher.Flush()
		case result, ok := <-resultChan:
			if ok {
				_, _ = w.Write([]byte(fmt.Sprintf("Hiddify result: %t\n", result)))
				flusher.Flush()
			}
			streamOpen = false
		case <-ctx.Done():
			_, _ = w.Write([]byte("Hiddify process stopped\n"))
			flusher.Flush()
			streamOpen = false
		}
	}

	_, _ = w.Write([]byte("Hiddify stream finished\n"))
	flusher.Flush()
}

func KillHiddifyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hiddify.KillHiddify(ctx)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("hiddify process kill triggered"))
	logger.Log(logger.INFO, "KillHiddify triggered")
}

