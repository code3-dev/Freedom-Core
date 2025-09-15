package core

import (
	"strconv"

	flags "github.com/Freedom-Guard/freedom-core/pkg/flag"
	"github.com/Freedom-Guard/freedom-core/internal/server"

)

func StartCore() {
	addr := ":" + strconv.Itoa(flags.AppConfig.Port)
	srv := &server.Server{Addr: addr}
	srv.ListenAndServe()
}
