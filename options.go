package k6lint

// Options contains settings that modify the linter's operation.
type Options struct {
	// Passed contains a list of checkers that have already been marked as successful.
	Passed []Checker
}

func passedChecks(checkers []Checker) map[Checker]Check {
	dict := make(map[Checker]Check, len(checkers))

	for _, checker := range checkers {
		dict[checker] = Check{ID: checker, Passed: true, Details: "marked as passed because it is a requirement"}
	}

	return dict
}
