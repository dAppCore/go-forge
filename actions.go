package forge

import (
	"context"
	"fmt"
	"iter"

	"forge.lthn.ai/core/go-forge/types"
)

// ActionsService handles CI/CD actions operations across repositories and
// organisations — secrets, variables, and workflow dispatches.
// No Resource embedding — heterogeneous endpoints across repo and org levels.
type ActionsService struct {
	client *Client
}

func newActionsService(c *Client) *ActionsService {
	return &ActionsService{client: c}
}

// ListRepoSecrets returns all secrets for a repository.
func (s *ActionsService) ListRepoSecrets(ctx context.Context, owner, repo string) ([]types.Secret, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/secrets", owner, repo)
	return ListAll[types.Secret](ctx, s.client, path, nil)
}

// IterRepoSecrets returns an iterator over all secrets for a repository.
func (s *ActionsService) IterRepoSecrets(ctx context.Context, owner, repo string) iter.Seq2[types.Secret, error] {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/secrets", owner, repo)
	return ListIter[types.Secret](ctx, s.client, path, nil)
}

// CreateRepoSecret creates or updates a secret in a repository.
// Forgejo expects a PUT with {"data": "secret-value"} body.
func (s *ActionsService) CreateRepoSecret(ctx context.Context, owner, repo, name string, data string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/secrets/%s", owner, repo, name)
	body := map[string]string{"data": data}
	return s.client.Put(ctx, path, body, nil)
}

// DeleteRepoSecret removes a secret from a repository.
func (s *ActionsService) DeleteRepoSecret(ctx context.Context, owner, repo, name string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/secrets/%s", owner, repo, name)
	return s.client.Delete(ctx, path)
}

// ListRepoVariables returns all action variables for a repository.
func (s *ActionsService) ListRepoVariables(ctx context.Context, owner, repo string) ([]types.ActionVariable, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/variables", owner, repo)
	return ListAll[types.ActionVariable](ctx, s.client, path, nil)
}

// IterRepoVariables returns an iterator over all action variables for a repository.
func (s *ActionsService) IterRepoVariables(ctx context.Context, owner, repo string) iter.Seq2[types.ActionVariable, error] {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/variables", owner, repo)
	return ListIter[types.ActionVariable](ctx, s.client, path, nil)
}

// CreateRepoVariable creates a new action variable in a repository.
// Forgejo expects a POST with {"value": "var-value"} body.
func (s *ActionsService) CreateRepoVariable(ctx context.Context, owner, repo, name, value string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/variables/%s", owner, repo, name)
	body := types.CreateVariableOption{Value: value}
	return s.client.Post(ctx, path, body, nil)
}

// DeleteRepoVariable removes an action variable from a repository.
func (s *ActionsService) DeleteRepoVariable(ctx context.Context, owner, repo, name string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/variables/%s", owner, repo, name)
	return s.client.Delete(ctx, path)
}

// ListOrgSecrets returns all secrets for an organisation.
func (s *ActionsService) ListOrgSecrets(ctx context.Context, org string) ([]types.Secret, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/actions/secrets", org)
	return ListAll[types.Secret](ctx, s.client, path, nil)
}

// IterOrgSecrets returns an iterator over all secrets for an organisation.
func (s *ActionsService) IterOrgSecrets(ctx context.Context, org string) iter.Seq2[types.Secret, error] {
	path := fmt.Sprintf("/api/v1/orgs/%s/actions/secrets", org)
	return ListIter[types.Secret](ctx, s.client, path, nil)
}

// ListOrgVariables returns all action variables for an organisation.
func (s *ActionsService) ListOrgVariables(ctx context.Context, org string) ([]types.ActionVariable, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/actions/variables", org)
	return ListAll[types.ActionVariable](ctx, s.client, path, nil)
}

// IterOrgVariables returns an iterator over all action variables for an organisation.
func (s *ActionsService) IterOrgVariables(ctx context.Context, org string) iter.Seq2[types.ActionVariable, error] {
	path := fmt.Sprintf("/api/v1/orgs/%s/actions/variables", org)
	return ListIter[types.ActionVariable](ctx, s.client, path, nil)
}

// DispatchWorkflow triggers a workflow run.
func (s *ActionsService) DispatchWorkflow(ctx context.Context, owner, repo, workflow string, opts map[string]any) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/actions/workflows/%s/dispatches", owner, repo, workflow)
	return s.client.Post(ctx, path, opts, nil)
}
