package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// PullService handles pull request operations within a repository.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Pulls.ListAll(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
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
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/merge", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	body := map[string]string{"Do": method}
	return s.client.Post(ctx, path, body, nil)
}

// CancelScheduledAutoMerge cancels the scheduled auto merge for a pull request.
func (s *PullService) CancelScheduledAutoMerge(ctx context.Context, owner, repo string, index int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/merge", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Delete(ctx, path)
}

// Update updates a pull request branch with the base branch.
func (s *PullService) Update(ctx context.Context, owner, repo string, index int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/update", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Post(ctx, path, nil, nil)
}

// ListReviews returns all reviews on a pull request.
func (s *PullService) ListReviews(ctx context.Context, owner, repo string, index int64) ([]types.PullReview, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.PullReview](ctx, s.client, path, nil)
}

// IterReviews returns an iterator over all reviews on a pull request.
func (s *PullService) IterReviews(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.PullReview, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.PullReview](ctx, s.client, path, nil)
}

// ListFiles returns all changed files on a pull request.
func (s *PullService) ListFiles(ctx context.Context, owner, repo string, index int64) ([]types.ChangedFile, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/files", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.ChangedFile](ctx, s.client, path, nil)
}

// IterFiles returns an iterator over all changed files on a pull request.
func (s *PullService) IterFiles(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.ChangedFile, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/files", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.ChangedFile](ctx, s.client, path, nil)
}

// ListReviewers returns all users who can be requested to review a pull request.
func (s *PullService) ListReviewers(ctx context.Context, owner, repo string) ([]types.User, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/reviewers", pathParams("owner", owner, "repo", repo))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterReviewers returns an iterator over all users who can be requested to review a pull request.
func (s *PullService) IterReviewers(ctx context.Context, owner, repo string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/reviewers", pathParams("owner", owner, "repo", repo))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// RequestReviewers creates review requests for a pull request.
func (s *PullService) RequestReviewers(ctx context.Context, owner, repo string, index int64, opts *types.PullReviewRequestOptions) ([]types.PullReview, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/requested_reviewers", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	var out []types.PullReview
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CancelReviewRequests cancels review requests for a pull request.
func (s *PullService) CancelReviewRequests(ctx context.Context, owner, repo string, index int64, opts *types.PullReviewRequestOptions) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/requested_reviewers", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.DeleteWithBody(ctx, path, opts)
}

// SubmitReview creates a new review on a pull request.
func (s *PullService) SubmitReview(ctx context.Context, owner, repo string, index int64, review map[string]any) (*types.PullReview, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	var out types.PullReview
	if err := s.client.Post(ctx, path, review, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DismissReview dismisses a pull request review.
func (s *PullService) DismissReview(ctx context.Context, owner, repo string, index, reviewID int64, msg string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{reviewID}/dismissals", pathParams("owner", owner, "repo", repo, "index", int64String(index), "reviewID", int64String(reviewID)))
	body := map[string]string{"message": msg}
	return s.client.Post(ctx, path, body, nil)
}

// UndismissReview undismisses a pull request review.
func (s *PullService) UndismissReview(ctx context.Context, owner, repo string, index, reviewID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{reviewID}/undismissals", pathParams("owner", owner, "repo", repo, "index", int64String(index), "reviewID", int64String(reviewID)))
	return s.client.Post(ctx, path, nil, nil)
}
