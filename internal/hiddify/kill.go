package hiddify

import (
	"context"
	"os/exec"
	"runtime"
	"syscall"
)

func KillHiddify(ctx context.Context) {
	if runtime.GOOS == "windows" {
		_ = killHiddifyWindows(ctx)
	} else {
		_ = killHiddifyUnix(ctx)
	}
}

func killHiddifyUnix(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "pkill", "-f", "hiddify")
	setpgid(cmd)
	return cmd.Run()
}

func killHiddifyWindows(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "taskkill", "/F", "/IM", "hiddify.exe", "/T")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
	return cmd.Run()
}
