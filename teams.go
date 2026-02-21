package forge

import (
	"context"
	"fmt"

	"forge.lthn.ai/core/go-forge/types"
)

// TeamService handles team operations.
type TeamService struct {
	Resource[types.Team, types.CreateTeamOption, types.EditTeamOption]
}

func newTeamService(c *Client) *TeamService {
	return &TeamService{
		Resource: *NewResource[types.Team, types.CreateTeamOption, types.EditTeamOption](
			c, "/api/v1/teams/{id}",
		),
	}
}

// ListMembers returns all members of a team.
func (s *TeamService) ListMembers(ctx context.Context, teamID int64) ([]types.User, error) {
	path := fmt.Sprintf("/api/v1/teams/%d/members", teamID)
	return ListAll[types.User](ctx, s.client, path, nil)
}

// AddMember adds a user to a team.
func (s *TeamService) AddMember(ctx context.Context, teamID int64, username string) error {
	path := fmt.Sprintf("/api/v1/teams/%d/members/%s", teamID, username)
	return s.client.Put(ctx, path, nil, nil)
}

// RemoveMember removes a user from a team.
func (s *TeamService) RemoveMember(ctx context.Context, teamID int64, username string) error {
	path := fmt.Sprintf("/api/v1/teams/%d/members/%s", teamID, username)
	return s.client.Delete(ctx, path)
}

// ListRepos returns all repositories managed by a team.
func (s *TeamService) ListRepos(ctx context.Context, teamID int64) ([]types.Repository, error) {
	path := fmt.Sprintf("/api/v1/teams/%d/repos", teamID)
	return ListAll[types.Repository](ctx, s.client, path, nil)
}

// AddRepo adds a repository to a team.
func (s *TeamService) AddRepo(ctx context.Context, teamID int64, org, repo string) error {
	path := fmt.Sprintf("/api/v1/teams/%d/repos/%s/%s", teamID, org, repo)
	return s.client.Put(ctx, path, nil, nil)
}

// RemoveRepo removes a repository from a team.
func (s *TeamService) RemoveRepo(ctx context.Context, teamID int64, org, repo string) error {
	path := fmt.Sprintf("/api/v1/teams/%d/repos/%s/%s", teamID, org, repo)
	return s.client.Delete(ctx, path)
}

// ListOrgTeams returns all teams in an organisation.
func (s *TeamService) ListOrgTeams(ctx context.Context, org string) ([]types.Team, error) {
	path := fmt.Sprintf("/api/v1/orgs/%s/teams", org)
	return ListAll[types.Team](ctx, s.client, path, nil)
}
