package k6lint

import (
	"context"
	"os"
	"regexp"
)

var reREADME = regexp.MustCompile(
	`(?i)^readme\.(?:markdown|mdown|mkdn|md|textile|rdoc|org|creole|mediawiki|wiki|rst|asciidoc|adoc|asc|pod|txt)`,
)

func checkerReadme(_ context.Context, dir string) *checkResult {
	entries, err := os.ReadDir(dir) //nolint:forbidigo
	if err != nil {
		return checkFailed("")
	}

	for _, entry := range entries {
		if reREADME.Match([]byte(entry.Name())) {
			return checkPassed("found `%s` as README file", entry.Name())
		}
	}

	return checkFailed("no README file found")
}
