package forge

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"

	core "dappco.re/go/core"
	"dappco.re/go/core/forge/types"
)

// NotificationListOptions controls filtering for notification listings.
//
// Usage:
//
//	opts := forge.NotificationListOptions{All: true, StatusTypes: []string{"unread"}}
type NotificationListOptions struct {
	All          bool
	StatusTypes  []string
	SubjectTypes []string
	Since        *time.Time
	Before       *time.Time
}

// String returns a safe summary of the notification filters.
func (o NotificationListOptions) String() string {
	return optionString("forge.NotificationListOptions",
		"all", o.All,
		"status_types", o.StatusTypes,
		"subject_types", o.SubjectTypes,
		"since", o.Since,
		"before", o.Before,
	)
}

// GoString returns a safe Go-syntax summary of the notification filters.
func (o NotificationListOptions) GoString() string { return o.String() }

func (o NotificationListOptions) addQuery(values url.Values) {
	if o.All {
		values.Set("all", "true")
	}
	for _, status := range o.StatusTypes {
		if status != "" {
			values.Add("status-types", status)
		}
	}
	for _, subjectType := range o.SubjectTypes {
		if subjectType != "" {
			values.Add("subject-type", subjectType)
		}
	}
	if o.Since != nil {
		values.Set("since", o.Since.Format(time.RFC3339))
	}
	if o.Before != nil {
		values.Set("before", o.Before.Format(time.RFC3339))
	}
}

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
//
// Usage:
//
//	opts := forge.NotificationRepoMarkOptions{All: true, ToStatus: "read"}
type NotificationRepoMarkOptions struct {
	All         bool
	StatusTypes []string
	ToStatus    string
	LastReadAt  *time.Time
}

// String returns a safe summary of the repository notification mark options.
func (o NotificationRepoMarkOptions) String() string {
	return optionString("forge.NotificationRepoMarkOptions",
		"all", o.All,
		"status_types", o.StatusTypes,
		"to_status", o.ToStatus,
		"last_read_at", o.LastReadAt,
	)
}

// GoString returns a safe Go-syntax summary of the repository notification mark options.
func (o NotificationRepoMarkOptions) GoString() string { return o.String() }

// NotificationMarkOptions controls how authenticated-user notifications are marked.
//
// Usage:
//
//	opts := forge.NotificationMarkOptions{All: true, ToStatus: "read"}
type NotificationMarkOptions struct {
	All         bool
	StatusTypes []string
	ToStatus    string
	LastReadAt  *time.Time
}

// String returns a safe summary of the authenticated-user notification mark options.
func (o NotificationMarkOptions) String() string {
	return optionString("forge.NotificationMarkOptions",
		"all", o.All,
		"status_types", o.StatusTypes,
		"to_status", o.ToStatus,
		"last_read_at", o.LastReadAt,
	)
}

// GoString returns a safe Go-syntax summary of the authenticated-user notification mark options.
func (o NotificationMarkOptions) GoString() string { return o.String() }

func newNotificationService(c *Client) *NotificationService {
	return &NotificationService{client: c}
}

func notificationMarkQueryString(all bool, statusTypes []string, toStatus string, lastReadAt *time.Time) string {
	values := url.Values{}
	if all {
		values.Set("all", "true")
	}
	for _, status := range statusTypes {
		if status != "" {
			values.Add("status-types", status)
		}
	}
	if toStatus != "" {
		values.Set("to-status", toStatus)
	}
	if lastReadAt != nil {
		values.Set("last_read_at", lastReadAt.Format(time.RFC3339))
	}
	return values.Encode()
}

func (o NotificationRepoMarkOptions) queryString() string {
	return notificationMarkQueryString(o.All, o.StatusTypes, o.ToStatus, o.LastReadAt)
}

func (o NotificationMarkOptions) queryString() string {
	return notificationMarkQueryString(o.All, o.StatusTypes, o.ToStatus, o.LastReadAt)
}

// List returns all notifications for the authenticated user.
func (s *NotificationService) List(ctx context.Context, filters ...NotificationListOptions) ([]types.NotificationThread, error) {
	return s.listAll(ctx, "/api/v1/notifications", filters...)
}

// Iter returns an iterator over all notifications for the authenticated user.
func (s *NotificationService) Iter(ctx context.Context, filters ...NotificationListOptions) iter.Seq2[types.NotificationThread, error] {
	return s.listIter(ctx, "/api/v1/notifications", filters...)
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
func (s *NotificationService) ListRepo(ctx context.Context, owner, repo string, filters ...NotificationListOptions) ([]types.NotificationThread, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/notifications", pathParams("owner", owner, "repo", repo))
	return s.listAll(ctx, path, filters...)
}

// IterRepo returns an iterator over all notifications for a specific repository.
func (s *NotificationService) IterRepo(ctx context.Context, owner, repo string, filters ...NotificationListOptions) iter.Seq2[types.NotificationThread, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/notifications", pathParams("owner", owner, "repo", repo))
	return s.listIter(ctx, path, filters...)
}

// MarkNotifications marks authenticated-user notification threads as read, pinned, or unread.
func (s *NotificationService) MarkNotifications(ctx context.Context, opts *NotificationMarkOptions) ([]types.NotificationThread, error) {
	path := "/api/v1/notifications"
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

func (s *NotificationService) listAll(ctx context.Context, path string, filters ...NotificationListOptions) ([]types.NotificationThread, error) {
	var all []types.NotificationThread
	page := 1

	for {
		result, err := s.listPage(ctx, path, ListOptions{Page: page, Limit: defaultPageLimit}, filters...)
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

func (s *NotificationService) listIter(ctx context.Context, path string, filters ...NotificationListOptions) iter.Seq2[types.NotificationThread, error] {
	return func(yield func(types.NotificationThread, error) bool) {
		page := 1
		for {
			result, err := s.listPage(ctx, path, ListOptions{Page: page, Limit: defaultPageLimit}, filters...)
			if err != nil {
				yield(*new(types.NotificationThread), err)
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

func (s *NotificationService) listPage(ctx context.Context, path string, opts ListOptions, filters ...NotificationListOptions) (*PagedResult[types.NotificationThread], error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 {
		opts.Limit = defaultPageLimit
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, core.E("NotificationService.listPage", "forge: parse path", err)
	}

	values := u.Query()
	values.Set("page", strconv.Itoa(opts.Page))
	values.Set("limit", strconv.Itoa(opts.Limit))
	for _, filter := range filters {
		filter.addQuery(values)
	}
	u.RawQuery = values.Encode()

	var items []types.NotificationThread
	resp, err := s.client.doJSON(ctx, http.MethodGet, u.String(), nil, &items)
	if err != nil {
		return nil, err
	}

	totalCount, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))
	return &PagedResult[types.NotificationThread]{
		Items:      items,
		TotalCount: totalCount,
		Page:       opts.Page,
		HasMore: (totalCount > 0 && (opts.Page-1)*opts.Limit+len(items) < totalCount) ||
			(totalCount == 0 && len(items) >= opts.Limit),
	}, nil
}
