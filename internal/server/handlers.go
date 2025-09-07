package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
	"github.com/Freedom-Guard/freedom-core/internal/xray"
	"github.com/Freedom-Guard/freedom-core/internal/hiddify"
)

func XRayStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	logger.Log(logger.INFO, "X-Ray streaming started")

	for i := 0; i < 5; i++ {
		part := fmt.Sprintf("XRay step %d: %s\n", i+1, xray.RunXRay("test"))
		_, _ = w.Write([]byte(part))
		flusher.Flush()
		time.Sleep(1 * time.Second)
	}

	_, _ = w.Write([]byte("XRay stream finished\n"))
	flusher.Flush()
}

func HiddifyStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	logger.Log(logger.INFO, "Hiddify streaming started")

	part := fmt.Sprintf("Hiddify result: %t\n", hiddify.RunHiddify("test"))
	_, _ = w.Write([]byte(part))
	flusher.Flush()

	_, _ = w.Write([]byte("Hiddify stream finished\n"))
	flusher.Flush()
}
