package k6lint

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"golang.org/x/mod/modfile"
)

type moduleChecker struct {
	file *modfile.File
	exe  string
	js   bool
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

func (mc *moduleChecker) hasNoReplace(_ context.Context, _ string) *checkResult {
	if mc.file == nil {
		return checkFailed("missing go.mod")
	}

	if len(mc.file.Replace) != 0 {
		return checkFailed("the `go.mod` file contains `replace` directive(s)")
	}

	return checkPassed("no `replace` directive in the `go.mod` file")
}

func (mc *moduleChecker) canBuild(ctx context.Context, dir string) *checkResult {
	if mc.file == nil {
		return checkFailed("missing go.mod")
	}

	exe, err := build(ctx, mc.file.Module.Mod.Path, dir)
	if err != nil {
		return checkError(err)
	}

	out, err := exec.CommandContext(ctx, exe, "version").CombinedOutput() //nolint:gosec
	if err != nil {
		return checkError(err)
	}

	rex, err := regexp.Compile("(?i)  " + mc.file.Module.Mod.String() + "[^,]+, [^ ]+ \\[(?P<type>[a-z]+)\\]")
	if err != nil {
		return checkError(err)
	}

	subs := rex.FindAllSubmatch(out, -1)
	if subs == nil {
		return checkFailed(mc.file.Module.Mod.String() + " is not in the version command's output")
	}

	for _, one := range subs {
		if string(one[rex.SubexpIndex("type")]) == "js" {
			mc.js = true
			break
		}
	}

	mc.exe = exe

	return checkPassed("can be built with the latest k6 version")
}

var reSmoke = regexp.MustCompile(`(?i)^smoke(\.test)?\.(?:js|ts)`)

//nolint:forbidigo
func (mc *moduleChecker) smoke(ctx context.Context, dir string) *checkResult {
	if mc.exe == "" {
		return checkFailed("can't build")
	}

	filename, shortname, err := findFile(reSmoke,
		dir,
		filepath.Join(dir, "test"),
		filepath.Join(dir, "tests"),
		filepath.Join(dir, "examples"),
	)
	if err != nil {
		return checkError(err)
	}

	if len(shortname) == 0 {
		return checkFailed("no smoke test file found")
	}

	cmd := exec.CommandContext(ctx, mc.exe, "run", "--no-usage-report", "--no-summary", "--quiet", filename) //nolint:gosec

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))

		return checkError(err)
	}

	return checkPassed("`%s` successfully run with k6", shortname)
}

var reIndexDTS = regexp.MustCompile("^index.d.ts$")

func (mc *moduleChecker) types(_ context.Context, dir string) *checkResult {
	if mc.exe == "" {
		return checkFailed("can't build")
	}

	if !mc.js {
		return checkPassed("skipped due to output extension")
	}

	_, shortname, err := findFile(reIndexDTS,
		dir,
		filepath.Join(dir, "docs"),
	)
	if err != nil {
		return checkError(err)
	}

	if len(shortname) > 0 {
		return checkPassed("found `index.d.ts` file")
	}

	return checkFailed("no `index.d.ts` file found")
}
