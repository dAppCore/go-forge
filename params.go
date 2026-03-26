package forge

import (
	"net/url"

	core "dappco.re/go/core"
)

// Params maps path variable names to values.
// Example: Params{"owner": "core", "repo": "go-forge"}
type Params map[string]string

// ResolvePath substitutes {placeholders} in path with values from params.
func ResolvePath(path string, params Params) string {
	for k, v := range params {
		path = core.Replace(path, "{"+k+"}", url.PathEscape(v))
	}
	return path
}
