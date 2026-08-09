package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"theseus-target/license-checker/internal/cli"
	"theseus-target/license-checker/internal/ordered"
)

func Format(m *ordered.Map, opts cli.Options) (string, error) {
	if opts.JSON {
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(m); err != nil {
			return "", err
		}
		return b.String(), nil
	}
	if opts.CSV {
		return CSV(m, opts.CustomFormat, opts.CSVComponentPrefix), nil
	}
	if opts.Markdown {
		return Markdown(m, opts.CustomFormat) + "\n", nil
	}
	if opts.Summary {
		return Summary(m), nil
	}
	return Tree(m), nil
}
func Emit(m *ordered.Map, formatted string, opts cli.Options, stdout, stderr io.Writer) error {
	if opts.Files != "" {
		return Files(m, opts.Files, stderr)
	}
	if opts.Out != "" {
		if err := os.MkdirAll(filepath.Dir(opts.Out), 0755); err != nil {
			return err
		}
		return os.WriteFile(opts.Out, []byte(formatted), 0644)
	}
	_, err := fmt.Fprintln(stdout, formatted)
	return err
}
func q(s string) string { return strconv.Quote(s) }
func CSV(m *ordered.Map, fields []cli.Field, prefix string) string {
	lines := []string{}
	hdr := []string{}
	if prefix != "" {
		hdr = append(hdr, "component")
	}
	hdr = append(hdr, "module name")
	if len(fields) == 0 {
		hdr = append(hdr, "license", "repository")
	} else {
		for _, f := range fields {
			hdr = append(hdr, f.Key)
		}
	}
	lines = append(lines, quoteRow(hdr))
	for _, k := range m.Keys() {
		r, _ := m.Get(k)
		row := []string{}
		if prefix != "" {
			row = append(row, prefix)
		}
		row = append(row, k)
		if len(fields) == 0 {
			row = append(row, r.String("licenses"), r.String("repository"))
		} else {
			for _, f := range fields {
				row = append(row, fmt.Sprint(value(r, f.Key)))
			}
		}
		lines = append(lines, quoteRow(row))
	}
	return strings.Join(lines, "\n")
}
func value(r *ordered.Record, k string) any { v, _ := r.Get(k); return v }
func quoteRow(row []string) string {
	qd := make([]string, len(row))
	for i, s := range row {
		qd[i] = q(s)
	}
	return strings.Join(qd, ",")
}
func Markdown(m *ordered.Map, fields []cli.Field) string {
	var lines []string
	for _, k := range m.Keys() {
		r, _ := m.Get(k)
		repo := r.String("repository")
		if len(fields) == 0 {
			lines = append(lines, fmt.Sprintf("[%s](%s) - %s", k, repo, r.String("licenses")))
		} else {
			lines = append(lines, fmt.Sprintf(" - **[%s](%s)**", k, repo))
			for _, f := range fields {
				lines = append(lines, fmt.Sprintf("    - %s: %v", f.Key, value(r, f.Key)))
			}
		}
	}
	return strings.Join(lines, "\n")
}
func Tree(m *ordered.Map) string {
	var b strings.Builder
	keys := m.Keys()
	for i, k := range keys {
		last := i == len(keys)-1
		pref := "├─ "
		child := "│  "
		if last {
			pref = "└─ "
			child = "   "
		}
		b.WriteString(pref + k + "\n")
		r, _ := m.Get(k)
		rks := r.Keys()
		for j, f := range rks {
			lp := "├─ "
			if j == len(rks)-1 {
				lp = "└─ "
			}
			b.WriteString(child + lp + f + ": " + fmt.Sprint(value(r, f)))
			if !(i == len(keys)-1 && j == len(rks)-1) {
				b.WriteByte('\n')
			}
		}
	}
	return b.String()
}
func Summary(m *ordered.Map) string {
	counts := map[string]int{}
	for _, k := range m.Keys() {
		r, _ := m.Get(k)
		counts[r.String("licenses")]++
	}
	type p struct {
		lic string
		c   int
	}
	var ps []p
	for l, c := range counts {
		ps = append(ps, p{l, c})
	}
	sort.SliceStable(ps, func(i, j int) bool { return ps[i].c > ps[j].c })
	om := ordered.NewMap()
	for _, x := range ps {
		r := ordered.NewRecord()
		r.Set("count", x.c)
		om.Set(x.lic, r)
	}
	return Tree(om)
}
func Files(m *ordered.Map, outDir string, stderr io.Writer) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	for _, k := range m.Keys() {
		r, _ := m.Get(k)
		lf := r.String("licenseFile")
		b, err := os.ReadFile(lf)
		if lf == "" || err != nil {
			fmt.Fprintln(stderr, "no license file found for: "+k)
			continue
		}
		dest := filepath.Join(outDir, k+"-LICENSE.txt")
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, b, 0644); err != nil {
			return err
		}
	}
	return nil
}
