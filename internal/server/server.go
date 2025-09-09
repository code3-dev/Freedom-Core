package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

type Server struct {
	Addr string
}

func (s *Server) ListenAndServe() {
	logger.Log(logger.INFO, "Server listening on "+s.Addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/hiddify/stream", HiddifyStreamHandler)

	srv := &http.Server{
		Addr:    s.Addr,
		Handler: mux,
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
}
