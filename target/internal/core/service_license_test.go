package core

import "testing"

type stubLicenseStage struct{}

func (stubLicenseStage) Resolve(input ModuleLicenseInput) (NormalizedLicenseValue, string) {
	if input.License == "MIT" {
		return NormalizedLicenseValue("MIT"), ""
	}
	return NormalizedLicenseValue("UNKNOWN"), ""
}

func TestServiceResolveLicenses_AssignsSingleValuePerModule(t *testing.T) {
	svc := Service{LicenseStage: stubLicenseStage{}}
	in := ModuleInventoryMap{
		"a@1.0.0": {Name: "a", Version: "1.0.0", LicenseInput: ModuleLicenseInput{License: "MIT"}},
	}

	out := svc.ResolveLicenses(in)
	if got := out["a@1.0.0"].Licenses; got != "MIT" {
		t.Fatalf("got %q want MIT", got)
	}
}
