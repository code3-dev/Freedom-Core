package masqueplus

import (
	"fmt"
	"net/http"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func MasquePlusStreamHandler(w http.ResponseWriter, r *http.Request) {
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

	ctx := r.Context()
	lines := make(chan string, 100)
	resultChan := make(chan bool, 1)

	go func() {
		found := RunMasquePlusStream(ctx, args, func(line string) {
			select {
			case lines <- line:
			case <-ctx.Done():
			}
		})
		resultChan <- found
		close(resultChan)
		close(lines)
	}()

	for {
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
				_, _ = w.Write([]byte(fmt.Sprintf("Masque-plus result: %t\n", result)))
				flusher.Flush()
			}
			return
		case <-ctx.Done():
			_, _ = w.Write([]byte("Masque-plus stopped\n"))
			flusher.Flush()
			return
		}
	}
}

func KillMasquePlusHandler(w http.ResponseWriter, r *http.Request) {
	KillMasquePlus()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Masque-plus process kill triggered"))
	logger.Log(logger.INFO, "KillMasquePlus triggered")
}
