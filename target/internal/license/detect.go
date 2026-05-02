package license

import (
	"fmt"

	"theseus/target/internal/core"
)

type Detector struct{}

func (Detector) Resolve(input core.ModuleLicenseInput) (core.NormalizedLicenseValue, string) {
	if v := resolveFromMetadata(input); v != "" {
		return core.NormalizedLicenseValue(v), ""
	}

	if v, _ := normalize(input.ReadmeText); v != "" {
		return core.NormalizedLicenseValue(v), "README"
	}

	for _, candidate := range orderedCandidates(input.CandidateFiles) {
		if v, _ := normalize(candidate.Content); v != "" {
			return core.NormalizedLicenseValue(v), candidate.Name
		}
	}

	// OD-002 NEEDS CLARIFICATION: keep explicit sentinel until legacy parity is confirmed.
	return core.NormalizedLicenseValue("UNKNOWN"), ""
}

func resolveFromMetadata(input core.ModuleLicenseInput) string {
	if v := firstNormalizedFromAny(input.License); v != "" {
		return v
	}
	// OD-001 NEEDS CLARIFICATION: in conflicts, current PoC prioritizes `license` over `licenses`.
	if v := firstNormalizedFromAny(input.Licenses); v != "" {
		return v
	}
	return ""
}

func firstNormalizedFromAny(raw any) string {
	switch v := raw.(type) {
	case string:
		out, _ := normalize(v)
		return out
	case map[string]any:
		if t, ok := v["type"].(string); ok {
			out, _ := normalize(t)
			return out
		}
		if n, ok := v["name"].(string); ok {
			out, _ := normalize(n)
			return out
		}
	case []any:
		for _, item := range v {
			if out := firstNormalizedFromAny(item); out != "" {
				return out
			}
		}
	default:
		if v != nil {
			out, _ := normalize(fmt.Sprint(v))
			return out
		}
	}
	return ""
}
