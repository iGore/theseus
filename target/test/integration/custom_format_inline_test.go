package integration

import (
	"testing"

	"theseus/target/internal/config"
	"theseus/target/internal/core"
)

func TestCustomFormatInline_AddsConfiguredKeys(t *testing.T) {
	module := core.ModuleEntry{Name: "alpha", Version: "1.0.0"}
	cfg := config.CustomFormatConfig{
		"name":    "fallback",
		"missing": "default-missing",
	}

	fields := config.ApplyCustomFields(module, cfg, false)
	if fields["name"] != "alpha" {
		t.Fatalf("expected name from module, got %q", fields["name"])
	}
	if fields["missing"] != "default-missing" {
		t.Fatalf("expected default fallback, got %q", fields["missing"])
	}
}
