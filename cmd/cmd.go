// Package cmd contains lint cobra command factory function.
package cmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/grafana/k6lint"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
)

//go:embed help.md
var help string

type options struct {
	out       string
	compact   bool
	quiet     bool
	json      bool
	passing   k6lint.Grade
	passedStr []string
	passed    []k6lint.Checker
	official  bool
}

// New creates new cobra command for exec command.
func New() (*cobra.Command, error) {
	opts := new(options)
	opts.passing = k6lint.GradeC

	root := &cobra.Command{
		Use:               "k6lint [flags] [directory]",
		Short:             "Linter for k6 extensions",
		Long:              help,
		SilenceUsage:      true,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
		Args:              cobra.MaximumNArgs(1),
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PreRunE: func(_ *cobra.Command, _ []string) error {
			opts.json = opts.json || opts.compact

			for _, str := range opts.passedStr {
				checker, err := k6lint.ParseChecker(str)
				if err != nil {
					return err
				}

				opts.passed = append(opts.passed, checker)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), args, opts)
		},
	}

	root.SetContext(context.TODO())

	flags := root.Flags()

	flags.SortFlags = false

	flags.Var(&opts.passing, "passing", "set lowest passing grade")
	flags.BoolVar(&opts.official, "official", false, "enable extra checks for official extensions")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "no output, only validation")
	flags.StringVarP(&opts.out, "out", "o", "", "write output to file instead of stdout")
	flags.BoolVar(&opts.json, "json", false, "generate JSON output")
	flags.BoolVarP(&opts.compact, "compact", "c", false, "compact instead of pretty-printed JSON output")
	flags.StringSliceVar(&opts.passedStr, "passed", []string{}, "set checker(s) to passed")

	flags.MarkHidden("passed") //nolint:errcheck,gosec
	root.MarkFlagsMutuallyExclusive("compact", "quiet")
	root.MarkFlagsMutuallyExclusive("json", "quiet")

	flags.BoolP("version", "V", false, "print version")

	return root, nil
}

func run(ctx context.Context, args []string, opts *options) (result error) {
	var dir string

	if len(args) == 0 {
		cwd, err := os.Getwd() //nolint:forbidigo
		if err != nil {
			return err
		}

		dir = cwd
	} else {
		dir = args[0]
	}

	info, err := os.Stat(dir) //nolint:forbidigo
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: %s", errNotDirectory, dir)
	}

	output := colorable.NewColorableStdout()

	if len(opts.out) > 0 {
		file, err := os.Create(opts.out) //nolint:forbidigo
		if err != nil {
			return err
		}

		defer func() {
			err := file.Close()
			if result == nil && err != nil {
				result = err
			}
		}()

		output = file
	}

	compliance, err := k6lint.Lint(ctx, dir, &k6lint.Options{Passed: opts.passed, Official: opts.official})
	if err != nil {
		return err
	}

	if isGitHubAction() {
		emitOutput(compliance)
	}

	if opts.quiet {
		return nil
	}

	if opts.json {
		err = jsonOutput(compliance, output, opts.compact)
		if err != nil {
			return err
		}

		return result
	}

	textOutput(compliance, output)

	if string(compliance.Grade) > string(opts.passing) {
		return fmt.Errorf("%w: %s", errFailingGrade, compliance.Grade)
	}

	return nil
}

func jsonOutput(compliance *k6lint.Compliance, output io.Writer, compact bool) error {
	encoder := json.NewEncoder(output)

	if !compact {
		encoder.SetIndent("", "  ")
	}

	return encoder.Encode(compliance)
}

func textOutput(compliance *k6lint.Compliance, output io.Writer) {
	heading := color.New(color.FgHiWhite, color.Bold).FprintfFunc()

	var c color.Attribute

	switch compliance.Grade {
	case k6lint.GradeA:
		c = color.BgHiGreen
	case k6lint.GradeB:
		c = color.BgGreen
	case k6lint.GradeC:
		c = color.BgHiCyan
	case k6lint.GradeD:
		c = color.BgYellow
	case k6lint.GradeE:
		c = color.BgMagenta
	case k6lint.GradeF:
		c = color.BgRed
	default:
		c = color.BgBlue
	}

	value := color.New(c).FprintfFunc()
	label := color.New(color.FgBlack, color.BgWhite).FprintfFunc()
	details := color.New(color.Italic).FprintfFunc()
	failed := color.New(color.FgRed).FprintfFunc()
	passed := color.New(color.FgGreen).FprintfFunc()
	plain := color.New(color.FgWhite).FprintfFunc()

	heading(output, "k6 extension compliance\n──────────┬─────────────\n")
	label(output, " grade ")
	value(output, " %s ", string(compliance.Grade))
	heading(output, "│")
	label(output, " level ")
	value(output, " %3d%% ", compliance.Level)
	plain(output, "\n\n")

	heading(output, "Details\n───────\n")

	for _, check := range compliance.Checks {
		fprintf := failed
		symbol := "✗"

		if check.Passed {
			fprintf = passed
			symbol = "✔"
		}

		fprintf(output, "%s %-20s\n", symbol, check.ID)

		if len(check.Details) != 0 {
			details(output, "  %s\n", check.Details)
		}
	}

	plain(output, "\n")
}

var (
	errNotDirectory = errors.New("not a directory")
	errFailingGrade = errors.New("failing grade")
)
