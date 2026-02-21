package forge

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go-forge/types"
)

// OrgService handles organisation operations.
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
	path := fmt.Sprintf("/api/v1/orgs/%s/members", org)
	return ListAll[types.User](ctx, s.client, path, nil)
}

// AddMember adds a user to an organisation.
func (s *OrgService) AddMember(ctx context.Context, org, username string) error {
	path := fmt.Sprintf("/api/v1/orgs/%s/members/%s", org, username)
	return s.client.Put(ctx, path, nil, nil)
}

// RemoveMember removes a user from an organisation.
func (s *OrgService) RemoveMember(ctx context.Context, org, username string) error {
	path := fmt.Sprintf("/api/v1/orgs/%s/members/%s", org, username)
	return s.client.Delete(ctx, path)
}

// ListUserOrgs returns all organisations for a user.
func (s *OrgService) ListUserOrgs(ctx context.Context, username string) ([]types.Organization, error) {
	path := fmt.Sprintf("/api/v1/users/%s/orgs", username)
	return ListAll[types.Organization](ctx, s.client, path, nil)
}

// ListMyOrgs returns all organisations for the authenticated user.
func (s *OrgService) ListMyOrgs(ctx context.Context) ([]types.Organization, error) {
	return ListAll[types.Organization](ctx, s.client, "/api/v1/user/orgs", nil)
}
