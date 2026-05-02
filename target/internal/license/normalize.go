package license

import "strings"

func normalize(raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false
	}

	if isCustomReference(value) {
		return "Custom: " + value, false
	}

	if isSPDXExpression(value) {
		return value, false
	}

	if inferred := inferLicense(value); inferred != "" {
		return inferred + "*", true
	}

	return "", false
}

func isCustomReference(value string) bool {
	lower := strings.ToLower(value)
	return strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "file:") ||
		strings.HasPrefix(value, "./") ||
		strings.HasPrefix(value, "../") ||
		strings.HasPrefix(value, "/")
}

func isSPDXExpression(value string) bool {
	replacer := strings.NewReplacer("(", " ", ")", " ", "+", " ")
	cleaned := replacer.Replace(value)
	tokens := strings.Fields(cleaned)
	if len(tokens) == 0 {
		return false
	}

	for _, token := range tokens {
		t := strings.TrimSpace(token)
		if t == "" || t == "AND" || t == "OR" || t == "WITH" {
			continue
		}
		if _, ok := knownSPDX[strings.ToUpper(t)]; !ok {
			return false
		}
	}

	return true
}

func inferLicense(value string) string {
	lower := strings.ToLower(value)
	switch {
	case strings.Contains(lower, "mit"):
		return "MIT"
	case strings.Contains(lower, "apache"):
		return "Apache-2.0"
	case strings.Contains(lower, "bsd"):
		return "BSD"
	case strings.Contains(lower, "gpl"):
		return "GPL"
	case strings.Contains(lower, "unlicense"):
		return "UNLICENSED"
	default:
		return ""
	}
}

var knownSPDX = map[string]struct{}{
	"MIT":          {},
	"APACHE-2.0":   {},
	"BSD-2-CLAUSE": {},
	"BSD-3-CLAUSE": {},
	"GPL-2.0":      {},
	"GPL-3.0":      {},
	"ISC":          {},
	"MPL-2.0":      {},
	"UNLICENSED":   {},
}
