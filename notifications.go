package forge

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go-forge/types"
)

// NotificationService handles notification operations via the Forgejo API.
// No Resource embedding — varied endpoint shapes.
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

// ListRepo returns all notifications for a specific repository.
func (s *NotificationService) ListRepo(ctx context.Context, owner, repo string) ([]types.NotificationThread, error) {
	path := fmt.Sprintf("/api/v1/repos/%s/%s/notifications", owner, repo)
	return ListAll[types.NotificationThread](ctx, s.client, path, nil)
}

// MarkRead marks all notifications as read.
func (s *NotificationService) MarkRead(ctx context.Context) error {
	return s.client.Put(ctx, "/api/v1/notifications", nil, nil)
}

// GetThread returns a single notification thread by ID.
func (s *NotificationService) GetThread(ctx context.Context, id int64) (*types.NotificationThread, error) {
	path := fmt.Sprintf("/api/v1/notifications/threads/%d", id)
	var out types.NotificationThread
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MarkThreadRead marks a single notification thread as read.
func (s *NotificationService) MarkThreadRead(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/v1/notifications/threads/%d", id)
	return s.client.Patch(ctx, path, nil, nil)
}
