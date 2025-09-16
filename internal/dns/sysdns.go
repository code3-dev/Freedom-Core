package sysdns

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func DNSStreamHandler(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	lines := make(chan string, 100)
	done := make(chan struct{})

	go func() {
		manager := NewDNSManager()
		for _, arg := range args {
			if strings.HasPrefix(arg, "set:") {
				paramStr := strings.TrimPrefix(arg, "set:")
				cfg := DNSConfig{}
				for _, kv := range strings.Split(paramStr, ",") {
					parts := strings.SplitN(kv, "=", 2)
					if len(parts) != 2 {
						lines <- "invalid parameter: " + kv
						continue
					}
					key := strings.ToLower(strings.TrimSpace(parts[0]))
					val := strings.TrimSpace(parts[1])
					switch key {
					case "primary":
						cfg.Primary = val
					case "secondary":
						cfg.Secondary = val
					case "enable":
						cfg.Enable = val == "true" || val == "1"
					default:
						lines <- "unknown key: " + key
					}
				}
				if err := manager.SetDNS(&cfg); err != nil {
					lines <- "set failed: " + err.Error()
					logger.Log(logger.DEBUG, "SetDNS failed: "+err.Error())
				} else {
					lines <- fmt.Sprintf("DNS set: %+v", cfg)
					logger.Log(logger.DEBUG, "DNS set: "+cfg.Primary+" "+cfg.Secondary)
				}

			} else if arg == "get" {
				cfg, err := manager.GetDNS()
				if err != nil {
					lines <- "get failed: " + err.Error()
					logger.Log(logger.DEBUG, "GetDNS failed: "+err.Error())
				} else {
					lines <- fmt.Sprintf("current DNS: %+v", cfg)
					logger.Log(logger.DEBUG, "GetDNS: "+cfg.Primary+" "+cfg.Secondary)
				}

			} else if arg == "clear" {
				if err := manager.ClearDNS(); err != nil {
					lines <- "clear failed: " + err.Error()
					logger.Log(logger.DEBUG, "ClearDNS failed: "+err.Error())
				} else {
					lines <- "DNS cleared"
					logger.Log(logger.DEBUG, "DNS cleared")
				}

			} else {
				lines <- "unknown command: " + arg
				logger.Log(logger.DEBUG, "Unknown command: "+arg)
			}
			time.Sleep(100 * time.Millisecond)
		}
		close(done)
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
		case <-done:
			_, _ = w.Write([]byte("done\n"))
			flusher.Flush()
			return
		case <-ctx.Done():
			_, _ = w.Write([]byte("aborted\n"))
			flusher.Flush()
			return
		}
	}
}
