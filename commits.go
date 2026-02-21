package forge

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go-forge/types"
)

// CommitService handles commit-related operations such as commit statuses
// and git notes.
// No Resource embedding — heterogeneous endpoints across status and note paths.
type CommitService struct {
	client *Client
}

func newCommitService(c *Client) *CommitService {
	return &CommitService{client: c}
}

// GetCombinedStatus returns the combined status for a given ref (branch, tag, or SHA).
func (s *CommitService) GetCombinedStatus(ctx context.Context, owner, repo, ref string) (*types.CombinedStatus, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/statuses/%s", owner, repo, ref)
	var out types.CombinedStatus
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListStatuses returns all commit statuses for a given ref.
func (s *CommitService) ListStatuses(ctx context.Context, owner, repo, ref string) ([]types.CommitStatus, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/commits/%s/statuses", owner, repo, ref)
	var out []types.CommitStatus
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateStatus creates a new commit status for the given SHA.
func (s *CommitService) CreateStatus(ctx context.Context, owner, repo, sha string, opts *types.CreateStatusOption) (*types.CommitStatus, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/statuses/%s", owner, repo, sha)
	var out types.CommitStatus
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNote returns the git note for a given commit SHA.
func (s *CommitService) GetNote(ctx context.Context, owner, repo, sha string) (*types.Note, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/git/notes/%s", owner, repo, sha)
	var out types.Note
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
