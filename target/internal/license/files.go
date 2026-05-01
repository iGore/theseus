package license

import (
	"sort"
	"strings"

	"theseus/target/internal/core"
)

func orderedCandidates(in []core.LicenseCandidateFile) []core.LicenseCandidateFile {
	out := append([]core.LicenseCandidateFile(nil), in...)
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		ra, rb := fileRank(a.Name), fileRank(b.Name)
		if ra != rb {
			return ra < rb
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
	return out
}

func fileRank(name string) int {
	upper := strings.ToUpper(strings.TrimSpace(name))
	switch {
	case strings.HasPrefix(upper, "LICENSE"):
		return 0
	case strings.HasPrefix(upper, "LICENCE"):
		return 1
	case upper == "COPYING":
		return 2
	case strings.HasPrefix(upper, "README"):
		return 3
	default:
		return 4
	}
}
