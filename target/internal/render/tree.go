package render

import (
	"fmt"
	"strings"
)

func Tree(modules ModuleResultMap, colorize bool) (string, error) {
	var b strings.Builder
	for _, key := range sortedKeys(modules) {
		m := modules[key]
		name := key
		if colorize {
			name = "\x1b[32m" + key + "\x1b[0m"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", name, m.License))
	}
	return b.String(), nil
}
