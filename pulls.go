package forge

import (
	"context"
	"iter"
	"net/url"
	"strconv"

	core "dappco.re/go/core"
	"dappco.re/go/forge/types"
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

// PullListOptions controls filtering for repository pull request listings.
//
// Usage:
//
//	opts := forge.PullListOptions{State: "open", Labels: []int64{1, 2}}
type PullListOptions struct {
	State     string
	Sort      string
	Milestone int64
	Labels    []int64
	Poster    string
}

// String returns a safe summary of the pull request list filters.
func (o PullListOptions) String() string {
	return optionString("forge.PullListOptions",
		"state", o.State,
		"sort", o.Sort,
		"milestone", o.Milestone,
		"labels", o.Labels,
		"poster", o.Poster,
	)
}

// GoString returns a safe Go-syntax summary of the pull request list filters.
func (o PullListOptions) GoString() string { return o.String() }

func (o PullListOptions) addQuery(values url.Values) {
	if o.State != "" {
		values.Set("state", o.State)
	}
	if o.Sort != "" {
		values.Set("sort", o.Sort)
	}
	if o.Milestone != 0 {
		values.Set("milestone", strconv.FormatInt(o.Milestone, 10))
	}
	for _, label := range o.Labels {
		if label != 0 {
			values.Add("labels", strconv.FormatInt(label, 10))
		}
	}
	if o.Poster != "" {
		values.Set("poster", o.Poster)
	}
}

func newPullService(c *Client) *PullService {
	return &PullService{
		Resource: *NewResource[types.PullRequest, types.CreatePullRequestOption, types.EditPullRequestOption](
			c, "/api/v1/repos/{owner}/{repo}/pulls/{index}",
		),
	}
}

// ListPullRequestsPage returns a single page of pull requests in a repository.
func (s *PullService) ListPullRequestsPage(ctx context.Context, owner, repo string, opts ListOptions, filters ...any) (*PagedResult[types.PullRequest], error) {
	return s.listPage(ctx, owner, repo, opts, filters...)
}

// ListPullRequests returns all pull requests in a repository.
func (s *PullService) ListPullRequests(ctx context.Context, owner, repo string, filters ...any) ([]types.PullRequest, error) {
	if pageOpts, ok := compatPullListPageOptions(filters...); ok {
		page, err := s.listPage(ctx, owner, repo, pageOpts, filters...)
		if err != nil {
			return nil, err
		}
		return page.Items, nil
	}
	return s.listAll(ctx, owner, repo, filters...)
}

// IterPullRequests returns an iterator over all pull requests in a repository.
func (s *PullService) IterPullRequests(ctx context.Context, owner, repo string, filters ...any) iter.Seq2[types.PullRequest, error] {
	if pageOpts, ok := compatPullListPageOptions(filters...); ok {
		return func(yield func(types.PullRequest, error) bool) {
			page, err := s.listPage(ctx, owner, repo, pageOpts, filters...)
			if err != nil {
				yield(*new(types.PullRequest), err)
				return
			}
			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}
		}
	}
	return s.listIter(ctx, owner, repo, filters...)
}

// CreatePullRequest creates a pull request in a repository.
func (s *PullService) CreatePullRequest(ctx context.Context, owner, repo string, opts *types.CreatePullRequestOption) (*types.PullRequest, error) {
	var out types.PullRequest
	if err := s.client.Post(ctx, ResolvePath("/api/v1/repos/{owner}/{repo}/pulls", pathParams("owner", owner, "repo", repo)), opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPullRequest returns a single pull request by index.
func (s *PullService) GetPullRequest(ctx context.Context, owner, repo string, index int64) (*types.PullRequest, error) {
	return s.Get(ctx, pathParams("owner", owner, "repo", repo, "index", int64String(index)))
}

// EditPullRequest updates an existing pull request.
func (s *PullService) EditPullRequest(ctx context.Context, owner, repo string, index int64, opts *types.EditPullRequestOption) (*types.PullRequest, error) {
	return s.Resource.Update(ctx, pathParams("owner", owner, "repo", repo, "index", int64String(index)), opts)
}

// DeletePullRequest deletes a pull request.
func (s *PullService) DeletePullRequest(ctx context.Context, owner, repo string, index int64) error {
	return s.Delete(ctx, pathParams("owner", owner, "repo", repo, "index", int64String(index)))
}

// Merge merges a pull request. Method is one of "merge", "rebase", "rebase-merge", "squash", "fast-forward-only", "manually-merged".
func (s *PullService) Merge(ctx context.Context, owner, repo string, index int64, method string) error {
	return s.MergePullRequest(ctx, owner, repo, index, &types.MergePullRequestOption{Do: method})
}

// MergePullRequest merges a pull request.
func (s *PullService) MergePullRequest(ctx context.Context, owner, repo string, index int64, opts *types.MergePullRequestOption) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/merge", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	if opts != nil {
		body := *opts
		if body.Do == "" {
			body.Do = body.MergeStyle
		}
		opts = &body
	}
	return s.client.Post(ctx, path, opts, nil)
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

// GetDiffOrPatch returns a pull request diff or patch as raw bytes.
func (s *PullService) GetDiffOrPatch(ctx context.Context, owner, repo string, index int64, diffType string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}.{diffType}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "diffType", diffType))
	return s.client.GetRaw(ctx, path)
}

// ListCommits returns all commits for a pull request.
func (s *PullService) ListCommits(ctx context.Context, owner, repo string, index int64) ([]types.Commit, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/commits", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.Commit](ctx, s.client, path, nil)
}

// IterCommits returns an iterator over all commits for a pull request.
func (s *PullService) IterCommits(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Commit, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/commits", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.Commit](ctx, s.client, path, nil)
}

// ListReviews returns all reviews on a pull request.
func (s *PullService) ListReviews(ctx context.Context, owner, repo string, index int64) ([]types.PullReview, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.PullReview](ctx, s.client, path, nil)
}

// ListPullReviews returns all reviews on a pull request.
func (s *PullService) ListPullReviews(ctx context.Context, owner, repo string, index int64) ([]types.PullReview, error) {
	return s.ListReviews(ctx, owner, repo, index)
}

// IterReviews returns an iterator over all reviews on a pull request.
func (s *PullService) IterReviews(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.PullReview, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.PullReview](ctx, s.client, path, nil)
}

// IterPullReviews returns an iterator over all reviews on a pull request.
func (s *PullService) IterPullReviews(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.PullReview, error] {
	return s.IterReviews(ctx, owner, repo, index)
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

// GetByBaseHead returns a pull request for a given base and head branch pair.
func (s *PullService) GetByBaseHead(ctx context.Context, owner, repo, base, head string) (*types.PullRequest, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{base}/{head}", pathParams(
		"owner", owner,
		"repo", repo,
		"base", base,
		"head", head,
	))
	var out types.PullRequest
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
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
func (s *PullService) SubmitReview(ctx context.Context, owner, repo string, index int64, review *types.SubmitPullReviewOptions) (*types.PullReview, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	var out types.PullReview
	if err := s.client.Post(ctx, path, review, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetReview returns a single pull request review.
func (s *PullService) GetReview(ctx context.Context, owner, repo string, index, reviewID int64) (*types.PullReview, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID)))
	var out types.PullReview
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPullReview returns a single pull request review.
func (s *PullService) GetPullReview(ctx context.Context, owner, repo string, index, reviewID int64) (*types.PullReview, error) {
	return s.GetReview(ctx, owner, repo, index, reviewID)
}

// DeleteReview deletes a pull request review.
func (s *PullService) DeleteReview(ctx context.Context, owner, repo string, index, reviewID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID)))
	return s.client.Delete(ctx, path)
}

// DeletePullReview deletes a pull request review.
func (s *PullService) DeletePullReview(ctx context.Context, owner, repo string, index, reviewID int64) error {
	return s.DeleteReview(ctx, owner, repo, index, reviewID)
}

func (s *PullService) listPage(ctx context.Context, owner, repo string, opts ListOptions, filters ...any) (*PagedResult[types.PullRequest], error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	pageSize := opts.PageSize
	if pageSize < 1 {
		pageSize = opts.Limit
	}
	if pageSize < 1 {
		pageSize = defaultPageLimit
	}

	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls", pathParams("owner", owner, "repo", repo))
	u, err := url.Parse(path)
	if err != nil {
		return nil, core.E("PullService.listPage", "forge: parse path", err)
	}

	values := u.Query()
	values.Set("page", strconv.Itoa(opts.Page))
	values.Set("limit", strconv.Itoa(pageSize))
	addPullFilters(values, filters...)
	u.RawQuery = values.Encode()

	var items []types.PullRequest
	resp, err := s.client.doJSON(ctx, "GET", u.String(), nil, &items)
	if err != nil {
		return nil, err
	}

	totalCount, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))
	return &PagedResult[types.PullRequest]{
		Items:      items,
		TotalCount: totalCount,
		Page:       opts.Page,
		HasMore: (totalCount > 0 && (opts.Page-1)*pageSize+len(items) < totalCount) ||
			(totalCount == 0 && len(items) >= pageSize),
	}, nil
}

func (s *PullService) listAll(ctx context.Context, owner, repo string, filters ...any) ([]types.PullRequest, error) {
	var all []types.PullRequest
	page := 1

	for {
		result, err := s.listPage(ctx, owner, repo, ListOptions{Page: page, PageSize: defaultPageLimit}, filters...)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Items...)
		if !result.HasMore {
			break
		}
		page++
	}

	return all, nil
}

func (s *PullService) listIter(ctx context.Context, owner, repo string, filters ...any) iter.Seq2[types.PullRequest, error] {
	return func(yield func(types.PullRequest, error) bool) {
		page := 1
		for {
			result, err := s.listPage(ctx, owner, repo, ListOptions{Page: page, PageSize: defaultPageLimit}, filters...)
			if err != nil {
				yield(*new(types.PullRequest), err)
				return
			}
			for _, item := range result.Items {
				if !yield(item, nil) {
					return
				}
			}
			if !result.HasMore {
				break
			}
			page++
		}
	}
}

func addPullFilters(values url.Values, filters ...any) {
	for _, filter := range filters {
		switch v := filter.(type) {
		case PullListOptions:
			v.addQuery(values)
		case *PullListOptions:
			if v != nil {
				v.addQuery(values)
			}
		case types.ListPullRequestsOption:
			addCompatPullFilter(values, v)
		case *types.ListPullRequestsOption:
			if v != nil {
				addCompatPullFilter(values, *v)
			}
		}
	}
}

func addCompatPullFilter(values url.Values, filter types.ListPullRequestsOption) {
	if filter.State != "" {
		values.Set("state", filter.State)
	}
	if filter.Sort != "" {
		values.Set("sort", filter.Sort)
	}
	if filter.Milestone != 0 {
		values.Set("milestone", strconv.FormatInt(filter.Milestone, 10))
	}
	for _, label := range filter.Labels {
		if label != 0 {
			values.Add("labels", strconv.FormatInt(label, 10))
		}
	}
	if filter.Poster != "" {
		values.Set("poster", filter.Poster)
	}
}

func compatPullListPageOptions(filters ...any) (ListOptions, bool) {
	for _, filter := range filters {
		switch v := filter.(type) {
		case types.ListPullRequestsOption:
			if opts, ok := compatListOptions(v.Page, v.PageSize, v.Limit); ok {
				return opts, true
			}
		case *types.ListPullRequestsOption:
			if v != nil {
				if opts, ok := compatListOptions(v.Page, v.PageSize, v.Limit); ok {
					return opts, true
				}
			}
		}
	}
	return ListOptions{}, false
}

// ListReviewComments returns all comments on a pull request review.
func (s *PullService) ListReviewComments(ctx context.Context, owner, repo string, index, reviewID int64) ([]types.PullReviewComment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID)))
	return ListAll[types.PullReviewComment](ctx, s.client, path, nil)
}

// IterReviewComments returns an iterator over all comments on a pull request review.
func (s *PullService) IterReviewComments(ctx context.Context, owner, repo string, index, reviewID int64) iter.Seq2[types.PullReviewComment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID)))
	return ListIter[types.PullReviewComment](ctx, s.client, path, nil)
}

// GetReviewComment returns a single comment on a pull request review.
func (s *PullService) GetReviewComment(ctx context.Context, owner, repo string, index, reviewID, commentID int64) (*types.PullReviewComment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments/{comment}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID), "comment", int64String(commentID)))
	var out types.PullReviewComment
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateReviewComment creates a new comment on a pull request review.
func (s *PullService) CreateReviewComment(ctx context.Context, owner, repo string, index, reviewID int64, opts *types.CreatePullReviewCommentOptions) (*types.PullReviewComment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID)))
	var out types.PullReviewComment
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteReviewComment deletes a comment on a pull request review.
func (s *PullService) DeleteReviewComment(ctx context.Context, owner, repo string, index, reviewID, commentID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/comments/{comment}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID), "comment", int64String(commentID)))
	return s.client.Delete(ctx, path)
}

// DismissReview dismisses a pull request review.
func (s *PullService) DismissReview(ctx context.Context, owner, repo string, index, reviewID int64, msg string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/dismissals", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID)))
	body := map[string]string{"message": msg}
	return s.client.Post(ctx, path, body, nil)
}

// UndismissReview undismisses a pull request review.
func (s *PullService) UndismissReview(ctx context.Context, owner, repo string, index, reviewID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/{index}/reviews/{id}/undismissals", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(reviewID)))
	return s.client.Post(ctx, path, nil, nil)
}
