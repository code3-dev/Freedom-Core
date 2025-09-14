package server

import (
	"fmt"
	"net/http"

	"github.com/Freedom-Guard/freedom-core/internal/hiddify"
	"github.com/Freedom-Guard/freedom-core/internal/singbox"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
	helpers "github.com/Freedom-Guard/freedom-core/pkg/utils"
)

func HiddifyStreamHandler(w http.ResponseWriter, r *http.Request) {
	if helpers.AllowDialog("آیا به هسته هیدیفای اجازه اجرا می‌دهید؟") {
	} else {
		logger.Log(logger.ERROR, "Core blocked dialog")
		return
	}
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
	hiddify.KillHiddify()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("hiddify process kill triggered"))
	logger.Log(logger.INFO, "KillHiddify triggered")
}

func SingBoxStreamHandler(w http.ResponseWriter, r *http.Request) {
	if helpers.AllowDialog("آیا به هسته Sing-Box اجازه اجرا می‌دهید؟") {
	} else {
		logger.Log(logger.ERROR, "Sing-Box blocked dialog")
		return
	}

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

	logger.Log(logger.INFO, fmt.Sprintf("Sing-Box streaming started with args: %v", args))

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
				_, _ = w.Write([]byte(fmt.Sprintf("Sing-Box result: %t\n", result)))
				flusher.Flush()
			}
			streamOpen = false
		case <-ctx.Done():
			_, _ = w.Write([]byte("Sing-Box process stopped\n"))
			flusher.Flush()
			streamOpen = false
		}
	}

	_, _ = w.Write([]byte("Sing-Box stream finished\n"))
	flusher.Flush()
}

func KillSingBoxHandler(w http.ResponseWriter, r *http.Request) {
	singbox.KillSingBox()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Sing-Box process kill triggered"))
	logger.Log(logger.INFO, "KillSingBox triggered")
}
