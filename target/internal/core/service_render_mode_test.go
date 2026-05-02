package core

import "testing"

func TestResolveRenderModePrecedence(t *testing.T) {
	svc := NewRenderService()
	tests := []struct {
		name string
		opts RenderOptions
		want RenderMode
	}{
		{name: "default tree", opts: RenderOptions{}, want: RenderModeTree},
		{name: "summary only", opts: RenderOptions{Summary: true}, want: RenderModeSummary},
		{name: "markdown over summary", opts: RenderOptions{Markdown: true, Summary: true}, want: RenderModeMarkdown},
		{name: "csv over markdown", opts: RenderOptions{CSV: true, Markdown: true, Summary: true}, want: RenderModeCSV},
		{name: "json highest priority", opts: RenderOptions{JSON: true, CSV: true, Markdown: true, Summary: true}, want: RenderModeJSON},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := svc.ResolveRenderMode(tc.opts)
			if got != tc.want {
				t.Fatalf("ResolveRenderMode() = %q, want %q", got, tc.want)
			}
		})
	}
}
