package app

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"

	flags "github.com/Freedom-Guard/freedom-core/pkg/flag"
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

	mLogs := systray.AddMenuItem("Open Logs", "Show logs in browser")
	mExit := systray.AddMenuItem("Exit", "Quit the application")

	go func() {
		for {
			select {
			case <-mLogs.ClickedCh:
				openLogsWindow()
			case <-mExit.ClickedCh:
				logger.Log(logger.INFO, "Exiting application")
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func openLogsWindow() {

	url := "http://localhost:" + strconv.Itoa(flags.AppConfig.Port) + "/logs"
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	}
}

func onExit() {
	logger.Log(logger.INFO, "Freedom Core exited")
}
