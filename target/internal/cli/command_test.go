package cli

import (
	"bytes"
	"strings"
	"testing"

	"theseus/target/internal/core"
)

func TestExecute_GuardrailsAndWarnings(t *testing.T) {
	t.Run("conflicting failOn and onlyAllow returns non-zero", func(t *testing.T) {
		out := &bytes.Buffer{}
		errOut := &bytes.Buffer{}
		result, err := Execute(
			[]string{"--failOn", "MIT", "--onlyAllow", "Apache-2.0"},
			CommandConfig{Version: "v0.0.0", VersionExitMode: core.VersionExitSuccess, Out: out, Err: errOut},
			func() (string, error) { return "/cwd", nil },
		)
		if err == nil {
			t.Fatalf("expected error")
		}
		if result.ExitCode == 0 {
			t.Fatalf("expected non-zero exit code")
		}
	})

	t.Run("comma delimiter emits warning", func(t *testing.T) {
		out := &bytes.Buffer{}
		errOut := &bytes.Buffer{}
		result, err := Execute(
			[]string{"--failOn", "MIT,Apache-2.0"},
			CommandConfig{Version: "v0.0.0", VersionExitMode: core.VersionExitSuccess, Out: out, Err: errOut},
			func() (string, error) { return "/cwd", nil },
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Warnings) == 0 {
			t.Fatalf("expected warnings")
		}
		if !strings.Contains(errOut.String(), "semicolon-delimited") {
			t.Fatalf("expected delimiter warning output, got %q", errOut.String())
		}
	})
}

func TestExecute_HelpAndVersionDeterministic(t *testing.T) {
	for i := 0; i < 2; i++ {
		out := &bytes.Buffer{}
		errOut := &bytes.Buffer{}
		result, err := Execute(
			[]string{"--help"},
			CommandConfig{Version: "v1.2.3", VersionExitMode: core.VersionExitSuccess, Out: out, Err: errOut},
			func() (string, error) { return "/cwd", nil },
		)
		if err != nil {
			t.Fatalf("help run %d unexpected error: %v", i, err)
		}
		if result.ExitCode != 0 {
			t.Fatalf("help run %d expected exit 0, got %d", i, result.ExitCode)
		}
		if !strings.Contains(out.String(), "Usage:") {
			t.Fatalf("help run %d missing usage output", i)
		}
	}

	for i := 0; i < 2; i++ {
		out := &bytes.Buffer{}
		errOut := &bytes.Buffer{}
		result, err := Execute(
			[]string{"--version"},
			CommandConfig{Version: "v1.2.3", VersionExitMode: core.VersionExitSuccess, Out: out, Err: errOut},
			func() (string, error) { return "/cwd", nil },
		)
		if err != nil {
			t.Fatalf("version success run %d unexpected error: %v", i, err)
		}
		if result.ExitCode != 0 {
			t.Fatalf("version success run %d expected exit 0, got %d", i, result.ExitCode)
		}
		if strings.TrimSpace(out.String()) != "v1.2.3" {
			t.Fatalf("version success run %d unexpected output: %q", i, out.String())
		}
	}
}

func TestExecute_VersionLegacyCompatibilityBranch(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	result, err := Execute(
		[]string{"--version"},
		CommandConfig{Version: "v1.2.3", VersionExitMode: core.VersionExitLegacyNonZero, Out: out, Err: errOut},
		func() (string, error) { return "/cwd", nil },
	)
	if err == nil {
		t.Fatalf("expected legacy non-zero exit error")
	}
	if result.ExitCode != 1 {
		t.Fatalf("expected exit 1 for legacy mode, got %d", result.ExitCode)
	}
	if strings.TrimSpace(out.String()) != "v1.2.3" {
		t.Fatalf("unexpected version output: %q", out.String())
	}
}
