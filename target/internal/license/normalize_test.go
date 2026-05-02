package license

import "testing"

func TestNormalize_SPDXFirstAndHeuristicMarker(t *testing.T) {
	spdx, inferred := normalize("MIT")
	if spdx != "MIT" || inferred {
		t.Fatalf("expected direct SPDX MIT, got %q inferred=%v", spdx, inferred)
	}

	heuristic, inferred := normalize("mit license")
	if heuristic != "MIT*" || !inferred {
		t.Fatalf("expected inferred MIT*, got %q inferred=%v", heuristic, inferred)
	}
}

func TestNormalize_CustomReference(t *testing.T) {
	got, inferred := normalize("https://licenses.example/custom")
	if got != "Custom: https://licenses.example/custom" || inferred {
		t.Fatalf("unexpected custom normalization: %q inferred=%v", got, inferred)
	}
}
