package forge

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	core "dappco.re/go"
)

// Params maps path variable names to values.
// Example: Params{"owner": "core", "repo": "go-forge"}
//
// Usage:
//
//	params := forge.Params{"owner": "core", "repo": "go-forge"}
//	_ = params
type Params map[string]string

// String returns a safe summary of the path parameters.
//
// Usage:
//
//	_ = forge.Params{"owner": "core"}.String()
func (p Params) String() string {
	if p == nil {
		return "forge.Params{<nil>}"
	}

	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("forge.Params{")
	for i, k := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(strconv.Quote(p[k]))
	}
	b.WriteString("}")
	return b.String()
}

// GoString returns a safe Go-syntax summary of the path parameters.
//
// Usage:
//
//	_ = fmt.Sprintf("%#v", forge.Params{"owner": "core"})
func (p Params) GoString() string { return p.String() }

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
