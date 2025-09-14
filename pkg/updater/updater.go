package updater

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

const (
	githubBaseURL = "https://github.com/Freedom-Guard/Freedom-Core/releases/latest/download/"
)

func Update() error {
	var binaryName string
	switch runtime.GOOS {
	case "linux":
		binaryName = "freedom-core-linux-x64"
	case "windows":
		binaryName = "freedom-core-windows-x64.exe"
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	url := githubBaseURL + binaryName
	fmt.Println("Downloading:", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(binaryName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Update completed successfully!")
	return nil
}
