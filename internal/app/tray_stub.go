//go:build !windows,!linux,!darwin

package app

// RunTray is a stub implementation for platforms that don't support systray
func RunTray() {
	// No-op for unsupported platforms
}