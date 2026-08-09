package policy

import (
	"fmt"
	"io"
	"strings"
	"theseus-target/license-checker/internal/cli"
	"theseus-target/license-checker/internal/diagnostics"
	"theseus-target/license-checker/internal/ordered"
	"theseus-target/license-checker/internal/spdxcompat"
)

func Apply(in *ordered.Map, opts cli.Options, stderr io.Writer) (*ordered.Map, int) {
	sorted := ordered.NewMap()
	for _, k := range in.Keys() {
		r, _ := in.Get(k)
		if opts.OnlyUnknown {
			l := r.String("licenses")
			if !strings.Contains(l, "*") && !strings.Contains(l, "UNKNOWN") {
				continue
			}
		}
		sorted.Set(k, r)
	}
	filtered := sorted
	if opts.Exclude != "" {
		filtered = exclude(sorted, opts.Exclude)
	}
	restricted := filtered
	if opts.Packages != "" {
		restricted = ordered.NewMap()
		list := strings.Split(opts.Packages, ";")
		for _, k := range filtered.Keys() {
			if contains(list, k) {
				r, _ := filtered.Get(k)
				restricted.Set(k, r)
			}
		}
	}
	if opts.ExcludePackages != "" {
		restricted = ordered.NewMap()
		list := strings.Split(opts.ExcludePackages, ";")
		for _, k := range filtered.Keys() {
			if !contains(list, k) {
				r, _ := filtered.Get(k)
				restricted.Set(k, r)
			}
		}
	}
	if opts.ExcludePrivatePackages {
		for _, k := range restricted.Keys() {
			r, _ := restricted.Get(k)
			if r.Bool("private") {
				restricted.Delete(k)
			}
		}
	}
	if opts.FailOn != "" {
		toks := policyTokens(opts.FailOn)
		for _, k := range restricted.Keys() {
			r, _ := restricted.Get(k)
			for _, t := range toks {
				if r.String("licenses") == t {
					fmt.Fprintln(stderr, diagnostics.FailOn(t))
					return restricted, 1
				}
			}
		}
	}
	if opts.OnlyAllow != "" {
		toks := policyTokens(opts.OnlyAllow)
		for _, k := range restricted.Keys() {
			r, _ := restricted.Get(k)
			l := r.String("licenses")
			ok := false
			for _, t := range toks {
				if strings.Contains(l, t) {
					ok = true
				}
			}
			if !ok {
				fmt.Fprintln(stderr, diagnostics.OnlyAllow(k, l))
				return restricted, 1
			}
		}
	}
	return restricted, 0
}
func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
func policyTokens(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ";") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
func splitExclude(s string) []string {
	var out []string
	var b strings.Builder
	esc := false
	for _, r := range s {
		if esc {
			if r == ',' {
				b.WriteRune(',')
			} else {
				b.WriteRune('\\')
				b.WriteRune(r)
			}
			esc = false
			continue
		}
		if r == '\\' {
			esc = true
			continue
		}
		if r == ',' {
			if strings.TrimSpace(b.String()) != "" {
				out = append(out, strings.TrimSpace(b.String()))
			}
			b.Reset()
			continue
		}
		b.WriteRune(r)
	}
	if strings.TrimSpace(b.String()) != "" {
		out = append(out, strings.TrimSpace(b.String()))
	}
	return out
}
func exclude(in *ordered.Map, spec string) *ordered.Map {
	toks := splitExclude(spec)
	out := ordered.NewMap()
	for _, k := range in.Keys() {
		r, _ := in.Get(k)
		l := r.String("licenses")
		keep := true
		if strings.Contains(l, "UNKNOWN") {
			keep = true
		} else {
			cand := l
			if strings.Contains(cand, "*") {
				cand = cand[:len(cand)-1]
			}
			cand = spdxcompat.TransformBSD(cand)
			for _, t := range toks {
				tr := spdxcompat.TransformBSD(t)
				if cand == tr || (!strings.HasPrefix(cand, "Custom:") && spdxcompat.Satisfies(cand, tr)) || (spdxcompat.Correct(cand) != "" && spdxcompat.Satisfies(spdxcompat.Correct(cand), tr)) {
					keep = false
					break
				}
			}
		}
		if keep {
			out.Set(k, r)
		}
	}
	return out
}
