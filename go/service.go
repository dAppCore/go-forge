// SPDX-License-Identifier: EUPL-1.2

// Service registration for the forge package — exposes the canonical
// `NewService(opts)` + `Register(c)` shape per Mantis #1336, wrapping
// the existing NewForge(url, token, opts...) constructor in a
// Core-registerable factory.
//
//	c, _ := core.New(
//	    core.WithService(forge.NewService(forge.ServiceOptions{
//	        URL:   "https://forge.lthn.sh",
//	        Token: bearerToken,
//	    })),
//	)
//	svc := core.MustServiceFor[*forge.Service](c, "forge")
//	repo, err := svc.Forge.Repos.Get(ctx, forge.Params{"owner": "core", "repo": "go-scm"})
//
// The *Forge type and the NewForge constructor (with all 20+ sub-
// services like Repos, Issues, Pulls, etc.) remain the source of truth
// — Service is a thin Core-side handle that gives the package a
// registerable identity the rest of the framework can discover via
// core.ServiceFor.

package forge

import (
	core "dappco.re/go"
)

// ServiceOptions configures the forge service. URL + Token are required
// for runtime use — empty fields register a Service with a nil Forge,
// suitable for late wiring.
//
// Named ServiceOptions (not Options) to avoid colliding with the
// existing `type Option` from the Forge constructor's option-funcs API.
type ServiceOptions struct {
	// URL is the Forgejo / Gitea API base URL.
	// Empty → no Forge wired; svc.Forge stays nil for late injection.
	URL string
	// Token is the bearer token for API auth. Empty + URL set → still
	// constructs a Forge but unauthenticated calls will fail.
	Token string
	// Options are the full forge.Option func slice (WithHTTPClient,
	// WithRetry, etc.) passed straight through to NewForge.
	Options []Option
}

// Service is the registerable handle for the forge package — embeds
// *core.ServiceRuntime[ServiceOptions] for typed options access and
// holds a constructed *Forge ready for the full sub-service surface
// (Repos, Issues, Pulls, etc.).
//
// Usage example: `svc := core.MustServiceFor[*forge.Service](c, "forge"); _ = svc.Forge.Repos.Get(ctx, params)`
type Service struct {
	*core.ServiceRuntime[ServiceOptions]
	// Forge is the live *Forge the service was constructed with.
	// nil if NewService was called with empty URL — callers wire one
	// later via NewForge(...) and assign to svc.Forge.
	Forge *Forge
}

// NewService returns a factory that constructs a *Forge from the
// supplied options and wraps it as a Core-registerable *Service. Empty
// URL registers a Service with a nil Forge for late wiring.
//
//	core.WithService(forge.NewService(forge.ServiceOptions{
//	    URL:   "https://forge.lthn.sh",
//	    Token: bearerToken,
//	}))
func NewService(opts ServiceOptions) func(*core.Core) core.Result {
	return func(c *core.Core) core.Result {
		svc := &Service{
			ServiceRuntime: core.NewServiceRuntime(c, opts),
		}
		if opts.URL == "" {
			return core.Ok(svc)
		}
		svc.Forge = NewForge(opts.URL, opts.Token, opts.Options...)
		return core.Ok(svc)
	}
}

// Register wires the forge service into the Core with empty
// ServiceOptions — the imperative-style alternative to NewService.
// The resulting *Service holds a nil Forge, so consumers must wire
// one via NewForge(...) and assign to svc.Forge before use.
//
//	c := core.New()
//	if r := forge.Register(c); !r.OK { return r }
//	svc := core.MustServiceFor[*forge.Service](c, "forge")
//	svc.Forge = forge.NewForge("https://forge.lthn.sh", token)
func Register(c *core.Core) core.Result {
	return NewService(ServiceOptions{})(c)
}
