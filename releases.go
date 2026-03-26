package forge

import (
	"context"
	"iter"

	core "dappco.re/go/core"
	"dappco.re/go/core/forge/types"
)

// ReleaseService handles release operations within a repository.
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
	path := core.Sprintf("/api/v1/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	var out types.Release
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteByTag deletes a release by its tag name.
func (s *ReleaseService) DeleteByTag(ctx context.Context, owner, repo, tag string) error {
	path := core.Sprintf("/api/v1/repos/%s/%s/releases/tags/%s", owner, repo, tag)
	return s.client.Delete(ctx, path)
}

// ListAssets returns all assets for a release.
func (s *ReleaseService) ListAssets(ctx context.Context, owner, repo string, releaseID int64) ([]types.Attachment, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/releases/%d/assets", owner, repo, releaseID)
	return ListAll[types.Attachment](ctx, s.client, path, nil)
}

// IterAssets returns an iterator over all assets for a release.
func (s *ReleaseService) IterAssets(ctx context.Context, owner, repo string, releaseID int64) iter.Seq2[types.Attachment, error] {
	path := core.Sprintf("/api/v1/repos/%s/%s/releases/%d/assets", owner, repo, releaseID)
	return ListIter[types.Attachment](ctx, s.client, path, nil)
}

// GetAsset returns a single asset for a release.
func (s *ReleaseService) GetAsset(ctx context.Context, owner, repo string, releaseID, assetID int64) (*types.Attachment, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/releases/%d/assets/%d", owner, repo, releaseID, assetID)
	var out types.Attachment
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAsset deletes a single asset from a release.
func (s *ReleaseService) DeleteAsset(ctx context.Context, owner, repo string, releaseID, assetID int64) error {
	path := core.Sprintf("/api/v1/repos/%s/%s/releases/%d/assets/%d", owner, repo, releaseID, assetID)
	return s.client.Delete(ctx, path)
}
