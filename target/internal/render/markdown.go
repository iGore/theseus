package render

import "strings"

func Markdown(modules ModuleResultMap) (string, error) {
	var b strings.Builder
	b.WriteString("| key | license |\n")
	b.WriteString("| --- | --- |\n")
	for _, key := range sortedKeys(modules) {
		m := modules[key]
		b.WriteString("| " + key + " | " + m.License + " |\n")
	}
	return b.String(), nil
}
