package core

import (
	"strconv"

	"github.com/Freedom-Guard/freedom-core/internal/server"
	flags "github.com/Freedom-Guard/freedom-core/pkg/flag"
)

func StartCore() {
	addr := ":" + strconv.Itoa(flags.AppConfig.Port)
	srv := &server.Server{Addr: addr}
	srv.ListenAndServe()
}
