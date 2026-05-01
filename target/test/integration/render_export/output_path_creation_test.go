package render_export

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"theseus/target/internal/core"
)

func TestOutputPathCreation(t *testing.T) {
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "nested", "deep", "report.txt")
	svc := core.NewRenderService()

	_, _, err := svc.Execute(sampleModules(), core.RenderOptions{OutPath: outPath}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("output file missing: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("output file was empty")
	}
}
