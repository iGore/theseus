package render_export

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"theseus/target/internal/core"
)

func TestOutputContentParity(t *testing.T) {
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "artifact", "report.md")
	svc := core.NewRenderService()

	got, _, err := svc.Execute(sampleModules(), core.RenderOptions{Markdown: true, OutPath: outPath}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if got != string(b) {
		t.Fatalf("file content mismatch with rendered output")
	}
}
