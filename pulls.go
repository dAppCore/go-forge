package forge

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go-forge/types"
)

// PullService handles pull request operations within a repository.
type PullService struct {
	Resource[types.PullRequest, types.CreatePullRequestOption, types.EditPullRequestOption]
}

func newPullService(c *Client) *PullService {
	return &PullService{
		Resource: *NewResource[types.PullRequest, types.CreatePullRequestOption, types.EditPullRequestOption](
			c, "/api/v1/repos/{owner}/{repo}/pulls/{index}",
		),
	}
}

// Merge merges a pull request. Method is one of "merge", "rebase", "rebase-merge", "squash", "fast-forward-only", "manually-merged".
func (s *PullService) Merge(ctx context.Context, owner, repo string, index int64, method string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/pulls/%d/merge", owner, repo, index)
	body := map[string]string{"Do": method}
	return s.client.Post(ctx, path, body, nil)
}

// Update updates a pull request branch with the base branch.
func (s *PullService) Update(ctx context.Context, owner, repo string, index int64) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/pulls/%d/update", owner, repo, index)
	return s.client.Post(ctx, path, nil, nil)
}

// ListReviews returns all reviews on a pull request.
func (s *PullService) ListReviews(ctx context.Context, owner, repo string, index int64) ([]types.PullReview, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/pulls/%d/reviews", owner, repo, index)
	return ListAll[types.PullReview](ctx, s.client, path, nil)
}

// SubmitReview creates a new review on a pull request.
func (s *PullService) SubmitReview(ctx context.Context, owner, repo string, index int64, review map[string]any) (*types.PullReview, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/pulls/%d/reviews", owner, repo, index)
	var out types.PullReview
	if err := s.client.Post(ctx, path, review, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DismissReview dismisses a pull request review.
func (s *PullService) DismissReview(ctx context.Context, owner, repo string, index, reviewID int64, msg string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/pulls/%d/reviews/%d/dismissals", owner, repo, index, reviewID)
	body := map[string]string{"message": msg}
	return s.client.Post(ctx, path, body, nil)
}

// UndismissReview undismisses a pull request review.
func (s *PullService) UndismissReview(ctx context.Context, owner, repo string, index, reviewID int64) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/pulls/%d/reviews/%d/undismissals", owner, repo, index, reviewID)
	return s.client.Post(ctx, path, nil, nil)
}
