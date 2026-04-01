package forge

import (
	"context"
	"iter"
	"strconv"
	"time"

	"dappco.re/go/core/forge/types"
)

// IssueService handles issue operations within a repository.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Issues.ListAll(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
type IssueService struct {
	Resource[types.Issue, types.CreateIssueOption, types.EditIssueOption]
}

func newIssueService(c *Client) *IssueService {
	return &IssueService{
		Resource: *NewResource[types.Issue, types.CreateIssueOption, types.EditIssueOption](
			c, "/api/v1/repos/{owner}/{repo}/issues/{index}",
		),
	}
}

// SearchIssuesOptions controls filtering for the global issue search endpoint.
type SearchIssuesOptions struct {
	State           string
	Labels          string
	Milestones      string
	Query           string
	PriorityRepoID  int64
	Type            string
	Since           *time.Time
	Before          *time.Time
	Assigned        bool
	Created         bool
	Mentioned       bool
	ReviewRequested bool
	Reviewed        bool
	Owner           string
	Team            string
}

func (o SearchIssuesOptions) queryParams() map[string]string {
	query := make(map[string]string, 12)
	if o.State != "" {
		query["state"] = o.State
	}
	if o.Labels != "" {
		query["labels"] = o.Labels
	}
	if o.Milestones != "" {
		query["milestones"] = o.Milestones
	}
	if o.Query != "" {
		query["q"] = o.Query
	}
	if o.PriorityRepoID != 0 {
		query["priority_repo_id"] = strconv.FormatInt(o.PriorityRepoID, 10)
	}
	if o.Type != "" {
		query["type"] = o.Type
	}
	if o.Since != nil {
		query["since"] = o.Since.Format(time.RFC3339)
	}
	if o.Before != nil {
		query["before"] = o.Before.Format(time.RFC3339)
	}
	if o.Assigned {
		query["assigned"] = strconv.FormatBool(true)
	}
	if o.Created {
		query["created"] = strconv.FormatBool(true)
	}
	if o.Mentioned {
		query["mentioned"] = strconv.FormatBool(true)
	}
	if o.ReviewRequested {
		query["review_requested"] = strconv.FormatBool(true)
	}
	if o.Reviewed {
		query["reviewed"] = strconv.FormatBool(true)
	}
	if o.Owner != "" {
		query["owner"] = o.Owner
	}
	if o.Team != "" {
		query["team"] = o.Team
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

// SearchIssuesPage returns a single page of issues matching the search filters.
func (s *IssueService) SearchIssuesPage(ctx context.Context, opts SearchIssuesOptions, pageOpts ListOptions) (*PagedResult[types.Issue], error) {
	return ListPage[types.Issue](ctx, s.client, "/api/v1/repos/issues/search", opts.queryParams(), pageOpts)
}

// SearchIssues returns all issues matching the search filters.
func (s *IssueService) SearchIssues(ctx context.Context, opts SearchIssuesOptions) ([]types.Issue, error) {
	return ListAll[types.Issue](ctx, s.client, "/api/v1/repos/issues/search", opts.queryParams())
}

// IterSearchIssues returns an iterator over issues matching the search filters.
func (s *IssueService) IterSearchIssues(ctx context.Context, opts SearchIssuesOptions) iter.Seq2[types.Issue, error] {
	return ListIter[types.Issue](ctx, s.client, "/api/v1/repos/issues/search", opts.queryParams())
}

// Pin pins an issue.
func (s *IssueService) Pin(ctx context.Context, owner, repo string, index int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/pin", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Post(ctx, path, nil, nil)
}

// MovePin moves a pinned issue to a new position.
func (s *IssueService) MovePin(ctx context.Context, owner, repo string, index, position int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/pin/{position}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "position", int64String(position)))
	return s.client.Patch(ctx, path, nil, nil)
}

// Unpin unpins an issue.
func (s *IssueService) Unpin(ctx context.Context, owner, repo string, index int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/pin", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Delete(ctx, path)
}

// SetDeadline sets or updates the deadline on an issue.
func (s *IssueService) SetDeadline(ctx context.Context, owner, repo string, index int64, deadline string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/deadline", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	body := map[string]string{"due_date": deadline}
	return s.client.Post(ctx, path, body, nil)
}

// AddReaction adds a reaction to an issue.
func (s *IssueService) AddReaction(ctx context.Context, owner, repo string, index int64, reaction string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/reactions", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	body := map[string]string{"content": reaction}
	return s.client.Post(ctx, path, body, nil)
}

// DeleteReaction removes a reaction from an issue.
func (s *IssueService) DeleteReaction(ctx context.Context, owner, repo string, index int64, reaction string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/reactions", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	body := map[string]string{"content": reaction}
	return s.client.DeleteWithBody(ctx, path, body)
}

// StartStopwatch starts the stopwatch on an issue.
func (s *IssueService) StartStopwatch(ctx context.Context, owner, repo string, index int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/stopwatch/start", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Post(ctx, path, nil, nil)
}

// StopStopwatch stops the stopwatch on an issue.
func (s *IssueService) StopStopwatch(ctx context.Context, owner, repo string, index int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/stopwatch/stop", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Post(ctx, path, nil, nil)
}

// DeleteStopwatch deletes an issue's existing stopwatch.
func (s *IssueService) DeleteStopwatch(ctx context.Context, owner, repo string, index int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/stopwatch/delete", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Delete(ctx, path)
}

// ListTimes returns all tracked times on an issue.
func (s *IssueService) ListTimes(ctx context.Context, owner, repo string, index int64, user string, since, before *time.Time) ([]types.TrackedTime, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/times", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	query := make(map[string]string, 3)
	if user != "" {
		query["user"] = user
	}
	if since != nil {
		query["since"] = since.Format(time.RFC3339)
	}
	if before != nil {
		query["before"] = before.Format(time.RFC3339)
	}
	if len(query) == 0 {
		query = nil
	}
	return ListAll[types.TrackedTime](ctx, s.client, path, query)
}

// AddTime adds tracked time to an issue.
func (s *IssueService) AddTime(ctx context.Context, owner, repo string, index int64, opts *types.AddTimeOption) (*types.TrackedTime, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/times", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	var out types.TrackedTime
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ResetTime removes all tracked time from an issue.
func (s *IssueService) ResetTime(ctx context.Context, owner, repo string, index int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/times", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Delete(ctx, path)
}

// AddLabels adds labels to an issue.
func (s *IssueService) AddLabels(ctx context.Context, owner, repo string, index int64, labelIDs []int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/labels", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	body := types.IssueLabelsOption{Labels: toAnySlice(labelIDs)}
	return s.client.Post(ctx, path, body, nil)
}

// RemoveLabel removes a single label from an issue.
func (s *IssueService) RemoveLabel(ctx context.Context, owner, repo string, index int64, labelID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/labels/{labelID}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "labelID", int64String(labelID)))
	return s.client.Delete(ctx, path)
}

// ListComments returns all comments on an issue.
func (s *IssueService) ListComments(ctx context.Context, owner, repo string, index int64) ([]types.Comment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/comments", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.Comment](ctx, s.client, path, nil)
}

// IterComments returns an iterator over all comments on an issue.
func (s *IssueService) IterComments(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Comment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/comments", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.Comment](ctx, s.client, path, nil)
}

// CreateComment creates a comment on an issue.
func (s *IssueService) CreateComment(ctx context.Context, owner, repo string, index int64, body string) (*types.Comment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/comments", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	opts := types.CreateIssueCommentOption{Body: body}
	var out types.Comment
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTimeline returns all comments and events on an issue.
func (s *IssueService) ListTimeline(ctx context.Context, owner, repo string, index int64, since, before *time.Time) ([]types.TimelineComment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/timeline", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	query := make(map[string]string, 2)
	if since != nil {
		query["since"] = since.Format(time.RFC3339)
	}
	if before != nil {
		query["before"] = before.Format(time.RFC3339)
	}
	if len(query) == 0 {
		query = nil
	}
	return ListAll[types.TimelineComment](ctx, s.client, path, query)
}

// IterTimeline returns an iterator over all comments and events on an issue.
func (s *IssueService) IterTimeline(ctx context.Context, owner, repo string, index int64, since, before *time.Time) iter.Seq2[types.TimelineComment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/timeline", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	query := make(map[string]string, 2)
	if since != nil {
		query["since"] = since.Format(time.RFC3339)
	}
	if before != nil {
		query["before"] = before.Format(time.RFC3339)
	}
	if len(query) == 0 {
		query = nil
	}
	return ListIter[types.TimelineComment](ctx, s.client, path, query)
}

// ListSubscriptions returns all users subscribed to an issue.
func (s *IssueService) ListSubscriptions(ctx context.Context, owner, repo string, index int64) ([]types.User, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/subscriptions", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterSubscriptions returns an iterator over all users subscribed to an issue.
func (s *IssueService) IterSubscriptions(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/subscriptions", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// CheckSubscription returns the authenticated user's subscription state for an issue.
func (s *IssueService) CheckSubscription(ctx context.Context, owner, repo string, index int64) (*types.WatchInfo, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/subscriptions/check", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	var out types.WatchInfo
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SubscribeUser subscribes a user to an issue.
func (s *IssueService) SubscribeUser(ctx context.Context, owner, repo string, index int64, user string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/subscriptions/{user}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "user", user))
	return s.client.Put(ctx, path, nil, nil)
}

// UnsubscribeUser unsubscribes a user from an issue.
func (s *IssueService) UnsubscribeUser(ctx context.Context, owner, repo string, index int64, user string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/subscriptions/{user}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "user", user))
	return s.client.Delete(ctx, path)
}

// toAnySlice converts a slice of int64 to a slice of any for IssueLabelsOption.
func toAnySlice(ids []int64) []any {
	out := make([]any, len(ids))
	for i, id := range ids {
		out[i] = id
	}
	return out
}
