package installer

import (
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed hiddify-core.exe
var CoreBinary []byte

//go:embed hiddify-core.dll
var CoreDLL []byte

func PrepareCore() (string, error) {
	var dir string
	var coreName string
	var dllName string

	if runtime.GOOS == "windows" {
		appData, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(appData, "Hiddify", "bin")
		coreName = "hiddify-core.exe"
		dllName = "hiddify-core.dll"
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".local", "bin")
		coreName = "hiddify-core"
		dllName = "hiddify-core.dll"
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	corePath := filepath.Join(dir, coreName)
	dllPath := filepath.Join(dir, dllName)

	if err := ensureFile(corePath, CoreBinary, 0755); err != nil {
		return "", err
	}
	if err := ensureFile(dllPath, CoreDLL, 0644); err != nil {
		return "", err
	}

	return corePath, nil
}

func ensureFile(path string, data []byte, perm os.FileMode) error {
	if fi, err := os.Stat(path); err == nil && fi.Mode().Perm()&0100 != 0 {
		return nil
	}
	return os.WriteFile(path, data, perm)
}
