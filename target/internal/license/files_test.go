package license

import (
	"testing"

	"theseus/target/internal/core"
)

func TestOrderedCandidates_Precedence(t *testing.T) {
	in := []core.LicenseCandidateFile{
		{Name: "README.md"},
		{Name: "COPYING"},
		{Name: "LICENCE.txt"},
		{Name: "LICENSE"},
	}
	got := orderedCandidates(in)
	want := []string{"LICENSE", "LICENCE.txt", "COPYING", "README.md"}
	for i := range want {
		if got[i].Name != want[i] {
			t.Fatalf("index %d got %q want %q", i, got[i].Name, want[i])
		}
	}
}
