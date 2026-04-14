package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// LabelService handles repository labels, organisation labels, and issue labels.
// No Resource embedding — paths are heterogeneous.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Labels.ListRepoLabels(ctx, "core", "go-forge")
type LabelService struct {
	client *Client
}

func newLabelService(c *Client) *LabelService {
	return &LabelService{client: c}
}

// ListRepoLabels returns all labels for a repository.
func (s *LabelService) ListRepoLabels(ctx context.Context, owner, repo string) ([]types.Label, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/labels", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Label](ctx, s.client, path, nil)
}

// ListLabels returns all labels for a repository.
func (s *LabelService) ListLabels(ctx context.Context, owner, repo string) ([]types.Label, error) {
	return s.ListRepoLabels(ctx, owner, repo)
}

// IterRepoLabels returns an iterator over all labels for a repository.
func (s *LabelService) IterRepoLabels(ctx context.Context, owner, repo string) iter.Seq2[types.Label, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/labels", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Label](ctx, s.client, path, nil)
}

// IterLabels returns an iterator over all labels for a repository.
func (s *LabelService) IterLabels(ctx context.Context, owner, repo string) iter.Seq2[types.Label, error] {
	return s.IterRepoLabels(ctx, owner, repo)
}

// GetRepoLabel returns a single label by ID.
func (s *LabelService) GetRepoLabel(ctx context.Context, owner, repo string, id int64) (*types.Label, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/labels/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.Label
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateRepoLabel creates a new label in a repository.
func (s *LabelService) CreateRepoLabel(ctx context.Context, owner, repo string, opts *types.CreateLabelOption) (*types.Label, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/labels", pathParams("owner", owner, "repo", repo))
	var out types.Label
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditRepoLabel updates an existing label in a repository.
func (s *LabelService) EditRepoLabel(ctx context.Context, owner, repo string, id int64, opts *types.EditLabelOption) (*types.Label, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/labels/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.Label
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRepoLabel deletes a label from a repository.
func (s *LabelService) DeleteRepoLabel(ctx context.Context, owner, repo string, id int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/labels/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListOrgLabels returns all labels for an organisation.
func (s *LabelService) ListOrgLabels(ctx context.Context, org string) ([]types.Label, error) {
	path := ResolvePath("/api/v1/orgs/{org}/labels", pathParams("org", org))
	return ListAll[types.Label](ctx, s.client, path, nil)
}

// IterOrgLabels returns an iterator over all labels for an organisation.
func (s *LabelService) IterOrgLabels(ctx context.Context, org string) iter.Seq2[types.Label, error] {
	path := ResolvePath("/api/v1/orgs/{org}/labels", pathParams("org", org))
	return ListIter[types.Label](ctx, s.client, path, nil)
}

// CreateOrgLabel creates a new label in an organisation.
func (s *LabelService) CreateOrgLabel(ctx context.Context, org string, opts *types.CreateLabelOption) (*types.Label, error) {
	path := ResolvePath("/api/v1/orgs/{org}/labels", pathParams("org", org))
	var out types.Label
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOrgLabel returns a single label for an organisation.
func (s *LabelService) GetOrgLabel(ctx context.Context, org string, id int64) (*types.Label, error) {
	path := ResolvePath("/api/v1/orgs/{org}/labels/{id}", pathParams("org", org, "id", int64String(id)))
	var out types.Label
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditOrgLabel updates an existing label in an organisation.
func (s *LabelService) EditOrgLabel(ctx context.Context, org string, id int64, opts *types.EditLabelOption) (*types.Label, error) {
	path := ResolvePath("/api/v1/orgs/{org}/labels/{id}", pathParams("org", org, "id", int64String(id)))
	var out types.Label
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteOrgLabel deletes a label from an organisation.
func (s *LabelService) DeleteOrgLabel(ctx context.Context, org string, id int64) error {
	path := ResolvePath("/api/v1/orgs/{org}/labels/{id}", pathParams("org", org, "id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListLabelTemplates returns all available label template names.
func (s *LabelService) ListLabelTemplates(ctx context.Context) ([]string, error) {
	var out []string
	if err := s.client.Get(ctx, "/api/v1/label/templates", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IterLabelTemplates returns an iterator over all available label template names.
func (s *LabelService) IterLabelTemplates(ctx context.Context) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		items, err := s.ListLabelTemplates(ctx)
		if err != nil {
			yield("", err)
			return
		}
		for _, item := range items {
			if !yield(item, nil) {
				return
			}
		}
	}
}

// GetLabelTemplate returns all labels for a label template.
func (s *LabelService) GetLabelTemplate(ctx context.Context, name string) ([]types.LabelTemplate, error) {
	path := ResolvePath("/api/v1/label/templates/{name}", pathParams("name", name))
	var out []types.LabelTemplate
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}
