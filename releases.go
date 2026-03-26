package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// ReleaseService handles release operations within a repository.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Releases.ListAll(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
type ReleaseService struct {
	Resource[types.Release, types.CreateReleaseOption, types.EditReleaseOption]
}

func newReleaseService(c *Client) *ReleaseService {
	return &ReleaseService{
		Resource: *NewResource[types.Release, types.CreateReleaseOption, types.EditReleaseOption](
			c, "/api/v1/repos/{owner}/{repo}/releases/{id}",
		),
	}
}

// GetByTag returns a release by its tag name.
func (s *ReleaseService) GetByTag(ctx context.Context, owner, repo, tag string) (*types.Release, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/tags/{tag}", pathParams("owner", owner, "repo", repo, "tag", tag))
	var out types.Release
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteByTag deletes a release by its tag name.
func (s *ReleaseService) DeleteByTag(ctx context.Context, owner, repo, tag string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/tags/{tag}", pathParams("owner", owner, "repo", repo, "tag", tag))
	return s.client.Delete(ctx, path)
}

// ListAssets returns all assets for a release.
func (s *ReleaseService) ListAssets(ctx context.Context, owner, repo string, releaseID int64) ([]types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{releaseID}/assets", pathParams("owner", owner, "repo", repo, "releaseID", int64String(releaseID)))
	return ListAll[types.Attachment](ctx, s.client, path, nil)
}

// IterAssets returns an iterator over all assets for a release.
func (s *ReleaseService) IterAssets(ctx context.Context, owner, repo string, releaseID int64) iter.Seq2[types.Attachment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{releaseID}/assets", pathParams("owner", owner, "repo", repo, "releaseID", int64String(releaseID)))
	return ListIter[types.Attachment](ctx, s.client, path, nil)
}

// GetAsset returns a single asset for a release.
func (s *ReleaseService) GetAsset(ctx context.Context, owner, repo string, releaseID, assetID int64) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{releaseID}/assets/{assetID}", pathParams("owner", owner, "repo", repo, "releaseID", int64String(releaseID), "assetID", int64String(assetID)))
	var out types.Attachment
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAsset deletes a single asset from a release.
func (s *ReleaseService) DeleteAsset(ctx context.Context, owner, repo string, releaseID, assetID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{releaseID}/assets/{assetID}", pathParams("owner", owner, "repo", repo, "releaseID", int64String(releaseID), "assetID", int64String(assetID)))
	return s.client.Delete(ctx, path)
}
