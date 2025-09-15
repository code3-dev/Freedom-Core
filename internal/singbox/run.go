package singbox

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Freedom-Guard/freedom-core/internal/run"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

var releaseVersion = "1.12.8"
var procMu sync.Mutex
var currentProc *exec.Cmd

func getCoreURL() (url, filename string, isTarGz bool, err error) {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("https://github.com/SagerNet/sing-box/releases/download/v%s/sing-box-%s-windows-amd64.zip", releaseVersion, releaseVersion),
			"sing-box.exe", false, nil
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return fmt.Sprintf("https://github.com/SagerNet/sing-box/releases/download/v%s/sing-box-%s-linux-x86_64.tar.gz", releaseVersion, releaseVersion),
				"sing-box", true, nil
		case "arm64":
			return fmt.Sprintf("https://github.com/SagerNet/sing-box/releases/download/v%s/sing-box-%s-linux-aarch64.tar.gz", releaseVersion, releaseVersion),
				"sing-box", true, nil
		}
	}
	return "", "", false, errors.New("unsupported OS/ARCH")
}

func corePath(subdir string) string {
	dir, _ := os.UserCacheDir()
	dir = filepath.Join(dir, "freedom-core", subdir)
	os.MkdirAll(dir, 0o755)
	return dir
}

func downloadAndExtract() (string, error) {
	url, filename, isTarGz, err := getCoreURL()
	if err != nil {
		return "", err
	}

	dir := corePath("singbox")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	dest := filepath.Join(dir, filename)
	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	}

	tmpFile := dest
	if isTarGz {
		tmpFile += ".tar.gz"
	} else {
		tmpFile += ".zip"
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(tmpFile)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return "", err
	}

	if isTarGz {
		err = extractTarGz(tmpFile, dest)
	} else {
		err = extractZip(tmpFile, dest)
	}
	if err != nil {
		return "", err
	}

	if runtime.GOOS != "windows" {
		os.Chmod(dest, 0o755)
	}

	return dest, nil
}

func extractZip(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	baseName := filepath.Base(dest)
	for _, f := range r.File {
		if filepath.Base(f.Name) == baseName {
			fr, _ := f.Open()
			defer fr.Close()
			out, _ := os.Create(dest)
			defer out.Close()
			_, err = io.Copy(out, fr)
			return err
		}
	}
	return errors.New("file not found in zip")
}

func extractTarGz(tarGzPath, dest string) error {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	baseName := filepath.Base(dest)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if filepath.Base(hdr.Name) == baseName {
			out, _ := os.Create(dest)
			defer out.Close()
			_, err = io.Copy(out, tr)
			return err
		}
	}
	return errors.New("file not found in tar.gz")
}

func PrepareCore() (string, error) {
	return downloadAndExtract()
}

func RunSingBoxStream(ctx context.Context, args []string, callback func(string)) bool {

	path, err := PrepareCore()
	if err != nil {
		callback("Sing-box core not installed: " + err.Error())
		return false
	}

	cmd := exec.CommandContext(ctx, path, args...)

	run.SetupCmd(cmd)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		callback("Failed to start sing-box: " + err.Error())
		return false
	}

	procMu.Lock()
	currentProc = cmd
	procMu.Unlock()

	found := false
	done := make(chan struct{}, 2)

	scan := func(r io.Reader, isErr bool) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			callback(line)
			if isErr {
				logger.Log(logger.ERROR, "sing-box stderr: "+line)
			} else {
				logger.Log(logger.INFO, "sing-box stdout: "+line)
				if strings.Contains(strings.ToLower(line), "started") {
					found = true
				}
			}
		}
		done <- struct{}{}
	}

	go scan(stdout, false)
	go scan(stderr, true)
	go func() { <-ctx.Done(); KillSingBox() }()

	<-done
	<-done
	cmd.Wait()

	procMu.Lock()
	currentProc = nil
	procMu.Unlock()
	return found
}

func KillSingBox() {
	procMu.Lock()
	defer procMu.Unlock()
	if currentProc != nil && currentProc.Process != nil {
		_ = currentProc.Process.Kill()
		logger.Log(logger.INFO, "Sing-box process killed")
		currentProc = nil
	}
}
