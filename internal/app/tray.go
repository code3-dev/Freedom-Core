//go:build (windows && amd64) || (linux && amd64) || (darwin && (amd64 || arm64))

package app

import (
	"os"
	"strconv"

	flags "github.com/Freedom-Guard/freedom-core/pkg/flag"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
	"github.com/getlantern/systray"

	_ "embed"
)

//go:embed icon.ico
var iconData []byte

// RunTray starts the system tray
func RunTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Freedom Core")
	systray.SetTooltip("Freedom Core")

	mLogs := systray.AddMenuItem("Open Logs", "Show logs in browser")
	mFWeb := systray.AddMenuItem("Open FCORE WEB", "Open FCORE WEB in browser")
	mExit := systray.AddMenuItem("Exit", "Quit the application")

	go func() {
		for {
			select {
			case <-mLogs.ClickedCh:
				openLogsWindow()
			case <-mFWeb.ClickedCh:
				OpenUrl("https://freedom-guard.github.io/FCORE-WEB-CLI/")
			case <-mExit.ClickedCh:
				logger.Log(logger.INFO, "Exiting application")
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func openLogsWindow() {
	url := "http://localhost:" + strconv.Itoa(flags.AppConfig.Port) + "/logs/stream"
	OpenUrl(url)
}

func onExit() {
	logger.Log(logger.INFO, "Freedom Core exited")
}