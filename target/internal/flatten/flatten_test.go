package flatten

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"theseus/target/internal/core"
)

func TestFlatten_NestedToFlat(t *testing.T) {
	root := &core.DependencyNode{
		Name:    "root",
		Version: "1.0.0",
		Dependencies: map[string]*core.DependencyNode{
			"a": {
				Name:    "a",
				Version: "1.0.0",
				Dependencies: map[string]*core.DependencyNode{
					"b": {Name: "b", Version: "2.0.0"},
				},
			},
		},
	}

	got, err := (Flattener{}).Flatten(context.Background(), root, core.FlattenOptionsView{IncludeDev: true})
	if err != nil {
		t.Fatalf("Flatten() error = %v", err)
	}

	wantKeys := []string{"a@1.0.0", "b@2.0.0", "root@1.0.0"}
	if !reflect.DeepEqual(sortedKeys(got), wantKeys) {
		t.Fatalf("keys mismatch: got=%v want=%v", sortedKeys(got), wantKeys)
	}
}

func TestFlatten_DedupAndCycle(t *testing.T) {
	a := &core.DependencyNode{Name: "a", Version: "1.0.0", Dependencies: map[string]*core.DependencyNode{}}
	b := &core.DependencyNode{Name: "b", Version: "1.0.0", Dependencies: map[string]*core.DependencyNode{}}
	c := &core.DependencyNode{Name: "c", Version: "1.0.0", Dependencies: map[string]*core.DependencyNode{}}

	a.Dependencies["b"] = b
	b.Dependencies["c"] = c
	c.Dependencies["a"] = a // cycle

	root := &core.DependencyNode{
		Name:    "root",
		Version: "1.0.0",
		Dependencies: map[string]*core.DependencyNode{
			"a1": a,
			"a2": a, // duplicate path
		},
	}

	got, err := (Flattener{}).Flatten(context.Background(), root, core.FlattenOptionsView{IncludeDev: true})
	if err != nil {
		t.Fatalf("Flatten() error = %v", err)
	}

	wantKeys := []string{"a@1.0.0", "b@1.0.0", "c@1.0.0", "root@1.0.0"}
	if !reflect.DeepEqual(sortedKeys(got), wantKeys) {
		t.Fatalf("keys mismatch: got=%v want=%v", sortedKeys(got), wantKeys)
	}
}

func TestFlatten_GatingAndMetadata(t *testing.T) {
	root := &core.DependencyNode{
		Name:       "root",
		Version:    "1.0.0",
		Repository: "https://example.com/root",
		Author:     "Root Author",
		URL:        "https://pkg/root",
		Path:       "/tmp/root",
		Private:    true,
		Dependencies: map[string]*core.DependencyNode{
			"dev": {
				Name:    "devmod",
				Version: "0.1.0",
				Dev:     true,
			},
			"prod": {
				Name:    "prodmod",
				Version: "2.3.4",
			},
		},
	}

	got, err := (Flattener{}).Flatten(context.Background(), root, core.FlattenOptionsView{IncludeDev: false})
	if err != nil {
		t.Fatalf("Flatten() error = %v", err)
	}

	if _, exists := got["devmod@0.1.0"]; exists {
		t.Fatalf("dev module should be excluded")
	}
	if _, exists := got["prodmod@2.3.4"]; !exists {
		t.Fatalf("prod module should be included")
	}

	rootEntry, ok := got["root@1.0.0"]
	if !ok {
		t.Fatalf("root entry missing")
	}

	if rootEntry.Repository != "https://example.com/root" ||
		rootEntry.Author != "Root Author" ||
		rootEntry.URL != "https://pkg/root" ||
		rootEntry.Path != "/tmp/root" ||
		!rootEntry.Private {
		t.Fatalf("metadata mismatch: %+v", rootEntry)
	}
}

func TestFlatten_MissingIdentitySkipped(t *testing.T) {
	root := &core.DependencyNode{
		Name:    "root",
		Version: "1.0.0",
		Dependencies: map[string]*core.DependencyNode{
			"missing-name": {Name: "", Version: "1.0.0"},
			"missing-ver":  {Name: "x", Version: ""},
		},
	}

	got, err := (Flattener{}).Flatten(context.Background(), root, core.FlattenOptionsView{IncludeDev: true})
	if err != nil {
		t.Fatalf("Flatten() error = %v", err)
	}

	wantKeys := []string{"root@1.0.0"}
	if !reflect.DeepEqual(sortedKeys(got), wantKeys) {
		t.Fatalf("keys mismatch: got=%v want=%v", sortedKeys(got), wantKeys)
	}
}

func sortedKeys(m core.ModuleInventoryMap) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
