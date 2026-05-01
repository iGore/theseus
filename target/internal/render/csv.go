package render

import (
	"bytes"
	"encoding/csv"
)

func CSV(modules ModuleResultMap) (string, error) {
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	if err := w.Write([]string{"key", "name", "version", "license", "licenseFile"}); err != nil {
		return "", err
	}
	for _, key := range sortedKeys(modules) {
		m := modules[key]
		if err := w.Write([]string{key, m.Name, m.Version, m.License, m.LicenseFile}); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return b.String(), nil
}
