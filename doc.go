// Package forge provides a full-coverage Go client for the Forgejo API.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "your-token")
//	repos, err := f.Repos.List(ctx, forge.Params{"org": "core"}, forge.DefaultList)
//
// Types are generated from Forgejo's swagger.v1.json spec via cmd/forgegen/.
// Run `go generate ./types/...` to regenerate after a Forgejo upgrade.
package forge
