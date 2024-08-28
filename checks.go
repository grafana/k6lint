package k6lint

import (
	"context"
	"fmt"
	"os"
)

type checkFunc func(ctx context.Context, dir string) *checkResult

type checkResult struct {
	passed  bool
	details string
}

func checkFailed(details string) *checkResult {
	return &checkResult{passed: false, details: fmt.Sprint(details)}
}

func checkPassed(details string, args ...any) *checkResult {
	return &checkResult{passed: true, details: fmt.Sprintf(details, args...)}
}

func checkError(err error) *checkResult {
	return &checkResult{passed: false, details: "error: " + err.Error()}
}

type checkDefinition struct {
	id    Checker
	fn    checkFunc
	score int
}

func checkDefinitions() []checkDefinition {
	modCheck := newModuleChecker()
	gitCheck := newGitChecker()

	defs := []checkDefinition{
		{id: CheckerModule, score: 2, fn: modCheck.hasGoModule},
		{id: CheckerReplace, score: 2, fn: modCheck.hasNoReplace},
		{id: CheckerReadme, score: 5, fn: checkerReadme},
		{id: CheckerExamples, score: 5, fn: checkerExamples},
		{id: CheckerLicense, score: 5, fn: checkerLicense},
		{id: CheckerGit, score: 1, fn: gitCheck.isWorkDir},
		{id: CheckerVersions, score: 5, fn: gitCheck.hasVersions},
	}

	return defs
}

func runChecks(ctx context.Context, dir string, opts *Options) ([]Check, int) {
	checkDefs := checkDefinitions()
	results := make([]Check, 0, len(checkDefs))
	passed := passedChecks(opts.Passed)

	var total, sum float64

	for _, checker := range checkDefs {
		total += float64(checker.score)

		var check Check

		if c, found := passed[checker.id]; found {
			check = c
		} else {
			res := checker.fn(ctx, dir)

			check.ID = checker.id
			check.Passed = res.passed
			check.Details = res.details
		}

		if check.Passed {
			sum += float64(checker.score)
		}

		results = append(results, check)
	}

	return results, int((sum / total) * 100.0)
}

// ParseChecker parses checker name from string.
func ParseChecker(val string) (Checker, error) {
	v := Checker(val)

	switch v {
	case
		CheckerModule,
		CheckerReplace,
		CheckerReadme,
		CheckerExamples,
		CheckerLicense,
		CheckerGit,
		CheckerVersions:

		return v, nil
	default:
		return "", fmt.Errorf("%w: %s", os.ErrInvalid, val) //nolint:forbidigo
	}
}
