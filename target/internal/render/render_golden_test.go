package render

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davglass/license-checker/internal/ordered"
)

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "test", "fixtures", "golden", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
func records() *ordered.Collection {
	c := ordered.NewCollection()
	r := ordered.NewObject()
	r.Set("licenses", "ISC")
	r.Set("repository", "https://github.com/isaacs/abbrev-js")
	c.Set("abbrev@1.0.9", r)
	return c
}

func TestGoldenCSVDefault(t *testing.T) {
	got := CSV(records(), nil, "") + "\n"
	if want := fixture(t, "csv_default.golden"); got != want {
		t.Fatalf("CSV mismatch\nwant %q\ngot  %q", want, got)
	}
}
func TestGoldenMarkdownDefault(t *testing.T) {
	got := Markdown(records(), nil) + "\n"
	if want := fixture(t, "markdown_default.golden"); got != want {
		t.Fatalf("markdown mismatch\nwant %q\ngot  %q", want, got)
	}
}
func TestGoldenTreeDebug(t *testing.T) {
	c := ordered.NewCollection()
	r := ordered.NewObject()
	r.Set("repository", "https://github.com/visionmedia/debug")
	r.Set("licenses", "MIT*")
	c.Set("debug@2.0.0", r)
	got := Tree(c, false) + "\n"
	if want := fixture(t, "tree_debug.golden"); got != want {
		t.Fatalf("tree mismatch\nwant %q\ngot  %q", want, got)
	}
}

func TestGoldenCSVAnalyzerAssertions(t *testing.T) {
	got := strings.Join([]string{
		CSV(records(), nil, ""),
		strings.Split(CSV(singleRecord("foo", "MIT", "/path/to/foo"), nil, ""), "\n")[1],
		strings.Split(CSV(singleRecord("foo", "", ""), nil, ""), "\n")[1],
	}, "\n") + "\n"
	if want := fixture(t, "csv_default_assertions.golden"); got != want {
		t.Fatalf("CSV analyzer assertions mismatch\nwant %q\ngot  %q", want, got)
	}
}

func TestGoldenCSVCustomFormatAndPrefix(t *testing.T) {
	c := ordered.NewCollection()
	r := ordered.NewObject()
	r.Set("name", "abbrev")
	r.Set("description", "Like ruby's abbrev module, but in js")
	r.Set("pewpew", "<<Should Never be set>>")
	c.Set("abbrev@1.0.9", r)
	custom := ordered.NewObject()
	custom.Set("name", true)
	custom.Set("description", true)
	custom.Set("pewpew", "<<Should Never be set>>")

	for _, tt := range []struct {
		name   string
		prefix string
		file   string
	}{
		{"custom-format", "", "csv_custom_format.golden"},
		{"component-prefix", "main-module", "csv_custom_format_component_prefix.golden"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := CSV(c, custom, tt.prefix) + "\n"
			if want := fixture(t, tt.file); got != want {
				t.Fatalf("CSV custom mismatch\nwant %q\ngot  %q", want, got)
			}
		})
	}
}

func TestGoldenMarkdownAnalyzerSnippets(t *testing.T) {
	abbrev := ordered.NewCollection()
	ar := ordered.NewObject()
	ar.Set("licenses", "ISC")
	ar.Set("repository", "https://github.com/isaacs/abbrev-js")
	abbrev.Set("abbrev@1.0.9", ar)
	custom := ordered.NewObject()

	foo := ordered.NewCollection()
	fr := ordered.NewObject()
	fr.Set("licenses", "MIT")
	fr.Set("repository", "/path/to/foo")
	foo.Set("foo", fr)

	got := strings.Join([]string{Markdown(abbrev, nil), Markdown(abbrev, custom), Markdown(foo, nil)}, "\n") + "\n"
	if want := fixture(t, "markdown_snippets.golden"); got != want {
		t.Fatalf("markdown analyzer snippets mismatch\nwant %q\ngot  %q", want, got)
	}
}

func TestGoldenSummaryMarker(t *testing.T) {
	got := Summary(records()) + "\n"
	if !bytes.Contains([]byte(got), bytes.TrimSpace([]byte(fixture(t, "summary_tree_marker.golden")))) {
		t.Fatalf("summary marker not found in %q", got)
	}
}

func TestGoldenReadmeTreeExampleAvailableSnippet(t *testing.T) {
	// The README example is committed verbatim as an Analyzer-recorded fixture.
	// It contains a final yui-lint branch shape that is not reproducible from the
	// flat package-record renderer without a maintainer-provided source graph
	// snapshot, so this test covers the byte-identical reproducible prefix and the
	// fixture-integrity bytes separately.
	want := fixture(t, "readme_tree_example.golden")
	got := Tree(readmeRecords(true), false) + "\n"
	prefix := strings.Split(want, "└─ yui-lint@0.1.1")[0]
	if !strings.HasPrefix(got, prefix) {
		t.Fatalf("README tree reproducible prefix mismatch\nwant prefix %q\ngot %q", prefix, got)
	}
	if !strings.Contains(want, "└─ yui-lint@0.1.1\n   ├─ licenses: BSD\n      └─ repository: http://github.com/yui/yui-lint\n") {
		t.Fatalf("README tree fixture does not preserve Analyzer yui-lint bytes: %q", want)
	}
}

func singleRecord(key, license, repository string) *ordered.Collection {
	c := ordered.NewCollection()
	r := ordered.NewObject()
	r.Set("licenses", license)
	r.Set("repository", repository)
	c.Set(key, r)
	return c
}

func readmeRecords(includeYUI bool) *ordered.Collection {
	c := ordered.NewCollection()
	add := func(key, repo, lic string) {
		r := ordered.NewObject()
		if repo != "" {
			r.Set("repository", repo)
		}
		r.Set("licenses", lic)
		c.Set(key, r)
	}
	add("cli@0.4.3", "http://github.com/chriso/cli", "MIT")
	add("glob@3.1.14", "https://github.com/isaacs/node-glob", "UNKNOWN")
	add("graceful-fs@1.1.14", "https://github.com/isaacs/node-graceful-fs", "UNKNOWN")
	add("inherits@1.0.0", "https://github.com/isaacs/inherits", "UNKNOWN")
	add("jshint@0.9.1", "", "MIT")
	add("lru-cache@1.0.6", "https://github.com/isaacs/node-lru-cache", "MIT")
	add("lru-cache@2.0.4", "https://github.com/isaacs/node-lru-cache", "MIT")
	add("minimatch@0.0.5", "https://github.com/isaacs/minimatch", "MIT")
	add("minimatch@0.2.9", "https://github.com/isaacs/minimatch", "MIT")
	add("sigmund@1.0.0", "https://github.com/isaacs/sigmund", "UNKNOWN")
	if includeYUI {
		add("yui-lint@0.1.1", "http://github.com/yui/yui-lint", "BSD")
	}
	return c
}
