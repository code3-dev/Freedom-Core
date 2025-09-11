//go:build !windows

package hiddify

import (
	"context"
	"os/exec"
)

func KillHiddify(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "pkill", "-f", "hiddify")
	setpgid(cmd)
	return cmd.Run()
}
