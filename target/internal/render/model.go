package render

// ModuleResult is render-local to avoid package cycles.
type ModuleResult struct {
	Name        string
	Version     string
	License     string
	LicenseFile string
}

type ModuleResultMap map[string]ModuleResult
