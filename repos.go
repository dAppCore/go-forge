package forge

import (
	"context"
	"iter"

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

func newRepoService(c *Client) *RepoService {
	return &RepoService{
		Resource: *NewResource[types.Repository, types.CreateRepoOption, types.EditRepoOption](
			c, "/api/v1/repos/{owner}/{repo}",
		),
	}
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

// GetArchive returns a repository archive as raw bytes.
func (s *RepoService) GetArchive(ctx context.Context, owner, repo, archive string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/archive/{archive}", pathParams("owner", owner, "repo", repo, "archive", archive))
	return s.client.GetRaw(ctx, path)
}

// GetRawFile returns the raw content of a repository file as bytes.
func (s *RepoService) GetRawFile(ctx context.Context, owner, repo, filepath string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/raw/{filepath}", pathParams("owner", owner, "repo", repo, "filepath", filepath))
	return s.client.GetRaw(ctx, path)
}

// ListIssueTemplates returns all issue templates available for a repository.
func (s *RepoService) ListIssueTemplates(ctx context.Context, owner, repo string) ([]types.IssueTemplate, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/issue_templates", pathParams("owner", owner, "repo", repo))
	return ListAll[types.IssueTemplate](ctx, s.client, path, nil)
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

// GetNewPinAllowed returns whether new issue pins are allowed for a repository.
func (s *RepoService) GetNewPinAllowed(ctx context.Context, owner, repo string) (*types.NewIssuePinsAllowed, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/new_pin_allowed", pathParams("owner", owner, "repo", repo))
	var out types.NewIssuePinsAllowed
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
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

// Fork forks a repository. If org is non-empty, forks into that organisation.
func (s *RepoService) Fork(ctx context.Context, owner, repo, org string) (*types.Repository, error) {
	body := map[string]string{}
	if org != "" {
		body["organization"] = org
	}
	var out types.Repository
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/forks", pathParams("owner", owner, "repo", repo))
	err := s.client.Post(ctx, path, body, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Transfer initiates a repository transfer.
func (s *RepoService) Transfer(ctx context.Context, owner, repo string, opts map[string]any) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/transfer", pathParams("owner", owner, "repo", repo))
	return s.client.Post(ctx, path, opts, nil)
}

// AcceptTransfer accepts a pending repository transfer.
func (s *RepoService) AcceptTransfer(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/transfer/accept", pathParams("owner", owner, "repo", repo))
	return s.client.Post(ctx, path, nil, nil)
}

// RejectTransfer rejects a pending repository transfer.
func (s *RepoService) RejectTransfer(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/transfer/reject", pathParams("owner", owner, "repo", repo))
	return s.client.Post(ctx, path, nil, nil)
}

// MirrorSync triggers a mirror sync.
func (s *RepoService) MirrorSync(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/mirror-sync", pathParams("owner", owner, "repo", repo))
	return s.client.Post(ctx, path, nil, nil)
}

// SyncPushMirrors triggers a sync across all push mirrors configured for a repository.
func (s *RepoService) SyncPushMirrors(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/push_mirrors-sync", pathParams("owner", owner, "repo", repo))
	return s.client.Post(ctx, path, nil, nil)
}
