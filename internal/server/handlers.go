package server

import (
	"fmt"
	"net/http"

	"github.com/Freedom-Guard/freedom-core/internal/hiddify"
	"github.com/Freedom-Guard/freedom-core/internal/singbox"
	"github.com/Freedom-Guard/freedom-core/internal/xray"
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
				_, _ = w.Write([]byte(fmt.Sprintf("Hiddify result: %t\n", result)))
				flusher.Flush()
			}
			return
		case <-ctx.Done():
			_, _ = w.Write([]byte("Hiddify stopped\n"))
			flusher.Flush()
			return
		}
	}
}

func KillHiddifyHandler(w http.ResponseWriter, r *http.Request) {
	hiddify.KillHiddify()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hiddify process kill triggered"))
	logger.Log(logger.INFO, "KillHiddify triggered")
}

func SingBoxStreamHandler(w http.ResponseWriter, r *http.Request) {
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
		found := singbox.RunSingBoxStream(ctx, args, func(line string) {
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
				_, _ = w.Write([]byte(fmt.Sprintf("Sing-box result: %t\n", result)))
				flusher.Flush()
			}
			return
		case <-ctx.Done():
			_, _ = w.Write([]byte("Sing-box stopped\n"))
			flusher.Flush()
			return
		}
	}
}

func KillSingBoxHandler(w http.ResponseWriter, r *http.Request) {
	singbox.KillSingBox()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Sing-box process kill triggered"))
	logger.Log(logger.INFO, "KillSingBox triggered")
}

func XrayStreamHandler(w http.ResponseWriter, r *http.Request) {
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
		found := xray.RunXrayStream(ctx, args, func(line string) {
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
				_, _ = w.Write([]byte(fmt.Sprintf("Xray result: %t\n", result)))
				flusher.Flush()
			}
			return
		case <-ctx.Done():
			_, _ = w.Write([]byte("Xray stopped\n"))
			flusher.Flush()
			return
		}
	}
}

func KillXrayHandler(w http.ResponseWriter, r *http.Request) {
	xray.KillXray()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Xray process kill triggered"))
	logger.Log(logger.INFO, "KillXray triggered")
}
