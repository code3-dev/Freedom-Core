package hiddify

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/Freedom-Guard/freedom-core/internal/installer"
	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func RunHiddify(input string) bool {
	path, err := installer.PrepareCore()
	if err != nil {
		logger.Log(logger.ERROR, "Hiddify core not installed at: " + path)
	}
	logger.Log(logger.INFO, "Hiddify core installed at: "+path)

	cmd := exec.Command(path, input)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	found := false
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		logger.Log(logger.INFO, "Hiddify output: "+line)
		if strings.Contains(line, "CORE STARTED") {
			found = true
		}
	}
	errScanner := bufio.NewScanner(stderr)
	for errScanner.Scan() {
		line := errScanner.Text()
		logger.Log(logger.ERROR, "Hiddify error: "+line)
	}
	cmd.Wait()
	return found
}
