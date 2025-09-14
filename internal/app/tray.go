package app

import (
	"os"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
	"github.com/getlantern/systray"

	_ "embed"
)

//go:embed icon.ico
var iconData []byte

func RunTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Freedom Core")
	systray.SetTooltip("Freedom Core")

	mExit := systray.AddMenuItem("Exit", "Quit the application")

	go func() {
		for {
			select {
			case <-mExit.ClickedCh:
				logger.Log(logger.INFO, "Exiting application")
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func onExit() {
	logger.Log(logger.INFO, "Freedom Core exited")
}
