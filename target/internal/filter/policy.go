package filter

import (
	"fmt"
	"strings"
)

type Matcher func(license string, rules []string) bool

// onlyAllowMatcher preserves current compatibility semantics.
// NEEDS CLARIFICATION (OD-001): strict SPDX/exact behavior is unresolved.
var onlyAllowMatcher Matcher = matchesAnyRule

func EvaluatePolicies(modules map[string]ModuleEntry, opts Options) error {
	failOn := parseCommaListWithEscapes(opts.FailOn)
	onlyAllow := parseCommaListWithEscapes(opts.OnlyAllow)

	if len(failOn) == 0 && len(onlyAllow) == 0 {
		return nil
	}

	for pkg, entry := range modules {
		if len(failOn) > 0 && matchesAnyRule(entry.Licenses, failOn) {
			return fmt.Errorf("policy violation: failOn matched package=%s license=%s", pkg, entry.Licenses)
		}

		if len(onlyAllow) > 0 && !onlyAllowMatcher(entry.Licenses, onlyAllow) {
			return fmt.Errorf("policy violation: onlyAllow mismatch package=%s license=%s", pkg, entry.Licenses)
		}
	}

	return nil
}

func ValidatePolicyCombination(opts Options) error {
	if strings.TrimSpace(opts.FailOn) != "" && strings.TrimSpace(opts.OnlyAllow) != "" {
		return fmt.Errorf("invalid policy combination: failOn and onlyAllow are mutually exclusive")
	}
	return nil
}
