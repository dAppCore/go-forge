package forge

import (
	"syscall"

	core "dappco.re/go/core"
)

const (
	// DefaultURL is the fallback Forgejo instance URL when neither flag nor
	// environment variable is set.
	//
	// Usage:
	//  cfgURL, _, _ := forge.ResolveConfig("", "")
	//  _ = cfgURL == forge.DefaultURL
	DefaultURL = "http://localhost:3000"
)

// ResolveConfig resolves the Forgejo URL and API token from flags, environment
// variables, and built-in defaults. Priority order: flags > env > defaults.
//
// Environment variables:
//   - FORGE_URL   — base URL of the Forgejo instance
//   - FORGE_TOKEN — API token for authentication
//
// Usage:
//
//	url, token, err := forge.ResolveConfig("", "")
//	_ = url
//	_ = token
func ResolveConfig(flagURL, flagToken string) (url, token string, err error) {
	url, _ = syscall.Getenv("FORGE_URL")
	token, _ = syscall.Getenv("FORGE_TOKEN")

	if flagURL != "" {
		url = flagURL
	}
	if flagToken != "" {
		token = flagToken
	}
	if url == "" {
		url = DefaultURL
	}
	return url, token, nil
}

// NewForgeFromConfig creates a new Forge client using resolved configuration.
// It returns an error if no API token is available from flags or environment.
//
// Usage:
//
//	f, err := forge.NewForgeFromConfig("", "")
//	_ = f
func NewForgeFromConfig(flagURL, flagToken string, opts ...Option) (*Forge, error) {
	url, token, err := ResolveConfig(flagURL, flagToken)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return nil, core.E("NewForgeFromConfig", "forge: no API token configured (set FORGE_TOKEN or pass --token)", nil)
	}
	return NewForge(url, token, opts...), nil
}
