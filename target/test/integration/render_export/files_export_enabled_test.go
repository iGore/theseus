package render_export

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"theseus/target/internal/core"
)

func TestFilesExportEnabledAndCombinedWithOut(t *testing.T) {
	tmp := t.TempDir()
	licenseSrc := filepath.Join(tmp, "src", "ALPHA_LICENSE")
	if err := os.MkdirAll(filepath.Dir(licenseSrc), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(licenseSrc, []byte("MIT text"), 0o644); err != nil {
		t.Fatal(err)
	}

	modules := sampleModules()
	m := modules["alpha@1.0.0"]
	m.LicenseFile = licenseSrc
	modules["alpha@1.0.0"] = m

	outPath := filepath.Join(tmp, "out", "report.txt")
	exportPath := filepath.Join(tmp, "export")
	svc := core.NewRenderService()

	_, _, err := svc.Execute(modules, core.RenderOptions{OutPath: outPath, FilesPath: exportPath}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("expected --out artifact: %v", err)
	}
	exported := filepath.Join(exportPath, "alpha@1.0.0", "ALPHA_LICENSE")
	if _, err := os.Stat(exported); err != nil {
		t.Fatalf("expected exported license file: %v", err)
	}
}
