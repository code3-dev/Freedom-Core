package core

import (
	"strconv"

	"github.com/Freedom-Guard/freedom-core/pkg/config"
	"github.com/Freedom-Guard/freedom-core/internal/server"
)

func StartCore(cfg *config.Config) {
	addr := ":" + strconv.Itoa(cfg.Port)
	srv := &server.Server{Addr: addr}
	srv.ListenAndServe()
}
