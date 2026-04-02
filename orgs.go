package forge

import (
	"context"
	"iter"
	"net/http"
	"time"

	"dappco.re/go/core/forge/types"
)

// OrgService handles organisation operations.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Orgs.ListMembers(ctx, "core")
type OrgService struct {
	Resource[types.Organization, types.CreateOrgOption, types.EditOrgOption]
}

// OrgActivityFeedListOptions controls filtering for organisation activity feeds.
type OrgActivityFeedListOptions struct {
	Date *time.Time
}

func (o OrgActivityFeedListOptions) queryParams() map[string]string {
	if o.Date == nil {
		return nil
	}
	return map[string]string{
		"date": o.Date.Format("2006-01-02"),
	}
}

func newOrgService(c *Client) *OrgService {
	return &OrgService{
		Resource: *NewResource[types.Organization, types.CreateOrgOption, types.EditOrgOption](
			c, "/api/v1/orgs/{org}",
		),
	}
}

// ListMembers returns all members of an organisation.
func (s *OrgService) ListMembers(ctx context.Context, org string) ([]types.User, error) {
	path := ResolvePath("/api/v1/orgs/{org}/members", pathParams("org", org))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterMembers returns an iterator over all members of an organisation.
func (s *OrgService) IterMembers(ctx context.Context, org string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/orgs/{org}/members", pathParams("org", org))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// AddMember adds a user to an organisation.
func (s *OrgService) AddMember(ctx context.Context, org, username string) error {
	path := ResolvePath("/api/v1/orgs/{org}/members/{username}", pathParams("org", org, "username", username))
	return s.client.Put(ctx, path, nil, nil)
}

// RemoveMember removes a user from an organisation.
func (s *OrgService) RemoveMember(ctx context.Context, org, username string) error {
	path := ResolvePath("/api/v1/orgs/{org}/members/{username}", pathParams("org", org, "username", username))
	return s.client.Delete(ctx, path)
}

// IsMember reports whether a user is a member of an organisation.
func (s *OrgService) IsMember(ctx context.Context, org, username string) (bool, error) {
	path := ResolvePath("/api/v1/orgs/{org}/members/{username}", pathParams("org", org, "username", username))
	resp, err := s.client.doJSON(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == http.StatusNoContent, nil
}

// ListBlockedUsers returns all users blocked by an organisation.
func (s *OrgService) ListBlockedUsers(ctx context.Context, org string) ([]types.BlockedUser, error) {
	path := ResolvePath("/api/v1/orgs/{org}/list_blocked", pathParams("org", org))
	return ListAll[types.BlockedUser](ctx, s.client, path, nil)
}

// IterBlockedUsers returns an iterator over all users blocked by an organisation.
func (s *OrgService) IterBlockedUsers(ctx context.Context, org string) iter.Seq2[types.BlockedUser, error] {
	path := ResolvePath("/api/v1/orgs/{org}/list_blocked", pathParams("org", org))
	return ListIter[types.BlockedUser](ctx, s.client, path, nil)
}

// IsBlocked reports whether a user is blocked by an organisation.
func (s *OrgService) IsBlocked(ctx context.Context, org, username string) (bool, error) {
	path := ResolvePath("/api/v1/orgs/{org}/block/{username}", pathParams("org", org, "username", username))
	resp, err := s.client.doJSON(ctx, "GET", path, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == http.StatusNoContent, nil
}

// ListPublicMembers returns all public members of an organisation.
func (s *OrgService) ListPublicMembers(ctx context.Context, org string) ([]types.User, error) {
	path := ResolvePath("/api/v1/orgs/{org}/public_members", pathParams("org", org))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterPublicMembers returns an iterator over all public members of an organisation.
func (s *OrgService) IterPublicMembers(ctx context.Context, org string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/orgs/{org}/public_members", pathParams("org", org))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// IsPublicMember reports whether a user is a public member of an organisation.
func (s *OrgService) IsPublicMember(ctx context.Context, org, username string) (bool, error) {
	path := ResolvePath("/api/v1/orgs/{org}/public_members/{username}", pathParams("org", org, "username", username))
	resp, err := s.client.doJSON(ctx, "GET", path, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == http.StatusNoContent, nil
}

// PublicizeMember makes a user's membership public within an organisation.
func (s *OrgService) PublicizeMember(ctx context.Context, org, username string) error {
	path := ResolvePath("/api/v1/orgs/{org}/public_members/{username}", pathParams("org", org, "username", username))
	return s.client.Put(ctx, path, nil, nil)
}

// ConcealMember hides a user's public membership within an organisation.
func (s *OrgService) ConcealMember(ctx context.Context, org, username string) error {
	path := ResolvePath("/api/v1/orgs/{org}/public_members/{username}", pathParams("org", org, "username", username))
	return s.client.Delete(ctx, path)
}

// Block blocks a user within an organisation.
func (s *OrgService) Block(ctx context.Context, org, username string) error {
	path := ResolvePath("/api/v1/orgs/{org}/block/{username}", pathParams("org", org, "username", username))
	return s.client.Put(ctx, path, nil, nil)
}

// Unblock unblocks a user within an organisation.
func (s *OrgService) Unblock(ctx context.Context, org, username string) error {
	path := ResolvePath("/api/v1/orgs/{org}/unblock/{username}", pathParams("org", org, "username", username))
	return s.client.Delete(ctx, path)
}

// GetQuota returns the quota information for an organisation.
func (s *OrgService) GetQuota(ctx context.Context, org string) (*types.QuotaInfo, error) {
	path := ResolvePath("/api/v1/orgs/{org}/quota", pathParams("org", org))
	var out types.QuotaInfo
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CheckQuota reports whether an organisation is over quota for the current subject.
func (s *OrgService) CheckQuota(ctx context.Context, org string) (bool, error) {
	path := ResolvePath("/api/v1/orgs/{org}/quota/check", pathParams("org", org))
	var out bool
	if err := s.client.Get(ctx, path, &out); err != nil {
		return false, err
	}
	return out, nil
}

// ListQuotaArtifacts returns all artefacts counting towards an organisation's quota.
func (s *OrgService) ListQuotaArtifacts(ctx context.Context, org string) ([]types.QuotaUsedArtifact, error) {
	path := ResolvePath("/api/v1/orgs/{org}/quota/artifacts", pathParams("org", org))
	return ListAll[types.QuotaUsedArtifact](ctx, s.client, path, nil)
}

// IterQuotaArtifacts returns an iterator over all artefacts counting towards an organisation's quota.
func (s *OrgService) IterQuotaArtifacts(ctx context.Context, org string) iter.Seq2[types.QuotaUsedArtifact, error] {
	path := ResolvePath("/api/v1/orgs/{org}/quota/artifacts", pathParams("org", org))
	return ListIter[types.QuotaUsedArtifact](ctx, s.client, path, nil)
}

// ListQuotaAttachments returns all attachments counting towards an organisation's quota.
func (s *OrgService) ListQuotaAttachments(ctx context.Context, org string) ([]types.QuotaUsedAttachment, error) {
	path := ResolvePath("/api/v1/orgs/{org}/quota/attachments", pathParams("org", org))
	return ListAll[types.QuotaUsedAttachment](ctx, s.client, path, nil)
}

// IterQuotaAttachments returns an iterator over all attachments counting towards an organisation's quota.
func (s *OrgService) IterQuotaAttachments(ctx context.Context, org string) iter.Seq2[types.QuotaUsedAttachment, error] {
	path := ResolvePath("/api/v1/orgs/{org}/quota/attachments", pathParams("org", org))
	return ListIter[types.QuotaUsedAttachment](ctx, s.client, path, nil)
}

// ListQuotaPackages returns all packages counting towards an organisation's quota.
func (s *OrgService) ListQuotaPackages(ctx context.Context, org string) ([]types.QuotaUsedPackage, error) {
	path := ResolvePath("/api/v1/orgs/{org}/quota/packages", pathParams("org", org))
	return ListAll[types.QuotaUsedPackage](ctx, s.client, path, nil)
}

// IterQuotaPackages returns an iterator over all packages counting towards an organisation's quota.
func (s *OrgService) IterQuotaPackages(ctx context.Context, org string) iter.Seq2[types.QuotaUsedPackage, error] {
	path := ResolvePath("/api/v1/orgs/{org}/quota/packages", pathParams("org", org))
	return ListIter[types.QuotaUsedPackage](ctx, s.client, path, nil)
}

// GetRunnerRegistrationToken returns an organisation actions runner registration token.
func (s *OrgService) GetRunnerRegistrationToken(ctx context.Context, org string) (string, error) {
	path := ResolvePath("/api/v1/orgs/{org}/actions/runners/registration-token", pathParams("org", org))
	resp, err := s.client.doJSON(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return "", err
	}
	return resp.Header.Get("token"), nil
}

// UpdateAvatar updates an organisation avatar.
func (s *OrgService) UpdateAvatar(ctx context.Context, org string, opts *types.UpdateUserAvatarOption) error {
	path := ResolvePath("/api/v1/orgs/{org}/avatar", pathParams("org", org))
	return s.client.Post(ctx, path, opts, nil)
}

// DeleteAvatar deletes an organisation avatar.
func (s *OrgService) DeleteAvatar(ctx context.Context, org string) error {
	path := ResolvePath("/api/v1/orgs/{org}/avatar", pathParams("org", org))
	return s.client.Delete(ctx, path)
}

// SearchTeams searches for teams within an organisation.
func (s *OrgService) SearchTeams(ctx context.Context, org, q string) ([]types.Team, error) {
	path := ResolvePath("/api/v1/orgs/{org}/teams/search", pathParams("org", org))
	return ListAll[types.Team](ctx, s.client, path, map[string]string{"q": q})
}

// IterSearchTeams returns an iterator over teams within an organisation.
func (s *OrgService) IterSearchTeams(ctx context.Context, org, q string) iter.Seq2[types.Team, error] {
	path := ResolvePath("/api/v1/orgs/{org}/teams/search", pathParams("org", org))
	return ListIter[types.Team](ctx, s.client, path, map[string]string{"q": q})
}

// GetUserPermissions returns a user's permissions in an organisation.
func (s *OrgService) GetUserPermissions(ctx context.Context, username, org string) (*types.OrganizationPermissions, error) {
	path := ResolvePath("/api/v1/users/{username}/orgs/{org}/permissions", pathParams("username", username, "org", org))
	var out types.OrganizationPermissions
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListActivityFeeds returns the organisation's activity feed entries.
func (s *OrgService) ListActivityFeeds(ctx context.Context, org string, filters ...OrgActivityFeedListOptions) ([]types.Activity, error) {
	path := ResolvePath("/api/v1/orgs/{org}/activities/feeds", pathParams("org", org))
	return ListAll[types.Activity](ctx, s.client, path, orgActivityFeedQuery(filters...))
}

// IterActivityFeeds returns an iterator over the organisation's activity feed entries.
func (s *OrgService) IterActivityFeeds(ctx context.Context, org string, filters ...OrgActivityFeedListOptions) iter.Seq2[types.Activity, error] {
	path := ResolvePath("/api/v1/orgs/{org}/activities/feeds", pathParams("org", org))
	return ListIter[types.Activity](ctx, s.client, path, orgActivityFeedQuery(filters...))
}

// ListUserOrgs returns all organisations for a user.
func (s *OrgService) ListUserOrgs(ctx context.Context, username string) ([]types.Organization, error) {
	path := ResolvePath("/api/v1/users/{username}/orgs", pathParams("username", username))
	return ListAll[types.Organization](ctx, s.client, path, nil)
}

// IterUserOrgs returns an iterator over all organisations for a user.
func (s *OrgService) IterUserOrgs(ctx context.Context, username string) iter.Seq2[types.Organization, error] {
	path := ResolvePath("/api/v1/users/{username}/orgs", pathParams("username", username))
	return ListIter[types.Organization](ctx, s.client, path, nil)
}

// ListMyOrgs returns all organisations for the authenticated user.
func (s *OrgService) ListMyOrgs(ctx context.Context) ([]types.Organization, error) {
	return ListAll[types.Organization](ctx, s.client, "/api/v1/user/orgs", nil)
}

// IterMyOrgs returns an iterator over all organisations for the authenticated user.
func (s *OrgService) IterMyOrgs(ctx context.Context) iter.Seq2[types.Organization, error] {
	return ListIter[types.Organization](ctx, s.client, "/api/v1/user/orgs", nil)
}

func orgActivityFeedQuery(filters ...OrgActivityFeedListOptions) map[string]string {
	if len(filters) == 0 {
		return nil
	}

	query := make(map[string]string, 1)
	for _, filter := range filters {
		if filter.Date != nil {
			query["date"] = filter.Date.Format("2006-01-02")
		}
	}
	if len(query) == 0 {
		return nil
	}
	return query
}
