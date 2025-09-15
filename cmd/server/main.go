package main

import (
	"github.com/Freedom-Guard/freedom-core/internal/app"
	"github.com/Freedom-Guard/freedom-core/internal/core"
	flags "github.com/Freedom-Guard/freedom-core/pkg/flag"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func main() {
	flags.Parse()
	logger.Log(logger.INFO, "Freedom-Core is starting... 🚀")
	go app.RunTray()
	core.StartCore()
}
