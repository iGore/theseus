package cli

import (
	"strings"

	"theseus/target/internal/core"
	"theseus/target/internal/filter"
)

func NormalizeOptions(raw core.RawArgs, getwd func() (string, error)) (core.Options, error) {
	start := raw.Start
	if strings.TrimSpace(start) == "" {
		cwd, err := getwd()
		if err != nil {
			return core.Options{}, err
		}
		start = cwd
	}

	depth := core.FullTraversalDepth
	if raw.Direct {
		depth = 0
	}

	color := true
	if raw.ColorSet {
		color = raw.Color
	}
	if raw.JSON || raw.CSV || raw.Markdown {
		color = false
	}

	return core.Options{
		Start:      start,
		Depth:      depth,
		CustomPath: raw.CustomPath,
		Color:      color,
		FailOn:     raw.FailOn,
		OnlyAllow:  raw.OnlyAllow,
	}, nil
}

func ValidateOptions(opts core.Options) error {
	return filter.ValidatePolicyCombination(filter.Options{FailOn: opts.FailOn, OnlyAllow: opts.OnlyAllow})
}
