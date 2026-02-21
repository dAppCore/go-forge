package forge

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go-forge/types"
)

// BranchService handles branch operations within a repository.
type BranchService struct {
	Resource[types.Branch, types.CreateBranchRepoOption, struct{}]
}

func newBranchService(c *Client) *BranchService {
	return &BranchService{
		Resource: *NewResource[types.Branch, types.CreateBranchRepoOption, struct{}](
			c, "/api/v1/repos/{owner}/{repo}/branches/{branch}",
		),
	}
}

// ListBranchProtections returns all branch protections for a repository.
func (s *BranchService) ListBranchProtections(ctx context.Context, owner, repo string) ([]types.BranchProtection, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/branch_protections", owner, repo)
	return ListAll[types.BranchProtection](ctx, s.client, path, nil)
}

// GetBranchProtection returns a single branch protection by name.
func (s *BranchService) GetBranchProtection(ctx context.Context, owner, repo, name string) (*types.BranchProtection, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/branch_protections/%s", owner, repo, name)
	var out types.BranchProtection
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateBranchProtection creates a new branch protection rule.
func (s *BranchService) CreateBranchProtection(ctx context.Context, owner, repo string, opts *types.CreateBranchProtectionOption) (*types.BranchProtection, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/branch_protections", owner, repo)
	var out types.BranchProtection
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditBranchProtection updates an existing branch protection rule.
func (s *BranchService) EditBranchProtection(ctx context.Context, owner, repo, name string, opts *types.EditBranchProtectionOption) (*types.BranchProtection, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/branch_protections/%s", owner, repo, name)
	var out types.BranchProtection
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBranchProtection deletes a branch protection rule.
func (s *BranchService) DeleteBranchProtection(ctx context.Context, owner, repo, name string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/branch_protections/%s", owner, repo, name)
	return s.client.Delete(ctx, path)
}
