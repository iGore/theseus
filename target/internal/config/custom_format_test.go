package config

import (
	"os"
	"path/filepath"
	"testing"

	"theseus/target/internal/core"
)

func TestParseJSON_Matrix(t *testing.T) {
	t.Run("non string", func(t *testing.T) {
		if _, ok := ParseJSON(42).(error); !ok {
			t.Fatalf("expected error")
		}
	})

	t.Run("missing file", func(t *testing.T) {
		if _, ok := ParseJSON(filepath.Join(t.TempDir(), "missing.json")).(error); !ok {
			t.Fatalf("expected error")
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		d := t.TempDir()
		p := filepath.Join(d, "bad.json")
		if err := os.WriteFile(p, []byte("{"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, ok := ParseJSON(p).(error); !ok {
			t.Fatalf("expected error")
		}
	})

	t.Run("valid json", func(t *testing.T) {
		d := t.TempDir()
		p := filepath.Join(d, "ok.json")
		if err := os.WriteFile(p, []byte(`{"name":"fallback","unknown":"x"}`), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, ok := ParseJSON(p).(CustomFormatConfig); !ok {
			t.Fatalf("expected config")
		}
	})
}

func TestApplyCustomFields_DefaultFallbackAndFalseExclusion(t *testing.T) {
	module := core.ModuleEntry{Name: "pkg-a", Version: "1.2.3", LicenseText: "Line1\n\"Line2\""}
	cfg := CustomFormatConfig{
		"name":        "default-name",
		"description": "fallback-desc",
		"skipMe":      false,
		"licenseText": "",
	}

	got := ApplyCustomFields(module, cfg, true)
	if got["name"] != "pkg-a" {
		t.Fatalf("name mismatch: %q", got["name"])
	}
	if got["description"] != "fallback-desc" {
		t.Fatalf("description fallback mismatch: %q", got["description"])
	}
	if _, exists := got["skipMe"]; exists {
		t.Fatalf("skipMe should be excluded")
	}
	if got["licenseText"] != "Line1\\n'Line2'" {
		t.Fatalf("licenseText normalization mismatch: %q", got["licenseText"])
	}
}

func TestLoadCustomFormat_FromCustomPath(t *testing.T) {
	result := LoadCustomFormat(nil, "../../test/fixtures/custom-format/valid_custom_path.json")
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Config) == 0 {
		t.Fatalf("expected config entries")
	}
	if _, ok := result.Config["description"]; !ok {
		t.Fatalf("expected description key in loaded config")
	}
}
