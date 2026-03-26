package forge

import (
	"context"
	"iter"

	core "dappco.re/go/core"
	"dappco.re/go/core/forge/types"
)

// CommitService handles commit-related operations such as commit statuses
// and git notes.
// No Resource embedding — collection and item commit paths differ, and the
// remaining endpoints are heterogeneous across status and note paths.
type CommitService struct {
	client *Client
}

const (
	commitCollectionPath = "/api/v1/repos/{owner}/{repo}/commits"
	commitItemPath       = "/api/v1/repos/{owner}/{repo}/git/commits/{sha}"
)

func newCommitService(c *Client) *CommitService {
	return &CommitService{client: c}
}

// List returns a single page of commits for a repository.
func (s *CommitService) List(ctx context.Context, params Params, opts ListOptions) (*PagedResult[types.Commit], error) {
	return ListPage[types.Commit](ctx, s.client, ResolvePath(commitCollectionPath, params), nil, opts)
}

// ListAll returns all commits for a repository.
func (s *CommitService) ListAll(ctx context.Context, params Params) ([]types.Commit, error) {
	return ListAll[types.Commit](ctx, s.client, ResolvePath(commitCollectionPath, params), nil)
}

// Iter returns an iterator over all commits for a repository.
func (s *CommitService) Iter(ctx context.Context, params Params) iter.Seq2[types.Commit, error] {
	return ListIter[types.Commit](ctx, s.client, ResolvePath(commitCollectionPath, params), nil)
}

// Get returns a single commit by SHA or ref.
func (s *CommitService) Get(ctx context.Context, params Params) (*types.Commit, error) {
	var out types.Commit
	if err := s.client.Get(ctx, ResolvePath(commitItemPath, params), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCombinedStatus returns the combined status for a given ref (branch, tag, or SHA).
func (s *CommitService) GetCombinedStatus(ctx context.Context, owner, repo, ref string) (*types.CombinedStatus, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/statuses/%s", owner, repo, ref)
	var out types.CombinedStatus
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListStatuses returns all commit statuses for a given ref.
func (s *CommitService) ListStatuses(ctx context.Context, owner, repo, ref string) ([]types.CommitStatus, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/commits/%s/statuses", owner, repo, ref)
	var out []types.CommitStatus
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateStatus creates a new commit status for the given SHA.
func (s *CommitService) CreateStatus(ctx context.Context, owner, repo, sha string, opts *types.CreateStatusOption) (*types.CommitStatus, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/statuses/%s", owner, repo, sha)
	var out types.CommitStatus
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNote returns the git note for a given commit SHA.
func (s *CommitService) GetNote(ctx context.Context, owner, repo, sha string) (*types.Note, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/git/notes/%s", owner, repo, sha)
	var out types.Note
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
