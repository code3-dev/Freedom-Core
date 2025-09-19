package masqueplus

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
	helpers "github.com/Freedom-Guard/freedom-core/pkg/utils"
)

var releaseVersion = "v2.8.0"
var procMu sync.Mutex
var currentProc *exec.Cmd

func getCoreURL() (url string, isTarGz bool, err error) {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-windows_amd64.zip", releaseVersion), false, nil
		case "arm64":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-windows_arm64.zip", releaseVersion), false, nil
		}
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-linux_amd64.tar.gz", releaseVersion), true, nil
		case "arm64":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-linux_arm64.tar.gz", releaseVersion), true, nil
		case "armv6":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-linux_armv6.zip", releaseVersion), false, nil
		case "armv7":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-linux_armv7.zip", releaseVersion), false, nil
		case "mips":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-linux_mips.zip", releaseVersion), false, nil
		case "mips64":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-linux_mips64.zip", releaseVersion), false, nil
		case "mips64le":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-linux_mips64le.zip", releaseVersion), false, nil
		case "mipsle":
			return fmt.Sprintf("https://github.com/ircfspace/masque-plus/releases/download/%s/masque-plus-linux_mipsle.zip", releaseVersion), false, nil
		}
	}
	return "", false, errors.New("unsupported OS/ARCH")
}

func corePath(subdir string) string {
	dir, _ := os.UserCacheDir()
	dir = filepath.Join(dir, "freedom-core", subdir)
	os.MkdirAll(dir, 0o755)
	return dir
}

func downloadAndExtract() (string, error) {
	url, isTarGz, err := getCoreURL()
	if err != nil {
		return "", err
	}

	dir := corePath("masqueplus")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	tmpFile := filepath.Join(dir, "masqueplus_download")
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
		if err := extractTarGzAll(tmpFile, dir); err != nil {
			return "", err
		}
	} else {
		if err := extractZipAll(tmpFile, dir); err != nil {
			return "", err
		}
	}

	return dir, nil
}

func extractZipAll(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0o755); err != nil {
			return err
		}

		fr, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			fr.Close()
			return err
		}
		_, err = io.Copy(out, fr)
		fr.Close()
		out.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func extractTarGzAll(tarGzPath, destDir string) error {
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
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fpath := filepath.Join(destDir, hdr.Name)
		if hdr.FileInfo().IsDir() {
			os.MkdirAll(fpath, hdr.FileInfo().Mode())
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0o755); err != nil {
			return err
		}

		out, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, hdr.FileInfo().Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(out, tr)
		out.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func PrepareCore() (string, error) {
	return downloadAndExtract()
}

func RunMasquePlusStream(ctx context.Context, args []string, callback func(string)) bool {
	dir, err := PrepareCore()
	if err != nil {
		callback("Masque-plus core not installed: " + err.Error())
		return false
	}

	path := filepath.Join(dir, "masque-plus")
	if runtime.GOOS == "windows" {
		path += ".exe"
	}
	cmd := exec.CommandContext(ctx, path, args...)
	run.SetupCmd(cmd)

	helpers.ShowInfo("Masque Plus Status", "The masque-plus is Running.")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		callback("Failed to start masque-plus: " + err.Error())
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
				logger.Log(logger.ERROR, "masque-plus stderr: "+line)
			} else {
				logger.Log(logger.DEBUG, "masque-plus stdout: "+line)
				if strings.Contains(strings.ToLower(line), "started") {
					found = true
				}
			}
		}
		done <- struct{}{}
	}

	go scan(stdout, false)
	go scan(stderr, true)
	go func() { <-ctx.Done(); KillMasquePlus() }()

	<-done
	<-done
	cmd.Wait()

	procMu.Lock()
	currentProc = nil
	procMu.Unlock()
	return found
}

func KillMasquePlus() {
	procMu.Lock()
	defer procMu.Unlock()
	if currentProc != nil && currentProc.Process != nil {
		_ = currentProc.Process.Kill()
		logger.Log(logger.INFO, "Masque-plus process killed")
		currentProc = nil
	}
}
