package diagnostics

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"theseus-target/license-checker/internal/cli"
)

func golden(t *testing.T, name string) string {
	b, err := os.ReadFile(filepath.Join("..", "..", "testdata", "golden", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
func TestPreflightGolden(t *testing.T) {
	var b strings.Builder
	handled, code := Preflight(cli.Options{Version: true}, &b)
	if !handled || code != 1 || b.String() != golden(t, "version_stderr.golden") {
		t.Fatalf("version %d %q", code, b.String())
	}
	b.Reset()
	Preflight(cli.Options{FailOn: "MIT,ISC"}, &b)
	if b.String() != golden(t, "failon_warning.golden") {
		t.Fatalf("warn %q", b.String())
	}
	b.Reset()
	handled, code = Preflight(cli.Options{Help: true}, &b)
	if !handled || code != 0 || !bytes.Equal([]byte(b.String()), []byte(golden(t, "help_stderr.golden"))) {
		t.Fatalf("help %d %q", code, b.String())
	}
	b.Reset()
	handled, code = Preflight(cli.Options{FailOn: "MIT", OnlyAllow: "ISC"}, &b)
	if !handled || code != 1 || b.String() != MutualExclusion+"\n" {
		t.Fatalf("mutual %d %q", code, b.String())
	}
}
