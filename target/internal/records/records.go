package records

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"theseus-target/license-checker/internal/cli"
	"theseus-target/license-checker/internal/license"
	"theseus-target/license-checker/internal/npmgraph"
	"theseus-target/license-checker/internal/ordered"
)

type Builder struct{}

func Flatten(root *npmgraph.Node, opts cli.Options) *ordered.Map {
	data := map[string]*ordered.Record{}
	flatten(data, root, opts, opts.Start)
	out := ordered.NewMap()
	for _, k := range ordered.SortedKeys(data) {
		r := data[k]
		if r.Bool("private") {
			r.Set("licenses", "UNLICENSED")
		}
		if r.String("licenses") == "" {
			r.Set("licenses", "UNKNOWN")
		}
		if opts.Unknown {
			l := r.String("licenses")
			if l != "UNKNOWN" && strings.Contains(l, "*") {
				r.Set("licenses", "UNKNOWN")
			}
		}
		out.Set(k, r)
	}
	return out
}

func flatten(data map[string]*ordered.Record, n *npmgraph.Node, opts cli.Options, base string) {
	if n == nil {
		return
	}
	key := n.Name + "@" + n.Version
	if _, ok := data[key]; ok {
		return
	}
	if opts.Production && n.Extraneous {
		return
	}
	if opts.Development && !n.Extraneous && !n.Root {
		return
	}
	r := ordered.NewRecord()
	r.Set("licenses", "UNKNOWN")
	if n.Private {
		r.Set("private", true)
	}
	if n.Name != "" && n.Version != "" {
		data[key] = r
	}
	if n.Repository.URL != "" {
		r.Set("repository", normalizeRepo(n.Repository.URL))
	}
	if n.URL.Web != "" {
		r.Set("url", n.URL.Web)
	}
	if n.Author.Name != "" {
		r.Set("publisher", n.Author.Name)
	}
	if n.Author.Email != "" {
		r.Set("email", n.Author.Email)
	}
	if n.Author.URL != "" {
		r.Set("url", n.Author.URL)
	}
	if opts.Unknown {
		r.Set("dependencyPath", n.Path)
	}
	for _, f := range opts.CustomFormat {
		if f.Value == false {
			continue
		}
		if v, ok := n.Fields[f.Key]; ok {
			if s, ok := v.(string); ok && s != "" {
				r.Set(f.Key, s)
				continue
			}
			if ok {
				continue
			}
		}
		r.Set(f.Key, f.Value)
	}
	if n.Path != "" {
		r.Set("path", n.Path)
	}
	applyLicense(r, n, opts, base)
	if n.Name == "" || n.Version == "" {
		delete(data, key)
	}
	for _, c := range n.Dependencies {
		flatten(data, c, opts, base)
	}
}

func normalizeRepo(s string) string {
	s = strings.Replace(s, "git+ssh://git@", "git://", 1)
	s = strings.Replace(s, "git+https://github.com", "https://github.com", 1)
	s = strings.Replace(s, "git://github.com", "https://github.com", 1)
	s = strings.Replace(s, "git@github.com:", "https://github.com/", 1)
	return strings.TrimSuffix(s, ".git")
}

func applyLicense(r *ordered.Record, n *npmgraph.Node, opts cli.Options, base string) {
	set := func(v any) bool {
		switch x := v.(type) {
		case string:
			if out, ok := license.Classify(x, true); ok {
				r.Set("licenses", out)
				return true
			}
		case map[string]any:
			if t, _ := x["type"].(string); t != "" {
				if out, ok := license.Classify(t, true); ok {
					r.Set("licenses", out)
					return true
				}
			}
		case []any:
			var vals []string
			for _, e := range x {
				if s, ok := e.(string); ok {
					if out, ok := license.Classify(s, true); ok {
						vals = append(vals, out)
					}
				}
			}
			if len(vals) > 0 {
				r.Set("licenses", strings.Join(vals, ","))
				return true
			}
		}
		return false
	}
	if n.License != nil {
		set(n.License)
	} else if n.Licenses != nil {
		set(n.Licenses)
	} else if n.Readme != "" {
		set(n.Readme)
	}
	entries, _ := os.ReadDir(n.Path)
	names := []string{}
	for _, e := range entries {
		names = append(names, e.Name())
	}
	files := license.DetectFiles(names)
	if len(files) > 0 {
		lf := filepath.Join(n.Path, files[0])
		if opts.RelativeLicensePath {
			if rel, err := filepath.Rel(base, lf); err == nil {
				lf = rel
			}
		}
		r.Set("licenseFile", lf)
		if b, err := os.ReadFile(filepath.Join(n.Path, files[0])); err == nil {
			txt := string(b)
			cur := r.String("licenses")
			if cur == "" || strings.Contains(cur, "UNKNOWN") || strings.HasPrefix(cur, "Custom:") {
				if out, ok := license.Classify(txt, true); ok {
					r.Set("licenses", out)
				}
			}
			r.Set("licenseText", strings.TrimSpace(txt))
		}
	}
	for _, nm := range names {
		if strings.EqualFold(strings.TrimSuffix(nm, filepath.Ext(nm)), "NOTICE") {
			nf := filepath.Join(n.Path, nm)
			if opts.RelativeLicensePath {
				if rel, err := filepath.Rel(base, nf); err == nil {
					nf = rel
				}
			}
			r.Set("noticeFile", nf)
			break
		}
	}
	_ = fmt.Sprintf
}
