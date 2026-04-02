package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// WebhookService handles webhook (hook) operations within a repository.
// Embeds Resource for standard CRUD on /api/v1/repos/{owner}/{repo}/hooks/{id}.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Webhooks.ListAll(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
type WebhookService struct {
	Resource[types.Hook, types.CreateHookOption, types.EditHookOption]
}

func newWebhookService(c *Client) *WebhookService {
	return &WebhookService{
		Resource: *NewResource[types.Hook, types.CreateHookOption, types.EditHookOption](
			c, "/api/v1/repos/{owner}/{repo}/hooks/{id}",
		),
	}
}

// ListHooks returns all webhooks for a repository.
func (s *WebhookService) ListHooks(ctx context.Context, owner, repo string) ([]types.Hook, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/hooks", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Hook](ctx, s.client, path, nil)
}

// IterHooks returns an iterator over all webhooks for a repository.
func (s *WebhookService) IterHooks(ctx context.Context, owner, repo string) iter.Seq2[types.Hook, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/hooks", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Hook](ctx, s.client, path, nil)
}

// CreateHook creates a webhook for a repository.
func (s *WebhookService) CreateHook(ctx context.Context, owner, repo string, opts *types.CreateHookOption) (*types.Hook, error) {
	var out types.Hook
	if err := s.client.Post(ctx, ResolvePath("/api/v1/repos/{owner}/{repo}/hooks", pathParams("owner", owner, "repo", repo)), opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// TestHook triggers a test delivery for a webhook.
func (s *WebhookService) TestHook(ctx context.Context, owner, repo string, id int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/hooks/{id}/tests", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return s.client.Post(ctx, path, nil, nil)
}

// ListGitHooks returns all Git hooks for a repository.
func (s *WebhookService) ListGitHooks(ctx context.Context, owner, repo string) ([]types.GitHook, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/hooks/git", pathParams("owner", owner, "repo", repo))
	return ListAll[types.GitHook](ctx, s.client, path, nil)
}

// GetGitHook returns a single Git hook for a repository.
func (s *WebhookService) GetGitHook(ctx context.Context, owner, repo, id string) (*types.GitHook, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/hooks/git/{id}", pathParams("owner", owner, "repo", repo, "id", id))
	var out types.GitHook
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditGitHook updates an existing Git hook in a repository.
func (s *WebhookService) EditGitHook(ctx context.Context, owner, repo, id string, opts *types.EditGitHookOption) (*types.GitHook, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/hooks/git/{id}", pathParams("owner", owner, "repo", repo, "id", id))
	var out types.GitHook
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteGitHook deletes a Git hook from a repository.
func (s *WebhookService) DeleteGitHook(ctx context.Context, owner, repo, id string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/hooks/git/{id}", pathParams("owner", owner, "repo", repo, "id", id))
	return s.client.Delete(ctx, path)
}

// ListUserHooks returns all webhooks for the authenticated user.
func (s *WebhookService) ListUserHooks(ctx context.Context) ([]types.Hook, error) {
	return ListAll[types.Hook](ctx, s.client, "/api/v1/user/hooks", nil)
}

// IterUserHooks returns an iterator over all webhooks for the authenticated user.
func (s *WebhookService) IterUserHooks(ctx context.Context) iter.Seq2[types.Hook, error] {
	return ListIter[types.Hook](ctx, s.client, "/api/v1/user/hooks", nil)
}

// GetUserHook returns a single webhook for the authenticated user.
func (s *WebhookService) GetUserHook(ctx context.Context, id int64) (*types.Hook, error) {
	path := ResolvePath("/api/v1/user/hooks/{id}", pathParams("id", int64String(id)))
	var out types.Hook
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateUserHook creates a webhook for the authenticated user.
func (s *WebhookService) CreateUserHook(ctx context.Context, opts *types.CreateHookOption) (*types.Hook, error) {
	var out types.Hook
	if err := s.client.Post(ctx, "/api/v1/user/hooks", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditUserHook updates an existing authenticated-user webhook.
func (s *WebhookService) EditUserHook(ctx context.Context, id int64, opts *types.EditHookOption) (*types.Hook, error) {
	path := ResolvePath("/api/v1/user/hooks/{id}", pathParams("id", int64String(id)))
	var out types.Hook
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteUserHook deletes an authenticated-user webhook.
func (s *WebhookService) DeleteUserHook(ctx context.Context, id int64) error {
	path := ResolvePath("/api/v1/user/hooks/{id}", pathParams("id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListOrgHooks returns all webhooks for an organisation.
func (s *WebhookService) ListOrgHooks(ctx context.Context, org string) ([]types.Hook, error) {
	path := ResolvePath("/api/v1/orgs/{org}/hooks", pathParams("org", org))
	return ListAll[types.Hook](ctx, s.client, path, nil)
}

// IterOrgHooks returns an iterator over all webhooks for an organisation.
func (s *WebhookService) IterOrgHooks(ctx context.Context, org string) iter.Seq2[types.Hook, error] {
	path := ResolvePath("/api/v1/orgs/{org}/hooks", pathParams("org", org))
	return ListIter[types.Hook](ctx, s.client, path, nil)
}

// GetOrgHook returns a single webhook for an organisation.
func (s *WebhookService) GetOrgHook(ctx context.Context, org string, id int64) (*types.Hook, error) {
	path := ResolvePath("/api/v1/orgs/{org}/hooks/{id}", pathParams("org", org, "id", int64String(id)))
	var out types.Hook
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateOrgHook creates a webhook for an organisation.
func (s *WebhookService) CreateOrgHook(ctx context.Context, org string, opts *types.CreateHookOption) (*types.Hook, error) {
	path := ResolvePath("/api/v1/orgs/{org}/hooks", pathParams("org", org))
	var out types.Hook
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditOrgHook updates an existing organisation webhook.
func (s *WebhookService) EditOrgHook(ctx context.Context, org string, id int64, opts *types.EditHookOption) (*types.Hook, error) {
	path := ResolvePath("/api/v1/orgs/{org}/hooks/{id}", pathParams("org", org, "id", int64String(id)))
	var out types.Hook
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteOrgHook deletes an organisation webhook.
func (s *WebhookService) DeleteOrgHook(ctx context.Context, org string, id int64) error {
	path := ResolvePath("/api/v1/orgs/{org}/hooks/{id}", pathParams("org", org, "id", int64String(id)))
	return s.client.Delete(ctx, path)
}
