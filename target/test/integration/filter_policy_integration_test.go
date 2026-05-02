package integration

import (
	"strings"
	"testing"

	"theseus/target/internal/core"
)

type exitRecorder struct {
	code   int
	called bool
}

func (e *exitRecorder) Exit(code int) {
	e.called = true
	e.code = code
}

func TestFilterBeforePolicyAndExitCode(t *testing.T) {
	modules := map[string]core.ModuleEntry{
		"a@1.0.0": {Licenses: "MIT*"},
		"b@1.0.0": {Licenses: "Apache-2.0"},
	}

	rec := &exitRecorder{}
	svc := core.Service{ExitStrategy: rec}

	_, err := svc.ApplyFilteringAndPolicy(modules, core.Options{
		Unknown: true,
		FailOn:  "UNKNOWN",
	})
	if err == nil {
		t.Fatal("expected policy violation")
	}
	if !strings.Contains(err.Error(), "failOn") {
		t.Fatalf("expected failOn in error, got %q", err.Error())
	}
	if !rec.called || rec.code != 1 {
		t.Fatalf("expected exit strategy called with code 1, called=%v code=%d", rec.called, rec.code)
	}
}
