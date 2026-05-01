package filter

import (
	"testing"
)

// Traceability:
// - FR-001/SC-001: unknown transform
// - FR-002/SC-002: onlyunknown restriction
// - FR-003: exclude filtering with escaped comma + BSD compatibility
// - FR-004/SC-004: package include/exclude key restriction
// - FR-005: private exclusion
func TestApplyFilteringStages(t *testing.T) {
	base := map[string]ModuleEntry{
		"a@1.0.0": {Licenses: "MIT"},
		"b@1.0.0": {Licenses: "BSD-3-Clause"},
		"c@1.0.0": {Licenses: "GPL"},
		"d@1.0.0": {Licenses: "Foo,Bar"},
		"e@1.0.0": {Licenses: "UNKNOWN*"},
		"f@1.0.0": {Licenses: "UNLICENSED", Private: true},
	}

	t.Run("unknown rewrites guessed marker", func(t *testing.T) {
		out := Apply(base, Options{Unknown: true})
		if got := out["e@1.0.0"].Licenses; got != "UNKNOWN" {
			t.Fatalf("expected UNKNOWN, got %q", got)
		}
	})

	t.Run("onlyunknown retains only unknown or guessed", func(t *testing.T) {
		out := Apply(base, Options{OnlyUnknown: true})
		if len(out) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(out))
		}
		if _, ok := out["e@1.0.0"]; !ok {
			t.Fatalf("expected e@1.0.0 retained")
		}
	})

	t.Run("exclude supports BSD and escaped commas", func(t *testing.T) {
		out := Apply(base, Options{Exclude: "BSD,Foo\\,Bar"})
		if _, ok := out["b@1.0.0"]; ok {
			t.Fatalf("expected BSD package removed")
		}
		if _, ok := out["d@1.0.0"]; ok {
			t.Fatalf("expected escaped-comma package removed")
		}
	})

	t.Run("packages include-list keeps only explicit keys", func(t *testing.T) {
		out := Apply(base, Options{Packages: "a@1.0.0;c@1.0.0"})
		if len(out) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(out))
		}
		if _, ok := out["a@1.0.0"]; !ok {
			t.Fatalf("expected a@1.0.0 present")
		}
		if _, ok := out["c@1.0.0"]; !ok {
			t.Fatalf("expected c@1.0.0 present")
		}
	})

	t.Run("excludePackages removes explicit keys", func(t *testing.T) {
		out := Apply(base, Options{ExcludePackages: "a@1.0.0;c@1.0.0"})
		if _, ok := out["a@1.0.0"]; ok {
			t.Fatalf("expected a@1.0.0 removed")
		}
		if _, ok := out["c@1.0.0"]; ok {
			t.Fatalf("expected c@1.0.0 removed")
		}
	})

	t.Run("excludePrivatePackages removes private modules", func(t *testing.T) {
		out := Apply(base, Options{ExcludePrivatePackages: true})
		if _, ok := out["f@1.0.0"]; ok {
			t.Fatalf("expected private package removed")
		}
	})
}
