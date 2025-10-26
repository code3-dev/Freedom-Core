//go:build !windows && !linux && !darwin

package app

// RunTrayWrapper is a no-op on unsupported platforms
func RunTrayWrapper() {
	// Systray is not supported on this platform
}