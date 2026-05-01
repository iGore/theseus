package license

import (
	"testing"

	"theseus/target/internal/core"
)

func TestDetector_MetadataShapes(t *testing.T) {
	tests := []struct {
		name  string
		input core.ModuleLicenseInput
		want  string
	}{
		{name: "license string", input: core.ModuleLicenseInput{License: "MIT"}, want: "MIT"},
		{name: "license object", input: core.ModuleLicenseInput{License: map[string]any{"type": "Apache-2.0"}}, want: "Apache-2.0"},
		{name: "licenses array", input: core.ModuleLicenseInput{Licenses: []any{map[string]any{"type": "BSD-3-Clause"}}}, want: "BSD-3-Clause"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := (Detector{}).Resolve(tt.input)
			if string(got) != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestDetector_FallbackPrecedence(t *testing.T) {
	input := core.ModuleLicenseInput{
		ReadmeText: "MIT license text",
		CandidateFiles: []core.LicenseCandidateFile{
			{Name: "LICENSE", Content: "Apache License"},
		},
	}
	got, file := (Detector{}).Resolve(input)
	if string(got) != "MIT*" {
		t.Fatalf("expected README inferred license, got %q", got)
	}
	if file != "README" {
		t.Fatalf("expected README source marker, got %q", file)
	}
}

func TestDetector_ReadmeNoMatchContinuesToFiles(t *testing.T) {
	input := core.ModuleLicenseInput{
		ReadmeText: "this readme has no recognizable license",
		CandidateFiles: []core.LicenseCandidateFile{
			{Name: "COPYING", Content: "GPL-3.0"},
		},
	}
	got, file := (Detector{}).Resolve(input)
	if string(got) != "GPL-3.0" {
		t.Fatalf("got %q want GPL-3.0", got)
	}
	if file != "COPYING" {
		t.Fatalf("got file %q want COPYING", file)
	}
}
