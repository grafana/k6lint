package k6lint

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/grafana/k6foundry"
)

//nolint:forbidigo
func build(ctx context.Context, module string, dir string) (filename string, result error) {
	exe, err := os.CreateTemp("", "k6-*.exe")
	if err != nil {
		return "", err
	}

	if err = os.Chmod(exe.Name(), 0o700); err != nil { //nolint:gosec
		return "", err
	}

	var out bytes.Buffer

	defer func() {
		if result != nil {
			_, _ = io.Copy(os.Stderr, &out)
			fmt.Fprintln(os.Stderr)
		}
	}()

	builder, err := k6foundry.NewNativeBuilder(
		ctx,
		k6foundry.NativeBuilderOpts{
			Logger: slog.New(slog.NewTextHandler(&out, &slog.HandlerOptions{Level: slog.LevelError})),
			Stdout: &out,
			Stderr: &out,
			GoOpts: k6foundry.GoOpts{
				CopyGoEnv: true,
				Env:       map[string]string{"GOWORK": "off"},
			},
		},
	)
	if err != nil {
		result = err
		return "", result
	}

	_, result = builder.Build(
		ctx,
		k6foundry.NewPlatform(runtime.GOOS, runtime.GOARCH),
		"latest",
		[]k6foundry.Module{{Path: module, ReplacePath: dir}},
		nil,
		exe,
	)

	if result != nil {
		return "", result
	}

	if err = exe.Close(); err != nil {
		return "", err
	}

	return exe.Name(), nil
}

//nolint:forbidigo
func findFile(rex *regexp.Regexp, dirs ...string) (string, string, error) {
	for idx, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if idx == 0 {
				return "", "", err
			}

			continue
		}

		script := ""

		for _, entry := range entries {
			if rex.Match([]byte(entry.Name())) {
				script = entry.Name()

				break
			}
		}

		if len(script) > 0 {
			return filepath.Join(dir, script), script, nil
		}
	}

	return "", "", nil
}
