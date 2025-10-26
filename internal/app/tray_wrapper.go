//go:build (windows && amd64) || (linux && amd64) || (darwin && (amd64 || arm64))

package app

// RunTrayWrapper calls the actual systray implementation
func RunTrayWrapper() {
	// This will call RunTray which is only available on supported platforms
	// due to build constraints in tray.go
	RunTray()
}