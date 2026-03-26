package forge

import (
	"context"
	"iter"

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
