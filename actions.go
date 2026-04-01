package forge

import (
	"context"
	"iter"
	"net/url"
	"strconv"

	"dappco.re/go/core/forge/types"
)

// ActionsService handles CI/CD actions operations across repositories and
// organisations — secrets, variables, workflow dispatches, and tasks.
// No Resource embedding — heterogeneous endpoints across repo and org levels.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Actions.ListRepoSecrets(ctx, "core", "go-forge")
type ActionsService struct {
	client *Client
}

func newActionsService(c *Client) *ActionsService {
	return &ActionsService{client: c}
}

// ListRepoSecrets returns all secrets for a repository.
func (s *ActionsService) ListRepoSecrets(ctx context.Context, owner, repo string) ([]types.Secret, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/secrets", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Secret](ctx, s.client, path, nil)
}

// IterRepoSecrets returns an iterator over all secrets for a repository.
func (s *ActionsService) IterRepoSecrets(ctx context.Context, owner, repo string) iter.Seq2[types.Secret, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/secrets", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Secret](ctx, s.client, path, nil)
}

// CreateRepoSecret creates or updates a secret in a repository.
// Forgejo expects a PUT with {"data": "secret-value"} body.
func (s *ActionsService) CreateRepoSecret(ctx context.Context, owner, repo, name string, data string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/secrets/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	body := map[string]string{"data": data}
	return s.client.Put(ctx, path, body, nil)
}

// DeleteRepoSecret removes a secret from a repository.
func (s *ActionsService) DeleteRepoSecret(ctx context.Context, owner, repo, name string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/secrets/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	return s.client.Delete(ctx, path)
}

// ListRepoVariables returns all action variables for a repository.
func (s *ActionsService) ListRepoVariables(ctx context.Context, owner, repo string) ([]types.ActionVariable, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/variables", pathParams("owner", owner, "repo", repo))
	return ListAll[types.ActionVariable](ctx, s.client, path, nil)
}

// IterRepoVariables returns an iterator over all action variables for a repository.
func (s *ActionsService) IterRepoVariables(ctx context.Context, owner, repo string) iter.Seq2[types.ActionVariable, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/variables", pathParams("owner", owner, "repo", repo))
	return ListIter[types.ActionVariable](ctx, s.client, path, nil)
}

// CreateRepoVariable creates a new action variable in a repository.
// Forgejo expects a POST with {"value": "var-value"} body.
func (s *ActionsService) CreateRepoVariable(ctx context.Context, owner, repo, name, value string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/variables/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	body := types.CreateVariableOption{Value: value}
	return s.client.Post(ctx, path, body, nil)
}

// DeleteRepoVariable removes an action variable from a repository.
func (s *ActionsService) DeleteRepoVariable(ctx context.Context, owner, repo, name string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/variables/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	return s.client.Delete(ctx, path)
}

// ListOrgSecrets returns all secrets for an organisation.
func (s *ActionsService) ListOrgSecrets(ctx context.Context, org string) ([]types.Secret, error) {
	path := ResolvePath("/api/v1/orgs/{org}/actions/secrets", pathParams("org", org))
	return ListAll[types.Secret](ctx, s.client, path, nil)
}

// IterOrgSecrets returns an iterator over all secrets for an organisation.
func (s *ActionsService) IterOrgSecrets(ctx context.Context, org string) iter.Seq2[types.Secret, error] {
	path := ResolvePath("/api/v1/orgs/{org}/actions/secrets", pathParams("org", org))
	return ListIter[types.Secret](ctx, s.client, path, nil)
}

// ListOrgVariables returns all action variables for an organisation.
func (s *ActionsService) ListOrgVariables(ctx context.Context, org string) ([]types.ActionVariable, error) {
	path := ResolvePath("/api/v1/orgs/{org}/actions/variables", pathParams("org", org))
	return ListAll[types.ActionVariable](ctx, s.client, path, nil)
}

// IterOrgVariables returns an iterator over all action variables for an organisation.
func (s *ActionsService) IterOrgVariables(ctx context.Context, org string) iter.Seq2[types.ActionVariable, error] {
	path := ResolvePath("/api/v1/orgs/{org}/actions/variables", pathParams("org", org))
	return ListIter[types.ActionVariable](ctx, s.client, path, nil)
}

// DispatchWorkflow triggers a workflow run.
func (s *ActionsService) DispatchWorkflow(ctx context.Context, owner, repo, workflow string, opts map[string]any) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/workflows/{workflow}/dispatches", pathParams("owner", owner, "repo", repo, "workflow", workflow))
	return s.client.Post(ctx, path, opts, nil)
}

// ListRepoTasks returns a single page of action tasks for a repository.
func (s *ActionsService) ListRepoTasks(ctx context.Context, owner, repo string, opts ListOptions) (*types.ActionTaskResponse, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/tasks", pathParams("owner", owner, "repo", repo))

	if opts.Page > 0 || opts.Limit > 0 {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		if opts.Page > 0 {
			q.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			q.Set("limit", strconv.Itoa(opts.Limit))
		}
		u.RawQuery = q.Encode()
		path = u.String()
	}

	var out types.ActionTaskResponse
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IterRepoTasks returns an iterator over all action tasks for a repository.
func (s *ActionsService) IterRepoTasks(ctx context.Context, owner, repo string) iter.Seq2[types.ActionTask, error] {
	return func(yield func(types.ActionTask, error) bool) {
		const limit = 50
		var seen int64
		for page := 1; ; page++ {
			resp, err := s.ListRepoTasks(ctx, owner, repo, ListOptions{Page: page, Limit: limit})
			if err != nil {
				yield(*new(types.ActionTask), err)
				return
			}
			for _, item := range resp.Entries {
				if !yield(*item, nil) {
					return
				}
				seen++
			}
			if resp.TotalCount > 0 {
				if seen >= resp.TotalCount {
					return
				}
				continue
			}
			if len(resp.Entries) < limit {
				return
			}
		}
	}
}
