package render

import (
	"fmt"
	"strings"
)

func Summary(modules ModuleResultMap) (string, error) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("total: %d\n", len(modules)))
	for _, key := range sortedKeys(modules) {
		b.WriteString(key + "\n")
	}
	return b.String(), nil
}
