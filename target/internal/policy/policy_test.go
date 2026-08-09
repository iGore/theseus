package policy

import (
	"strings"
	"testing"
	"theseus-target/license-checker/internal/cli"
	"theseus-target/license-checker/internal/ordered"
)

func pm(items map[string]string) *ordered.Map {
	m := ordered.NewMap()
	keys := []string{"a@1", "b@1", "c@1", "d@1"}
	for _, k := range keys {
		if l, ok := items[k]; ok {
			r := ordered.NewRecord()
			r.Set("licenses", l)
			if k == "d@1" {
				r.Set("private", true)
			}
			m.Set(k, r)
		}
	}
	return m
}
func TestPolicyContracts(t *testing.T) {
	m, _ := Apply(pm(map[string]string{"a@1": "MIT", "b@1": "UNKNOWN", "c@1": "BSD-3-Clause"}), cli.Options{Exclude: "MIT, BSD"}, &strings.Builder{})
	if _, ok := m.Get("a@1"); ok {
		t.Fatal("MIT not excluded")
	}
	if _, ok := m.Get("c@1"); ok {
		t.Fatal("BSD not excluded")
	}
	if _, ok := m.Get("b@1"); !ok {
		t.Fatal("UNKNOWN excluded")
	}
	m, _ = Apply(pm(map[string]string{"a@1": "MIT*", "b@1": "UNKNOWN", "c@1": "MIT", "d@1": "UNLICENSED"}), cli.Options{OnlyUnknown: true}, &strings.Builder{})
	if m.Len() != 2 {
		t.Fatalf("onlyunknown len %d", m.Len())
	}
	var errb strings.Builder
	_, code := Apply(pm(map[string]string{"a@1": "MIT"}), cli.Options{FailOn: "ISC;MIT"}, &errb)
	if code != 1 || errb.String() != "Found license defined by the --failOn flag: \"MIT\". Exiting.\n" {
		t.Fatalf("%d %q", code, errb.String())
	}
}
