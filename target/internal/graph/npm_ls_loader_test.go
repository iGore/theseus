package graph

import (
	"context"
	"errors"
	"strings"
	"testing"

	"theseus/target/internal/core"
)

type fakeRunner struct {
	name string
	args []string
	dir  string
	out  []byte
	err  error
}

func (f *fakeRunner) Run(_ context.Context, name string, args []string, dir string) ([]byte, error) {
	f.name = name
	f.args = append([]string(nil), args...)
	f.dir = dir
	if f.err != nil {
		return nil, f.err
	}
	return f.out, nil
}

func TestNpmLSLoader_Load_UsesStartPathAndParsesTree(t *testing.T) {
	runner := &fakeRunner{out: []byte(`{"name":"root","version":"1.0.0","dependencies":{"left-pad":{"name":"left-pad","version":"1.3.0"}}}`)}
	loader := NpmLSLoader{runner: runner}

	root, err := loader.Load(context.Background(), "/workspace/proj", core.LoaderOptions{Depth: 0, Dev: false})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if runner.name != "npm" || runner.dir != "/workspace/proj" {
		t.Fatalf("command mismatch: name=%q dir=%q", runner.name, runner.dir)
	}
	if !contains(runner.args, "--depth=0") || !contains(runner.args, "--omit=dev") {
		t.Fatalf("expected depth/dev args, got %v", runner.args)
	}
	if root == nil || root.Name != "root" {
		t.Fatalf("unexpected root: %#v", root)
	}
	if root.Dependencies["left-pad"].Version != "1.3.0" {
		t.Fatalf("dependency parse mismatch: %#v", root.Dependencies["left-pad"])
	}
}

func TestNpmLSLoader_Load_ContextCancellationPropagates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	runner := &fakeRunner{err: context.Canceled}
	loader := NpmLSLoader{runner: runner}

	_, err := loader.Load(ctx, "/workspace/proj", core.LoaderOptions{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestNpmLSLoader_Load_WrapsCommandErrors(t *testing.T) {
	runner := &fakeRunner{err: errors.New("exit status 1")}
	loader := NpmLSLoader{runner: runner}

	_, err := loader.Load(context.Background(), "/workspace/proj", core.LoaderOptions{})
	if err == nil || !strings.Contains(err.Error(), "npm ls failed") {
		t.Fatalf("expected wrapped command error, got %v", err)
	}
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}
