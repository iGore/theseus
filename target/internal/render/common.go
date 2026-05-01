package render

import "sort"

func sortedKeys(modules ModuleResultMap) []string {
	keys := make([]string, 0, len(modules))
	for k := range modules {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
