package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// AdminService handles site administration operations.
// Unlike other services, AdminService does not embed Resource[T,C,U]
// because admin endpoints are heterogeneous.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Admin.ListUsers(ctx)
type AdminService struct {
	client *Client
}

func newAdminService(c *Client) *AdminService {
	return &AdminService{client: c}
}

// ListUsers returns all users (admin only).
func (s *AdminService) ListUsers(ctx context.Context) ([]types.User, error) {
	return ListAll[types.User](ctx, s.client, "/api/v1/admin/users", nil)
}

// IterUsers returns an iterator over all users (admin only).
func (s *AdminService) IterUsers(ctx context.Context) iter.Seq2[types.User, error] {
	return ListIter[types.User](ctx, s.client, "/api/v1/admin/users", nil)
}

// CreateUser creates a new user (admin only).
func (s *AdminService) CreateUser(ctx context.Context, opts *types.CreateUserOption) (*types.User, error) {
	var out types.User
	if err := s.client.Post(ctx, "/api/v1/admin/users", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditUser edits an existing user (admin only).
func (s *AdminService) EditUser(ctx context.Context, username string, opts map[string]any) error {
	path := ResolvePath("/api/v1/admin/users/{username}", Params{"username": username})
	return s.client.Patch(ctx, path, opts, nil)
}

// DeleteUser deletes a user (admin only).
func (s *AdminService) DeleteUser(ctx context.Context, username string) error {
	path := ResolvePath("/api/v1/admin/users/{username}", Params{"username": username})
	return s.client.Delete(ctx, path)
}

// RenameUser renames a user (admin only).
func (s *AdminService) RenameUser(ctx context.Context, username, newName string) error {
	path := ResolvePath("/api/v1/admin/users/{username}/rename", Params{"username": username})
	return s.client.Post(ctx, path, &types.RenameUserOption{NewName: newName}, nil)
}

// ListOrgs returns all organisations (admin only).
func (s *AdminService) ListOrgs(ctx context.Context) ([]types.Organization, error) {
	return ListAll[types.Organization](ctx, s.client, "/api/v1/admin/orgs", nil)
}

// IterOrgs returns an iterator over all organisations (admin only).
func (s *AdminService) IterOrgs(ctx context.Context) iter.Seq2[types.Organization, error] {
	return ListIter[types.Organization](ctx, s.client, "/api/v1/admin/orgs", nil)
}

// ListEmails returns all email addresses (admin only).
func (s *AdminService) ListEmails(ctx context.Context) ([]types.Email, error) {
	return ListAll[types.Email](ctx, s.client, "/api/v1/admin/emails", nil)
}

// IterEmails returns an iterator over all email addresses (admin only).
func (s *AdminService) IterEmails(ctx context.Context) iter.Seq2[types.Email, error] {
	return ListIter[types.Email](ctx, s.client, "/api/v1/admin/emails", nil)
}

// SearchEmails searches all email addresses by keyword (admin only).
func (s *AdminService) SearchEmails(ctx context.Context, q string) ([]types.Email, error) {
	return ListAll[types.Email](ctx, s.client, "/api/v1/admin/emails/search", map[string]string{"q": q})
}

// IterSearchEmails returns an iterator over all email addresses matching a keyword (admin only).
func (s *AdminService) IterSearchEmails(ctx context.Context, q string) iter.Seq2[types.Email, error] {
	return ListIter[types.Email](ctx, s.client, "/api/v1/admin/emails/search", map[string]string{"q": q})
}

// RunCron runs a cron task by name (admin only).
func (s *AdminService) RunCron(ctx context.Context, task string) error {
	path := ResolvePath("/api/v1/admin/cron/{task}", Params{"task": task})
	return s.client.Post(ctx, path, nil, nil)
}

// ListCron returns all cron tasks (admin only).
func (s *AdminService) ListCron(ctx context.Context) ([]types.Cron, error) {
	return ListAll[types.Cron](ctx, s.client, "/api/v1/admin/cron", nil)
}

// IterCron returns an iterator over all cron tasks (admin only).
func (s *AdminService) IterCron(ctx context.Context) iter.Seq2[types.Cron, error] {
	return ListIter[types.Cron](ctx, s.client, "/api/v1/admin/cron", nil)
}

// AdoptRepo adopts an unadopted repository (admin only).
func (s *AdminService) AdoptRepo(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/admin/unadopted/{owner}/{repo}", Params{"owner": owner, "repo": repo})
	return s.client.Post(ctx, path, nil, nil)
}

// GenerateRunnerToken generates an actions runner registration token.
func (s *AdminService) GenerateRunnerToken(ctx context.Context) (string, error) {
	var out struct {
		Token string `json:"token"`
	}
	if err := s.client.Get(ctx, "/api/v1/admin/runners/registration-token", &out); err != nil {
		return "", err
	}
	return out.Token, nil
}
