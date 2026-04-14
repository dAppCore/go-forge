package forge

import (
	"context"
	"iter"
	"net/url"

	core "dappco.re/go/core"
	"dappco.re/go/core/forge/types"
)

// ContentService handles file read/write operations via the Forgejo API.
// No Resource embedding — paths vary by operation.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Contents.GetFile(ctx, "core", "go-forge", "README.md")
type ContentService struct {
	client *Client
}

func newContentService(c *Client) *ContentService {
	return &ContentService{client: c}
}

// ListContents returns the entries in a repository directory.
// If ref is non-empty, the listing is resolved against that branch, tag, or commit.
func (s *ContentService) ListContents(ctx context.Context, owner, repo, ref string) ([]types.ContentsResponse, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/contents", pathParams("owner", owner, "repo", repo))
	if ref != "" {
		u, err := url.Parse(path)
		if err != nil {
			return nil, core.E("ContentService.ListContents", "forge: parse path", err)
		}
		q := u.Query()
		q.Set("ref", ref)
		u.RawQuery = q.Encode()
		path = u.String()
	}

	var out []types.ContentsResponse
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IterContents returns an iterator over the entries in a repository directory.
// If ref is non-empty, the listing is resolved against that branch, tag, or commit.
func (s *ContentService) IterContents(ctx context.Context, owner, repo, ref string) iter.Seq2[types.ContentsResponse, error] {
	return func(yield func(types.ContentsResponse, error) bool) {
		items, err := s.ListContents(ctx, owner, repo, ref)
		if err != nil {
			yield(*new(types.ContentsResponse), err)
			return
		}
		for _, item := range items {
			if !yield(item, nil) {
				return
			}
		}
	}
}

// GetFile returns metadata and content for a file in a repository.
func (s *ContentService) GetFile(ctx context.Context, owner, repo, filepath string) (*types.ContentsResponse, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/contents/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	var out types.ContentsResponse
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetContents returns metadata and content for a file in a repository.
func (s *ContentService) GetContents(ctx context.Context, owner, repo, filepath string) (*types.ContentsResponse, error) {
	return s.GetFile(ctx, owner, repo, filepath)
}

// CreateFile creates a new file in a repository.
func (s *ContentService) CreateFile(ctx context.Context, owner, repo, filepath string, opts *types.CreateFileOptions) (*types.FileResponse, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/contents/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	var out types.FileResponse
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateFile updates an existing file in a repository.
func (s *ContentService) UpdateFile(ctx context.Context, owner, repo, filepath string, opts *types.UpdateFileOptions) (*types.FileResponse, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/contents/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	var out types.FileResponse
	if err := s.client.Put(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteFile deletes a file from a repository. Uses DELETE with a JSON body.
func (s *ContentService) DeleteFile(ctx context.Context, owner, repo, filepath string, opts *types.DeleteFileOptions) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/contents/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	return s.client.DeleteWithBody(ctx, path, opts)
}

// GetRawFile returns the raw file content as bytes.
func (s *ContentService) GetRawFile(ctx context.Context, owner, repo, filepath string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/raw/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	return s.client.GetRaw(ctx, path)
}
