package integration

import (
	"context"
	"testing"

	"theseus/target/internal/core"
	"theseus/target/internal/flatten"
)

func TestFlattenStageContract_DownstreamCanConsumeMap(t *testing.T) {
	root := &core.DependencyNode{
		Name:    "root",
		Version: "1.0.0",
		Dependencies: map[string]*core.DependencyNode{
			"left":  {Name: "left", Version: "0.0.1"},
			"right": {Name: "right", Version: "0.0.2"},
		},
	}

	inventory, err := (flatten.Flattener{}).Flatten(context.Background(), root, core.FlattenOptionsView{IncludeDev: true})
	if err != nil {
		t.Fatalf("flatten error: %v", err)
	}

	// Mimic downstream consumer behavior that relies only on flat map.
	var count int
	for key, entry := range inventory {
		if key == "" || entry.Name == "" || entry.Version == "" {
			t.Fatalf("invalid downstream contract entry: key=%q entry=%+v", key, entry)
		}
		count++
	}

	if count != 3 {
		t.Fatalf("expected 3 flattened entries, got %d", count)
	}
}
