package spdxcompat

import "strings"

var ids = map[string]bool{"MIT": true, "ISC": true, "Apache-2.0": true, "BSD-2-Clause": true, "BSD-3-Clause": true, "BSD-4-Clause": true, "0BSD": true, "LGPL-2.0": true, "GPL-2.0+": true, "Bison-exception-2.2": true, "CC0-1.0": true}

func Valid(expr string) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false
	}
	rep := strings.NewReplacer("(", " ", ")", " ", "+", "+ ").Replace(expr)
	for _, tok := range strings.Fields(rep) {
		if tok == "AND" || tok == "OR" || tok == "WITH" {
			continue
		}
		if !ids[strings.TrimSpace(tok)] {
			return false
		}
	}
	return true
}
func Correct(s string) string {
	s = strings.TrimSpace(s)
	if Valid(s) {
		return s
	}
	switch strings.ToLower(s) {
	case "mit license":
		return "MIT"
	case "apache 2.0":
		return "Apache-2.0"
	}
	return ""
}
func Satisfies(candidate, allowed string) bool {
	candidate = TransformBSD(strings.TrimSpace(candidate))
	allowed = TransformBSD(allowed)
	if strings.Contains(allowed, " OR ") || strings.HasPrefix(allowed, "(") {
		for _, part := range strings.Split(strings.Trim(allowed, "() "), " OR ") {
			if Satisfies(candidate, part) {
				return true
			}
		}
		return false
	}
	return candidate == strings.TrimSpace(allowed)
}
func TransformBSD(s string) string {
	if strings.TrimSpace(s) == "BSD" {
		return "(0BSD OR BSD-2-Clause OR BSD-3-Clause OR BSD-4-Clause)"
	}
	return s
}
