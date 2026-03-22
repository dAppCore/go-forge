package forge

import (
	"context"
	"fmt"
	"iter"

	"dappco.re/go/core/forge/types"
)

// IssueService handles issue operations within a repository.
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

// Pin pins an issue.
func (s *IssueService) Pin(ctx context.Context, owner, repo string, index int64) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin", owner, repo, index)
	return s.client.Post(ctx, path, nil, nil)
}

// Unpin unpins an issue.
func (s *IssueService) Unpin(ctx context.Context, owner, repo string, index int64) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/pin", owner, repo, index)
	return s.client.Delete(ctx, path)
}

// SetDeadline sets or updates the deadline on an issue.
func (s *IssueService) SetDeadline(ctx context.Context, owner, repo string, index int64, deadline string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/deadline", owner, repo, index)
	body := map[string]string{"due_date": deadline}
	return s.client.Post(ctx, path, body, nil)
}

// AddReaction adds a reaction to an issue.
func (s *IssueService) AddReaction(ctx context.Context, owner, repo string, index int64, reaction string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/reactions", owner, repo, index)
	body := map[string]string{"content": reaction}
	return s.client.Post(ctx, path, body, nil)
}

// DeleteReaction removes a reaction from an issue.
func (s *IssueService) DeleteReaction(ctx context.Context, owner, repo string, index int64, reaction string) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/reactions", owner, repo, index)
	body := map[string]string{"content": reaction}
	return s.client.DeleteWithBody(ctx, path, body)
}

// StartStopwatch starts the stopwatch on an issue.
func (s *IssueService) StartStopwatch(ctx context.Context, owner, repo string, index int64) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/stopwatch/start", owner, repo, index)
	return s.client.Post(ctx, path, nil, nil)
}

// StopStopwatch stops the stopwatch on an issue.
func (s *IssueService) StopStopwatch(ctx context.Context, owner, repo string, index int64) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/stopwatch/stop", owner, repo, index)
	return s.client.Post(ctx, path, nil, nil)
}

// AddLabels adds labels to an issue.
func (s *IssueService) AddLabels(ctx context.Context, owner, repo string, index int64, labelIDs []int64) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/labels", owner, repo, index)
	body := types.IssueLabelsOption{Labels: toAnySlice(labelIDs)}
	return s.client.Post(ctx, path, body, nil)
}

// RemoveLabel removes a single label from an issue.
func (s *IssueService) RemoveLabel(ctx context.Context, owner, repo string, index int64, labelID int64) error {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/labels/%d", owner, repo, index, labelID)
	return s.client.Delete(ctx, path)
}

// ListComments returns all comments on an issue.
func (s *IssueService) ListComments(ctx context.Context, owner, repo string, index int64) ([]types.Comment, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/comments", owner, repo, index)
	return ListAll[types.Comment](ctx, s.client, path, nil)
}

// IterComments returns an iterator over all comments on an issue.
func (s *IssueService) IterComments(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Comment, error] {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/comments", owner, repo, index)
	return ListIter[types.Comment](ctx, s.client, path, nil)
}

// CreateComment creates a comment on an issue.
func (s *IssueService) CreateComment(ctx context.Context, owner, repo string, index int64, body string) (*types.Comment, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/comments", owner, repo, index)
	opts := types.CreateIssueCommentOption{Body: body}
	var out types.Comment
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// toAnySlice converts a slice of int64 to a slice of any for IssueLabelsOption.
func toAnySlice(ids []int64) []any {
	out := make([]any, len(ids))
	for i, id := range ids {
		out[i] = id
	}
	return out
}
