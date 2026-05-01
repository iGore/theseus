package cli

import (
	"testing"

	"theseus/target/internal/core"
)

func TestNormalizeOptions(t *testing.T) {
	tests := []struct {
		name       string
		raw        core.RawArgs
		wantStart  string
		wantDepth  int
		wantColor  bool
		wantErr    bool
		structured bool
	}{
		{
			name:      "defaults start to cwd when empty",
			raw:       core.RawArgs{},
			wantStart: "/tmp/work",
			wantDepth: core.FullTraversalDepth,
			wantColor: true,
		},
		{
			name:      "direct maps depth zero",
			raw:       core.RawArgs{Direct: true},
			wantStart: "/tmp/work",
			wantDepth: 0,
			wantColor: true,
		},
		{
			name:      "structured output disables color",
			raw:       core.RawArgs{JSON: true, Color: true, ColorSet: true},
			wantStart: "/tmp/work",
			wantDepth: core.FullTraversalDepth,
			wantColor: false,
		},
		{
			name:      "explicit start preserved",
			raw:       core.RawArgs{Start: "/repo"},
			wantStart: "/repo",
			wantDepth: core.FullTraversalDepth,
			wantColor: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizeOptions(tc.raw, func() (string, error) { return "/tmp/work", nil })
			if tc.wantErr && err == nil {
				t.Fatalf("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Start != tc.wantStart {
				t.Fatalf("start mismatch: got %q want %q", got.Start, tc.wantStart)
			}
			if got.Depth != tc.wantDepth {
				t.Fatalf("depth mismatch: got %d want %d", got.Depth, tc.wantDepth)
			}
			if got.Color != tc.wantColor {
				t.Fatalf("color mismatch: got %v want %v", got.Color, tc.wantColor)
			}
		})
	}
}
