package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"theseus/target/internal/core"
)

type CommandConfig struct {
	Version         string
	VersionExitMode core.VersionExitMode // NEEDS CLARIFICATION (FR-008)
	Out             io.Writer
	Err             io.Writer
}

type ExitCodeError struct {
	Code int
	Msg  string
}

func (e ExitCodeError) Error() string {
	if e.Msg == "" {
		return fmt.Sprintf("exit code %d", e.Code)
	}
	return e.Msg
}

func Execute(args []string, cfg CommandConfig, getwd func() (string, error)) (core.GuardrailResult, error) {
	if getwd == nil {
		getwd = os.Getwd
	}
	if cfg.Out == nil {
		cfg.Out = io.Discard
	}
	if cfg.Err == nil {
		cfg.Err = io.Discard
	}

	var raw core.RawArgs
	result := core.GuardrailResult{}

	cmd := &cobra.Command{
		Use:           "license-checker",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if raw.Version {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", cfg.Version)
				result.MetaCommand = "version"
				if cfg.VersionExitMode == core.VersionExitLegacyNonZero {
					result.ExitCode = 1
					return ExitCodeError{Code: 1, Msg: "legacy version exit behavior"}
				}
				return nil
			}

			normalized, err := NormalizeOptions(raw, getwd)
			if err != nil {
				return err
			}
			result.Options = &normalized

			if strings.TrimSpace(normalized.FailOn) != "" && strings.TrimSpace(normalized.OnlyAllow) != "" {
				result.ExitCode = 1
				return errors.New("cannot use --failOn and --onlyAllow together")
			}

			if strings.Contains(normalized.FailOn, ",") {
				msg := "warning: --failOn expects semicolon-delimited values"
				result.Warnings = append(result.Warnings, msg)
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), msg)
			}
			if strings.Contains(normalized.OnlyAllow, ",") {
				msg := "warning: --onlyAllow expects semicolon-delimited values"
				result.Warnings = append(result.Warnings, msg)
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), msg)
			}

			return nil
		},
	}

	cmd.SetOut(cfg.Out)
	cmd.SetErr(cfg.Err)
	cmd.SetArgs(args)

	cmd.Flags().StringVar(&raw.Start, "start", "", "scan start path")
	cmd.Flags().BoolVar(&raw.Direct, "direct", false, "direct dependencies only")
	cmd.Flags().StringVar(&raw.CustomPath, "customPath", "", "path to custom format JSON file")
	cmd.Flags().BoolVar(&raw.JSON, "json", false, "json output")
	cmd.Flags().BoolVar(&raw.CSV, "csv", false, "csv output")
	cmd.Flags().BoolVar(&raw.Markdown, "markdown", false, "markdown output")
	cmd.Flags().BoolVar(&raw.Color, "color", true, "color output")
	cmd.Flags().Lookup("color").NoOptDefVal = "true"
	cmd.Flags().StringVar(&raw.FailOn, "failOn", "", "fail-on policy")
	cmd.Flags().StringVar(&raw.OnlyAllow, "onlyAllow", "", "allow-only policy")
	cmd.Flags().BoolVar(&raw.Version, "version", false, "print version")

	_ = cmd.Flags().Lookup("color").Changed
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		raw.ColorSet = cmd.Flags().Changed("color")
	}

	err := cmd.Execute()
	if err != nil {
		var exitErr ExitCodeError
		if errors.As(err, &exitErr) {
			return result, exitErr
		}
		if result.ExitCode == 0 {
			result.ExitCode = 1
		}
		return result, err
	}
	return result, nil
}
