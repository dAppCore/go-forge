package forge

import (
	"context"
	"io"
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

// ReleaseAttachmentUploadOptions controls metadata sent when uploading a release attachment.
type ReleaseAttachmentUploadOptions struct {
	Name        string
	ExternalURL string
}

func releaseAttachmentUploadQuery(opts *ReleaseAttachmentUploadOptions) map[string]string {
	if opts == nil || opts.Name == "" {
		return nil
	}
	query := make(map[string]string, 1)
	if opts.Name != "" {
		query["name"] = opts.Name
	}
	return query
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

// GetLatest returns the most recent non-prerelease, non-draft release.
func (s *ReleaseService) GetLatest(ctx context.Context, owner, repo string) (*types.Release, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/latest", pathParams("owner", owner, "repo", repo))
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

// CreateAttachment uploads a new attachment to a release.
//
// If opts.ExternalURL is set, the upload uses the external_url form field and
// ignores filename/content.
func (s *ReleaseService) CreateAttachment(ctx context.Context, owner, repo string, releaseID int64, opts *ReleaseAttachmentUploadOptions, filename string, content io.Reader) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{releaseID}/assets", pathParams("owner", owner, "repo", repo, "releaseID", int64String(releaseID)))
	fields := make(map[string]string, 1)
	fieldName := "attachment"
	if opts != nil && opts.ExternalURL != "" {
		fields["external_url"] = opts.ExternalURL
		fieldName = ""
		filename = ""
		content = nil
	}
	var out types.Attachment
	if err := s.client.postMultipartJSON(ctx, path, releaseAttachmentUploadQuery(opts), fields, fieldName, filename, content, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateAsset uploads a new asset to a release.
func (s *ReleaseService) CreateAsset(ctx context.Context, owner, repo string, releaseID int64, opts *ReleaseAttachmentUploadOptions, filename string, content io.Reader) (*types.Attachment, error) {
	return s.CreateAttachment(ctx, owner, repo, releaseID, opts, filename, content)
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
