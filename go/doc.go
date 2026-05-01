// Package forge provides a full-coverage Go client for the Forgejo API.
//
// Usage:
//
//	ctx := context.Background()
//	f := forge.NewForge("https://forge.lthn.ai", "your-token")
//	repos, err := f.Repos.ListOrgRepos(ctx, "core")
//
// Types are generated from Forgejo's swagger.v1.json spec via cmd/forgegen/.
// Run `go generate ./types/...` to regenerate after a Forgejo upgrade.
package forge
