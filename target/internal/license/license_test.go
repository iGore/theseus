package license

import "testing"

func TestClassifySourceExamples(t *testing.T) {
	tests := []struct {
		in, want string
		ok       bool
	}{{"MIT", "MIT", true}, {"Apache-2.0 OR ISC OR MIT", "Apache-2.0 OR ISC OR MIT", true}, {"Permission is hereby granted, free of charge, to any person", "MIT*", true}, {"GNU GENERAL PUBLIC LICENSE Version 2", "GPL-2.0*", true}, {"SEE LICENSE IN LICENSE.md", "Custom: LICENSE.md", true}, {"gibberish", "", false}, {"", "Undefined", true}}
	for _, tt := range tests {
		got, ok := Classify(tt.in, true)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("%q got %q/%v", tt.in, got, ok)
		}
	}
}
func TestDetectFilesPrecedence(t *testing.T) {
	got := DetectFiles([]string{"README.txt", "COPYING", "LICENSE.md", "LICENSE-MIT", "LICENSE-APACHE", "licence.txt"})
	want := []string{"LICENSE.md", "LICENSE-MIT", "licence.txt", "COPYING", "README.txt"}
	if len(got) != len(want) {
		t.Fatalf("%v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%v", got)
		}
	}
}
