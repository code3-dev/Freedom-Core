package core

import (
	"strconv"

	"github.com/Freedom-Guard/freedom-core/pkg/config"
	"github.com/Freedom-Guard/freedom-core/internal/server"
)

func StartCore(cfg *config.Config) {
	srv := server.NewServer(":"+ strconv.Itoa(cfg.Port))
	srv.ListenAndServe()
}