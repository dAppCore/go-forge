package forge

import (
	"context"
	"iter"

	"dappco.re/go/forge/types"
)

// PackageService handles package registry operations via the Forgejo API.
// No Resource embedding — paths vary by operation.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Packages.List(ctx, "core")
type PackageService struct {
	client *Client
}

func newPackageService(c *Client) *PackageService {
	return &PackageService{client: c}
}

// List returns all packages for a given owner.
func (s *PackageService) List(ctx context.Context, owner string) ([]types.Package, error) {
	path := ResolvePath("/api/v1/packages/{owner}", pathParams("owner", owner))
	return ListAll[types.Package](ctx, s.client, path, nil)
}

// Iter returns an iterator over all packages for a given owner.
func (s *PackageService) Iter(ctx context.Context, owner string) iter.Seq2[types.Package, error] {
	path := ResolvePath("/api/v1/packages/{owner}", pathParams("owner", owner))
	return ListIter[types.Package](ctx, s.client, path, nil)
}

// Get returns a single package by owner, type, name, and version.
func (s *PackageService) Get(ctx context.Context, owner, pkgType, name, version string) (*types.Package, error) {
	path := ResolvePath("/api/v1/packages/{owner}/{type}/{name}/{version}", pathParams("owner", owner, "type", pkgType, "name", name, "version", version))
	var out types.Package
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a package by owner, type, name, and version.
func (s *PackageService) Delete(ctx context.Context, owner, pkgType, name, version string) error {
	path := ResolvePath("/api/v1/packages/{owner}/{type}/{name}/{version}", pathParams("owner", owner, "type", pkgType, "name", name, "version", version))
	return s.client.Delete(ctx, path)
}

// ListFiles returns all files for a specific package version.
func (s *PackageService) ListFiles(ctx context.Context, owner, pkgType, name, version string) ([]types.PackageFile, error) {
	path := ResolvePath("/api/v1/packages/{owner}/{type}/{name}/{version}/files", pathParams("owner", owner, "type", pkgType, "name", name, "version", version))
	return ListAll[types.PackageFile](ctx, s.client, path, nil)
}

// IterFiles returns an iterator over all files for a specific package version.
func (s *PackageService) IterFiles(ctx context.Context, owner, pkgType, name, version string) iter.Seq2[types.PackageFile, error] {
	path := ResolvePath("/api/v1/packages/{owner}/{type}/{name}/{version}/files", pathParams("owner", owner, "type", pkgType, "name", name, "version", version))
	return ListIter[types.PackageFile](ctx, s.client, path, nil)
}
