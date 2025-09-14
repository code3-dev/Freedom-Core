package hiddify

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
	"syscall"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
	helpers "github.com/Freedom-Guard/freedom-core/pkg/utils"
)

var releaseVersion = "3.2.0"
var procMu sync.Mutex
var currentProc *exec.Cmd

func getCoreURL() (url, filename string, isTarGz bool, err error) {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("https://github.com/hiddify/hiddify-core/releases/download/v%s/hiddify-core-windows-amd64.tar.gz", releaseVersion),
			"HiddifyCli.exe", true, nil
	case "linux":
		return fmt.Sprintf("https://github.com/hiddify/hiddify-core/releases/download/v%s/hiddify-core-linux-amd64.tar.gz", releaseVersion),
			"HiddifyCli", true, nil
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

	dir := corePath("hiddify")
	dest := filepath.Join(dir, filename)

	if _, err := os.Stat(dest); err == nil {
		logger.Log(logger.INFO, "Hiddify core already exists")
		return dest, nil
	}

	tmpFile := dest
	if isTarGz {
		tmpFile += ".tar.gz"
	} else {
		tmpFile += ".zip"
	}

	logger.Log(logger.INFO, "Downloading Hiddify core from "+url)
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

	var execPath string
	if isTarGz {
		execPath, err = extractTarGzAll(tmpFile, dir)
	} else {
		execPath, err = extractZipAll(tmpFile, dir)
	}
	if err != nil {
		return "", err
	}

	if runtime.GOOS != "windows" {
		os.Chmod(execPath, 0o755)
	}

	logger.Log(logger.INFO, "Hiddify core prepared at "+execPath)
	return execPath, nil
}

func extractZipAll(zipPath, destDir string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var execPath string
	for _, f := range r.File {
		outPath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0o755)
			continue
		}
		fr, _ := f.Open()
		defer fr.Close()
		os.MkdirAll(filepath.Dir(outPath), 0o755)
		out, _ := os.Create(outPath)
		defer out.Close()
		_, err = io.Copy(out, fr)
		if err != nil {
			return "", err
		}

		base := strings.ToLower(filepath.Base(f.Name))
		if base == "hiddifycli.exe" || base == "hiddifycli" {
			execPath = outPath
		}
	}

	if execPath == "" {
		return "", errors.New("Hiddify executable not found in zip")
	}
	return execPath, nil
}

func extractTarGzAll(tarGzPath, destDir string) (string, error) {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var execPath string

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		outPath := filepath.Join(destDir, hdr.Name)
		if hdr.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0o755)
			continue
		}

		os.MkdirAll(filepath.Dir(outPath), 0o755)
		out, _ := os.Create(outPath)
		defer out.Close()
		_, err = io.Copy(out, tr)
		if err != nil {
			return "", err
		}

		base := strings.ToLower(filepath.Base(hdr.Name))
		if base == "hiddifycli.exe" || base == "hiddifycli" {
			execPath = outPath
		}
	}

	if execPath == "" {
		return "", errors.New("Hiddify executable not found in tar.gz")
	}
	return execPath, nil
}

func PrepareCore() (string, error) {
	return downloadAndExtract()
}

func RunHiddifyStream(ctx context.Context, args []string, callback func(string)) bool {
	allowed := helpers.AllowDialog("آیا اجازه می‌دهید هسته هیدیفای اجرا شود؟ / Do you allow Hiddify Core to run?")
	if allowed {
		logger.Log(logger.INFO, "کاربر اجازه داد هسته اجرا شود / User allowed the Core to run")
	} else {
		logger.Log(logger.ERROR, "کاربر اجازه نداد هسته اجرا شود / User denied the Core")
		return false
	}

	path, err := PrepareCore()
	if err != nil {
		callback("Hiddify core not installed: " + err.Error())
		return false
	}

	cmd := exec.CommandContext(ctx, path, args...)

	const CREATE_NO_WINDOW = 0x08000000
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: CREATE_NO_WINDOW,
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		callback("Failed to start Hiddify: " + err.Error())
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
				logger.Log(logger.ERROR, "Hiddify stderr: "+line)
			} else {
				logger.Log(logger.INFO, "Hiddify stdout: "+line)
				if strings.Contains(line, "CORE STARTED") {
					found = true
				}
			}
		}
		done <- struct{}{}
	}

	go scan(stdout, false)
	go scan(stderr, true)
	go func() { <-ctx.Done(); KillHiddify() }()

	<-done
	<-done
	cmd.Wait()

	procMu.Lock()
	currentProc = nil
	procMu.Unlock()
	return found
}

func KillHiddify() {
	procMu.Lock()
	defer procMu.Unlock()
	if currentProc != nil && currentProc.Process != nil {
		_ = currentProc.Process.Kill()
		logger.Log(logger.INFO, "Hiddify process killed")
		currentProc = nil
	}
}
