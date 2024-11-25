package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/grafana/k6lint"
)

//nolint:forbidigo
func isGitHubAction() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

//nolint:forbidigo
func emitOutput(compliance *k6lint.Compliance) error {
	ghOutput := os.Getenv("GITHUB_OUTPUT")
	if len(ghOutput) == 0 {
		return nil
	}

	file, err := os.Create(filepath.Clean(ghOutput))
	if err != nil {
		return err
	}

	slog.Debug("Compliance", "grade", compliance.Grade, "level", compliance.Level)

	_, err = fmt.Fprintf(file, "grade=%s\n", compliance.Grade)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "level=%d\n", compliance.Level)
	if err != nil {
		return err
	}

	return file.Close()
}
