package app

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	helpers "github.com/Freedom-Guard/freedom-core/pkg/utils"
)

func RunAsAdminStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if runtime.GOOS == "windows" {
		if !IsAdmin() {
			fmt.Fprintln(w, "⚠️ The application is not running with administrator privileges. Restarting as admin...")
			if err := RunAsAdmin(); err != nil {
				http.Error(w, "Failed to restart as administrator: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		fmt.Fprintln(w, "✅ The application is already running with administrator privileges.")
	} else {
		fmt.Fprintln(w, "This feature is only available on Windows.")
	}
}

func IsAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if runtime.GOOS != "windows" {
		fmt.Fprint(w, `{"admin": false, "message": "This check is only available on Windows."}`)
		return
	}

	if IsAdmin() {
		fmt.Fprint(w, `{"admin": true, "message": "The application is running with administrator privileges."}`)
	} else {
		fmt.Fprint(w, `{"admin": false, "message": "The application is not running with administrator privileges."}`)
	}
}

func IsAdmin() bool {
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

func RunAsAdmin() error {
	helpers.ShowInfo(
		"Running as Administrator",
		"The program is running with administrator privileges.",
	)

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	psCmd := fmt.Sprintf("Start-Process -FilePath '%s' -Verb RunAs -WindowStyle Normal", exe)
	cmd := exec.Command("powershell", "-Command", psCmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Start()
	if err != nil {
		return err
	}

	os.Exit(0)
	return nil
}
