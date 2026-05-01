package flatten

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"theseus/target/internal/core"
)

// Flattener converts a nested dependency graph into a canonical module inventory map.
type Flattener struct{}

// Flatten transforms the graph rooted at root into a map keyed by name@version.
func (Flattener) Flatten(ctx context.Context, root *core.DependencyNode, opts core.FlattenOptionsView) (core.ModuleInventoryMap, error) {
	if root == nil {
		return core.ModuleInventoryMap{}, nil
	}

	out := make(core.ModuleInventoryMap)
	visited := make(map[string]struct{})

	if err := walk(ctx, root, opts, out, visited); err != nil {
		return nil, err
	}

	return out, nil
}

func walk(ctx context.Context, n *core.DependencyNode, opts core.FlattenOptionsView, out core.ModuleInventoryMap, visited map[string]struct{}) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if n == nil {
		return nil
	}

	key, ok := canonicalKey(n)
	if !ok {
		// OD-001 unresolved: skip nodes missing identity fields until parity behavior is clarified.
		return nil
	}

	if _, exists := visited[key]; exists {
		return nil
	}
	visited[key] = struct{}{}

	if includeNode(n, opts) {
		out[key] = toModuleEntry(n)
	}

	for _, child := range n.Dependencies {
		if err := walk(ctx, child, opts, out, visited); err != nil {
			return err
		}
	}

	return nil
}

func canonicalKey(n *core.DependencyNode) (string, bool) {
	if n.Name == "" || n.Version == "" {
		return "", false
	}
	return fmt.Sprintf("%s@%s", n.Name, n.Version), true
}

func includeNode(n *core.DependencyNode, opts core.FlattenOptionsView) bool {
	if n.Dev && !opts.IncludeDev {
		return false
	}
	return true
}

func toModuleEntry(n *core.DependencyNode) core.ModuleEntry {
	return core.ModuleEntry{
		Name:       n.Name,
		Version:    n.Version,
		Repository: n.Repository,
		Author:     n.Author,
		URL:        n.URL,
		Path:       n.Path,
		Private:    n.Private,
		Dev:        n.Dev,
	}
}

var copyrightPattern = regexp.MustCompile(`(?im)^.*copyright.*$`)

func EnrichLicenseDerivedFields(module core.ModuleEntry, csv bool) core.ModuleEntry {
	if module.LicenseText != "" && csv {
		normalized := strings.ReplaceAll(module.LicenseText, "\r\n", "\n")
		normalized = strings.ReplaceAll(normalized, "\n", "\\n")
		normalized = strings.ReplaceAll(normalized, `"`, `'`)
		module.LicenseText = normalized
	}

	if strings.TrimSpace(module.Copyright) == "" && strings.TrimSpace(module.LicenseText) != "" {
		module.Copyright = extractCopyright(module.LicenseText)
	}

	return module
}

func extractCopyright(licenseText string) string {
	lines := strings.Split(licenseText, "\n")
	out := make([]string, 0)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		if strings.Contains(lower, "copyright notice") || strings.Contains(lower, "copyright and related rights") {
			continue
		}
		if copyrightPattern.MatchString(trimmed) {
			out = append(out, trimmed)
		}
	}
	return strings.Join(out, "\n")
}
