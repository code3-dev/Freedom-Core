package hiddify

import (
	"bufio"
	"context"
	"os/exec"
	"strings"

	"github.com/Freedom-Guard/freedom-core/internal/installer"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func RunHiddifyStream(ctx context.Context, args []string, callback func(string)) bool {
	path, err := installer.PrepareCore()
	if err != nil {
		callback("Hiddify core not installed: " + err.Error())
		return false
	}

	cmd := exec.CommandContext(ctx, path, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		callback("Failed to start Hiddify: " + err.Error())
		return false
	}

	found := false
	done := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			callback(line)
			logger.Log(logger.INFO, "Hiddify stdout: "+line)
			if strings.Contains(line, "CORE STARTED") {
				found = true
			}
		}
		done <- struct{}{}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			callback(line)
			logger.Log(logger.ERROR, "Hiddify stderr: "+line)
		}
		done <- struct{}{}
	}()

	<-done
	<-done

	if err := cmd.Wait(); err != nil {
		callback("Hiddify process error: " + err.Error())
	}

	return found
}
