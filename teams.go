package forge

import (
	"context"
	"iter"

	"dappco.re/go/core/forge/types"
)

// TeamService handles team operations.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Teams.ListMembers(ctx, 42)
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
	path := ResolvePath("/api/v1/teams/{id}/members", pathParams("id", int64String(teamID)))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterMembers returns an iterator over all members of a team.
func (s *TeamService) IterMembers(ctx context.Context, teamID int64) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/teams/{id}/members", pathParams("id", int64String(teamID)))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// AddMember adds a user to a team.
func (s *TeamService) AddMember(ctx context.Context, teamID int64, username string) error {
	path := ResolvePath("/api/v1/teams/{id}/members/{username}", pathParams("id", int64String(teamID), "username", username))
	return s.client.Put(ctx, path, nil, nil)
}

// GetMember returns a particular member of a team.
func (s *TeamService) GetMember(ctx context.Context, teamID int64, username string) (*types.User, error) {
	path := ResolvePath("/api/v1/teams/{id}/members/{username}", pathParams("id", int64String(teamID), "username", username))
	var out types.User
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RemoveMember removes a user from a team.
func (s *TeamService) RemoveMember(ctx context.Context, teamID int64, username string) error {
	path := ResolvePath("/api/v1/teams/{id}/members/{username}", pathParams("id", int64String(teamID), "username", username))
	return s.client.Delete(ctx, path)
}

// ListRepos returns all repositories managed by a team.
func (s *TeamService) ListRepos(ctx context.Context, teamID int64) ([]types.Repository, error) {
	path := ResolvePath("/api/v1/teams/{id}/repos", pathParams("id", int64String(teamID)))
	return ListAll[types.Repository](ctx, s.client, path, nil)
}

// IterRepos returns an iterator over all repositories managed by a team.
func (s *TeamService) IterRepos(ctx context.Context, teamID int64) iter.Seq2[types.Repository, error] {
	path := ResolvePath("/api/v1/teams/{id}/repos", pathParams("id", int64String(teamID)))
	return ListIter[types.Repository](ctx, s.client, path, nil)
}

// AddRepo adds a repository to a team.
func (s *TeamService) AddRepo(ctx context.Context, teamID int64, org, repo string) error {
	path := ResolvePath("/api/v1/teams/{id}/repos/{org}/{repo}", pathParams("id", int64String(teamID), "org", org, "repo", repo))
	return s.client.Put(ctx, path, nil, nil)
}

// RemoveRepo removes a repository from a team.
func (s *TeamService) RemoveRepo(ctx context.Context, teamID int64, org, repo string) error {
	path := ResolvePath("/api/v1/teams/{id}/repos/{org}/{repo}", pathParams("id", int64String(teamID), "org", org, "repo", repo))
	return s.client.Delete(ctx, path)
}

// GetRepo returns a particular repository managed by a team.
func (s *TeamService) GetRepo(ctx context.Context, teamID int64, org, repo string) (*types.Repository, error) {
	path := ResolvePath("/api/v1/teams/{id}/repos/{org}/{repo}", pathParams("id", int64String(teamID), "org", org, "repo", repo))
	var out types.Repository
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListOrgTeams returns all teams in an organisation.
func (s *TeamService) ListOrgTeams(ctx context.Context, org string) ([]types.Team, error) {
	path := ResolvePath("/api/v1/orgs/{org}/teams", pathParams("org", org))
	return ListAll[types.Team](ctx, s.client, path, nil)
}

// IterOrgTeams returns an iterator over all teams in an organisation.
func (s *TeamService) IterOrgTeams(ctx context.Context, org string) iter.Seq2[types.Team, error] {
	path := ResolvePath("/api/v1/orgs/{org}/teams", pathParams("org", org))
	return ListIter[types.Team](ctx, s.client, path, nil)
}

// ListActivityFeeds returns a team's activity feed entries.
func (s *TeamService) ListActivityFeeds(ctx context.Context, teamID int64) ([]types.Activity, error) {
	path := ResolvePath("/api/v1/teams/{id}/activities/feeds", pathParams("id", int64String(teamID)))
	return ListAll[types.Activity](ctx, s.client, path, nil)
}

// IterActivityFeeds returns an iterator over a team's activity feed entries.
func (s *TeamService) IterActivityFeeds(ctx context.Context, teamID int64) iter.Seq2[types.Activity, error] {
	path := ResolvePath("/api/v1/teams/{id}/activities/feeds", pathParams("id", int64String(teamID)))
	return ListIter[types.Activity](ctx, s.client, path, nil)
}
