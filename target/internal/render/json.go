package render

import "encoding/json"

func JSON(modules ModuleResultMap) (string, error) {
	type item struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Version     string `json:"version"`
		License     string `json:"license"`
		LicenseFile string `json:"licenseFile,omitempty"`
	}
	out := struct {
		Modules []item `json:"modules"`
	}{Modules: make([]item, 0, len(modules))}

	for _, key := range sortedKeys(modules) {
		m := modules[key]
		out.Modules = append(out.Modules, item{Key: key, Name: m.Name, Version: m.Version, License: m.License, LicenseFile: m.LicenseFile})
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}
