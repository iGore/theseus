package cli

import "strings"

const UnboundedDepth = -1

// Options is the normalized source-compatible CLI option set.
type Options struct {
	Production, Development, JSON, CSV, Markdown, Unknown, OnlyUnknown, Version, Color, Help, RelativeLicensePath, Summary, ExcludePrivatePackages bool
	CSVComponentPrefix, Out, Start, Exclude, CustomPath, Files, FailOn, OnlyAllow, Packages, ExcludePackages                                       string
	Direct                                                                                                                                         int
	CustomFormat                                                                                                                                   []Field
}
type Field struct {
	Key   string
	Value any
}

func Parse(args []string, cwd string, supportsColor bool) Options {
	o := Options{Color: supportsColor, Start: cwd, Direct: UnboundedDepth}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "-v" {
			a = "--version"
		}
		if a == "-h" {
			a = "--help"
		}
		val := func() string {
			if j := strings.IndexByte(a, '='); j >= 0 {
				return a[j+1:]
			}
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				return args[i]
			}
			return ""
		}
		name := a
		if j := strings.IndexByte(a, '='); j >= 0 {
			name = a[:j]
		}
		switch name {
		case "--production":
			o.Production = true
		case "--development":
			o.Development = true
		case "--json":
			o.JSON = true
		case "--csv":
			o.CSV = true
		case "--markdown":
			o.Markdown = true
		case "--unknown":
			o.Unknown = true
		case "--onlyunknown":
			o.OnlyUnknown = true
		case "--version":
			o.Version = true
		case "--help":
			o.Help = true
		case "--relativeLicensePath":
			o.RelativeLicensePath = true
		case "--summary":
			o.Summary = true
		case "--excludePrivatePackages":
			o.ExcludePrivatePackages = true
		case "--direct":
			o.Direct = 0
		case "--color":
			o.Color = true
		case "--no-color":
			o.Color = false
		case "--csvComponentPrefix":
			o.CSVComponentPrefix = val()
		case "--out":
			o.Out = val()
		case "--start":
			if v := val(); v != "" {
				o.Start = v
			}
		case "--exclude":
			o.Exclude = val()
		case "--customPath":
			o.CustomPath = val()
		case "--files":
			o.Files = val()
		case "--failOn":
			o.FailOn = val()
		case "--onlyAllow":
			o.OnlyAllow = val()
		case "--packages":
			o.Packages = val()
		case "--excludePackages":
			o.ExcludePackages = val()
		}
	}
	if o.JSON || o.CSV || o.Markdown {
		o.Color = false
	}
	return o
}

func Has(args []string, name string) bool {
	p, n := "--"+name, "--no-"+name
	for _, a := range args {
		if a == p || a == n || strings.HasPrefix(a, p+"=") {
			return true
		}
	}
	return false
}
