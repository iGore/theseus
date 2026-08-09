package license

import (
	"path/filepath"
	"regexp"
	"strings"
	"theseus-target/license-checker/internal/spdxcompat"
)

func Classify(input string, ok bool) (string, bool) {
	if ok && spdxcompat.Valid(input) {
		return input, true
	}
	s := input
	if ok {
		s = strings.Replace(s, "\n", "", 1)
	} else {
		return "Undefined", true
	}
	if s == "" {
		return "Undefined", true
	}
	rules := []struct {
		re  *regexp.Regexp
		out string
	}{
		{regexp.MustCompile(`The ISC License`), "ISC*"}, {regexp.MustCompile(`ermission is hereby granted, free of charge, to any`), "MIT*"}, {regexp.MustCompile(`edistribution and use in source and binary forms, with or withou`), "BSD*"}, {regexp.MustCompile(`edistribution and use of this software in source and binary forms, with or withou`), "BSD-Source-Code*"}, {regexp.MustCompile(`DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE`), "WTFPL*"}, {regexp.MustCompile(`\bISC\b`), "ISC*"}, {regexp.MustCompile(`\bMIT\b`), "MIT*"}, {regexp.MustCompile(`\bBSD\b`), "BSD*"}, {regexp.MustCompile(`\bWTFPL\b`), "WTFPL*"}, {regexp.MustCompile(`\bApache License\b`), "Apache*"}, {regexp.MustCompile(`http://creativecommons.org/publicdomain/zero/1.0/`), "CC0-1.0*"}, {regexp.MustCompile(`[Pp]ublic [Dd]omain`), "Public Domain"},
	}
	for _, r := range rules {
		if r.re.MatchString(s) {
			return r.out, true
		}
	}
	if m := regexp.MustCompile(`(?i)\bGNU GENERAL PUBLIC LICENSE\s*Version ([^,]*)`).FindStringSubmatch(s); len(m) > 1 {
		v := strings.TrimSpace(m[1])
		if len(v) == 1 {
			v += ".0"
		}
		return "GPL-" + v + "*", true
	}
	if m := regexp.MustCompile(`(?i)(?:LESSER|LIBRARY) GENERAL PUBLIC LICENSE\s*Version ([^,]*)`).FindStringSubmatch(s); len(m) > 1 {
		v := strings.TrimSpace(m[1])
		if len(v) == 1 {
			v += ".0"
		}
		return "LGPL-" + v + "*", true
	}
	if m := regexp.MustCompile(`(?i)(https?://\S+|SEE LICENSE IN (.*))`).FindStringSubmatch(s); len(m) > 0 {
		cap := m[1]
		if len(m) > 2 && m[2] != "" {
			cap = m[2]
		}
		return "Custom: " + cap, true
	}
	return "", false
}

func DetectFiles(names []string) []string {
	pats := []*regexp.Regexp{regexp.MustCompile(`^LICENSE$`), regexp.MustCompile(`^LICENSE\-\w+$`), regexp.MustCompile(`^LICENCE$`), regexp.MustCompile(`^LICENCE\-\w+$`), regexp.MustCompile(`^COPYING$`), regexp.MustCompile(`^README$`)}
	out := []string{}
	for _, p := range pats {
		for _, n := range names {
			base := strings.ToUpper(strings.TrimSuffix(filepath.Base(n), filepath.Ext(n)))
			if p.MatchString(base) {
				out = append(out, n)
				break
			}
		}
	}
	return out
}
