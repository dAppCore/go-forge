package forge

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"

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

// UserSearchOptions controls filtering for user searches.
//
// Usage:
//
//	opts := forge.UserSearchOptions{UID: 1001}
type UserSearchOptions struct {
	UID int64
}

// String returns a safe summary of the user search filters.
func (o UserSearchOptions) String() string {
	return optionString("forge.UserSearchOptions", "uid", o.UID)
}

// GoString returns a safe Go-syntax summary of the user search filters.
func (o UserSearchOptions) GoString() string { return o.String() }

func (o UserSearchOptions) queryParams() map[string]string {
	if o.UID == 0 {
		return nil
	}
	return map[string]string{
		"uid": strconv.FormatInt(o.UID, 10),
	}
}

// UserKeyListOptions controls filtering for authenticated user public key listings.
//
// Usage:
//
//	opts := forge.UserKeyListOptions{Fingerprint: "AB:CD"}
type UserKeyListOptions struct {
	Fingerprint string
}

// String returns a safe summary of the user key filters.
func (o UserKeyListOptions) String() string {
	return optionString("forge.UserKeyListOptions", "fingerprint", o.Fingerprint)
}

// GoString returns a safe Go-syntax summary of the user key filters.
func (o UserKeyListOptions) GoString() string { return o.String() }

func (o UserKeyListOptions) queryParams() map[string]string {
	if o.Fingerprint == "" {
		return nil
	}
	return map[string]string{
		"fingerprint": o.Fingerprint,
	}
}

type userSearchResults struct {
	Data []*types.User `json:"data,omitempty"`
	OK   bool          `json:"ok,omitempty"`
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

// SearchUsersPage returns a single page of users matching the search filters.
func (s *UserService) SearchUsersPage(ctx context.Context, query string, pageOpts ListOptions, filters ...UserSearchOptions) (*PagedResult[types.User], error) {
	if pageOpts.Page < 1 {
		pageOpts.Page = 1
	}
	if pageOpts.Limit < 1 {
		pageOpts.Limit = 50
	}

	u, err := url.Parse("/api/v1/users/search")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", query)
	for _, filter := range filters {
		for key, value := range filter.queryParams() {
			q.Set(key, value)
		}
	}
	q.Set("page", strconv.Itoa(pageOpts.Page))
	q.Set("limit", strconv.Itoa(pageOpts.Limit))
	u.RawQuery = q.Encode()

	var out userSearchResults
	resp, err := s.client.doJSON(ctx, http.MethodGet, u.String(), nil, &out)
	if err != nil {
		return nil, err
	}

	totalCount, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))
	items := make([]types.User, 0, len(out.Data))
	for _, user := range out.Data {
		if user != nil {
			items = append(items, *user)
		}
	}

	return &PagedResult[types.User]{
		Items:      items,
		TotalCount: totalCount,
		Page:       pageOpts.Page,
		HasMore: (totalCount > 0 && (pageOpts.Page-1)*pageOpts.Limit+len(items) < totalCount) ||
			(totalCount == 0 && len(items) >= pageOpts.Limit),
	}, nil
}

// SearchUsers returns all users matching the search filters.
func (s *UserService) SearchUsers(ctx context.Context, query string, filters ...UserSearchOptions) ([]types.User, error) {
	var all []types.User
	page := 1

	for {
		result, err := s.SearchUsersPage(ctx, query, ListOptions{Page: page, Limit: 50}, filters...)
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

// IterSearchUsers returns an iterator over users matching the search filters.
func (s *UserService) IterSearchUsers(ctx context.Context, query string, filters ...UserSearchOptions) iter.Seq2[types.User, error] {
	return func(yield func(types.User, error) bool) {
		page := 1
		for {
			result, err := s.SearchUsersPage(ctx, query, ListOptions{Page: page, Limit: 50}, filters...)
			if err != nil {
				yield(*new(types.User), err)
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

// ListKeys returns all public keys owned by the authenticated user.
func (s *UserService) ListKeys(ctx context.Context, filters ...UserKeyListOptions) ([]types.PublicKey, error) {
	query := make(map[string]string, len(filters))
	for _, filter := range filters {
		for key, value := range filter.queryParams() {
			query[key] = value
		}
	}
	if len(query) == 0 {
		query = nil
	}
	return ListAll[types.PublicKey](ctx, s.client, "/api/v1/user/keys", query)
}

// IterKeys returns an iterator over all public keys owned by the authenticated user.
func (s *UserService) IterKeys(ctx context.Context, filters ...UserKeyListOptions) iter.Seq2[types.PublicKey, error] {
	query := make(map[string]string, len(filters))
	for _, filter := range filters {
		for key, value := range filter.queryParams() {
			query[key] = value
		}
	}
	if len(query) == 0 {
		query = nil
	}
	return ListIter[types.PublicKey](ctx, s.client, "/api/v1/user/keys", query)
}

// CreateKey creates a public key for the authenticated user.
func (s *UserService) CreateKey(ctx context.Context, opts *types.CreateKeyOption) (*types.PublicKey, error) {
	var out types.PublicKey
	if err := s.client.Post(ctx, "/api/v1/user/keys", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetKey returns a single public key owned by the authenticated user.
func (s *UserService) GetKey(ctx context.Context, id int64) (*types.PublicKey, error) {
	path := ResolvePath("/api/v1/user/keys/{id}", pathParams("id", int64String(id)))
	var out types.PublicKey
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteKey removes a public key owned by the authenticated user.
func (s *UserService) DeleteKey(ctx context.Context, id int64) error {
	path := ResolvePath("/api/v1/user/keys/{id}", pathParams("id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListUserKeys returns all public keys for a user.
func (s *UserService) ListUserKeys(ctx context.Context, username string, filters ...UserKeyListOptions) ([]types.PublicKey, error) {
	path := ResolvePath("/api/v1/users/{username}/keys", pathParams("username", username))
	query := make(map[string]string, len(filters))
	for _, filter := range filters {
		for key, value := range filter.queryParams() {
			query[key] = value
		}
	}
	if len(query) == 0 {
		query = nil
	}
	return ListAll[types.PublicKey](ctx, s.client, path, query)
}

// IterUserKeys returns an iterator over all public keys for a user.
func (s *UserService) IterUserKeys(ctx context.Context, username string, filters ...UserKeyListOptions) iter.Seq2[types.PublicKey, error] {
	path := ResolvePath("/api/v1/users/{username}/keys", pathParams("username", username))
	query := make(map[string]string, len(filters))
	for _, filter := range filters {
		for key, value := range filter.queryParams() {
			query[key] = value
		}
	}
	if len(query) == 0 {
		query = nil
	}
	return ListIter[types.PublicKey](ctx, s.client, path, query)
}

// ListGPGKeys returns all GPG keys owned by the authenticated user.
func (s *UserService) ListGPGKeys(ctx context.Context) ([]types.GPGKey, error) {
	return ListAll[types.GPGKey](ctx, s.client, "/api/v1/user/gpg_keys", nil)
}

// IterGPGKeys returns an iterator over all GPG keys owned by the authenticated user.
func (s *UserService) IterGPGKeys(ctx context.Context) iter.Seq2[types.GPGKey, error] {
	return ListIter[types.GPGKey](ctx, s.client, "/api/v1/user/gpg_keys", nil)
}

// CreateGPGKey adds a GPG key for the authenticated user.
func (s *UserService) CreateGPGKey(ctx context.Context, opts *types.CreateGPGKeyOption) (*types.GPGKey, error) {
	var out types.GPGKey
	if err := s.client.Post(ctx, "/api/v1/user/gpg_keys", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetGPGKey returns a single GPG key owned by the authenticated user.
func (s *UserService) GetGPGKey(ctx context.Context, id int64) (*types.GPGKey, error) {
	path := ResolvePath("/api/v1/user/gpg_keys/{id}", pathParams("id", int64String(id)))
	var out types.GPGKey
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteGPGKey removes a GPG key owned by the authenticated user.
func (s *UserService) DeleteGPGKey(ctx context.Context, id int64) error {
	path := ResolvePath("/api/v1/user/gpg_keys/{id}", pathParams("id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListUserGPGKeys returns all GPG keys for a user.
func (s *UserService) ListUserGPGKeys(ctx context.Context, username string) ([]types.GPGKey, error) {
	path := ResolvePath("/api/v1/users/{username}/gpg_keys", pathParams("username", username))
	return ListAll[types.GPGKey](ctx, s.client, path, nil)
}

// IterUserGPGKeys returns an iterator over all GPG keys for a user.
func (s *UserService) IterUserGPGKeys(ctx context.Context, username string) iter.Seq2[types.GPGKey, error] {
	path := ResolvePath("/api/v1/users/{username}/gpg_keys", pathParams("username", username))
	return ListIter[types.GPGKey](ctx, s.client, path, nil)
}

// GetGPGKeyVerificationToken returns the token used to verify a GPG key.
func (s *UserService) GetGPGKeyVerificationToken(ctx context.Context) (string, error) {
	data, err := s.client.GetRaw(ctx, "/api/v1/user/gpg_key_token")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// VerifyGPGKey verifies a GPG key for the authenticated user.
func (s *UserService) VerifyGPGKey(ctx context.Context) (*types.GPGKey, error) {
	var out types.GPGKey
	if err := s.client.Post(ctx, "/api/v1/user/gpg_key_verify", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTokens returns all access tokens for a user.
func (s *UserService) ListTokens(ctx context.Context, username string) ([]types.AccessToken, error) {
	path := ResolvePath("/api/v1/users/{username}/tokens", pathParams("username", username))
	return ListAll[types.AccessToken](ctx, s.client, path, nil)
}

// IterTokens returns an iterator over all access tokens for a user.
func (s *UserService) IterTokens(ctx context.Context, username string) iter.Seq2[types.AccessToken, error] {
	path := ResolvePath("/api/v1/users/{username}/tokens", pathParams("username", username))
	return ListIter[types.AccessToken](ctx, s.client, path, nil)
}

// CreateToken creates an access token for a user.
func (s *UserService) CreateToken(ctx context.Context, username string, opts *types.CreateAccessTokenOption) (*types.AccessToken, error) {
	path := ResolvePath("/api/v1/users/{username}/tokens", pathParams("username", username))
	var out types.AccessToken
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteToken deletes an access token for a user.
func (s *UserService) DeleteToken(ctx context.Context, username, token string) error {
	path := ResolvePath("/api/v1/users/{username}/tokens/{token}", pathParams("username", username, "token", token))
	return s.client.Delete(ctx, path)
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

// ListMyFollowers returns all followers of the authenticated user.
func (s *UserService) ListMyFollowers(ctx context.Context) ([]types.User, error) {
	return ListAll[types.User](ctx, s.client, "/api/v1/user/followers", nil)
}

// IterMyFollowers returns an iterator over all followers of the authenticated user.
func (s *UserService) IterMyFollowers(ctx context.Context) iter.Seq2[types.User, error] {
	return ListIter[types.User](ctx, s.client, "/api/v1/user/followers", nil)
}

// ListMyFollowing returns all users followed by the authenticated user.
func (s *UserService) ListMyFollowing(ctx context.Context) ([]types.User, error) {
	return ListAll[types.User](ctx, s.client, "/api/v1/user/following", nil)
}

// IterMyFollowing returns an iterator over all users followed by the authenticated user.
func (s *UserService) IterMyFollowing(ctx context.Context) iter.Seq2[types.User, error] {
	return ListIter[types.User](ctx, s.client, "/api/v1/user/following", nil)
}

// ListMyTeams returns all teams the authenticated user belongs to.
func (s *UserService) ListMyTeams(ctx context.Context) ([]types.Team, error) {
	return ListAll[types.Team](ctx, s.client, "/api/v1/user/teams", nil)
}

// IterMyTeams returns an iterator over all teams the authenticated user belongs to.
func (s *UserService) IterMyTeams(ctx context.Context) iter.Seq2[types.Team, error] {
	return ListIter[types.Team](ctx, s.client, "/api/v1/user/teams", nil)
}

// ListMyTrackedTimes returns all tracked times logged by the authenticated user.
func (s *UserService) ListMyTrackedTimes(ctx context.Context) ([]types.TrackedTime, error) {
	return ListAll[types.TrackedTime](ctx, s.client, "/api/v1/user/times", nil)
}

// IterMyTrackedTimes returns an iterator over all tracked times logged by the authenticated user.
func (s *UserService) IterMyTrackedTimes(ctx context.Context) iter.Seq2[types.TrackedTime, error] {
	return ListIter[types.TrackedTime](ctx, s.client, "/api/v1/user/times", nil)
}

// CheckQuota reports whether the authenticated user is over quota.
func (s *UserService) CheckQuota(ctx context.Context) (bool, error) {
	var out bool
	if err := s.client.Get(ctx, "/api/v1/user/quota/check", &out); err != nil {
		return false, err
	}
	return out, nil
}

// GetRunnerRegistrationToken returns the authenticated user's actions runner registration token.
func (s *UserService) GetRunnerRegistrationToken(ctx context.Context) (string, error) {
	path := "/api/v1/user/actions/runners/registration-token"
	resp, err := s.client.doJSON(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return "", err
	}
	return resp.Header.Get("token"), nil
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

// ListActivityFeeds returns a user's activity feed entries.
func (s *UserService) ListActivityFeeds(ctx context.Context, username string) ([]types.Activity, error) {
	path := ResolvePath("/api/v1/users/{username}/activities/feeds", pathParams("username", username))
	return ListAll[types.Activity](ctx, s.client, path, nil)
}

// IterActivityFeeds returns an iterator over a user's activity feed entries.
func (s *UserService) IterActivityFeeds(ctx context.Context, username string) iter.Seq2[types.Activity, error] {
	path := ResolvePath("/api/v1/users/{username}/activities/feeds", pathParams("username", username))
	return ListIter[types.Activity](ctx, s.client, path, nil)
}

// ListRepos returns all repositories owned by a user.
func (s *UserService) ListRepos(ctx context.Context, username string) ([]types.Repository, error) {
	path := ResolvePath("/api/v1/users/{username}/repos", pathParams("username", username))
	return ListAll[types.Repository](ctx, s.client, path, nil)
}

// IterRepos returns an iterator over all repositories owned by a user.
func (s *UserService) IterRepos(ctx context.Context, username string) iter.Seq2[types.Repository, error] {
	path := ResolvePath("/api/v1/users/{username}/repos", pathParams("username", username))
	return ListIter[types.Repository](ctx, s.client, path, nil)
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
