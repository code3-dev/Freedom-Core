package updater

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Freedom-Guard/freedom-core/pkg/logger"
)

func DeleteCores() error {
	base, err := os.UserCacheDir()
	if err != nil {
		logger.Log(logger.ERROR, fmt.Sprintf("failed to get cache dir: %v", err))
		return fmt.Errorf("get cache dir: %w", err)
	}

	coreDir := filepath.Join(base, "freedom-core")
	logger.Log(logger.INFO, fmt.Sprintf("using core directory: %s", coreDir))

	names := []string{"xray", "singbox", "hiddify"}

	var errs []error
	for _, n := range names {
		dir := filepath.Join(coreDir, n)
		if err := os.RemoveAll(dir); err != nil {
			if os.IsNotExist(err) {
				logger.Log(logger.WARN, fmt.Sprintf("core folder not found: %s", dir))
			} else {
				logger.Log(logger.ERROR, fmt.Sprintf("failed to delete core folder %s: %v", n, err))
				errs = append(errs, err)
			}
		} else {
			logger.Log(logger.INFO, fmt.Sprintf("deleted core folder: %s", dir))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("some deletions failed: %v", errs)
	}

	logger.Log(logger.INFO, "all selected core folders deleted successfully")
	return nil
}
