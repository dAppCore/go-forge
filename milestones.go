package forge

import (
	"context"

	"dappco.re/go/core/forge/types"
)

// MilestoneService handles repository milestones.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Milestones.ListAll(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
type MilestoneService struct {
	client *Client
}

func newMilestoneService(c *Client) *MilestoneService {
	return &MilestoneService{client: c}
}

// ListAll returns all milestones for a repository.
func (s *MilestoneService) ListAll(ctx context.Context, params Params) ([]types.Milestone, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/milestones", params)
	return ListAll[types.Milestone](ctx, s.client, path, nil)
}

// Get returns a single milestone by ID.
func (s *MilestoneService) Get(ctx context.Context, owner, repo string, id int64) (*types.Milestone, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/milestones/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.Milestone
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Create creates a new milestone.
func (s *MilestoneService) Create(ctx context.Context, owner, repo string, opts *types.CreateMilestoneOption) (*types.Milestone, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/milestones", pathParams("owner", owner, "repo", repo))
	var out types.Milestone
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
