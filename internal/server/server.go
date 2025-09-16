package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Freedom-Guard/freedom-core/internal/logs"
	sysproxy "github.com/Freedom-Guard/freedom-core/internal/proxy"
	dns "github.com/Freedom-Guard/freedom-core/internal/dns"
	flags "github.com/Freedom-Guard/freedom-core/pkg/flag"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
	"github.com/getlantern/systray"
)

type Server struct {
	Addr string
}

func (s *Server) ListenAndServe() {
	logger.Log(logger.INFO, "Server listening on "+s.Addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/hiddify/start", HiddifyStreamHandler)
	mux.HandleFunc("/hiddify/stop", KillHiddifyHandler)
	mux.HandleFunc("/singbox/start", SingBoxStreamHandler)
	mux.HandleFunc("/singbox/stop", KillSingBoxHandler)
	mux.HandleFunc("/xray/start", XrayStreamHandler)
	mux.HandleFunc("/xray/stop", KillXrayHandler)
	mux.HandleFunc("/proxy/start", sysproxy.ProxyStreamHandler)
	mux.HandleFunc("/dns/start", dns.DNSStreamHandler)
	mux.HandleFunc("/logs", logs.LogPageHandler())
	mux.HandleFunc("/logs/stream", logs.LogStreamHandler())

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(flags.AppConfig.Version))
	})

	mux.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is shutting down"))
		go func() {
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		}()
	})

	handler := corsMiddleware(mux)

	srv := &http.Server{
		Addr:    s.Addr,
		Handler: handler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log(logger.ERROR, "Server error: "+err.Error())
		}
	}()

	<-stop
	logger.Log(logger.INFO, "Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log(logger.ERROR, "Server forced to shutdown: "+err.Error())
	}

	logger.Log(logger.INFO, "Server exited cleanly")

	systray.Quit()
	os.Exit(0)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
