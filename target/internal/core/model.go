package core

import (
	"context"
	"log/slog"
)

const FullTraversalDepth = -1

type VersionExitMode int

const (
	VersionExitSuccess VersionExitMode = iota
	VersionExitLegacyNonZero
)

type RawArgs struct {
	Start     string
	Direct    bool
	CustomPath string
	JSON      bool
	CSV       bool
	Markdown  bool
	Color     bool
	ColorSet  bool
	FailOn    string
	OnlyAllow string
	Version   bool
}

type GuardrailResult struct {
	ExitCode    int
	Warnings    []string
	MetaCommand string
	Options     *Options
}

type Options struct {
	Start                  string
	Depth                  int
	CustomPath             string
	CustomFormat           map[string]any
	Production             bool
	Development            bool
	Color                  bool
	Unknown                bool
	OnlyUnknown            bool
	Exclude                string
	Packages               string
	ExcludePackages        string
	ExcludePrivatePackages bool
	FailOn                 string
	OnlyAllow              string
}

type ModuleEntry struct {
	Name         string
	Version      string
	Repository   string
	Author       string
	URL          string
	Path         string
	Private      bool
	Dev          bool
	Licenses     string
	LicenseFile  string
	LicenseText  string
	Copyright    string
	LicenseInput ModuleLicenseInput
	CustomFields map[string]string
}

type ModuleInventoryMap map[string]ModuleEntry

type DependencyNode struct {
	Name         string
	Version      string
	Repository   string
	Author       string
	URL          string
	Path         string
	Private      bool
	Dev          bool
	Dependencies map[string]*DependencyNode
}

type LoaderOptions struct {
	Depth  int
	Dev    bool
	Logger *slog.Logger
}

type DependencyLoader interface {
	Load(ctx context.Context, start string, opts LoaderOptions) (*DependencyNode, error)
}

type FlattenOptionsView struct {
	IncludeDev bool
}

// ModuleLicenseInput carries all F-004 inputs needed to resolve one license value.
type ModuleLicenseInput struct {
	License        any
	Licenses       any
	ReadmeText     string
	CandidateFiles []LicenseCandidateFile
}

// LicenseCandidateFile is an ordered candidate for file-based fallback.
type LicenseCandidateFile struct {
	Name    string
	Content string
}

// NormalizedLicenseValue is the normalized single-value output of F-004.
type NormalizedLicenseValue string

type RenderMode string

const (
	RenderModeTree     RenderMode = "tree"
	RenderModeSummary  RenderMode = "summary"
	RenderModeMarkdown RenderMode = "markdown"
	RenderModeCSV      RenderMode = "csv"
	RenderModeJSON     RenderMode = "json"
)

type RenderOptions struct {
	JSON       bool
	CSV        bool
	Markdown   bool
	Summary    bool
	Colorize   bool
	IsTerminal bool
	OutPath    string
	FilesPath  string
}

type ModuleResultMap map[string]ModuleEntry
