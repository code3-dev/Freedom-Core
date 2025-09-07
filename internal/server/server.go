package server

import (
	"log"
	"net/http"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) ListenAndServe() {
	logger.Log(logger.INFO, "Server listening on "+s.addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/xray/stream", XRayStreamHandler)
	mux.HandleFunc("/hiddify/stream", HiddifyStreamHandler)

	if err := http.ListenAndServe(s.addr, mux); err != nil {
		log.Fatal(err)
	}
}
