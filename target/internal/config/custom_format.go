package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"theseus/target/internal/core"
)

type CustomFormatConfig map[string]any

// ParsedCustomConfigResult preserves legacy compatibility semantics:
// valid object on success, Error object value on failure.
type ParsedCustomConfigResult struct {
	Config CustomFormatConfig
	Err    error
}

func ParseJSON(path any) any {
	p, ok := path.(string)
	if !ok || strings.TrimSpace(p) == "" {
		return errors.New("did not specify a path")
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return err
	}
	return CustomFormatConfig(out)
}

func LoadCustomFormat(inline CustomFormatConfig, customPath string) ParsedCustomConfigResult {
	if strings.TrimSpace(customPath) == "" {
		if inline == nil {
			return ParsedCustomConfigResult{Config: CustomFormatConfig{}}
		}
		return ParsedCustomConfigResult{Config: inline}
	}

	parsed := ParseJSON(customPath)
	if err, ok := parsed.(error); ok {
		return ParsedCustomConfigResult{Config: nil, Err: err}
	}
	cfg, _ := parsed.(CustomFormatConfig)
	if cfg == nil {
		cfg = CustomFormatConfig{}
	}
	return ParsedCustomConfigResult{Config: cfg}
}

func ApplyCustomFields(module core.ModuleEntry, cfg CustomFormatConfig, csv bool) map[string]string {
	out := map[string]string{}
	for key, defaultValue := range cfg {
		if b, ok := defaultValue.(bool); ok && !b {
			continue
		}

		if key == "licenseText" {
			if module.LicenseText != "" {
				if csv {
					out[key] = normalizeCSVLicenseText(module.LicenseText)
				} else {
					out[key] = module.LicenseText
				}
				continue
			}
		}

		if key == "copyright" {
			if module.Copyright != "" {
				out[key] = module.Copyright
				continue
			}
		}

		if v, ok := moduleStringValue(module, key); ok && strings.TrimSpace(v) != "" {
			out[key] = v
			continue
		}

		out[key] = stringifyDefault(defaultValue)
	}
	return out
}

func moduleStringValue(module core.ModuleEntry, key string) (string, bool) {
	switch key {
	case "name":
		return module.Name, true
	case "version":
		return module.Version, true
	case "repository":
		return module.Repository, true
	case "author":
		return module.Author, true
	case "url":
		return module.URL, true
	case "path":
		return module.Path, true
	case "licenses":
		return module.Licenses, true
	case "licenseFile":
		return module.LicenseFile, true
	default:
		return "", false
	}
}

func stringifyDefault(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func normalizeCSVLicenseText(in string) string {
	replaced := strings.ReplaceAll(in, "\r\n", "\n")
	replaced = strings.ReplaceAll(replaced, "\n", "\\n")
	replaced = strings.ReplaceAll(replaced, "\"", "'")
	return replaced
}
