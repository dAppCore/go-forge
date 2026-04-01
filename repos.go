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

// GetArchive returns a repository archive as raw bytes.
func (s *RepoService) GetArchive(ctx context.Context, owner, repo, archive string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/archive/{archive}", pathParams("owner", owner, "repo", repo, "archive", archive))
	return s.client.GetRaw(ctx, path)
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
