package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/davglass/license-checker/internal/cli"
)

func TestGoldenHelpAndVersion(t *testing.T) {
	root := filepath.Join("..", "..", "test", "golden")
	for _, tt := range []struct {
		name string
		argv []string
		file string
		code int
	}{{"help", []string{"--help"}, "help.stderr", 0}, {"version", []string{"--version"}, "version.stderr", 1}} {
		t.Run(tt.name, func(t *testing.T) {
			want, err := os.ReadFile(filepath.Join(root, tt.file))
			if err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			code := cli.New(&stdout, &stderr).Run(context.Background(), tt.argv)
			if code != tt.code {
				t.Fatalf("code want %d got %d", tt.code, code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("stdout not empty: %q", stdout.String())
			}
			if !bytes.Equal(stderr.Bytes(), want) {
				t.Fatalf("stderr mismatch\nwant %q\ngot  %q", want, stderr.Bytes())
			}
		})
	}
}
