package render_export

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"theseus/target/internal/core"
)

func sampleModules() core.ModuleResultMap {
	return core.ModuleResultMap{
		"beta@2.0.0":  {Name: "beta", Version: "2.0.0", Licenses: "Apache-2.0"},
		"alpha@1.0.0": {Name: "alpha", Version: "1.0.0", Licenses: "MIT"},
	}
}

func TestRenderModesGolden(t *testing.T) {
	svc := core.NewRenderService()

	tests := []struct {
		name   string
		opts   core.RenderOptions
		golden string
	}{
		{name: "tree", opts: core.RenderOptions{}, golden: "tree.golden"},
		{name: "summary", opts: core.RenderOptions{Summary: true}, golden: "summary.golden"},
		{name: "markdown", opts: core.RenderOptions{Markdown: true}, golden: "markdown.golden"},
		{name: "csv", opts: core.RenderOptions{CSV: true}, golden: "csv.golden"},
		{name: "json", opts: core.RenderOptions{JSON: true}, golden: "json.golden"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var stdout bytes.Buffer
			got, _, err := svc.Execute(sampleModules(), tc.opts, &stdout)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			wantBytes, err := os.ReadFile(filepath.Join("..", "..", "fixtures", "render", "f006_modes", tc.golden))
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			want := string(wantBytes)
			if got != want {
				t.Fatalf("render mismatch\nwant:\n%s\n\ngot:\n%s", want, got)
			}
		})
	}
}
