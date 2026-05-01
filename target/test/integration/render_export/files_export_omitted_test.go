package render_export

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"theseus/target/internal/core"
)

func TestFilesExportOmitted(t *testing.T) {
	tmp := t.TempDir()
	licenseSrc := filepath.Join(tmp, "src", "BETA_LICENSE")
	if err := os.MkdirAll(filepath.Dir(licenseSrc), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(licenseSrc, []byte("Apache text"), 0o644); err != nil {
		t.Fatal(err)
	}

	modules := sampleModules()
	m := modules["beta@2.0.0"]
	m.LicenseFile = licenseSrc
	modules["beta@2.0.0"] = m

	svc := core.NewRenderService()
	_, _, err := svc.Execute(modules, core.RenderOptions{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmp, "export")); !os.IsNotExist(err) {
		t.Fatalf("unexpected export directory side effect")
	}
}
