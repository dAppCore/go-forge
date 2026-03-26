package forge

import (
	"net/url"

	core "dappco.re/go/core"
)

// Params maps path variable names to values.
// Example: Params{"owner": "core", "repo": "go-forge"}
//
// Usage:
//
//	params := forge.Params{"owner": "core", "repo": "go-forge"}
//	_ = params
type Params map[string]string

// ResolvePath substitutes {placeholders} in path with values from params.
//
// Usage:
//
//	path := forge.ResolvePath("/api/v1/repos/{owner}/{repo}", forge.Params{"owner": "core", "repo": "go-forge"})
//	_ = path
func ResolvePath(path string, params Params) string {
	for k, v := range params {
		path = core.Replace(path, "{"+k+"}", url.PathEscape(v))
	}
	return path
}
