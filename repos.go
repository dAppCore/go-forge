package forge

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"dappco.re/go/core/forge/types"
)

// RepoService handles repository operations.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Repos.ListOrgRepos(ctx, "core")
type RepoService struct {
	Resource[types.Repository, types.CreateRepoOption, types.EditRepoOption]
}

// RepoKeyListOptions controls filtering for repository key listings.
type RepoKeyListOptions struct {
	KeyID       int64
	Fingerprint string
}

func (o RepoKeyListOptions) queryParams() map[string]string {
	query := make(map[string]string, 2)
	if o.KeyID != 0 {
		query["key_id"] = strconv.FormatInt(o.KeyID, 10)
	}
	if o.Fingerprint != "" {
		query["fingerprint"] = o.Fingerprint
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

// ActivityFeedListOptions controls filtering for repository activity feeds.
type ActivityFeedListOptions struct {
	Date *time.Time
}

func (o ActivityFeedListOptions) queryParams() map[string]string {
	if o.Date == nil {
		return nil
	}
	return map[string]string{
		"date": o.Date.Format("2006-01-02"),
	}
}

// RepoTimeListOptions controls filtering for repository tracked times.
type RepoTimeListOptions struct {
	User   string
	Since  *time.Time
	Before *time.Time
}

func (o RepoTimeListOptions) queryParams() map[string]string {
	query := make(map[string]string, 3)
	if o.User != "" {
		query["user"] = o.User
	}
	if o.Since != nil {
		query["since"] = o.Since.Format(time.RFC3339)
	}
	if o.Before != nil {
		query["before"] = o.Before.Format(time.RFC3339)
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func newRepoService(c *Client) *RepoService {
	return &RepoService{
		Resource: *NewResource[types.Repository, types.CreateRepoOption, types.EditRepoOption](
			c, "/api/v1/repos/{owner}/{repo}",
		),
	}
}

// Migrate imports a remote git repository into Forgejo.
func (s *RepoService) Migrate(ctx context.Context, opts *types.MigrateRepoOptions) (*types.Repository, error) {
	var out types.Repository
	if err := s.client.Post(ctx, "/api/v1/repos/migrate", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateOrgRepo creates a repository in an organisation.
func (s *RepoService) CreateOrgRepo(ctx context.Context, org string, opts *types.CreateRepoOption) (*types.Repository, error) {
	path := ResolvePath("/api/v1/orgs/{org}/repos", pathParams("org", org))
	var out types.Repository
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListOrgRepos returns all repositories for an organisation.
func (s *RepoService) ListOrgRepos(ctx context.Context, org string) ([]types.Repository, error) {
	path := ResolvePath("/api/v1/orgs/{org}/repos", pathParams("org", org))
	return ListAll[types.Repository](ctx, s.client, path, nil)
}

// IterOrgRepos returns an iterator over all repositories for an organisation.
func (s *RepoService) IterOrgRepos(ctx context.Context, org string) iter.Seq2[types.Repository, error] {
	path := ResolvePath("/api/v1/orgs/{org}/repos", pathParams("org", org))
	return ListIter[types.Repository](ctx, s.client, path, nil)
}

// ListUserRepos returns all repositories for the authenticated user.
func (s *RepoService) ListUserRepos(ctx context.Context) ([]types.Repository, error) {
	return ListAll[types.Repository](ctx, s.client, "/api/v1/user/repos", nil)
}

// IterUserRepos returns an iterator over all repositories for the authenticated user.
func (s *RepoService) IterUserRepos(ctx context.Context) iter.Seq2[types.Repository, error] {
	return ListIter[types.Repository](ctx, s.client, "/api/v1/user/repos", nil)
}

// GetByID returns a repository by its numeric ID.
func (s *RepoService) GetByID(ctx context.Context, id int64) (*types.Repository, error) {
	path := ResolvePath("/api/v1/repositories/{id}", pathParams("id", int64String(id)))
	var out types.Repository
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTags returns all tags for a repository.
func (s *RepoService) ListTags(ctx context.Context, owner, repo string) ([]types.Tag, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tags", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Tag](ctx, s.client, path, nil)
}

// IterTags returns an iterator over all tags for a repository.
func (s *RepoService) IterTags(ctx context.Context, owner, repo string) iter.Seq2[types.Tag, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tags", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Tag](ctx, s.client, path, nil)
}

// GetTag returns a single tag by name.
func (s *RepoService) GetTag(ctx context.Context, owner, repo, tag string) (*types.Tag, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tags/{tag}", pathParams("owner", owner, "repo", repo, "tag", tag))
	var out types.Tag
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTag deletes a repository tag by name.
func (s *RepoService) DeleteTag(ctx context.Context, owner, repo, tag string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tags/{tag}", pathParams("owner", owner, "repo", repo, "tag", tag))
	return s.client.Delete(ctx, path)
}

// ListTagProtections returns all tag protections for a repository.
func (s *RepoService) ListTagProtections(ctx context.Context, owner, repo string) ([]types.TagProtection, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tag_protections", pathParams("owner", owner, "repo", repo))
	return ListAll[types.TagProtection](ctx, s.client, path, nil)
}

// IterTagProtections returns an iterator over all tag protections for a repository.
func (s *RepoService) IterTagProtections(ctx context.Context, owner, repo string) iter.Seq2[types.TagProtection, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tag_protections", pathParams("owner", owner, "repo", repo))
	return ListIter[types.TagProtection](ctx, s.client, path, nil)
}

// GetTagProtection returns a single tag protection by ID.
func (s *RepoService) GetTagProtection(ctx context.Context, owner, repo string, id int64) (*types.TagProtection, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tag_protections/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.TagProtection
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateTagProtection creates a new tag protection for a repository.
func (s *RepoService) CreateTagProtection(ctx context.Context, owner, repo string, opts *types.CreateTagProtectionOption) (*types.TagProtection, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tag_protections", pathParams("owner", owner, "repo", repo))
	var out types.TagProtection
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditTagProtection updates an existing tag protection for a repository.
func (s *RepoService) EditTagProtection(ctx context.Context, owner, repo string, id int64, opts *types.EditTagProtectionOption) (*types.TagProtection, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tag_protections/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.TagProtection
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTagProtection deletes a tag protection from a repository.
func (s *RepoService) DeleteTagProtection(ctx context.Context, owner, repo string, id int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/tag_protections/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListKeys returns all deploy keys for a repository.
func (s *RepoService) ListKeys(ctx context.Context, owner, repo string, filters ...RepoKeyListOptions) ([]types.DeployKey, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/keys", pathParams("owner", owner, "repo", repo))
	return ListAll[types.DeployKey](ctx, s.client, path, repoKeyQuery(filters...))
}

// IterKeys returns an iterator over all deploy keys for a repository.
func (s *RepoService) IterKeys(ctx context.Context, owner, repo string, filters ...RepoKeyListOptions) iter.Seq2[types.DeployKey, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/keys", pathParams("owner", owner, "repo", repo))
	return ListIter[types.DeployKey](ctx, s.client, path, repoKeyQuery(filters...))
}

// GetKey returns a single deploy key by ID.
func (s *RepoService) GetKey(ctx context.Context, owner, repo string, id int64) (*types.DeployKey, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/keys/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	var out types.DeployKey
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateKey adds a deploy key to a repository.
func (s *RepoService) CreateKey(ctx context.Context, owner, repo string, opts *types.CreateKeyOption) (*types.DeployKey, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/keys", pathParams("owner", owner, "repo", repo))
	var out types.DeployKey
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteKey removes a deploy key from a repository by ID.
func (s *RepoService) DeleteKey(ctx context.Context, owner, repo string, id int64) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/keys/{id}", pathParams("owner", owner, "repo", repo, "id", int64String(id)))
	return s.client.Delete(ctx, path)
}

// ListStargazers returns all users who starred a repository.
func (s *RepoService) ListStargazers(ctx context.Context, owner, repo string) ([]types.User, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/stargazers", pathParams("owner", owner, "repo", repo))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterStargazers returns an iterator over all users who starred a repository.
func (s *RepoService) IterStargazers(ctx context.Context, owner, repo string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/stargazers", pathParams("owner", owner, "repo", repo))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// ListSubscribers returns all users watching a repository.
func (s *RepoService) ListSubscribers(ctx context.Context, owner, repo string) ([]types.User, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/subscribers", pathParams("owner", owner, "repo", repo))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterSubscribers returns an iterator over all users watching a repository.
func (s *RepoService) IterSubscribers(ctx context.Context, owner, repo string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/subscribers", pathParams("owner", owner, "repo", repo))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// ListAssignees returns all users that can be assigned to issues in a repository.
func (s *RepoService) ListAssignees(ctx context.Context, owner, repo string) ([]types.User, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/assignees", pathParams("owner", owner, "repo", repo))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterAssignees returns an iterator over all users that can be assigned to issues in a repository.
func (s *RepoService) IterAssignees(ctx context.Context, owner, repo string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/assignees", pathParams("owner", owner, "repo", repo))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// ListCollaborators returns all collaborators on a repository.
func (s *RepoService) ListCollaborators(ctx context.Context, owner, repo string) ([]types.User, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/collaborators", pathParams("owner", owner, "repo", repo))
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterCollaborators returns an iterator over all collaborators on a repository.
func (s *RepoService) IterCollaborators(ctx context.Context, owner, repo string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/collaborators", pathParams("owner", owner, "repo", repo))
	return ListIter[types.User](ctx, s.client, path, nil)
}

// ListRepoTeams returns all teams assigned to a repository.
func (s *RepoService) ListRepoTeams(ctx context.Context, owner, repo string) ([]types.Team, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/teams", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Team](ctx, s.client, path, nil)
}

// IterRepoTeams returns an iterator over all teams assigned to a repository.
func (s *RepoService) IterRepoTeams(ctx context.Context, owner, repo string) iter.Seq2[types.Team, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/teams", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Team](ctx, s.client, path, nil)
}

// GetRepoTeam returns a team assigned to a repository by name.
func (s *RepoService) GetRepoTeam(ctx context.Context, owner, repo, team string) (*types.Team, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/teams/{team}", pathParams("owner", owner, "repo", repo, "team", team))
	var out types.Team
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddRepoTeam assigns a team to a repository.
func (s *RepoService) AddRepoTeam(ctx context.Context, owner, repo, team string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/teams/{team}", pathParams("owner", owner, "repo", repo, "team", team))
	return s.client.Put(ctx, path, nil, nil)
}

// DeleteRepoTeam removes a team from a repository.
func (s *RepoService) DeleteRepoTeam(ctx context.Context, owner, repo, team string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/teams/{team}", pathParams("owner", owner, "repo", repo, "team", team))
	return s.client.Delete(ctx, path)
}

// CheckCollaborator reports whether a user is a collaborator on a repository.
func (s *RepoService) CheckCollaborator(ctx context.Context, owner, repo, collaborator string) (bool, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/collaborators/{collaborator}", pathParams("owner", owner, "repo", repo, "collaborator", collaborator))
	resp, err := s.client.doJSON(ctx, "GET", path, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == 204, nil
}

// AddCollaborator adds a user as a collaborator on a repository.
func (s *RepoService) AddCollaborator(ctx context.Context, owner, repo, collaborator string, opts *types.AddCollaboratorOption) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/collaborators/{collaborator}", pathParams("owner", owner, "repo", repo, "collaborator", collaborator))
	return s.client.Put(ctx, path, opts, nil)
}

// DeleteCollaborator removes a user from a repository's collaborators.
func (s *RepoService) DeleteCollaborator(ctx context.Context, owner, repo, collaborator string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/collaborators/{collaborator}", pathParams("owner", owner, "repo", repo, "collaborator", collaborator))
	return s.client.Delete(ctx, path)
}

// GetCollaboratorPermission returns repository permissions for a collaborator.
func (s *RepoService) GetCollaboratorPermission(ctx context.Context, owner, repo, collaborator string) (*types.RepoCollaboratorPermission, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/collaborators/{collaborator}/permission", pathParams("owner", owner, "repo", repo, "collaborator", collaborator))
	var out types.RepoCollaboratorPermission
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRepoPermissions returns repository permissions for a user.
func (s *RepoService) GetRepoPermissions(ctx context.Context, owner, repo, collaborator string) (*types.RepoCollaboratorPermission, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/collaborators/{collaborator}/permission", pathParams("owner", owner, "repo", repo, "collaborator", collaborator))
	var out types.RepoCollaboratorPermission
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetArchive returns a repository archive as raw bytes.
func (s *RepoService) GetArchive(ctx context.Context, owner, repo, archive string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/archive/{archive}", pathParams("owner", owner, "repo", repo, "archive", archive))
	return s.client.GetRaw(ctx, path)
}

// Compare returns commit comparison information between two branches or commits.
func (s *RepoService) Compare(ctx context.Context, owner, repo, basehead string) (*types.Compare, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/compare/{basehead}", pathParams("owner", owner, "repo", repo, "basehead", basehead))
	var out types.Compare
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetRawFile returns the raw content of a repository file as bytes.
func (s *RepoService) GetRawFile(ctx context.Context, owner, repo, filepath string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/raw/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	return s.client.GetRaw(ctx, path)
}

// GetRawFileOrLFS returns the raw content or LFS object for a repository file as bytes.
func (s *RepoService) GetRawFileOrLFS(ctx context.Context, owner, repo, filepath, ref string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/media/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	if ref != "" {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		q.Set("ref", ref)
		u.RawQuery = q.Encode()
		path = u.String()
	}
	return s.client.GetRaw(ctx, path)
}

// GetEditorConfig returns the EditorConfig definitions for a repository file.
func (s *RepoService) GetEditorConfig(ctx context.Context, owner, repo, filepath, ref string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/editorconfig/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	if ref != "" {
		u, err := url.Parse(path)
		if err != nil {
			return err
		}
		q := u.Query()
		q.Set("ref", ref)
		u.RawQuery = q.Encode()
		path = u.String()
	}
	return s.client.Get(ctx, path, nil)
}

// ApplyDiffPatch applies a diff patch to a repository.
func (s *RepoService) ApplyDiffPatch(ctx context.Context, owner, repo string, opts *types.UpdateFileOptions) (*types.FileResponse, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/diffpatch", pathParams("owner", owner, "repo", repo))
	var out types.FileResponse
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLanguages returns the byte counts per language for a repository.
func (s *RepoService) GetLanguages(ctx context.Context, owner, repo string) (map[string]int64, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/languages", pathParams("owner", owner, "repo", repo))
	var out map[string]int64
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListFlags returns all flags for a repository.
func (s *RepoService) ListFlags(ctx context.Context, owner, repo string) ([]string, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/flags", pathParams("owner", owner, "repo", repo))
	var out []string
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IterFlags returns an iterator over all flags for a repository.
func (s *RepoService) IterFlags(ctx context.Context, owner, repo string) iter.Seq2[string, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/flags", pathParams("owner", owner, "repo", repo))
	return ListIter[string](ctx, s.client, path, nil)
}

// ReplaceFlags replaces all flags for a repository.
func (s *RepoService) ReplaceFlags(ctx context.Context, owner, repo string, opts *types.ReplaceFlagsOption) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/flags", pathParams("owner", owner, "repo", repo))
	return s.client.Put(ctx, path, opts, nil)
}

// DeleteFlags removes all flags from a repository.
func (s *RepoService) DeleteFlags(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/flags", pathParams("owner", owner, "repo", repo))
	return s.client.Delete(ctx, path)
}

// GetSigningKey returns the repository signing key as ASCII-armoured text.
func (s *RepoService) GetSigningKey(ctx context.Context, owner, repo string) (string, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/signing-key.gpg", pathParams("owner", owner, "repo", repo))
	data, err := s.client.GetRaw(ctx, path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ListIssueTemplates returns all issue templates available for a repository.
func (s *RepoService) ListIssueTemplates(ctx context.Context, owner, repo string) ([]types.IssueTemplate, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issue_templates", pathParams("owner", owner, "repo", repo))
	return ListAll[types.IssueTemplate](ctx, s.client, path, nil)
}

// IterIssueTemplates returns an iterator over all issue templates available for a repository.
func (s *RepoService) IterIssueTemplates(ctx context.Context, owner, repo string) iter.Seq2[types.IssueTemplate, error] {
	return func(yield func(types.IssueTemplate, error) bool) {
		templates, err := s.ListIssueTemplates(ctx, owner, repo)
		if err != nil {
			yield(*new(types.IssueTemplate), err)
			return
		}
		for _, template := range templates {
			if !yield(template, nil) {
				return
			}
		}
	}
}

// GetIssueConfig returns the issue config for a repository.
func (s *RepoService) GetIssueConfig(ctx context.Context, owner, repo string) (*types.IssueConfig, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issue_config", pathParams("owner", owner, "repo", repo))
	var out types.IssueConfig
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ValidateIssueConfig returns the validation information for a repository's issue config.
func (s *RepoService) ValidateIssueConfig(ctx context.Context, owner, repo string) (*types.IssueConfigValidation, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issue_config/validate", pathParams("owner", owner, "repo", repo))
	var out types.IssueConfigValidation
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListActivityFeeds returns the repository's activity feed entries.
func (s *RepoService) ListActivityFeeds(ctx context.Context, owner, repo string, filters ...ActivityFeedListOptions) ([]types.Activity, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/activities/feeds", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Activity](ctx, s.client, path, activityFeedQuery(filters...))
}

// IterActivityFeeds returns an iterator over the repository's activity feed entries.
func (s *RepoService) IterActivityFeeds(ctx context.Context, owner, repo string, filters ...ActivityFeedListOptions) iter.Seq2[types.Activity, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/activities/feeds", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Activity](ctx, s.client, path, activityFeedQuery(filters...))
}

// ListTopics returns the topics assigned to a repository.
func (s *RepoService) ListTopics(ctx context.Context, owner, repo string) ([]string, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/topics", pathParams("owner", owner, "repo", repo))
	var out types.TopicName
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out.TopicNames, nil
}

// IterTopics returns an iterator over the topics assigned to a repository.
func (s *RepoService) IterTopics(ctx context.Context, owner, repo string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		topics, err := s.ListTopics(ctx, owner, repo)
		if err != nil {
			yield("", err)
			return
		}
		for _, topic := range topics {
			if !yield(topic, nil) {
				return
			}
		}
	}
}

// SearchTopics searches topics by keyword.
func (s *RepoService) SearchTopics(ctx context.Context, query string) ([]types.TopicResponse, error) {
	return ListAll[types.TopicResponse](ctx, s.client, "/api/v1/topics/search", map[string]string{"q": query})
}

// IterSearchTopics returns an iterator over topic search results.
func (s *RepoService) IterSearchTopics(ctx context.Context, query string) iter.Seq2[types.TopicResponse, error] {
	return ListIter[types.TopicResponse](ctx, s.client, "/api/v1/topics/search", map[string]string{"q": query})
}

// SearchRepositoriesPage returns a single page of repository search results.
func (s *RepoService) SearchRepositoriesPage(ctx context.Context, query string, pageOpts ListOptions) (*PagedResult[types.Repository], error) {
	if pageOpts.Page < 1 {
		pageOpts.Page = 1
	}
	if pageOpts.Limit < 1 {
		pageOpts.Limit = 50
	}

	u, err := url.Parse("/api/v1/repos/search")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("q", query)
	q.Set("page", strconv.Itoa(pageOpts.Page))
	q.Set("limit", strconv.Itoa(pageOpts.Limit))
	u.RawQuery = q.Encode()

	var out types.SearchResults
	resp, err := s.client.doJSON(ctx, http.MethodGet, u.String(), nil, &out)
	if err != nil {
		return nil, err
	}

	totalCount, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))
	items := make([]types.Repository, 0, len(out.Data))
	for _, repo := range out.Data {
		if repo != nil {
			items = append(items, *repo)
		}
	}

	return &PagedResult[types.Repository]{
		Items:      items,
		TotalCount: totalCount,
		Page:       pageOpts.Page,
		HasMore: (totalCount > 0 && (pageOpts.Page-1)*pageOpts.Limit+len(items) < totalCount) ||
			(totalCount == 0 && len(items) >= pageOpts.Limit),
	}, nil
}

// SearchRepositories returns all repositories matching the search query.
func (s *RepoService) SearchRepositories(ctx context.Context, query string) ([]types.Repository, error) {
	var all []types.Repository
	page := 1

	for {
		result, err := s.SearchRepositoriesPage(ctx, query, ListOptions{Page: page, Limit: 50})
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

// IterSearchRepositories returns an iterator over all repositories matching the search query.
func (s *RepoService) IterSearchRepositories(ctx context.Context, query string) iter.Seq2[types.Repository, error] {
	return func(yield func(types.Repository, error) bool) {
		page := 1
		for {
			result, err := s.SearchRepositoriesPage(ctx, query, ListOptions{Page: page, Limit: 50})
			if err != nil {
				yield(*new(types.Repository), err)
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

// UpdateTopics replaces the topics assigned to a repository.
func (s *RepoService) UpdateTopics(ctx context.Context, owner, repo string, topics []string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/topics", pathParams("owner", owner, "repo", repo))
	return s.client.Put(ctx, path, types.RepoTopicOptions{Topics: topics}, nil)
}

// AddTopic adds a topic to a repository.
func (s *RepoService) AddTopic(ctx context.Context, owner, repo, topic string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/topics/{topic}", pathParams("owner", owner, "repo", repo, "topic", topic))
	return s.client.Put(ctx, path, nil, nil)
}

// DeleteTopic removes a topic from a repository.
func (s *RepoService) DeleteTopic(ctx context.Context, owner, repo, topic string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/topics/{topic}", pathParams("owner", owner, "repo", repo, "topic", topic))
	return s.client.Delete(ctx, path)
}

// AddFlag adds a flag to a repository.
func (s *RepoService) AddFlag(ctx context.Context, owner, repo, flag string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/flags/{flag}", pathParams("owner", owner, "repo", repo, "flag", flag))
	return s.client.Put(ctx, path, nil, nil)
}

// HasFlag reports whether a repository has a given flag.
func (s *RepoService) HasFlag(ctx context.Context, owner, repo, flag string) (bool, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/flags/{flag}", pathParams("owner", owner, "repo", repo, "flag", flag))
	resp, err := s.client.doJSON(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		if IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return resp.StatusCode == http.StatusNoContent, nil
}

// RemoveFlag removes a flag from a repository.
func (s *RepoService) RemoveFlag(ctx context.Context, owner, repo, flag string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/flags/{flag}", pathParams("owner", owner, "repo", repo, "flag", flag))
	return s.client.Delete(ctx, path)
}

// GetNewPinAllowed returns whether new issue pins are allowed for a repository.
func (s *RepoService) GetNewPinAllowed(ctx context.Context, owner, repo string) (*types.NewIssuePinsAllowed, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/new_pin_allowed", pathParams("owner", owner, "repo", repo))
	var out types.NewIssuePinsAllowed
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListPinnedPullRequests returns all pinned pull requests in a repository.
func (s *RepoService) ListPinnedPullRequests(ctx context.Context, owner, repo string) ([]types.PullRequest, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/pinned", pathParams("owner", owner, "repo", repo))
	return ListAll[types.PullRequest](ctx, s.client, path, nil)
}

// IterPinnedPullRequests returns an iterator over all pinned pull requests in a repository.
func (s *RepoService) IterPinnedPullRequests(ctx context.Context, owner, repo string) iter.Seq2[types.PullRequest, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/pulls/pinned", pathParams("owner", owner, "repo", repo))
	return ListIter[types.PullRequest](ctx, s.client, path, nil)
}

// UpdateAvatar updates a repository avatar.
func (s *RepoService) UpdateAvatar(ctx context.Context, owner, repo string, opts *types.UpdateRepoAvatarOption) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/avatar", pathParams("owner", owner, "repo", repo))
	return s.client.Post(ctx, path, opts, nil)
}

// DeleteAvatar deletes a repository avatar.
func (s *RepoService) DeleteAvatar(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/avatar", pathParams("owner", owner, "repo", repo))
	return s.client.Delete(ctx, path)
}

// ListPushMirrors returns all push mirrors configured for a repository.
func (s *RepoService) ListPushMirrors(ctx context.Context, owner, repo string) ([]types.PushMirror, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/push_mirrors", pathParams("owner", owner, "repo", repo))
	return ListAll[types.PushMirror](ctx, s.client, path, nil)
}

// IterPushMirrors returns an iterator over all push mirrors configured for a repository.
func (s *RepoService) IterPushMirrors(ctx context.Context, owner, repo string) iter.Seq2[types.PushMirror, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/push_mirrors", pathParams("owner", owner, "repo", repo))
	return ListIter[types.PushMirror](ctx, s.client, path, nil)
}

// GetPushMirror returns a push mirror by its remote name.
func (s *RepoService) GetPushMirror(ctx context.Context, owner, repo, name string) (*types.PushMirror, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/push_mirrors/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	var out types.PushMirror
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreatePushMirror adds a push mirror to a repository.
func (s *RepoService) CreatePushMirror(ctx context.Context, owner, repo string, opts *types.CreatePushMirrorOption) (*types.PushMirror, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/push_mirrors", pathParams("owner", owner, "repo", repo))
	var out types.PushMirror
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeletePushMirror removes a push mirror from a repository by remote name.
func (s *RepoService) DeletePushMirror(ctx context.Context, owner, repo, name string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/push_mirrors/{name}", pathParams("owner", owner, "repo", repo, "name", name))
	return s.client.Delete(ctx, path)
}

// GetSubscription returns the current user's watch state for a repository.
func (s *RepoService) GetSubscription(ctx context.Context, owner, repo string) (*types.WatchInfo, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/subscription", pathParams("owner", owner, "repo", repo))
	var out types.WatchInfo
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Watch subscribes the current user to repository notifications.
func (s *RepoService) Watch(ctx context.Context, owner, repo string) (*types.WatchInfo, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/subscription", pathParams("owner", owner, "repo", repo))
	var out types.WatchInfo
	if err := s.client.Put(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Unwatch unsubscribes the current user from repository notifications.
func (s *RepoService) Unwatch(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/subscription", pathParams("owner", owner, "repo", repo))
	return s.client.Delete(ctx, path)
}

// Fork forks a repository into the authenticated user's namespace or the
// optional organisation.
func (s *RepoService) Fork(ctx context.Context, owner, repo, org string) (*types.Repository, error) {
	opts := &types.CreateForkOption{Organization: org}
	return s.ForkWithOptions(ctx, owner, repo, opts)
}

// ForkWithOptions forks a repository with full control over the fork target.
func (s *RepoService) ForkWithOptions(ctx context.Context, owner, repo string, opts *types.CreateForkOption) (*types.Repository, error) {
	var out types.Repository
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/forks", pathParams("owner", owner, "repo", repo))
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Generate creates a repository from a template repository.
func (s *RepoService) Generate(ctx context.Context, templateOwner, templateRepo string, opts *types.GenerateRepoOption) (*types.Repository, error) {
	path := ResolvePath("/api/v1/repos/{template_owner}/{template_repo}/generate", pathParams("template_owner", templateOwner, "template_repo", templateRepo))
	var out types.Repository
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListForks returns all forks of a repository.
func (s *RepoService) ListForks(ctx context.Context, owner, repo string) ([]types.Repository, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/forks", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Repository](ctx, s.client, path, nil)
}

// IterForks returns an iterator over all forks of a repository.
func (s *RepoService) IterForks(ctx context.Context, owner, repo string) iter.Seq2[types.Repository, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/forks", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Repository](ctx, s.client, path, nil)
}

// Transfer initiates a repository transfer.
func (s *RepoService) Transfer(ctx context.Context, owner, repo string, opts *types.TransferRepoOption) (*types.Repository, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/transfer", pathParams("owner", owner, "repo", repo))
	var out types.Repository
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AcceptTransfer accepts a pending repository transfer.
func (s *RepoService) AcceptTransfer(ctx context.Context, owner, repo string) (*types.Repository, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/transfer/accept", pathParams("owner", owner, "repo", repo))
	var out types.Repository
	if err := s.client.Post(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RejectTransfer rejects a pending repository transfer.
func (s *RepoService) RejectTransfer(ctx context.Context, owner, repo string) (*types.Repository, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/transfer/reject", pathParams("owner", owner, "repo", repo))
	var out types.Repository
	if err := s.client.Post(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MirrorSync triggers a mirror sync.
func (s *RepoService) MirrorSync(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/mirror-sync", pathParams("owner", owner, "repo", repo))
	return s.client.Post(ctx, path, nil, nil)
}

// GetRunnerRegistrationToken returns a repository actions runner registration token.
func (s *RepoService) GetRunnerRegistrationToken(ctx context.Context, owner, repo string) (string, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/actions/runners/registration-token", pathParams("owner", owner, "repo", repo))
	resp, err := s.client.doJSON(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return "", err
	}
	return resp.Header.Get("token"), nil
}

// SyncPushMirrors triggers a sync across all push mirrors configured for a repository.
func (s *RepoService) SyncPushMirrors(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/push_mirrors-sync", pathParams("owner", owner, "repo", repo))
	return s.client.Post(ctx, path, nil, nil)
}

// GetBlob returns the blob content for a repository object.
func (s *RepoService) GetBlob(ctx context.Context, owner, repo, sha string) (*types.GitBlobResponse, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/blobs/{sha}", pathParams("owner", owner, "repo", repo, "sha", sha))
	var out types.GitBlobResponse
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListGitRefs returns all git references for a repository.
func (s *RepoService) ListGitRefs(ctx context.Context, owner, repo string) ([]types.Reference, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/refs", pathParams("owner", owner, "repo", repo))
	return ListAll[types.Reference](ctx, s.client, path, nil)
}

// IterGitRefs returns an iterator over all git references for a repository.
func (s *RepoService) IterGitRefs(ctx context.Context, owner, repo string) iter.Seq2[types.Reference, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/refs", pathParams("owner", owner, "repo", repo))
	return ListIter[types.Reference](ctx, s.client, path, nil)
}

// ListGitRefsByRef returns all git references matching a ref prefix.
func (s *RepoService) ListGitRefsByRef(ctx context.Context, owner, repo, ref string) ([]types.Reference, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/refs/{ref}", pathParams("owner", owner, "repo", repo, "ref", ref))
	return ListAll[types.Reference](ctx, s.client, path, nil)
}

// IterGitRefsByRef returns an iterator over all git references matching a ref prefix.
func (s *RepoService) IterGitRefsByRef(ctx context.Context, owner, repo, ref string) iter.Seq2[types.Reference, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/refs/{ref}", pathParams("owner", owner, "repo", repo, "ref", ref))
	return ListIter[types.Reference](ctx, s.client, path, nil)
}

// GetAnnotatedTag returns the annotated tag object for a tag SHA.
func (s *RepoService) GetAnnotatedTag(ctx context.Context, owner, repo, sha string) (*types.AnnotatedTag, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/tags/{sha}", pathParams("owner", owner, "repo", repo, "sha", sha))
	var out types.AnnotatedTag
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTree returns the git tree for a repository object.
func (s *RepoService) GetTree(ctx context.Context, owner, repo, sha string) (*types.GitTreeResponse, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/trees/{sha}", pathParams("owner", owner, "repo", repo, "sha", sha))
	var out types.GitTreeResponse
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTimes returns all tracked times for a repository.
func (s *RepoService) ListTimes(ctx context.Context, owner, repo string, filters ...RepoTimeListOptions) ([]types.TrackedTime, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/times", pathParams("owner", owner, "repo", repo))
	return ListAll[types.TrackedTime](ctx, s.client, path, repoTimeQuery(filters...))
}

// IterTimes returns an iterator over all tracked times for a repository.
func (s *RepoService) IterTimes(ctx context.Context, owner, repo string, filters ...RepoTimeListOptions) iter.Seq2[types.TrackedTime, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/times", pathParams("owner", owner, "repo", repo))
	return ListIter[types.TrackedTime](ctx, s.client, path, repoTimeQuery(filters...))
}

// ListUserTimes returns all tracked times for a user in a repository.
func (s *RepoService) ListUserTimes(ctx context.Context, owner, repo, username string) ([]types.TrackedTime, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/times/{user}", pathParams("owner", owner, "repo", repo, "user", username))
	return ListAll[types.TrackedTime](ctx, s.client, path, nil)
}

// IterUserTimes returns an iterator over all tracked times for a user in a repository.
func (s *RepoService) IterUserTimes(ctx context.Context, owner, repo, username string) iter.Seq2[types.TrackedTime, error] {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/times/{user}", pathParams("owner", owner, "repo", repo, "user", username))
	return ListIter[types.TrackedTime](ctx, s.client, path, nil)
}

func repoKeyQuery(filters ...RepoKeyListOptions) map[string]string {
	if len(filters) == 0 {
		return nil
	}

	query := make(map[string]string, 2)
	for _, filter := range filters {
		if filter.KeyID != 0 {
			query["key_id"] = strconv.FormatInt(filter.KeyID, 10)
		}
		if filter.Fingerprint != "" {
			query["fingerprint"] = filter.Fingerprint
		}
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func repoTimeQuery(filters ...RepoTimeListOptions) map[string]string {
	if len(filters) == 0 {
		return nil
	}

	query := make(map[string]string, 3)
	for _, filter := range filters {
		for key, value := range filter.queryParams() {
			query[key] = value
		}
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func activityFeedQuery(filters ...ActivityFeedListOptions) map[string]string {
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
