package forge

import (
	"context"
	"iter"

	core "dappco.re/go/core"
	"dappco.re/go/core/forge/types"
)

// LabelService handles repository labels, organisation labels, and issue labels.
// No Resource embedding — paths are heterogeneous.
type LabelService struct {
	client *Client
}

func newLabelService(c *Client) *LabelService {
	return &LabelService{client: c}
}

// ListRepoLabels returns all labels for a repository.
func (s *LabelService) ListRepoLabels(ctx context.Context, owner, repo string) ([]types.Label, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/labels", owner, repo)
	return ListAll[types.Label](ctx, s.client, path, nil)
}

// IterRepoLabels returns an iterator over all labels for a repository.
func (s *LabelService) IterRepoLabels(ctx context.Context, owner, repo string) iter.Seq2[types.Label, error] {
	path := core.Sprintf("/api/v1/repos/%s/%s/labels", owner, repo)
	return ListIter[types.Label](ctx, s.client, path, nil)
}

// GetRepoLabel returns a single label by ID.
func (s *LabelService) GetRepoLabel(ctx context.Context, owner, repo string, id int64) (*types.Label, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/labels/%d", owner, repo, id)
	var out types.Label
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRepoLabel creates a new label in a repository.
func (s *LabelService) CreateRepoLabel(ctx context.Context, owner, repo string, opts *types.CreateLabelOption) (*types.Label, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/labels", owner, repo)
	var out types.Label
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditRepoLabel updates an existing label in a repository.
func (s *LabelService) EditRepoLabel(ctx context.Context, owner, repo string, id int64, opts *types.EditLabelOption) (*types.Label, error) {
	path := core.Sprintf("/api/v1/repos/%s/%s/labels/%d", owner, repo, id)
	var out types.Label
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRepoLabel deletes a label from a repository.
func (s *LabelService) DeleteRepoLabel(ctx context.Context, owner, repo string, id int64) error {
	path := core.Sprintf("/api/v1/repos/%s/%s/labels/%d", owner, repo, id)
	return s.client.Delete(ctx, path)
}

// ListOrgLabels returns all labels for an organisation.
func (s *LabelService) ListOrgLabels(ctx context.Context, org string) ([]types.Label, error) {
	path := core.Sprintf("/api/v1/orgs/%s/labels", org)
	return ListAll[types.Label](ctx, s.client, path, nil)
}

// IterOrgLabels returns an iterator over all labels for an organisation.
func (s *LabelService) IterOrgLabels(ctx context.Context, org string) iter.Seq2[types.Label, error] {
	path := core.Sprintf("/api/v1/orgs/%s/labels", org)
	return ListIter[types.Label](ctx, s.client, path, nil)
}

// CreateOrgLabel creates a new label in an organisation.
func (s *LabelService) CreateOrgLabel(ctx context.Context, org string, opts *types.CreateLabelOption) (*types.Label, error) {
	path := core.Sprintf("/api/v1/orgs/%s/labels", org)
	var out types.Label
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
