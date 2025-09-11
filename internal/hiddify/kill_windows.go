//go:build windows
package hiddify

import (
	"context"
	"os/exec"
	"syscall"
)

func KillHiddify(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "taskkill", "/F", "/IM", "hiddify.exe", "/T")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
	return cmd.Run()
}
