//go:build !windows && !linux && !darwin

package app

// RunTray is a no-op on unsupported platforms
func RunTray() {
	// Systray is not supported on this platform
}