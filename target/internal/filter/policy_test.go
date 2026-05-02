package filter

import (
	"strings"
	"testing"
)

func TestEvaluatePoliciesFailFast(t *testing.T) {
	modules := map[string]ModuleEntry{
		"a@1.0.0": {Licenses: "MIT"},
	}

	err := EvaluatePolicies(modules, Options{FailOn: "MIT"})
	if err == nil {
		t.Fatal("expected failOn violation error")
	}
	if !strings.Contains(err.Error(), "failOn") {
		t.Fatalf("expected failOn violation text, got %q", err.Error())
	}
}

func TestEvaluatePoliciesOnlyAllowViolation(t *testing.T) {
	modules := map[string]ModuleEntry{
		"a@1.0.0": {Licenses: "GPL"},
	}

	err := EvaluatePolicies(modules, Options{OnlyAllow: "MIT"})
	if err == nil {
		t.Fatal("expected onlyAllow violation error")
	}
	if !strings.Contains(err.Error(), "onlyAllow") {
		t.Fatalf("expected onlyAllow violation text, got %q", err.Error())
	}
}

func TestValidatePolicyCombination(t *testing.T) {
	err := ValidatePolicyCombination(Options{FailOn: "MIT", OnlyAllow: "Apache-2.0"})
	if err == nil {
		t.Fatal("expected invalid combination error")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("unexpected error text: %q", err.Error())
	}
}
