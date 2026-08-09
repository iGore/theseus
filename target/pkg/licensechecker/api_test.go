package licensechecker_test

import (
	"os"
	"path/filepath"
	"testing"

	lc "github.com/davglass/license-checker/pkg/licensechecker"
)

func fixturePackages() *lc.Collection {
	c := lc.NewCollection()
	a := lc.NewObject()
	a.Set("licenses", "ISC")
	a.Set("repository", "https://github.com/isaacs/abbrev-js")
	c.Set("abbrev@1.0.9", a)
	f := lc.NewObject()
	f.Set("licenses", "MIT")
	f.Set("repository", "/path/to/foo")
	c.Set("foo", f)
	return c
}

func TestGoldenCSVAndMarkdownSnippets(t *testing.T) {
	c := fixturePackages()
	wantCSV := "\"module name\",\"license\",\"repository\"\n\"abbrev@1.0.9\",\"ISC\",\"https://github.com/isaacs/abbrev-js\"\n\"foo\",\"MIT\",\"/path/to/foo\""
	if got := lc.AsCSV(c, nil, ""); got != wantCSV {
		t.Fatalf("csv mismatch\nwant %q\ngot  %q", wantCSV, got)
	}
	wantMD := "[abbrev@1.0.9](https://github.com/isaacs/abbrev-js) - ISC\n[foo](/path/to/foo) - MIT"
	if got := lc.AsMarkDown(c, nil); got != wantMD {
		t.Fatalf("markdown mismatch\nwant %q\ngot  %q", wantMD, got)
	}
}

func TestParseJSONAndAsFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := filepath.Join(dir, "custom.json")
	if err := os.WriteFile(cfg, []byte(`{"name":true,"licenseText":"none"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	obj, err := lc.ParseJSON(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if obj.String("licenseText") != "none" {
		t.Fatalf("custom default not preserved")
	}
	lic := filepath.Join(dir, "LICENSE")
	if err := os.WriteFile(lic, []byte("MIT license bytes"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := lc.NewCollection()
	rec := lc.NewObject()
	rec.Set("licenses", "MIT")
	rec.Set("licenseFile", lic)
	c.Set("foo@1.0.0", rec)
	out := filepath.Join(dir, "out")
	if err := lc.AsFiles(c, out, os.Stderr); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(out, "foo-LICENSE.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "MIT license bytes" {
		t.Fatalf("copied bytes mismatch: %q", b)
	}
}
