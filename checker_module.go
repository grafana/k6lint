package k6lint

import (
	"context"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type moduleChecker struct {
	file *modfile.File
}

func newModuleChecker() *moduleChecker {
	return new(moduleChecker)
}

func (mc *moduleChecker) hasGoModule(_ context.Context, dir string) *checkResult {
	filename := filepath.Join(dir, "go.mod")

	data, err := os.ReadFile(filepath.Clean(filename)) //nolint:forbidigo
	if err != nil {
		return checkError(err)
	}

	mc.file, err = modfile.Parse(filename, data, nil)
	if err != nil {
		return checkError(err)
	}

	return checkPassed("found `%s` as go module", mc.file.Module.Mod.String())
}

func (mc *moduleChecker) hasNoReplace(ctx context.Context, dir string) *checkResult {
	if mc.file == nil {
		res := mc.hasGoModule(ctx, dir)
		if !res.passed {
			return res
		}
	}

	if len(mc.file.Replace) != 0 {
		return checkFailed("the `go.mod` file contains `replace` directive(s)")
	}

	return checkPassed("no `replace` directive in the `go.mod` file")
}
