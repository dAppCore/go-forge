package forge

import (
	"context"
	"iter"
	"net/http"

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

// GetSettings returns the authenticated user's settings.
func (s *UserService) GetSettings(ctx context.Context) (*types.UserSettings, error) {
	var out types.UserSettings
	if err := s.client.Get(ctx, "/api/v1/user/settings", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateSettings updates the authenticated user's settings.
func (s *UserService) UpdateSettings(ctx context.Context, opts *types.UserSettingsOptions) (*types.UserSettings, error) {
	var out types.UserSettings
	if err := s.client.Patch(ctx, "/api/v1/user/settings", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetQuota returns the authenticated user's quota information.
func (s *UserService) GetQuota(ctx context.Context) (*types.QuotaInfo, error) {
	var out types.QuotaInfo
	if err := s.client.Get(ctx, "/api/v1/user/quota", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListQuotaArtifacts returns all artifacts affecting the authenticated user's quota.
func (s *UserService) ListQuotaArtifacts(ctx context.Context) ([]types.QuotaUsedArtifact, error) {
	return ListAll[types.QuotaUsedArtifact](ctx, s.client, "/api/v1/user/quota/artifacts", nil)
}

// IterQuotaArtifacts returns an iterator over all artifacts affecting the authenticated user's quota.
func (s *UserService) IterQuotaArtifacts(ctx context.Context) iter.Seq2[types.QuotaUsedArtifact, error] {
	return ListIter[types.QuotaUsedArtifact](ctx, s.client, "/api/v1/user/quota/artifacts", nil)
}

// ListQuotaAttachments returns all attachments affecting the authenticated user's quota.
func (s *UserService) ListQuotaAttachments(ctx context.Context) ([]types.QuotaUsedAttachment, error) {
	return ListAll[types.QuotaUsedAttachment](ctx, s.client, "/api/v1/user/quota/attachments", nil)
}

// IterQuotaAttachments returns an iterator over all attachments affecting the authenticated user's quota.
func (s *UserService) IterQuotaAttachments(ctx context.Context) iter.Seq2[types.QuotaUsedAttachment, error] {
	return ListIter[types.QuotaUsedAttachment](ctx, s.client, "/api/v1/user/quota/attachments", nil)
}

// ListQuotaPackages returns all packages affecting the authenticated user's quota.
func (s *UserService) ListQuotaPackages(ctx context.Context) ([]types.QuotaUsedPackage, error) {
	return ListAll[types.QuotaUsedPackage](ctx, s.client, "/api/v1/user/quota/packages", nil)
}

// IterQuotaPackages returns an iterator over all packages affecting the authenticated user's quota.
func (s *UserService) IterQuotaPackages(ctx context.Context) iter.Seq2[types.QuotaUsedPackage, error] {
	return ListIter[types.QuotaUsedPackage](ctx, s.client, "/api/v1/user/quota/packages", nil)
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

// UpdateAvatar updates the authenticated user's avatar.
func (s *UserService) UpdateAvatar(ctx context.Context, opts *types.UpdateUserAvatarOption) error {
	return s.client.Post(ctx, "/api/v1/user/avatar", opts, nil)
}

// DeleteAvatar deletes the authenticated user's avatar.
func (s *UserService) DeleteAvatar(ctx context.Context) error {
	return s.client.Delete(ctx, "/api/v1/user/avatar")
}

// ListOAuth2Applications returns all OAuth2 applications owned by the authenticated user.
func (s *UserService) ListOAuth2Applications(ctx context.Context) ([]types.OAuth2Application, error) {
	return ListAll[types.OAuth2Application](ctx, s.client, "/api/v1/user/applications/oauth2", nil)
}

// IterOAuth2Applications returns an iterator over all OAuth2 applications owned by the authenticated user.
func (s *UserService) IterOAuth2Applications(ctx context.Context) iter.Seq2[types.OAuth2Application, error] {
	return ListIter[types.OAuth2Application](ctx, s.client, "/api/v1/user/applications/oauth2", nil)
}

// CreateOAuth2Application creates a new OAuth2 application for the authenticated user.
func (s *UserService) CreateOAuth2Application(ctx context.Context, opts *types.CreateOAuth2ApplicationOptions) (*types.OAuth2Application, error) {
	var out types.OAuth2Application
	if err := s.client.Post(ctx, "/api/v1/user/applications/oauth2", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOAuth2Application returns a single OAuth2 application owned by the authenticated user.
func (s *UserService) GetOAuth2Application(ctx context.Context, id int64) (*types.OAuth2Application, error) {
	path := ResolvePath("/api/v1/user/applications/oauth2/{id}", pathParams("id", int64String(id)))
	var out types.OAuth2Application
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateOAuth2Application updates an OAuth2 application owned by the authenticated user.
func (s *UserService) UpdateOAuth2Application(ctx context.Context, id int64, opts *types.CreateOAuth2ApplicationOptions) (*types.OAuth2Application, error) {
	path := ResolvePath("/api/v1/user/applications/oauth2/{id}", pathParams("id", int64String(id)))
	var out types.OAuth2Application
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteOAuth2Application deletes an OAuth2 application owned by the authenticated user.
func (s *UserService) DeleteOAuth2Application(ctx context.Context, id int64) error {
	path := ResolvePath("/api/v1/user/applications/oauth2/{id}", pathParams("id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListStopwatches returns all existing stopwatches for the authenticated user.
func (s *UserService) ListStopwatches(ctx context.Context) ([]types.StopWatch, error) {
	return ListAll[types.StopWatch](ctx, s.client, "/api/v1/user/stopwatches", nil)
}

// IterStopwatches returns an iterator over all existing stopwatches for the authenticated user.
func (s *UserService) IterStopwatches(ctx context.Context) iter.Seq2[types.StopWatch, error] {
	return ListIter[types.StopWatch](ctx, s.client, "/api/v1/user/stopwatches", nil)
}

// ListBlockedUsers returns all users blocked by the authenticated user.
func (s *UserService) ListBlockedUsers(ctx context.Context) ([]types.BlockedUser, error) {
	return ListAll[types.BlockedUser](ctx, s.client, "/api/v1/user/list_blocked", nil)
}

// IterBlockedUsers returns an iterator over all users blocked by the authenticated user.
func (s *UserService) IterBlockedUsers(ctx context.Context) iter.Seq2[types.BlockedUser, error] {
	return ListIter[types.BlockedUser](ctx, s.client, "/api/v1/user/list_blocked", nil)
}

// Block blocks a user as the authenticated user.
func (s *UserService) Block(ctx context.Context, username string) error {
	path := ResolvePath("/api/v1/user/block/{username}", pathParams("username", username))
	return s.client.Put(ctx, path, nil, nil)
}

// Unblock unblocks a user as the authenticated user.
func (s *UserService) Unblock(ctx context.Context, username string) error {
	path := ResolvePath("/api/v1/user/unblock/{username}", pathParams("username", username))
	return s.client.Put(ctx, path, nil, nil)
}

// ListMySubscriptions returns all repositories watched by the authenticated user.
func (s *UserService) ListMySubscriptions(ctx context.Context) ([]types.Repository, error) {
	return ListAll[types.Repository](ctx, s.client, "/api/v1/user/subscriptions", nil)
}

// IterMySubscriptions returns an iterator over all repositories watched by the authenticated user.
func (s *UserService) IterMySubscriptions(ctx context.Context) iter.Seq2[types.Repository, error] {
	return ListIter[types.Repository](ctx, s.client, "/api/v1/user/subscriptions", nil)
}

// ListMyStarred returns all repositories starred by the authenticated user.
func (s *UserService) ListMyStarred(ctx context.Context) ([]types.Repository, error) {
	return ListAll[types.Repository](ctx, s.client, "/api/v1/user/starred", nil)
}

// IterMyStarred returns an iterator over all repositories starred by the authenticated user.
func (s *UserService) IterMyStarred(ctx context.Context) iter.Seq2[types.Repository, error] {
	return ListIter[types.Repository](ctx, s.client, "/api/v1/user/starred", nil)
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

// CheckFollowing reports whether one user is following another user.
func (s *UserService) CheckFollowing(ctx context.Context, username, target string) (bool, error) {
	path := ResolvePath("/api/v1/users/{username}/following/{target}", pathParams("username", username, "target", target))
	resp, err := s.client.doJSON(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == http.StatusNoContent, nil
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

// GetHeatmap returns a user's contribution heatmap data.
func (s *UserService) GetHeatmap(ctx context.Context, username string) ([]types.UserHeatmapData, error) {
	path := ResolvePath("/api/v1/users/{username}/heatmap", pathParams("username", username))
	var out []types.UserHeatmapData
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
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

// CheckStarring reports whether the authenticated user is starring a repository.
func (s *UserService) CheckStarring(ctx context.Context, owner, repo string) (bool, error) {
	path := ResolvePath("/api/v1/user/starred/{owner}/{repo}", pathParams("owner", owner, "repo", repo))
	resp, err := s.client.doJSON(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == http.StatusNoContent, nil
}
