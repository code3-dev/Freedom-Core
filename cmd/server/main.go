package main

import (
	
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
	"github.com/Freedom-Guard/freedom-core/internal/core"
	"github.com/Freedom-Guard/freedom-core/pkg/config"
)

func main() {
	cfg := config.Load()
	logger.Log(logger.INFO, "Freedom-Core is starting... 🚀")
	core.StartCore(cfg)
}
