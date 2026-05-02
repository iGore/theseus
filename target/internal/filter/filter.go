package filter

import (
	"strings"
)

type ModuleEntry struct {
	Licenses string
	Private  bool
}

type Options struct {
	Unknown                bool
	OnlyUnknown            bool
	Exclude                string
	Packages               string
	ExcludePackages        string
	ExcludePrivatePackages bool
	FailOn                 string
	OnlyAllow              string
}

// Apply runs the F-005 filter stage in required order.
func Apply(input map[string]ModuleEntry, opts Options) map[string]ModuleEntry {
	result := clone(input)

	if opts.Unknown {
		for key, entry := range result {
			if strings.HasSuffix(entry.Licenses, "*") {
				entry.Licenses = "UNKNOWN"
				result[key] = entry
			}
		}
	}

	if opts.OnlyUnknown {
		for key, entry := range result {
			if entry.Licenses != "UNKNOWN" && !strings.HasSuffix(entry.Licenses, "*") {
				delete(result, key)
			}
		}
	}

	excludeRules := parseCommaListWithEscapes(opts.Exclude)
	if len(excludeRules) > 0 {
		for key, entry := range result {
			if matchesAnyRule(entry.Licenses, excludeRules) {
				delete(result, key)
			}
		}
	}

	includePackages := parseSemicolonList(opts.Packages)
	if len(includePackages) > 0 {
		for key := range result {
			if _, ok := includePackages[key]; !ok {
				delete(result, key)
			}
		}
	}

	excludePackages := parseSemicolonList(opts.ExcludePackages)
	if len(excludePackages) > 0 {
		for key := range result {
			if _, ok := excludePackages[key]; ok {
				delete(result, key)
			}
		}
	}

	if opts.ExcludePrivatePackages {
		for key, entry := range result {
			if entry.Private {
				delete(result, key)
			}
		}
	}

	return result
}

func clone(input map[string]ModuleEntry) map[string]ModuleEntry {
	result := make(map[string]ModuleEntry, len(input))
	for k, v := range input {
		result[k] = v
	}
	return result
}

func parseSemicolonList(raw string) map[string]struct{} {
	set := map[string]struct{}{}
	for _, part := range strings.Split(raw, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		set[part] = struct{}{}
	}
	return set
}

func parseCommaListWithEscapes(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	var (
		out     []string
		current strings.Builder
		escaped bool
	)

	for _, r := range raw {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
		case r == ',':
			item := strings.TrimSpace(current.String())
			if item != "" {
				out = append(out, item)
			}
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	item := strings.TrimSpace(current.String())
	if item != "" {
		out = append(out, item)
	}
	return out
}

func matchesAnyRule(license string, rules []string) bool {
	licenseLower := strings.ToLower(license)
	for _, rule := range rules {
		ruleLower := strings.ToLower(strings.TrimSpace(rule))
		if ruleLower == "" {
			continue
		}
		// Compatibility behavior for exclude=BSD matching BSD variants.
		if ruleLower == "bsd" {
			if strings.Contains(licenseLower, "bsd") {
				return true
			}
			continue
		}
		if strings.Contains(licenseLower, ruleLower) {
			return true
		}
	}
	return false
}
