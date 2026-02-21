package forge

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go-forge/types"
)

// WikiService handles wiki page operations for a repository.
// No Resource embedding — custom endpoints for wiki CRUD.
type WikiService struct {
	client *Client
}

func newWikiService(c *Client) *WikiService {
	return &WikiService{client: c}
}

// ListPages returns all wiki page metadata for a repository.
func (s *WikiService) ListPages(ctx context.Context, owner, repo string) ([]types.WikiPageMetaData, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/pages", owner, repo)
	var out []types.WikiPageMetaData
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetPage returns a single wiki page by name.
func (s *WikiService) GetPage(ctx context.Context, owner, repo, pageName string) (*types.WikiPage, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, pageName)
	var out types.WikiPage
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreatePage creates a new wiki page.
func (s *WikiService) CreatePage(ctx context.Context, owner, repo string, opts *types.CreateWikiPageOptions) (*types.WikiPage, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/new", owner, repo)
	var out types.WikiPage
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditPage updates an existing wiki page.
func (s *WikiService) EditPage(ctx context.Context, owner, repo, pageName string, opts *types.CreateWikiPageOptions) (*types.WikiPage, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, pageName)
	var out types.WikiPage
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeletePage removes a wiki page.
func (s *WikiService) DeletePage(ctx context.Context, owner, repo, pageName string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/wiki/page/%s", owner, repo, pageName)
	return s.client.Delete(ctx, path)
}
