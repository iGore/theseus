package integration

import (
	"testing"

	"theseus/target/internal/config"
	"theseus/target/internal/core"
)

func TestCustomFormat_FalsePropertyExcluded(t *testing.T) {
	module := core.ModuleEntry{Name: "alpha"}
	cfg := config.CustomFormatConfig{
		"name":   "",
		"private": false,
	}

	fields := config.ApplyCustomFields(module, cfg, false)
	if _, ok := fields["private"]; ok {
		t.Fatalf("private should be excluded when false")
	}
}
