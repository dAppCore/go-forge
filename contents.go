package forge

import (
	"context"

	core "dappco.re/go/core"
	"dappco.re/go/core/forge/types"
)

// ContentService handles file read/write operations via the Forgejo API.
// No Resource embedding — paths vary by operation.
type ContentService struct {
	client *Client
}

func newContentService(c *Client) *ContentService {
	return &ContentService{client: c}
}

// GetFile returns metadata and content for a file in a repository.
func (s *ContentService) GetFile(ctx context.Context, owner, repo, filepath string) (*types.ContentsResponse, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/contents/%s", owner, repo, filepath)
	var out types.ContentsResponse
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateFile creates a new file in a repository.
func (s *ContentService) CreateFile(ctx context.Context, owner, repo, filepath string, opts *types.CreateFileOptions) (*types.FileResponse, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/contents/%s", owner, repo, filepath)
	var out types.FileResponse
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateFile updates an existing file in a repository.
func (s *ContentService) UpdateFile(ctx context.Context, owner, repo, filepath string, opts *types.UpdateFileOptions) (*types.FileResponse, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/contents/%s", owner, repo, filepath)
	var out types.FileResponse
	if err := s.client.Put(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteFile deletes a file from a repository. Uses DELETE with a JSON body.
func (s *ContentService) DeleteFile(ctx context.Context, owner, repo, filepath string, opts *types.DeleteFileOptions) error {
	path := core.Sprintf("/api/v1/repos/%s/%s/contents/%s", owner, repo, filepath)
	return s.client.DeleteWithBody(ctx, path, opts)
}

// GetRawFile returns the raw file content as bytes.
func (s *ContentService) GetRawFile(ctx context.Context, owner, repo, filepath string) ([]byte, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/raw/%s", owner, repo, filepath)
	return s.client.GetRaw(ctx, path)
}
