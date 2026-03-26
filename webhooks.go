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

// TestHook triggers a test delivery for a webhook.
func (s *WebhookService) TestHook(ctx context.Context, owner, repo string, id int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/hooks/{id}/tests", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return s.client.Post(ctx, path, nil, nil)
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
