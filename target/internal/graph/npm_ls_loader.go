package graph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"theseus/target/internal/core"
)

type commandRunner interface {
	Run(ctx context.Context, name string, args []string, dir string) ([]byte, error)
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, name string, args []string, dir string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	return cmd.Output()
}

type NpmLSLoader struct {
	runner commandRunner
}

func NewNpmLSLoader() NpmLSLoader {
	return NpmLSLoader{runner: execRunner{}}
}

func (l NpmLSLoader) Load(ctx context.Context, start string, opts core.LoaderOptions) (*core.DependencyNode, error) {
	if l.runner == nil {
		l.runner = execRunner{}
	}

	args := []string{"ls", "--json"}
	if opts.Depth == 0 {
		args = append(args, "--depth=0")
	}
	if !opts.Dev {
		args = append(args, "--omit=dev")
	}

	out, err := l.runner.Run(ctx, "npm", args, start)
	if err != nil {
		if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("npm ls failed: %w", err)
	}

	var root npmLSNode
	if err := json.Unmarshal(out, &root); err != nil {
		return nil, fmt.Errorf("parse npm ls output: %w", err)
	}

	return root.toCoreNode(), nil
}

type npmLSNode struct {
	Name         string                `json:"name"`
	Version      string                `json:"version"`
	Path         string                `json:"path"`
	Private      bool                  `json:"private"`
	Dependencies map[string]*npmLSNode `json:"dependencies"`
}

func (n *npmLSNode) toCoreNode() *core.DependencyNode {
	if n == nil {
		return nil
	}

	out := &core.DependencyNode{
		Name:         n.Name,
		Version:      n.Version,
		Path:         n.Path,
		Private:      n.Private,
		Dependencies: make(map[string]*core.DependencyNode, len(n.Dependencies)),
	}
	for k, v := range n.Dependencies {
		out.Dependencies[k] = v.toCoreNode()
	}
	return out
}
