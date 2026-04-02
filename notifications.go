package forge

import (
	"context"
	"iter"
	"net/url"
	"time"

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

// NotificationRepoMarkOptions controls how repository notifications are marked.
type NotificationRepoMarkOptions struct {
	All         bool
	StatusTypes []string
	ToStatus    string
	LastReadAt  *time.Time
}

func newNotificationService(c *Client) *NotificationService {
	return &NotificationService{client: c}
}

func (o NotificationRepoMarkOptions) queryString() string {
	values := url.Values{}
	if o.All {
		values.Set("all", "true")
	}
	for _, status := range o.StatusTypes {
		if status != "" {
			values.Add("status-types", status)
		}
	}
	if o.ToStatus != "" {
		values.Set("to-status", o.ToStatus)
	}
	if o.LastReadAt != nil {
		values.Set("last_read_at", o.LastReadAt.Format(time.RFC3339))
	}
	return values.Encode()
}

// List returns all notifications for the authenticated user.
func (s *NotificationService) List(ctx context.Context) ([]types.NotificationThread, error) {
	return ListAll[types.NotificationThread](ctx, s.client, "/api/v1/notifications", nil)
}

// Iter returns an iterator over all notifications for the authenticated user.
func (s *NotificationService) Iter(ctx context.Context) iter.Seq2[types.NotificationThread, error] {
	return ListIter[types.NotificationThread](ctx, s.client, "/api/v1/notifications", nil)
}

// NewAvailable returns the count of unread notifications for the authenticated user.
func (s *NotificationService) NewAvailable(ctx context.Context) (*types.NotificationCount, error) {
	var out types.NotificationCount
	if err := s.client.Get(ctx, "/api/v1/notifications/new", &out); err != nil {
		return nil, err
	}
	return &out, nil
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

// MarkRepoNotifications marks repository notification threads as read, unread, or pinned.
func (s *NotificationService) MarkRepoNotifications(ctx context.Context, owner, repo string, opts *NotificationRepoMarkOptions) ([]types.NotificationThread, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/notifications", pathParams("owner", owner, "repo", repo))
	if opts != nil {
		if query := opts.queryString(); query != "" {
			path += "?" + query
		}
	}
	var out []types.NotificationThread
	if err := s.client.Put(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
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
