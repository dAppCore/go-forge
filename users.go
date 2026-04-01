package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// UserService handles user operations.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Users.GetCurrent(ctx)
type UserService struct {
	Resource[types.User, struct{}, struct{}]
}

func newUserService(c *Client) *UserService {
	return &UserService{
		Resource: *NewResource[types.User, struct{}, struct{}](
			c, "/api/v1/users/{username}",
		),
	}
}

// GetCurrent returns the authenticated user.
func (s *UserService) GetCurrent(ctx context.Context) (*types.User, error) {
	var out types.User
	if err := s.client.Get(ctx, "/api/v1/user", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListEmails returns all email addresses for the authenticated user.
func (s *UserService) ListEmails(ctx context.Context) ([]types.Email, error) {
	return ListAll[types.Email](ctx, s.client, "/api/v1/user/emails", nil)
}

// IterEmails returns an iterator over all email addresses for the authenticated user.
func (s *UserService) IterEmails(ctx context.Context) iter.Seq2[types.Email, error] {
	return ListIter[types.Email](ctx, s.client, "/api/v1/user/emails", nil)
}

// AddEmails adds email addresses for the authenticated user.
func (s *UserService) AddEmails(ctx context.Context, emails ...string) ([]types.Email, error) {
	var out []types.Email
	if err := s.client.Post(ctx, "/api/v1/user/emails", types.CreateEmailOption{Emails: emails}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteEmails deletes email addresses for the authenticated user.
func (s *UserService) DeleteEmails(ctx context.Context, emails ...string) error {
	return s.client.DeleteWithBody(ctx, "/api/v1/user/emails", types.DeleteEmailOption{Emails: emails})
}

// ListStopwatches returns all existing stopwatches for the authenticated user.
func (s *UserService) ListStopwatches(ctx context.Context) ([]types.StopWatch, error) {
	return ListAll[types.StopWatch](ctx, s.client, "/api/v1/user/stopwatches", nil)
}

// IterStopwatches returns an iterator over all existing stopwatches for the authenticated user.
func (s *UserService) IterStopwatches(ctx context.Context) iter.Seq2[types.StopWatch, error] {
	return ListIter[types.StopWatch](ctx, s.client, "/api/v1/user/stopwatches", nil)
}

// ListMySubscriptions returns all repositories watched by the authenticated user.
func (s *UserService) ListMySubscriptions(ctx context.Context) ([]types.Repository, error) {
	return ListAll[types.Repository](ctx, s.client, "/api/v1/user/subscriptions", nil)
}

// IterMySubscriptions returns an iterator over all repositories watched by the authenticated user.
func (s *UserService) IterMySubscriptions(ctx context.Context) iter.Seq2[types.Repository, error] {
	return ListIter[types.Repository](ctx, s.client, "/api/v1/user/subscriptions", nil)
}

// ListFollowers returns all followers of a user.
func (s *UserService) ListFollowers(ctx context.Context, username string) ([]types.User, error) {
	path := ResolvePath("/api/v1/users/{username}/followers", pathParams("username", username))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterFollowers returns an iterator over all followers of a user.
func (s *UserService) IterFollowers(ctx context.Context, username string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/users/{username}/followers", pathParams("username", username))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// ListSubscriptions returns all repositories watched by a user.
func (s *UserService) ListSubscriptions(ctx context.Context, username string) ([]types.Repository, error) {
	path := ResolvePath("/api/v1/users/{username}/subscriptions", pathParams("username", username))
	return ListAll[types.Repository](ctx, s.client, path, nil)
}

// IterSubscriptions returns an iterator over all repositories watched by a user.
func (s *UserService) IterSubscriptions(ctx context.Context, username string) iter.Seq2[types.Repository, error] {
	path := ResolvePath("/api/v1/users/{username}/subscriptions", pathParams("username", username))
	return ListIter[types.Repository](ctx, s.client, path, nil)
}

// ListFollowing returns all users that a user is following.
func (s *UserService) ListFollowing(ctx context.Context, username string) ([]types.User, error) {
	path := ResolvePath("/api/v1/users/{username}/following", pathParams("username", username))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterFollowing returns an iterator over all users that a user is following.
func (s *UserService) IterFollowing(ctx context.Context, username string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/users/{username}/following", pathParams("username", username))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// Follow follows a user as the authenticated user.
func (s *UserService) Follow(ctx context.Context, username string) error {
	path := ResolvePath("/api/v1/user/following/{username}", pathParams("username", username))
	return s.client.Put(ctx, path, nil, nil)
}

// Unfollow unfollows a user as the authenticated user.
func (s *UserService) Unfollow(ctx context.Context, username string) error {
	path := ResolvePath("/api/v1/user/following/{username}", pathParams("username", username))
	return s.client.Delete(ctx, path)
}

// ListStarred returns all repositories starred by a user.
func (s *UserService) ListStarred(ctx context.Context, username string) ([]types.Repository, error) {
	path := ResolvePath("/api/v1/users/{username}/starred", pathParams("username", username))
	return ListAll[types.Repository](ctx, s.client, path, nil)
}

// IterStarred returns an iterator over all repositories starred by a user.
func (s *UserService) IterStarred(ctx context.Context, username string) iter.Seq2[types.Repository, error] {
	path := ResolvePath("/api/v1/users/{username}/starred", pathParams("username", username))
	return ListIter[types.Repository](ctx, s.client, path, nil)
}

// Star stars a repository as the authenticated user.
func (s *UserService) Star(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/user/starred/{owner}/{repo}", pathParams("owner", owner, "repo", repo))
	return s.client.Put(ctx, path, nil, nil)
}

// Unstar unstars a repository as the authenticated user.
func (s *UserService) Unstar(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/user/starred/{owner}/{repo}", pathParams("owner", owner, "repo", repo))
	return s.client.Delete(ctx, path)
}
