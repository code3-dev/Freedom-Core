package main

import (
	"github.com/Freedom-Guard/freedom-core/internal/app"
	"github.com/Freedom-Guard/freedom-core/internal/core"
	"github.com/Freedom-Guard/freedom-core/pkg/config"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func main() {
	cfg := config.Load()
	logger.Log(logger.INFO, "Freedom-Core is starting... 🚀")
	go app.RunTray()
	core.StartCore(cfg)
}
