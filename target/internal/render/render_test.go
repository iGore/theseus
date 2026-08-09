package render

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"theseus-target/license-checker/internal/cli"
	"theseus-target/license-checker/internal/npmgraph"
	"theseus-target/license-checker/internal/ordered"
	"theseus-target/license-checker/internal/policy"
	"theseus-target/license-checker/internal/records"
)

func abbrevMap() *ordered.Map {
	m := ordered.NewMap()
	r := ordered.NewRecord()
	r.Set("licenses", "ISC")
	r.Set("repository", "https://github.com/isaacs/abbrev-js")
	r.Set("name", "abbrev")
	r.Set("description", "Like ruby's abbrev module, but in js")
	r.Set("pewpew", "<<Should Never be set>>")
	m.Set("abbrev@1.0.9", r)
	return m
}
func goldenBytes(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "testdata", "golden", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func golden(t *testing.T, name string) string { return string(goldenBytes(t, name)) }
func fields() []cli.Field                     { return []cli.Field{{Key: "name"}, {Key: "description"}, {Key: "pewpew"}} }

func TestGoldenRenderSnippets(t *testing.T) {
	tests := []struct{ name, got, file string }{
		{"csv default", CSV(abbrevMap(), nil, ""), "csv_default.golden"},
		{"csv custom", CSV(abbrevMap(), fields(), ""), "csv_custom.golden"},
		{"csv component", CSV(abbrevMap(), fields(), "main-module"), "csv_component.golden"},
		{"markdown default", Markdown(abbrevMap(), nil), "markdown_default.golden"},
		{"markdown custom", Markdown(abbrevMap(), fields()), "markdown_custom_row.golden"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !bytes.Equal([]byte(tt.got), goldenBytes(t, tt.file)) {
				t.Fatalf("%s mismatch\ngot  %q\nwant %q", tt.name, tt.got, golden(t, tt.file))
			}
		})
	}
}

func TestTreeGoldenSnippets(t *testing.T) {
	m := ordered.NewMap()
	for _, x := range []struct{ k, repo, lic string }{{"cli@0.4.3", "http://github.com/chriso/cli", "MIT"}, {"glob@3.1.14", "https://github.com/isaacs/node-glob", "UNKNOWN"}, {"yui-lint@0.1.1", "http://github.com/yui/yui-lint", "BSD"}} {
		r := ordered.NewRecord()
		if x.k == "yui-lint@0.1.1" {
			r.Set("licenses", x.lic)
			r.Set("repository", x.repo)
		} else {
			r.Set("repository", x.repo)
			r.Set("licenses", x.lic)
		}
		m.Set(x.k, r)
	}
	if got := Tree(m); !bytes.Equal([]byte(got), goldenBytes(t, "tree_excerpt.golden")) {
		t.Fatalf("tree mismatch\n%s", got)
	}

	guessed := ordered.NewMap()
	r := ordered.NewRecord()
	r.Set("repository", "https://github.com/visionmedia/debug")
	r.Set("licenses", "MIT*")
	guessed.Set("debug@2.0.0", r)
	if got := Tree(guessed); !bytes.Equal([]byte(got), goldenBytes(t, "guessed_tree.golden")) {
		t.Fatalf("guessed tree mismatch\ngot  %q\nwant %q", got, golden(t, "guessed_tree.golden"))
	}
}

func TestGoldenStructuredModesAndSideEffects(t *testing.T) {
	formatTests := []struct {
		name string
		opts cli.Options
		file string
	}{
		{"json", cli.Options{JSON: true}, "json_default.golden"},
		{"csv", cli.Options{CSV: true}, "csv_default.golden"},
		{"csv component prefix", cli.Options{CSV: true, CustomFormat: fields(), CSVComponentPrefix: "main-module"}, "csv_component.golden"},
		{"markdown", cli.Options{Markdown: true}, "markdown_cli.golden"},
	}
	for _, tt := range formatTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(abbrevMap(), tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal([]byte(got), goldenBytes(t, tt.file)) {
				t.Fatalf("format mismatch\ngot  %q\nwant %q", got, golden(t, tt.file))
			}
		})
	}

	t.Run("summary", func(t *testing.T) {
		m := ordered.NewMap()
		for _, x := range []struct{ k, lic string }{{"a@1.0.0", "MIT"}, {"b@1.0.0", "MIT"}, {"c@1.0.0", "ISC"}} {
			r := ordered.NewRecord()
			r.Set("licenses", x.lic)
			m.Set(x.k, r)
		}
		if got := Summary(m); !bytes.Equal([]byte(got), goldenBytes(t, "summary_counts.golden")) {
			t.Fatalf("summary mismatch\ngot  %q\nwant %q", got, golden(t, "summary_counts.golden"))
		}
	})

	t.Run("out writes raw bytes and no stdout", func(t *testing.T) {
		tmp := t.TempDir()
		out := filepath.Join(tmp, "licenses.csv")
		var stdout, stderr strings.Builder
		formatted := CSV(abbrevMap(), nil, "")
		if err := Emit(abbrevMap(), formatted, cli.Options{Out: out}, &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(b, goldenBytes(t, "out_csv_file.golden")) {
			t.Fatalf("out file mismatch\ngot  %q\nwant %q", string(b), golden(t, "out_csv_file.golden"))
		}
		if !bytes.Equal([]byte(stdout.String()), goldenBytes(t, "empty_stdout.golden")) {
			t.Fatalf("stdout mismatch %q", stdout.String())
		}
	})

	t.Run("stdout adds line ending to csv and markdown formatted bytes", func(t *testing.T) {
		for _, tt := range []struct {
			name      string
			formatted string
			file      string
		}{
			{"csv stdout", CSV(abbrevMap(), nil, ""), "csv_stdout.golden"},
			{"markdown stdout", Markdown(abbrevMap(), nil) + "\n", "markdown_stdout.golden"},
		} {
			t.Run(tt.name, func(t *testing.T) {
				var stdout, stderr strings.Builder
				if err := Emit(abbrevMap(), tt.formatted, cli.Options{}, &stdout, &stderr); err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal([]byte(stdout.String()), goldenBytes(t, tt.file)) {
					t.Fatalf("stdout mismatch\ngot  %q\nwant %q", stdout.String(), golden(t, tt.file))
				}
			})
		}
	})
}

func TestGoldenConformanceRowsFromSourceContracts(t *testing.T) {
	t.Run("production prunes extraneous", func(t *testing.T) {
		m := records.Flatten(graphFixture(), cli.Options{Production: true, Start: "/fixture/root"})
		if got := Tree(m); !bytes.Equal([]byte(got), goldenBytes(t, "production_tree.golden")) {
			t.Fatalf("production mismatch\ngot  %q\nwant %q", got, golden(t, "production_tree.golden"))
		}
	})
	t.Run("development keeps extraneous dependency", func(t *testing.T) {
		m := records.Flatten(graphFixture(), cli.Options{Development: true, Start: "/fixture/root"})
		if got := Tree(m); !bytes.Equal([]byte(got), goldenBytes(t, "development_tree.golden")) {
			t.Fatalf("development mismatch\ngot  %q\nwant %q", got, golden(t, "development_tree.golden"))
		}
	})
	t.Run("start path propagates into output", func(t *testing.T) {
		m := records.Flatten(&npmgraph.Node{Name: "abbrev", Version: "1.0.9", License: "ISC", Path: "/workspace/project/node_modules/abbrev"}, cli.Options{Start: "/workspace/project"})
		if got := Tree(m); !bytes.Equal([]byte(got), goldenBytes(t, "start_path_tree.golden")) {
			t.Fatalf("start mismatch\ngot  %q\nwant %q", got, golden(t, "start_path_tree.golden"))
		}
	})
	t.Run("unknown rewrites guessed and includes dependencyPath", func(t *testing.T) {
		m := records.Flatten(&npmgraph.Node{Name: "debug", Version: "2.0.0", License: "Permission is hereby granted, free of charge, to any", Path: "/workspace/project/node_modules/debug"}, cli.Options{Unknown: true, Start: "/workspace/project"})
		if got := Tree(m); !bytes.Equal([]byte(got), goldenBytes(t, "unknown_dependency_path_tree.golden")) {
			t.Fatalf("unknown mismatch\ngot  %q\nwant %q", got, golden(t, "unknown_dependency_path_tree.golden"))
		}
	})
	t.Run("direct depth zero provider contract output", func(t *testing.T) {
		m := records.Flatten(&npmgraph.Node{Name: "root", Version: "1.0.0", License: "MIT", Path: "/fixture/root"}, cli.Options{Direct: 0, Start: "/fixture/root"})
		if got := Tree(m); !bytes.Equal([]byte(got), goldenBytes(t, "direct_tree.golden")) {
			t.Fatalf("direct mismatch\ngot  %q\nwant %q", got, golden(t, "direct_tree.golden"))
		}
	})
	t.Run("policy filters render byte fixtures", func(t *testing.T) {
		tests := []struct {
			name string
			opts cli.Options
			file string
		}{
			{"onlyunknown", cli.Options{OnlyUnknown: true}, "onlyunknown_tree.golden"},
			{"exclude", cli.Options{Exclude: "MIT,ISC"}, "exclude_tree.golden"},
			{"packages", cli.Options{Packages: "abbrev@1.0.9"}, "packages_tree.golden"},
			{"excludePackages", cli.Options{ExcludePackages: "abbrev@1.0.9;custom@1.0.0;mystery@1.0.0;private@1.0.0"}, "exclude_packages_tree.golden"},
			{"excludePrivatePackages", cli.Options{ExcludePrivatePackages: true}, "exclude_private_empty_tree.golden"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				m, code := policy.Apply(policyFixture(), tt.opts, &strings.Builder{})
				if code != 0 {
					t.Fatalf("unexpected policy code %d", code)
				}
				if got := Tree(m); !bytes.Equal([]byte(got), goldenBytes(t, tt.file)) {
					t.Fatalf("%s mismatch\ngot  %q\nwant %q", tt.name, got, golden(t, tt.file))
				}
			})
		}
	})
	t.Run("onlyAllow stderr", func(t *testing.T) {
		var stderr strings.Builder
		_, code := policy.Apply(policyFixture(), cli.Options{OnlyAllow: "ISC"}, &stderr)
		if code != 1 || !bytes.Equal([]byte(stderr.String()), goldenBytes(t, "onlyallow_violation_stderr.golden")) {
			t.Fatalf("onlyAllow code=%d stderr=%q", code, stderr.String())
		}
	})
	t.Run("failOn stderr", func(t *testing.T) {
		var stderr strings.Builder
		_, code := policy.Apply(policyFixture(), cli.Options{FailOn: "MIT*"}, &stderr)
		if code != 1 || !bytes.Equal([]byte(stderr.String()), goldenBytes(t, "failon_violation_stderr.golden")) {
			t.Fatalf("failOn code=%d stderr=%q", code, stderr.String())
		}
	})
	t.Run("relativeLicensePath output", func(t *testing.T) {
		root := t.TempDir()
		pkg := filepath.Join(root, "node_modules", "pkg")
		if err := os.MkdirAll(pkg, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(pkg, "LICENSE"), []byte("Permission is hereby granted, free of charge, to any"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(pkg, "NOTICE"), []byte("notice"), 0644); err != nil {
			t.Fatal(err)
		}
		m := records.Flatten(&npmgraph.Node{Name: "pkg", Version: "1.0.0", Path: pkg}, cli.Options{Start: root, RelativeLicensePath: true})
		got := strings.ReplaceAll(Tree(m), root, "/fixture/root")
		if !bytes.Equal([]byte(got), goldenBytes(t, "relative_license_path_tree.golden")) {
			t.Fatalf("relative path mismatch\ngot  %q\nwant %q", got, golden(t, "relative_license_path_tree.golden"))
		}
	})
}

func graphFixture() *npmgraph.Node {
	return &npmgraph.Node{Root: true, Path: "/fixture/root", Dependencies: []*npmgraph.Node{
		{Name: "prod", Version: "1.0.0", License: "MIT", Path: "/fixture/root/node_modules/prod"},
		{Name: "dev", Version: "1.0.0", License: "ISC", Extraneous: true, Path: "/fixture/root/node_modules/dev"},
	}}
}

func policyFixture() *ordered.Map {
	m := ordered.NewMap()
	for _, x := range []struct {
		k       string
		lic     string
		private bool
	}{
		{"abbrev@1.0.9", "ISC", false},
		{"debug@2.0.0", "MIT*", false},
		{"custom@1.0.0", "MIT*", false},
		{"mystery@1.0.0", "UNKNOWN", false},
		{"private@1.0.0", "UNLICENSED", true},
	} {
		r := ordered.NewRecord()
		r.Set("licenses", x.lic)
		if x.private {
			r.Set("private", true)
		}
		m.Set(x.k, r)
	}
	return m
}

func TestFilesSideEffectAndWarning(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "LICENSE")
	if err := os.WriteFile(src, []byte("MIT text"), 0644); err != nil {
		t.Fatal(err)
	}
	m := ordered.NewMap()
	r := ordered.NewRecord()
	r.Set("licenseFile", src)
	m.Set("foo", r)
	r2 := ordered.NewRecord()
	m.Set("bar", r2)
	var errb strings.Builder
	out := filepath.Join(tmp, "out")
	if err := Files(m, out, &errb); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(out, "foo-LICENSE.txt"))
	if string(b) != "MIT text" {
		t.Fatal("copy mismatch")
	}
	if !bytes.Equal([]byte(errb.String()), goldenBytes(t, "files_warning.golden")) {
		t.Fatalf("warn %q", errb.String())
	}
}
