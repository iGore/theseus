package npmgraph

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type Node struct {
	Name, Version, Description, Readme, Path string
	Private, Extraneous, Root                bool
	License                                  any `json:"license"`
	Licenses                                 any `json:"licenses"`
	Repository                               struct {
		URL string `json:"url"`
	} `json:"repository"`
	URL struct {
		Web string `json:"web"`
	} `json:"url"`
	Author struct {
		Name, Email, URL string
	}
	Dependencies []*Node
	Fields       map[string]any
}
type Options struct {
	Dev   bool
	Depth int
}

type Scanner struct{}

func (Scanner) Scan(root string, opts Options) (*Node, error) {
	return readPackage(root, true, opts.Depth)
}

func readPackage(dir string, root bool, depth int) (*Node, error) {
	b, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil, err
	}
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	n := &Node{Root: root, Path: dir, Fields: raw}
	n.Name, _ = raw["name"].(string)
	n.Version, _ = raw["version"].(string)
	n.Description, _ = raw["description"].(string)
	n.Readme, _ = raw["readme"].(string)
	n.Private, _ = raw["private"].(bool)
	n.License = raw["license"]
	n.Licenses = raw["licenses"]
	if r, ok := raw["repository"].(map[string]any); ok {
		n.Repository.URL, _ = r["url"].(string)
	}
	if u, ok := raw["url"].(map[string]any); ok {
		n.URL.Web, _ = u["web"].(string)
	}
	if a, ok := raw["author"].(map[string]any); ok {
		n.Author.Name, _ = a["name"].(string)
		n.Author.Email, _ = a["email"].(string)
		n.Author.URL, _ = a["url"].(string)
	}
	if n.Readme == "" {
		if rb, err := os.ReadFile(filepath.Join(dir, "README.md")); err == nil {
			n.Readme = string(rb)
		}
	}
	if depth == 0 {
		return n, nil
	}
	entries, err := os.ReadDir(filepath.Join(dir, "node_modules"))
	if err != nil {
		return n, nil
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	nextDepth := depth
	if nextDepth > 0 {
		nextDepth--
	}
	for _, name := range names {
		p := filepath.Join(dir, "node_modules", name)
		if name != "" && name[0] == '@' {
			scoped, _ := os.ReadDir(p)
			for _, se := range scoped {
				if child, err := readPackage(filepath.Join(p, se.Name()), false, nextDepth); err == nil {
					n.Dependencies = append(n.Dependencies, child)
				}
			}
			continue
		}
		if child, err := readPackage(p, false, nextDepth); err == nil {
			n.Dependencies = append(n.Dependencies, child)
		}
	}
	return n, nil
}
