package k6lint

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

//nolint:forbidigo
func checkerExamples(_ context.Context, dir string) *checkResult {
	dir = filepath.Join(dir, "examples")

	info, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return checkFailed("missing `examples` directory")
		}

		return checkError(err)
	}

	if !info.IsDir() {
		return checkFailed("`examples` is not a directory")
	}

	hasRegular := false

	err = filepath.WalkDir(dir, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.Type().IsRegular() {
			hasRegular = true
		}

		return nil
	})
	if err != nil {
		return checkError(err)
	}

	if hasRegular {
		return checkPassed("found `examples` as examples directory")
	}

	return checkFailed("no examples found in the `examples` directory")
}
