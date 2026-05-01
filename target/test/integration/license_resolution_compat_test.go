package integration

import (
	"testing"

	"theseus/target/internal/core"
	"theseus/target/internal/license"
)

func TestLicenseCompatibility_OD001_Placeholder(t *testing.T) {
	t.Skip("NEEDS CLARIFICATION (OD-001): confirm legacy tie-break when license and licenses conflict")

	_ = core.ModuleLicenseInput{License: "MIT", Licenses: []any{"GPL-3.0"}}
	_, _ = (license.Detector{}).Resolve(core.ModuleLicenseInput{})
}

func TestLicenseCompatibility_OD002_Placeholder(t *testing.T) {
	t.Skip("NEEDS CLARIFICATION (OD-002): confirm fallback value when no recognizable signal exists")

	_, _ = (license.Detector{}).Resolve(core.ModuleLicenseInput{})
}
