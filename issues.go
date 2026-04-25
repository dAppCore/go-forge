package forge

import (
	"context"
	"iter"
	"time"

	// Note: AX-6 intrinsic — upload APIs must expose the structural request body type; coreio Medium is used inside Client multipart handling.
	goio "io"

	"dappco.re/go/forge/types"
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

// IssueListOptions controls filtering for repository issue listings.
//
// Usage:
//
//	opts := forge.IssueListOptions{State: "open", Labels: "bug"}
type IssueListOptions struct {
	State       string
	Sort        string
	Labels      string
	Query       string
	Type        string
	Milestones  string
	Since       *time.Time
	Before      *time.Time
	CreatedBy   string
	AssignedBy  string
	MentionedBy string
}

// String returns a safe summary of the issue list filters.
func (o IssueListOptions) String() string {
	return optionString("forge.IssueListOptions",
		"state", o.State,
		"sort", o.Sort,
		"labels", o.Labels,
		"q", o.Query,
		"type", o.Type,
		"milestones", o.Milestones,
		"since", o.Since,
		"before", o.Before,
		"created_by", o.CreatedBy,
		"assigned_by", o.AssignedBy,
		"mentioned_by", o.MentionedBy,
	)
}

// GoString returns a safe Go-syntax summary of the issue list filters.
func (o IssueListOptions) GoString() string { return o.String() }

func (o IssueListOptions) queryParams() map[string]string {
	query := make(map[string]string, 10)
	if o.State != "" {
		query["state"] = o.State
	}
	if o.Sort != "" {
		query["sort"] = o.Sort
	}
	if o.Labels != "" {
		query["labels"] = o.Labels
	}
	if o.Query != "" {
		query["q"] = o.Query
	}
	if o.Type != "" {
		query["type"] = o.Type
	}
	if o.Milestones != "" {
		query["milestones"] = o.Milestones
	}
	if o.Since != nil {
		query["since"] = o.Since.Format(time.RFC3339)
	}
	if o.Before != nil {
		query["before"] = o.Before.Format(time.RFC3339)
	}
	if o.CreatedBy != "" {
		query["created_by"] = o.CreatedBy
	}
	if o.AssignedBy != "" {
		query["assigned_by"] = o.AssignedBy
	}
	if o.MentionedBy != "" {
		query["mentioned_by"] = o.MentionedBy
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

// AttachmentUploadOptions controls metadata sent when uploading an attachment.
//
// Usage:
//
//	opts := forge.AttachmentUploadOptions{Name: "screenshot.png"}
type AttachmentUploadOptions struct {
	Name      string
	UpdatedAt *time.Time
}

// String returns a safe summary of the attachment upload metadata.
func (o AttachmentUploadOptions) String() string {
	return optionString("forge.AttachmentUploadOptions",
		"name", o.Name,
		"updated_at", o.UpdatedAt,
	)
}

// GoString returns a safe Go-syntax summary of the attachment upload metadata.
func (o AttachmentUploadOptions) GoString() string { return o.String() }

// RepoCommentListOptions controls filtering for repository-wide issue comment listings.
//
// Usage:
//
//	opts := forge.RepoCommentListOptions{Page: 1, Limit: 50}
type RepoCommentListOptions struct {
	Since  *time.Time
	Before *time.Time
}

// String returns a safe summary of the repository comment filters.
func (o RepoCommentListOptions) String() string {
	return optionString("forge.RepoCommentListOptions",
		"since", o.Since,
		"before", o.Before,
	)
}

// GoString returns a safe Go-syntax summary of the repository comment filters.
func (o RepoCommentListOptions) GoString() string { return o.String() }

func (o RepoCommentListOptions) queryParams() map[string]string {
	query := make(map[string]string, 2)
	if o.Since != nil {
		query["since"] = o.Since.Format(time.RFC3339)
	}
	if o.Before != nil {
		query["before"] = o.Before.Format(time.RFC3339)
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func newIssueService(c *Client) *IssueService {
	return &IssueService{
		Resource: *NewResource[types.Issue, types.CreateIssueOption, types.EditIssueOption](
			c, "/api/v1/repos/{owner}/{repo}/issues/{index}",
		),
	}
}

// GetIssue returns a single issue by index.
func (s *IssueService) GetIssue(ctx context.Context, owner, repo string, index int64) (*types.Issue, error) {
	return s.Get(ctx, pathParams("owner", owner, "repo", repo, "index", int64String(index)))
}

// EditIssue updates an existing issue.
func (s *IssueService) EditIssue(ctx context.Context, owner, repo string, index int64, opts *types.EditIssueOption) (*types.Issue, error) {
	return s.Update(ctx, pathParams("owner", owner, "repo", repo, "index", int64String(index)), opts)
}

// DeleteIssue deletes an issue.
func (s *IssueService) DeleteIssue(ctx context.Context, owner, repo string, index int64) error {
	return s.Delete(ctx, pathParams("owner", owner, "repo", repo, "index", int64String(index)))
}

// SearchIssuesOptions controls filtering for the global issue search endpoint.
//
// Usage:
//
//	opts := forge.SearchIssuesOptions{State: "open"}
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

// String returns a safe summary of the issue search filters.
func (o SearchIssuesOptions) String() string {
	return optionString("forge.SearchIssuesOptions",
		"state", o.State,
		"labels", o.Labels,
		"milestones", o.Milestones,
		"q", o.Query,
		"priority_repo_id", o.PriorityRepoID,
		"type", o.Type,
		"since", o.Since,
		"before", o.Before,
		"assigned", o.Assigned,
		"created", o.Created,
		"mentioned", o.Mentioned,
		"review_requested", o.ReviewRequested,
		"reviewed", o.Reviewed,
		"owner", o.Owner,
		"team", o.Team,
	)
}

// GoString returns a safe Go-syntax summary of the issue search filters.
func (o SearchIssuesOptions) GoString() string { return o.String() }

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
		query["priority_repo_id"] = int64String(o.PriorityRepoID)
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
		query["assigned"] = "true"
	}
	if o.Created {
		query["created"] = "true"
	}
	if o.Mentioned {
		query["mentioned"] = "true"
	}
	if o.ReviewRequested {
		query["review_requested"] = "true"
	}
	if o.Reviewed {
		query["reviewed"] = "true"
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

// ListIssuesPage returns a single page of issues in a repository.
func (s *IssueService) ListIssuesPage(ctx context.Context, owner, repo string, opts ListOptions, filters ...any) (*PagedResult[types.Issue], error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues", pathParams("owner", owner, "repo", repo))
	return ListPage[types.Issue](ctx, s.client, path, issueListQuery(filters...), opts)
}

// ListIssues returns all issues in a repository.
func (s *IssueService) ListIssues(ctx context.Context, owner, repo string, filters ...any) ([]types.Issue, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues", pathParams("owner", owner, "repo", repo))
	if pageOpts, ok := issueListPageOptions(filters...); ok {
		page, err := ListPage[types.Issue](ctx, s.client, path, issueListQuery(filters...), pageOpts)
		if err != nil {
			return nil, err
		}
		return page.Items, nil
	}
	return ListAll[types.Issue](ctx, s.client, path, issueListQuery(filters...))
}

// IterIssues returns an iterator over all issues in a repository.
func (s *IssueService) IterIssues(ctx context.Context, owner, repo string, filters ...any) iter.Seq2[types.Issue, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues", pathParams("owner", owner, "repo", repo))
	if pageOpts, ok := issueListPageOptions(filters...); ok {
		return func(yield func(types.Issue, error) bool) {
			page, err := ListPage[types.Issue](ctx, s.client, path, issueListQuery(filters...), pageOpts)
			if err != nil {
				yield(*new(types.Issue), err)
				return
			}
			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}
		}
	}
	return ListIter[types.Issue](ctx, s.client, path, issueListQuery(filters...))
}

// ListRepoIssues returns all issues in a repository.
func (s *IssueService) ListRepoIssues(ctx context.Context, owner, repo string, filters ...any) ([]types.Issue, error) {
	return s.ListIssues(ctx, owner, repo, filters...)
}

// ListRepoIssuesPage returns a single page of issues in a repository.
func (s *IssueService) ListRepoIssuesPage(ctx context.Context, owner, repo string, opts ListOptions, filters ...any) (*PagedResult[types.Issue], error) {
	return s.ListIssuesPage(ctx, owner, repo, opts, filters...)
}

// IterRepoIssues returns an iterator over all issues in a repository.
func (s *IssueService) IterRepoIssues(ctx context.Context, owner, repo string, filters ...any) iter.Seq2[types.Issue, error] {
	return s.IterIssues(ctx, owner, repo, filters...)
}

// CreateIssue creates a new issue in a repository.
func (s *IssueService) CreateIssue(ctx context.Context, owner, repo string, opts *types.CreateIssueOption) (*types.Issue, error) {
	var out types.Issue
	if err := s.client.Post(ctx, ResolvePath("/api/v1/repos/{owner}/{repo}/issues", pathParams("owner", owner, "repo", repo)), opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
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

// ListPinnedIssues returns all pinned issues in a repository.
func (s *IssueService) ListPinnedIssues(ctx context.Context, owner, repo string) ([]types.Issue, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/pinned", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Issue](ctx, s.client, path, nil)
}

// IterPinnedIssues returns an iterator over all pinned issues in a repository.
func (s *IssueService) IterPinnedIssues(ctx context.Context, owner, repo string) iter.Seq2[types.Issue, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/pinned", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Issue](ctx, s.client, path, nil)
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
	body := types.EditReactionOption{Reaction: reaction}
	return s.client.Post(ctx, path, body, nil)
}

// ListReactions returns all reactions on an issue.
func (s *IssueService) ListReactions(ctx context.Context, owner, repo string, index int64) ([]types.Reaction, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/reactions", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.Reaction](ctx, s.client, path, nil)
}

// IterReactions returns an iterator over all reactions on an issue.
func (s *IssueService) IterReactions(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Reaction, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/reactions", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.Reaction](ctx, s.client, path, nil)
}

// DeleteReaction removes a reaction from an issue.
func (s *IssueService) DeleteReaction(ctx context.Context, owner, repo string, index int64, reaction string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/reactions", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	body := types.EditReactionOption{Reaction: reaction}
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
	return ListAll[types.TrackedTime](ctx, s.client, path, issueTimeQuery(user, since, before))
}

// IterTimes returns an iterator over all tracked times on an issue.
func (s *IssueService) IterTimes(ctx context.Context, owner, repo string, index int64, user string, since, before *time.Time) iter.Seq2[types.TrackedTime, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/times", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.TrackedTime](ctx, s.client, path, issueTimeQuery(user, since, before))
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

// DeleteTime removes a specific tracked time entry from an issue.
func (s *IssueService) DeleteTime(ctx context.Context, owner, repo string, index, timeID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/times/{id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(timeID)))
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
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/labels/{id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(labelID)))
	return s.client.Delete(ctx, path)
}

// ListComments returns all comments on an issue.
func (s *IssueService) ListComments(ctx context.Context, owner, repo string, index int64) ([]types.Comment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/comments", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.Comment](ctx, s.client, path, nil)
}

// ListIssueComments returns all comments on an issue.
func (s *IssueService) ListIssueComments(ctx context.Context, owner, repo string, index int64) ([]types.Comment, error) {
	return s.ListComments(ctx, owner, repo, index)
}

// IterComments returns an iterator over all comments on an issue.
func (s *IssueService) IterComments(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Comment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/comments", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.Comment](ctx, s.client, path, nil)
}

// IterIssueComments returns an iterator over all comments on an issue.
func (s *IssueService) IterIssueComments(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Comment, error] {
	return s.IterComments(ctx, owner, repo, index)
}

// GetIssueComment returns a single comment on an issue.
func (s *IssueService) GetIssueComment(ctx context.Context, owner, repo string, index, id int64) (*types.Comment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/comments/{id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(id)))
	var out types.Comment
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditIssueComment updates an issue comment.
func (s *IssueService) EditIssueComment(ctx context.Context, owner, repo string, index, id int64, opts *types.EditIssueCommentOption) (*types.Comment, error) {
	return s.EditComment(ctx, owner, repo, index, id, opts)
}

// DeleteIssueComment deletes an issue comment.
func (s *IssueService) DeleteIssueComment(ctx context.Context, owner, repo string, index, id int64) error {
	return s.DeleteComment(ctx, owner, repo, index, id)
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

// EditComment updates an issue comment.
func (s *IssueService) EditComment(ctx context.Context, owner, repo string, index, id int64, opts *types.EditIssueCommentOption) (*types.Comment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/comments/{id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(id)))
	var out types.Comment
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteComment deletes an issue comment.
func (s *IssueService) DeleteComment(ctx context.Context, owner, repo string, index, id int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/comments/{id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListRepoComments returns all comments in a repository.
func (s *IssueService) ListRepoComments(ctx context.Context, owner, repo string, filters ...RepoCommentListOptions) ([]types.Comment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Comment](ctx, s.client, path, repoCommentQuery(filters...))
}

// IterRepoComments returns an iterator over all comments in a repository.
func (s *IssueService) IterRepoComments(ctx context.Context, owner, repo string, filters ...RepoCommentListOptions) iter.Seq2[types.Comment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Comment](ctx, s.client, path, repoCommentQuery(filters...))
}

// GetRepoComment returns a single comment in a repository.
func (s *IssueService) GetRepoComment(ctx context.Context, owner, repo string, id int64) (*types.Comment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.Comment
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditRepoComment updates a repository comment.
func (s *IssueService) EditRepoComment(ctx context.Context, owner, repo string, id int64, opts *types.EditIssueCommentOption) (*types.Comment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.Comment
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteRepoComment deletes a repository comment.
func (s *IssueService) DeleteRepoComment(ctx context.Context, owner, repo string, id int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListCommentReactions returns all reactions on an issue comment.
func (s *IssueService) ListCommentReactions(ctx context.Context, owner, repo string, id int64) ([]types.Reaction, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/reactions", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return ListAll[types.Reaction](ctx, s.client, path, nil)
}

// IterCommentReactions returns an iterator over all reactions on an issue comment.
func (s *IssueService) IterCommentReactions(ctx context.Context, owner, repo string, id int64) iter.Seq2[types.Reaction, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/reactions", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return ListIter[types.Reaction](ctx, s.client, path, nil)
}

// AddCommentReaction adds a reaction to an issue comment.
func (s *IssueService) AddCommentReaction(ctx context.Context, owner, repo string, id int64, reaction string) (*types.Reaction, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/reactions", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.Reaction
	if err := s.client.Post(ctx, path, types.EditReactionOption{Reaction: reaction}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCommentReaction removes a reaction from an issue comment.
func (s *IssueService) DeleteCommentReaction(ctx context.Context, owner, repo string, id int64, reaction string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/reactions", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return s.client.DeleteWithBody(ctx, path, types.EditReactionOption{Reaction: reaction})
}

func issueListQuery(filters ...any) map[string]string {
	query := make(map[string]string, len(filters))
	for _, filter := range filters {
		switch v := filter.(type) {
		case IssueListOptions:
			for key, value := range issueListQueryFromOption(v) {
				query[key] = value
			}
		case *IssueListOptions:
			if v != nil {
				for key, value := range issueListQueryFromOption(*v) {
					query[key] = value
				}
			}
		case types.ListIssueOption:
			for key, value := range issueListQueryFromCompat(v) {
				query[key] = value
			}
		case *types.ListIssueOption:
			if v != nil {
				for key, value := range issueListQueryFromCompat(*v) {
					query[key] = value
				}
			}
		}
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func issueListQueryFromOption(filter IssueListOptions) map[string]string {
	query := make(map[string]string, 10)
	if filter.State != "" {
		query["state"] = filter.State
	}
	if filter.Sort != "" {
		query["sort"] = filter.Sort
	}
	if filter.Labels != "" {
		query["labels"] = filter.Labels
	}
	if filter.Query != "" {
		query["q"] = filter.Query
	}
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	if filter.Milestones != "" {
		query["milestones"] = filter.Milestones
	}
	if filter.Since != nil {
		query["since"] = filter.Since.Format(time.RFC3339)
	}
	if filter.Before != nil {
		query["before"] = filter.Before.Format(time.RFC3339)
	}
	if filter.CreatedBy != "" {
		query["created_by"] = filter.CreatedBy
	}
	if filter.AssignedBy != "" {
		query["assigned_by"] = filter.AssignedBy
	}
	if filter.MentionedBy != "" {
		query["mentioned_by"] = filter.MentionedBy
	}
	return query
}

func issueListQueryFromCompat(filter types.ListIssueOption) map[string]string {
	query := make(map[string]string, 10)
	if filter.State != "" {
		query["state"] = filter.State
	}
	if filter.Sort != "" {
		query["sort"] = filter.Sort
	}
	if filter.Labels != "" {
		query["labels"] = filter.Labels
	}
	if filter.Query != "" {
		query["q"] = filter.Query
	}
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	if filter.Milestones != "" {
		query["milestones"] = filter.Milestones
	}
	if filter.Since != nil {
		query["since"] = filter.Since.Format(time.RFC3339)
	}
	if filter.Before != nil {
		query["before"] = filter.Before.Format(time.RFC3339)
	}
	if filter.CreatedBy != "" {
		query["created_by"] = filter.CreatedBy
	}
	if filter.AssignedBy != "" {
		query["assigned_by"] = filter.AssignedBy
	}
	if filter.MentionedBy != "" {
		query["mentioned_by"] = filter.MentionedBy
	}
	return query
}

func issueListPageOptions(filters ...any) (ListOptions, bool) {
	for _, filter := range filters {
		switch v := filter.(type) {
		case types.ListIssueOption:
			if opts, ok := compatListOptions(v.Page, v.PageSize, v.Limit); ok {
				return opts, true
			}
		case *types.ListIssueOption:
			if v != nil {
				if opts, ok := compatListOptions(v.Page, v.PageSize, v.Limit); ok {
					return opts, true
				}
			}
		}
	}
	return ListOptions{}, false
}

func attachmentUploadQuery(opts *AttachmentUploadOptions) map[string]string {
	if opts == nil {
		return nil
	}
	query := make(map[string]string, 2)
	if opts.Name != "" {
		query["name"] = opts.Name
	}
	if opts.UpdatedAt != nil {
		query["updated_at"] = opts.UpdatedAt.Format(time.RFC3339)
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func (s *IssueService) createAttachment(ctx context.Context, path string, opts *AttachmentUploadOptions, filename string, content goio.Reader) (*types.Attachment, error) {
	var out types.Attachment
	if err := s.client.postMultipartJSON(ctx, path, attachmentUploadQuery(opts), nil, "attachment", filename, content, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateAttachment uploads a new attachment to an issue.
func (s *IssueService) CreateAttachment(ctx context.Context, owner, repo string, index int64, opts *AttachmentUploadOptions, filename string, content goio.Reader) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/assets", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.createAttachment(ctx, path, opts, filename, content)
}

// ListAttachments returns all attachments on an issue.
func (s *IssueService) ListAttachments(ctx context.Context, owner, repo string, index int64) ([]types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/assets", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.Attachment](ctx, s.client, path, nil)
}

// IterAttachments returns an iterator over all attachments on an issue.
func (s *IssueService) IterAttachments(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Attachment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/assets", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.Attachment](ctx, s.client, path, nil)
}

// GetAttachment returns a single attachment on an issue.
func (s *IssueService) GetAttachment(ctx context.Context, owner, repo string, index, attachmentID int64) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "attachment_id", int64String(attachmentID)))
	var out types.Attachment
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditAttachment updates an issue attachment.
func (s *IssueService) EditAttachment(ctx context.Context, owner, repo string, index, attachmentID int64, opts *types.EditAttachmentOptions) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "attachment_id", int64String(attachmentID)))
	var out types.Attachment
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAttachment removes an issue attachment.
func (s *IssueService) DeleteAttachment(ctx context.Context, owner, repo string, index, attachmentID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/assets/{attachment_id}", pathParams("owner", owner, "repo", repo, "index", int64String(index), "attachment_id", int64String(attachmentID)))
	return s.client.Delete(ctx, path)
}

// ListCommentAttachments returns all attachments on an issue comment.
func (s *IssueService) ListCommentAttachments(ctx context.Context, owner, repo string, id int64) ([]types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/assets", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return ListAll[types.Attachment](ctx, s.client, path, nil)
}

// IterCommentAttachments returns an iterator over all attachments on an issue comment.
func (s *IssueService) IterCommentAttachments(ctx context.Context, owner, repo string, id int64) iter.Seq2[types.Attachment, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/assets", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return ListIter[types.Attachment](ctx, s.client, path, nil)
}

// GetCommentAttachment returns a single attachment on an issue comment.
func (s *IssueService) GetCommentAttachment(ctx context.Context, owner, repo string, id, attachmentID int64) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/assets/{attachment_id}", pathParams("owner", owner, "repo", repo, "id", int64String(id), "attachment_id", int64String(attachmentID)))
	var out types.Attachment
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCommentAttachment uploads a new attachment to an issue comment.
func (s *IssueService) CreateCommentAttachment(ctx context.Context, owner, repo string, id int64, opts *AttachmentUploadOptions, filename string, content goio.Reader) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/assets", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return s.createAttachment(ctx, path, opts, filename, content)
}

// EditCommentAttachment updates an issue comment attachment.
func (s *IssueService) EditCommentAttachment(ctx context.Context, owner, repo string, id, attachmentID int64, opts *types.EditAttachmentOptions) (*types.Attachment, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/assets/{attachment_id}", pathParams("owner", owner, "repo", repo, "id", int64String(id), "attachment_id", int64String(attachmentID)))
	var out types.Attachment
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteCommentAttachment removes an issue comment attachment.
func (s *IssueService) DeleteCommentAttachment(ctx context.Context, owner, repo string, id, attachmentID int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/comments/{id}/assets/{attachment_id}", pathParams("owner", owner, "repo", repo, "id", int64String(id), "attachment_id", int64String(attachmentID)))
	return s.client.Delete(ctx, path)
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

// ListDependencies returns all issues that block the given issue.
func (s *IssueService) ListDependencies(ctx context.Context, owner, repo string, index int64) ([]types.Issue, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/dependencies", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.Issue](ctx, s.client, path, nil)
}

// IterDependencies returns an iterator over all issues that block the given issue.
func (s *IssueService) IterDependencies(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Issue, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/dependencies", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.Issue](ctx, s.client, path, nil)
}

// AddDependency makes another issue block the issue at the given path.
func (s *IssueService) AddDependency(ctx context.Context, owner, repo string, index int64, dependency types.IssueMeta) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/dependencies", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Post(ctx, path, dependency, nil)
}

// RemoveDependency removes an issue dependency from the issue at the given path.
func (s *IssueService) RemoveDependency(ctx context.Context, owner, repo string, index int64, dependency types.IssueMeta) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/dependencies", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.DeleteWithBody(ctx, path, dependency)
}

// ListBlocks returns all issues blocked by the given issue.
func (s *IssueService) ListBlocks(ctx context.Context, owner, repo string, index int64) ([]types.Issue, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/blocks", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListAll[types.Issue](ctx, s.client, path, nil)
}

// IterBlocks returns an iterator over all issues blocked by the given issue.
func (s *IssueService) IterBlocks(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Issue, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/blocks", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return ListIter[types.Issue](ctx, s.client, path, nil)
}

// AddBlock makes the issue at the given path block another issue.
func (s *IssueService) AddBlock(ctx context.Context, owner, repo string, index int64, blockedIssue types.IssueMeta) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/blocks", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.Post(ctx, path, blockedIssue, nil)
}

// RemoveBlock removes an issue block from the issue at the given path.
func (s *IssueService) RemoveBlock(ctx context.Context, owner, repo string, index int64, blockedIssue types.IssueMeta) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}/blocks", pathParams("owner", owner, "repo", repo, "index", int64String(index)))
	return s.client.DeleteWithBody(ctx, path, blockedIssue)
}

// toAnySlice converts a slice of int64 to a slice of any for IssueLabelsOption.
func toAnySlice(ids []int64) []any {
	out := make([]any, len(ids))
	for i, id := range ids {
		out[i] = id
	}
	return out
}

func repoCommentQuery(filters ...RepoCommentListOptions) map[string]string {
	if len(filters) == 0 {
		return nil
	}

	query := make(map[string]string, 2)
	for _, filter := range filters {
		for key, value := range filter.queryParams() {
			query[key] = value
		}
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func issueTimeQuery(user string, since, before *time.Time) map[string]string {
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
		return nil
	}
	return query
}
