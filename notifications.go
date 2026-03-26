package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// NotificationService handles notification operations via the Forgejo API.
// No Resource embedding — varied endpoint shapes.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Notifications.List(ctx)
type NotificationService struct {
	client *Client
}

func newNotificationService(c *Client) *NotificationService {
	return &NotificationService{client: c}
}

// List returns all notifications for the authenticated user.
func (s *NotificationService) List(ctx context.Context) ([]types.NotificationThread, error) {
	return ListAll[types.NotificationThread](ctx, s.client, "/api/v1/notifications", nil)
}

// Iter returns an iterator over all notifications for the authenticated user.
func (s *NotificationService) Iter(ctx context.Context) iter.Seq2[types.NotificationThread, error] {
	return ListIter[types.NotificationThread](ctx, s.client, "/api/v1/notifications", nil)
}

// ListRepo returns all notifications for a specific repository.
func (s *NotificationService) ListRepo(ctx context.Context, owner, repo string) ([]types.NotificationThread, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/notifications", pathParams("owner", owner, "repo", repo))
	return ListAll[types.NotificationThread](ctx, s.client, path, nil)
}

// IterRepo returns an iterator over all notifications for a specific repository.
func (s *NotificationService) IterRepo(ctx context.Context, owner, repo string) iter.Seq2[types.NotificationThread, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/notifications", pathParams("owner", owner, "repo", repo))
	return ListIter[types.NotificationThread](ctx, s.client, path, nil)
}

// MarkRead marks all notifications as read.
func (s *NotificationService) MarkRead(ctx context.Context) error {
	return s.client.Put(ctx, "/api/v1/notifications", nil, nil)
}

// GetThread returns a single notification thread by ID.
func (s *NotificationService) GetThread(ctx context.Context, id int64) (*types.NotificationThread, error) {
	path := ResolvePath("/api/v1/notifications/threads/{id}", pathParams("id", int64String(id)))
	var out types.NotificationThread
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MarkThreadRead marks a single notification thread as read.
func (s *NotificationService) MarkThreadRead(ctx context.Context, id int64) error {
	path := ResolvePath("/api/v1/notifications/threads/{id}", pathParams("id", int64String(id)))
	return s.client.Patch(ctx, path, nil, nil)
}
