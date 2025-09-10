//go:build windows

package hiddify

import "os/exec"

func setpgid(cmd *exec.Cmd) {}
