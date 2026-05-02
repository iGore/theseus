package integration

import (
	"testing"

	"theseus/target/internal/core"
	"theseus/target/internal/license"
)

func TestLicenseFallback_PrecedenceAndDeterminism(t *testing.T) {
	input := core.ModuleLicenseInput{
		CandidateFiles: []core.LicenseCandidateFile{
			{Name: "README", Content: "MIT license text"},
			{Name: "LICENCE", Content: "Apache License"},
			{Name: "COPYING", Content: "GPL-3.0"},
			{Name: "LICENSE", Content: "BSD-3-Clause"},
		},
	}

	firstValue, firstFile := (license.Detector{}).Resolve(input)
	for i := 0; i < 5; i++ {
		got, file := (license.Detector{}).Resolve(input)
		if got != firstValue || file != firstFile {
			t.Fatalf("non-deterministic run %d: got=(%q,%q) want=(%q,%q)", i, got, file, firstValue, firstFile)
		}
	}

	if firstFile != "LICENSE" {
		t.Fatalf("expected LICENSE precedence, got %q", firstFile)
	}
	if string(firstValue) != "BSD-3-Clause" {
		t.Fatalf("expected BSD-3-Clause from LICENSE content, got %q", firstValue)
	}
}
