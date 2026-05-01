package core

import (
	"fmt"
	"io"

	"theseus/target/internal/render"
)

type RenderService struct{}

func NewRenderService() RenderService { return RenderService{} }

func (s RenderService) ResolveRenderMode(opts RenderOptions) RenderMode {
	if opts.JSON {
		return RenderModeJSON
	}
	if opts.CSV {
		return RenderModeCSV
	}
	if opts.Markdown {
		return RenderModeMarkdown
	}
	if opts.Summary {
		return RenderModeSummary
	}
	return RenderModeTree
}

func (s RenderService) Execute(modules ModuleResultMap, opts RenderOptions, stdout io.Writer) (string, RenderMode, error) {
	mode := s.ResolveRenderMode(opts)
	rm := toRenderModules(modules)

	var out string
	var err error
	switch mode {
	case RenderModeJSON:
		out, err = render.JSON(rm)
	case RenderModeCSV:
		out, err = render.CSV(rm)
	case RenderModeMarkdown:
		out, err = render.Markdown(rm)
	case RenderModeSummary:
		out, err = render.Summary(rm)
	case RenderModeTree:
		out, err = render.Tree(rm, opts.Colorize && opts.IsTerminal)
	default:
		err = fmt.Errorf("unsupported mode: %s", mode)
	}
	if err != nil {
		return "", mode, err
	}

	if opts.OutPath != "" {
		if err := render.WriteOutputFile(opts.OutPath, out); err != nil {
			return "", mode, err
		}
	} else {
		if _, err := io.WriteString(stdout, out); err != nil {
			return "", mode, err
		}
	}

	if opts.FilesPath != "" {
		if err := render.ExportLicenseFiles(rm, opts.FilesPath); err != nil {
			return "", mode, err
		}
	}

	return out, mode, nil
}

func toRenderModules(modules ModuleResultMap) render.ModuleResultMap {
	out := make(render.ModuleResultMap, len(modules))
	for key, v := range modules {
		out[key] = render.ModuleResult{Name: v.Name, Version: v.Version, License: v.Licenses, LicenseFile: v.LicenseFile}
	}
	return out
}
