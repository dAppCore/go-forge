package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// BranchService handles branch operations within a repository.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Branches.ListBranchProtections(ctx, "core", "go-forge")
type BranchService struct {
	Resource[types.Branch, types.CreateBranchRepoOption, types.UpdateBranchRepoOption]
}

func newBranchService(c *Client) *BranchService {
	return &BranchService{
		Resource: *NewResource[types.Branch, types.CreateBranchRepoOption, types.UpdateBranchRepoOption](
			c, "/api/v1/repos/{owner}/{repo}/branches/{branch}",
		),
	}
}

// ListBranchesPage returns a single page of branches for a repository.
func (s *BranchService) ListBranchesPage(ctx context.Context, owner, repo string, opts ListOptions) (*PagedResult[types.Branch], error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branches", pathParams("owner", owner, "repo", repo))
	return ListPage[types.Branch](ctx, s.client, path, nil, opts)
}

// ListBranches returns all branches for a repository.
func (s *BranchService) ListBranches(ctx context.Context, owner, repo string) ([]types.Branch, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branches", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Branch](ctx, s.client, path, nil)
}

// IterBranches returns an iterator over all branches for a repository.
func (s *BranchService) IterBranches(ctx context.Context, owner, repo string) iter.Seq2[types.Branch, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branches", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Branch](ctx, s.client, path, nil)
}

// CreateBranch creates a new branch in a repository.
func (s *BranchService) CreateBranch(ctx context.Context, owner, repo string, opts *types.CreateBranchRepoOption) (*types.Branch, error) {
	var out types.Branch
	if err := s.client.Post(ctx, ResolvePath("/api/v1/repos/{owner}/{repo}/branches", pathParams("owner", owner, "repo", repo)), opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBranch returns a single branch by name.
func (s *BranchService) GetBranch(ctx context.Context, owner, repo, branch string) (*types.Branch, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branches/{branch}", pathParams("owner", owner, "repo", repo, "branch", branch))
	var out types.Branch
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateBranch renames a branch in a repository.
func (s *BranchService) UpdateBranch(ctx context.Context, owner, repo, branch string, opts *types.UpdateBranchRepoOption) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branches/{branch}", pathParams("owner", owner, "repo", repo, "branch", branch))
	return s.client.Patch(ctx, path, opts, nil)
}

// DeleteBranch removes a branch from a repository.
func (s *BranchService) DeleteBranch(ctx context.Context, owner, repo, branch string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branches/{branch}", pathParams("owner", owner, "repo", repo, "branch", branch))
	return s.client.Delete(ctx, path)
}

// ListBranchProtections returns all branch protections for a repository.
func (s *BranchService) ListBranchProtections(ctx context.Context, owner, repo string) ([]types.BranchProtection, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branch_protections", pathParams("owner", owner, "repo", repo))
	return ListAll[types.BranchProtection](ctx, s.client, path, nil)
}

// IterBranchProtections returns an iterator over all branch protections for a repository.
func (s *BranchService) IterBranchProtections(ctx context.Context, owner, repo string) iter.Seq2[types.BranchProtection, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branch_protections", pathParams("owner", owner, "repo", repo))
	return ListIter[types.BranchProtection](ctx, s.client, path, nil)
}

// GetBranchProtection returns a single branch protection by name.
func (s *BranchService) GetBranchProtection(ctx context.Context, owner, repo, name string) (*types.BranchProtection, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branch_protections/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	var out types.BranchProtection
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateBranchProtection creates a new branch protection rule.
func (s *BranchService) CreateBranchProtection(ctx context.Context, owner, repo string, opts *types.CreateBranchProtectionOption) (*types.BranchProtection, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branch_protections", pathParams("owner", owner, "repo", repo))
	var out types.BranchProtection
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditBranchProtection updates an existing branch protection rule.
func (s *BranchService) EditBranchProtection(ctx context.Context, owner, repo, name string, opts *types.EditBranchProtectionOption) (*types.BranchProtection, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branch_protections/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	var out types.BranchProtection
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBranchProtection deletes a branch protection rule.
func (s *BranchService) DeleteBranchProtection(ctx context.Context, owner, repo, name string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/branch_protections/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	return s.client.Delete(ctx, path)
}
