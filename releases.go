package forge

import (
	"context"
	"iter"

	// Note: AX-6 intrinsic — upload APIs must expose the structural request body type; coreio Medium is used inside Client multipart handling.
	goio "io"

	"dappco.re/go/forge/types"
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

// ReleaseListOptions controls filtering for repository release listings.
//
// Usage:
//
//	opts := forge.ReleaseListOptions{Draft: true, Query: "1.0"}
type ReleaseListOptions struct {
	Draft      bool
	PreRelease bool
	Query      string
}

// String returns a safe summary of the release list filters.
func (o ReleaseListOptions) String() string {
	return optionString("forge.ReleaseListOptions",
		"draft", o.Draft,
		"pre-release", o.PreRelease,
		"q", o.Query,
	)
}

// GoString returns a safe Go-syntax summary of the release list filters.
func (o ReleaseListOptions) GoString() string { return o.String() }

func (o ReleaseListOptions) queryParams() map[string]string {
	query := make(map[string]string, 3)
	if o.Draft {
		query["draft"] = "true"
	}
	if o.PreRelease {
		query["pre-release"] = "true"
	}
	if o.Query != "" {
		query["q"] = o.Query
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

// ReleaseAttachmentUploadOptions controls metadata sent when uploading a release attachment.
//
// Usage:
//
//	opts := forge.ReleaseAttachmentUploadOptions{Name: "release.zip"}
type ReleaseAttachmentUploadOptions struct {
	Name        string
	ExternalURL string
}

// String returns a safe summary of the release attachment upload metadata.
func (o ReleaseAttachmentUploadOptions) String() string {
	return optionString("forge.ReleaseAttachmentUploadOptions",
		"name", o.Name,
		"external_url", o.ExternalURL,
	)
}

// GoString returns a safe Go-syntax summary of the release attachment upload metadata.
func (o ReleaseAttachmentUploadOptions) GoString() string { return o.String() }

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

// ListReleasesPage returns a single page of releases in a repository.
func (s *ReleaseService) ListReleasesPage(ctx context.Context, owner, repo string, opts ListOptions, filters ...ReleaseListOptions) (*PagedResult[types.Release], error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases", pathParams("owner", owner, "repo", repo))
	return ListPage[types.Release](ctx, s.client, path, releaseListQuery(filters...), opts)
}

// ListReleases returns all releases in a repository.
func (s *ReleaseService) ListReleases(ctx context.Context, owner, repo string, filters ...ReleaseListOptions) ([]types.Release, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Release](ctx, s.client, path, releaseListQuery(filters...))
}

// IterReleases returns an iterator over all releases in a repository.
func (s *ReleaseService) IterReleases(ctx context.Context, owner, repo string, filters ...ReleaseListOptions) iter.Seq2[types.Release, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Release](ctx, s.client, path, releaseListQuery(filters...))
}

// CreateRelease creates a release in a repository.
func (s *ReleaseService) CreateRelease(ctx context.Context, owner, repo string, opts *types.CreateReleaseOption) (*types.Release, error) {
	var out types.Release
	if err := s.client.Post(ctx, ResolvePath("/api/v1/repos/{owner}/{repo}/releases", pathParams("owner", owner, "repo", repo)), opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
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

// GetRelease returns a release by its tag name.
func (s *ReleaseService) GetRelease(ctx context.Context, owner, repo, tag string) (*types.Release, error) {
	return s.GetByTag(ctx, owner, repo, tag)
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
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{id}/assets", pathParams("owner", owner, "repo", repo, "id", int64String(releaseID)))
	return ListAll[types.Attachment](ctx, s.client, path, nil)
}

// CreateAttachment uploads a new attachment to a release.
//
// If opts.ExternalURL is set, the upload uses the external_url form field and
// ignores filename/content.
func (s *ReleaseService) CreateAttachment(ctx context.Context, owner, repo string, releaseID int64, opts *ReleaseAttachmentUploadOptions, filename string, content goio.Reader) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{id}/assets", pathParams("owner", owner, "repo", repo, "id", int64String(releaseID)))
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

// EditAttachment updates a release attachment.
func (s *ReleaseService) EditAttachment(ctx context.Context, owner, repo string, releaseID, attachmentID int64, opts *types.EditAttachmentOptions) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{id}/assets/{attachment_id}", pathParams("owner", owner, "repo", repo, "id", int64String(releaseID), "attachment_id", int64String(attachmentID)))
	var out types.Attachment
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateAsset uploads a new asset to a release.
func (s *ReleaseService) CreateAsset(ctx context.Context, owner, repo string, releaseID int64, opts *ReleaseAttachmentUploadOptions, filename string, content goio.Reader) (*types.Attachment, error) {
	return s.CreateAttachment(ctx, owner, repo, releaseID, opts, filename, content)
}

// EditAsset updates a release asset.
func (s *ReleaseService) EditAsset(ctx context.Context, owner, repo string, releaseID, attachmentID int64, opts *types.EditAttachmentOptions) (*types.Attachment, error) {
	return s.EditAttachment(ctx, owner, repo, releaseID, attachmentID, opts)
}

// IterAssets returns an iterator over all assets for a release.
func (s *ReleaseService) IterAssets(ctx context.Context, owner, repo string, releaseID int64) iter.Seq2[types.Attachment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/releases/{id}/assets", pathParams("owner", owner, "repo", repo, "id", int64String(releaseID)))
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

func releaseListQuery(filters ...ReleaseListOptions) map[string]string {
	query := make(map[string]string, len(filters))
	for _, filter := range filters {
		for key, value := range filter.queryParams() {
			query[key] = value
		}
	}
	if len(query) == 0 {
		return nil
	}
	return query
}
