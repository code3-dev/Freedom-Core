package core

import (
	"strconv"

	"github.com/Freedom-Guard/freedom-core/internal/server"
	flags "github.com/Freedom-Guard/freedom-core/pkg/flag"
)

func StartCore(cfg *flags.Config) {
	addr := ":" + strconv.Itoa(cfg.Port)
	srv := &server.Server{Addr: addr}
	srv.ListenAndServe()
}
