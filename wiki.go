package forge

import (
	"context"
	"iter"
	"strconv"

	"dappco.re/go/forge/types"
)

// WikiService handles wiki page operations for a repository.
// No Resource embedding — custom endpoints for wiki CRUD.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Wiki.ListPages(ctx, "core", "go-forge")
type WikiService struct {
	client *Client
}

func newWikiService(c *Client) *WikiService {
	return &WikiService{client: c}
}

// ListPages returns all wiki page metadata for a repository.
func (s *WikiService) ListPages(ctx context.Context, owner, repo string) ([]types.WikiPageMetaData, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/wiki/pages", pathParams("owner", owner, "repo", repo))
	var out []types.WikiPageMetaData
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IterPages returns an iterator over all wiki page metadata for a repository.
func (s *WikiService) IterPages(ctx context.Context, owner, repo string) iter.Seq2[types.WikiPageMetaData, error] {
	return func(yield func(types.WikiPageMetaData, error) bool) {
		items, err := s.ListPages(ctx, owner, repo)
		if err != nil {
			yield(*new(types.WikiPageMetaData), err)
			return
		}
		for _, item := range items {
			if !yield(item, nil) {
				return
			}
		}
	}
}

// GetPage returns a single wiki page by name.
func (s *WikiService) GetPage(ctx context.Context, owner, repo, pageName string) (*types.WikiPage, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/wiki/page/{pageName}", pathParams("owner", owner, "repo", repo, "pageName", pageName))
	var out types.WikiPage
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPageRevisions returns the revision history for a wiki page.
// Page is optional; pass a value greater than zero to request a specific page of results.
func (s *WikiService) GetPageRevisions(ctx context.Context, owner, repo, pageName string, page int) (*types.WikiCommitList, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/wiki/revisions/{pageName}", pathParams("owner", owner, "repo", repo, "pageName", pageName))
	if page > 0 {
		path += "?page=" + strconv.Itoa(page)
	}
	var out types.WikiCommitList
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreatePage creates a new wiki page.
func (s *WikiService) CreatePage(ctx context.Context, owner, repo string, opts *types.CreateWikiPageOptions) (*types.WikiPage, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/wiki/new", pathParams("owner", owner, "repo", repo))
	var out types.WikiPage
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditPage updates an existing wiki page.
func (s *WikiService) EditPage(ctx context.Context, owner, repo, pageName string, opts *types.CreateWikiPageOptions) (*types.WikiPage, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/wiki/page/{pageName}", pathParams("owner", owner, "repo", repo, "pageName", pageName))
	var out types.WikiPage
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeletePage removes a wiki page.
func (s *WikiService) DeletePage(ctx context.Context, owner, repo, pageName string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/wiki/page/{pageName}", pathParams("owner", owner, "repo", repo, "pageName", pageName))
	return s.client.Delete(ctx, path)
}
