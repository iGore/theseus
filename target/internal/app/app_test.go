package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func appGolden(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "testdata", "golden", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestCustomPathCLIGolden(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "package.json"), `{"version":"1.0.0"}`)
	pkg := filepath.Join(root, "node_modules", "abbrev")
	writeFile(t, filepath.Join(pkg, "package.json"), `{"name":"abbrev","version":"1.0.9","license":"ISC","homepage":"https://example.test/abbrev"}`)
	custom := filepath.Join(root, "custom.json")
	writeFile(t, custom, `{"name":"","homepage":""}`)

	var stdout, stderr strings.Builder
	code := Run(context.Background(), Invocation{
		Args:   []string{"--start", root, "--csv", "--customPath", custom},
		Stdout: &stdout,
		Stderr: &stderr,
		Getwd:  func() (string, error) { return root, nil },
	})
	if code != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	if !bytes.Equal([]byte(stdout.String()), appGolden(t, "custom_path_csv.golden")) {
		t.Fatalf("customPath stdout mismatch\ngot  %q\nwant %q", stdout.String(), string(appGolden(t, "custom_path_csv.golden")))
	}
	if stderr.String() != "" {
		t.Fatalf("unexpected stderr %q", stderr.String())
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
