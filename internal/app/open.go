package app

import (
	"os/exec"
	"runtime"
)

func OpenUrl(URL string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", URL).Start()
	case "linux":
		exec.Command("xdg-open", URL).Start()
	case "darwin":
		exec.Command("open", URL).Start()
	}
}
