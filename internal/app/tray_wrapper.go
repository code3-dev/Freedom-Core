//go:build windows || linux || darwin

package app

// RunTrayWrapper calls the actual systray implementation
func RunTrayWrapper() {
	RunTray()
}