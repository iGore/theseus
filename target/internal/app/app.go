package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"theseus-target/license-checker/internal/cli"
	"theseus-target/license-checker/internal/customformat"
	"theseus-target/license-checker/internal/diagnostics"
	"theseus-target/license-checker/internal/npmgraph"
	"theseus-target/license-checker/internal/policy"
	"theseus-target/license-checker/internal/records"
	"theseus-target/license-checker/internal/render"
)

type Invocation struct {
	Args           []string
	Stdout, Stderr io.Writer
	Getwd          func() (string, error)
}

func Run(ctx context.Context, inv Invocation) int {
	_ = ctx
	cwd, err := inv.Getwd()
	if err != nil {
		fmt.Fprintln(inv.Stderr, "Found error")
		fmt.Fprintln(inv.Stderr, err)
		return 0
	}
	opts := cli.Parse(inv.Args, cwd, false)
	if handled, code := diagnostics.Preflight(opts, inv.Stderr); handled {
		return code
	}
	if opts.CustomPath != "" {
		cf, err := customformat.Load(opts.CustomPath)
		if err != nil {
			fmt.Fprintln(inv.Stderr, "Found error")
			fmt.Fprintln(inv.Stderr, err)
			return 0
		}
		opts.CustomFormat = cf
	}
	n, err := (npmgraph.Scanner{}).Scan(opts.Start, npmgraph.Options{Dev: !(opts.Production || opts.Development), Depth: opts.Direct})
	if err != nil {
		fmt.Fprintln(inv.Stderr, "Found error")
		fmt.Fprintln(inv.Stderr, err)
		return 0
	}
	m := records.Flatten(n, opts)
	m, code := policy.Apply(m, opts, inv.Stderr)
	if code != 0 {
		return code
	}
	f, err := render.Format(m, opts)
	if err != nil {
		fmt.Fprintln(inv.Stderr, "Found error")
		fmt.Fprintln(inv.Stderr, err)
		return 0
	}
	if err := render.Emit(m, f, opts, inv.Stdout, inv.Stderr); err != nil {
		fmt.Fprintln(inv.Stderr, "Found error")
		fmt.Fprintln(inv.Stderr, err)
		return 0
	}
	return 0
}
func DefaultInvocation(args []string) Invocation {
	return Invocation{Args: args, Stdout: os.Stdout, Stderr: os.Stderr, Getwd: os.Getwd}
}
