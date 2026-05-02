package core

import (
	"context"
	"errors"
	"log/slog"

	"theseus/target/internal/filter"
)

type ExitStrategy interface {
	Exit(code int)
}

// FlattenStage defines the F-003 stage boundary.
type FlattenStage interface {
	Flatten(ctx context.Context, root *DependencyNode, opts FlattenOptionsView) (ModuleInventoryMap, error)
}

// FilterStage decouples core from concrete filter package to avoid cycles.
type FilterStage interface {
	Apply(modules map[string]ModuleEntry, opts Options) map[string]ModuleEntry
	EvaluatePolicies(modules map[string]ModuleEntry, opts Options) error
}

// LicenseStage defines the F-004 stage boundary.
type LicenseStage interface {
	Resolve(input ModuleLicenseInput) (NormalizedLicenseValue, string)
}

type Service struct {
	DependencyLoader DependencyLoader
	FlattenStage     FlattenStage
	ExitStrategy     ExitStrategy
	FilterStage      FilterStage
	LicenseStage     LicenseStage
	Logger           *slog.Logger
}

func (s Service) BuildLoaderOptions(opts Options) LoaderOptions {
	depth := FullTraversalDepth
	if opts.Depth == 0 {
		depth = 0
	}

	dev := true
	if opts.Production || opts.Development {
		dev = false
	}

	return LoaderOptions{Depth: depth, Dev: dev, Logger: s.Logger}
}

func (s Service) LoadDependencyGraph(ctx context.Context, opts Options) (*DependencyNode, error) {
	if s.DependencyLoader == nil {
		return nil, errors.New("dependency loader is required")
	}

	root, err := s.DependencyLoader.Load(ctx, opts.Start, s.BuildLoaderOptions(opts))
	if err != nil {
		return nil, err
	}
	return root, nil
}

func (s Service) LoadAndFlatten(ctx context.Context, opts Options, flattenOpts FlattenOptionsView) (ModuleInventoryMap, error) {
	if s.FlattenStage == nil {
		return nil, errors.New("flatten stage is required")
	}

	root, err := s.LoadDependencyGraph(ctx, opts)
	if err != nil {
		return nil, err
	}

	return s.FlattenStage.Flatten(ctx, root, flattenOpts)
}

type defaultFilterStage struct{}

func (defaultFilterStage) Apply(modules map[string]ModuleEntry, opts Options) map[string]ModuleEntry {
	in := make(map[string]filter.ModuleEntry, len(modules))
	for k, v := range modules {
		in[k] = filter.ModuleEntry{Licenses: v.Licenses, Private: v.Private}
	}
	out := filter.Apply(in, filter.Options{
		Unknown:                opts.Unknown,
		OnlyUnknown:            opts.OnlyUnknown,
		Exclude:                opts.Exclude,
		Packages:               opts.Packages,
		ExcludePackages:        opts.ExcludePackages,
		ExcludePrivatePackages: opts.ExcludePrivatePackages,
		FailOn:                 opts.FailOn,
		OnlyAllow:              opts.OnlyAllow,
	})
	converted := make(map[string]ModuleEntry, len(out))
	for k, v := range out {
		converted[k] = ModuleEntry{Licenses: v.Licenses, Private: v.Private}
	}
	return converted
}

func (defaultFilterStage) EvaluatePolicies(modules map[string]ModuleEntry, opts Options) error {
	in := make(map[string]filter.ModuleEntry, len(modules))
	for k, v := range modules {
		in[k] = filter.ModuleEntry{Licenses: v.Licenses, Private: v.Private}
	}
	return filter.EvaluatePolicies(in, filter.Options{FailOn: opts.FailOn, OnlyAllow: opts.OnlyAllow})
}

func (s Service) ApplyFilteringAndPolicy(modules map[string]ModuleEntry, opts Options) (map[string]ModuleEntry, error) {
	if s.FilterStage == nil {
		s.FilterStage = defaultFilterStage{}
	}
	filtered := s.FilterStage.Apply(modules, opts)
	if err := s.FilterStage.EvaluatePolicies(filtered, opts); err != nil {
		if s.ExitStrategy != nil {
			s.ExitStrategy.Exit(1)
		}
		return filtered, err
	}
	return filtered, nil
}

// ResolveLicenses applies F-004 license resolution to each module and returns a copied map.
func (s Service) ResolveLicenses(modules ModuleInventoryMap) ModuleInventoryMap {
	if s.LicenseStage == nil {
		return modules
	}

	resolved := make(ModuleInventoryMap, len(modules))
	for key, entry := range modules {
		next := entry
		value, file := s.LicenseStage.Resolve(entry.LicenseInput)
		next.Licenses = string(value)
		next.LicenseFile = file
		resolved[key] = next
	}

	return resolved
}
