package core

import (
	"context"
	"errors"
	"testing"
)

type stubLoader struct {
	start string
	opts  LoaderOptions
	root  *DependencyNode
	err   error
}

func (s *stubLoader) Load(_ context.Context, start string, opts LoaderOptions) (*DependencyNode, error) {
	s.start = start
	s.opts = opts
	if s.err != nil {
		return nil, s.err
	}
	return s.root, nil
}

type stubFlattener struct {
	gotRoot *DependencyNode
	out     ModuleInventoryMap
}

func (s *stubFlattener) Flatten(_ context.Context, root *DependencyNode, _ FlattenOptionsView) (ModuleInventoryMap, error) {
	s.gotRoot = root
	return s.out, nil
}

func TestServiceLoadDependencyGraph_PassesExactStartPath(t *testing.T) {
	loader := &stubLoader{root: &DependencyNode{Name: "root", Version: "1.0.0"}}
	svc := Service{DependencyLoader: loader}

	_, err := svc.LoadDependencyGraph(context.Background(), Options{Start: "/tmp/project", Depth: 0})
	if err != nil {
		t.Fatalf("LoadDependencyGraph() error = %v", err)
	}

	if loader.start != "/tmp/project" {
		t.Fatalf("start path mismatch: got=%q", loader.start)
	}
}

func TestServiceLoadAndFlatten_ForwardsLoadedTree(t *testing.T) {
	root := &DependencyNode{Name: "root", Version: "1.0.0"}
	loader := &stubLoader{root: root}
	flattener := &stubFlattener{out: ModuleInventoryMap{"root@1.0.0": {Name: "root", Version: "1.0.0"}}}
	svc := Service{DependencyLoader: loader, FlattenStage: flattener}

	_, err := svc.LoadAndFlatten(context.Background(), Options{Start: "/tmp/project"}, FlattenOptionsView{IncludeDev: true})
	if err != nil {
		t.Fatalf("LoadAndFlatten() error = %v", err)
	}

	if flattener.gotRoot != root {
		t.Fatalf("flatten stage did not receive loaded root")
	}
}

func TestServiceBuildLoaderOptions_Mapping(t *testing.T) {
	tests := []struct {
		name string
		in   Options
		want LoaderOptions
	}{
		{name: "direct depth with default dev", in: Options{Depth: 0}, want: LoaderOptions{Depth: 0, Dev: true}},
		{name: "recursive depth default", in: Options{Depth: FullTraversalDepth}, want: LoaderOptions{Depth: FullTraversalDepth, Dev: true}},
		{name: "production disables dev", in: Options{Depth: FullTraversalDepth, Production: true}, want: LoaderOptions{Depth: FullTraversalDepth, Dev: false}},
		{name: "development disables dev", in: Options{Depth: FullTraversalDepth, Development: true}, want: LoaderOptions{Depth: FullTraversalDepth, Dev: false}},
	}

	svc := Service{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.BuildLoaderOptions(tt.in)
			if got.Depth != tt.want.Depth || got.Dev != tt.want.Dev {
				t.Fatalf("BuildLoaderOptions() = %+v want %+v", got, tt.want)
			}
		})
	}
}

func TestServiceLoadDependencyGraph_PropagatesLoaderError(t *testing.T) {
	loaderErr := errors.New("unresolvable start")
	loader := &stubLoader{err: loaderErr}
	svc := Service{DependencyLoader: loader}

	_, err := svc.LoadDependencyGraph(context.Background(), Options{Start: "/bad/path"})
	if !errors.Is(err, loaderErr) {
		t.Fatalf("expected loader error, got %v", err)
	}
}
