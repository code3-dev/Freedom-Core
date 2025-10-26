//go:build !((windows && amd64) || (linux && amd64) || (darwin && (amd64 || arm64)))

package app

// RunTrayWrapper is a no-op on unsupported platforms
func RunTrayWrapper() {
	// Systray is not supported on this platform
}