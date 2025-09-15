//go:build windows

package run

import (
	"os/exec"
	"syscall"
)

func SetupCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
