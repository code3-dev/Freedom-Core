package installer

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed hiddify-core
var CoreBinary []byte

func PrepareCore() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	corePath := filepath.Join(dir, "hiddify-core")
	if fi, err := os.Stat(corePath); err == nil && fi.Mode().Perm()&0100 != 0 {
		return corePath, nil
	}
	if err := os.WriteFile(corePath, CoreBinary, 0755); err != nil {
		return "", err
	}
	return corePath, nil
}
