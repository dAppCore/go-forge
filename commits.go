package forge

import (
	"context"
	"iter"

	"dappco.re/go/forge/types"
)

// CommitService handles commit-related operations such as commit statuses
// and git notes.
// No Resource embedding — collection and item commit paths differ, and the
// remaining endpoints are heterogeneous across status and note paths.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Commits.GetCombinedStatus(ctx, "core", "go-forge", "main")
type CommitService struct {
	client *Client
}

// CommitListOptions controls filtering for repository commit listings.
//
// Usage:
//
//	stat := false
//	opts := forge.CommitListOptions{Sha: "main", Stat: &stat}
type CommitListOptions struct {
	Sha          string
	Path         string
	Stat         *bool
	Verification *bool
	Files        *bool
	Not          string
}

// String returns a safe summary of the commit list filters.
func (o CommitListOptions) String() string {
	return optionString("forge.CommitListOptions",
		"sha", o.Sha,
		"path", o.Path,
		"stat", o.Stat,
		"verification", o.Verification,
		"files", o.Files,
		"not", o.Not,
	)
}

// GoString returns a safe Go-syntax summary of the commit list filters.
func (o CommitListOptions) GoString() string { return o.String() }

func (o CommitListOptions) queryParams() map[string]string {
	query := make(map[string]string, 6)
	if o.Sha != "" {
		query["sha"] = o.Sha
	}
	if o.Path != "" {
		query["path"] = o.Path
	}
	if o.Stat != nil {
		query["stat"] = boolString(*o.Stat)
	}
	if o.Verification != nil {
		query["verification"] = boolString(*o.Verification)
	}
	if o.Files != nil {
		query["files"] = boolString(*o.Files)
	}
	if o.Not != "" {
		query["not"] = o.Not
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

const (
	commitCollectionPath = "/api/v1/repos/{owner}/{repo}/commits"
	commitItemPath       = "/api/v1/repos/{owner}/{repo}/git/commits/{sha}"
)

func newCommitService(c *Client) *CommitService {
	return &CommitService{client: c}
}

// List returns a single page of commits for a repository.
func (s *CommitService) List(ctx context.Context, params Params, opts ListOptions, filters ...CommitListOptions) (*PagedResult[types.Commit], error) {
	return ListPage[types.Commit](ctx, s.client, ResolvePath(commitCollectionPath, params), commitListQuery(filters...), opts)
}

// ListAll returns all commits for a repository.
func (s *CommitService) ListAll(ctx context.Context, params Params, filters ...CommitListOptions) ([]types.Commit, error) {
	return ListAll[types.Commit](ctx, s.client, ResolvePath(commitCollectionPath, params), commitListQuery(filters...))
}

// Iter returns an iterator over all commits for a repository.
func (s *CommitService) Iter(ctx context.Context, params Params, filters ...CommitListOptions) iter.Seq2[types.Commit, error] {
	return ListIter[types.Commit](ctx, s.client, ResolvePath(commitCollectionPath, params), commitListQuery(filters...))
}

// Get returns a single commit by SHA or ref.
func (s *CommitService) Get(ctx context.Context, params Params) (*types.Commit, error) {
	var out types.Commit
	if err := s.client.Get(ctx, ResolvePath(commitItemPath, params), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListCommitsPage returns a single page of commits for a repository.
func (s *CommitService) ListCommitsPage(ctx context.Context, owner, repo string, opts ListOptions, filters ...any) (*PagedResult[types.Commit], error) {
	params := pathParams("owner", owner, "repo", repo)
	path := ResolvePath(commitCollectionPath, params)
	return ListPage[types.Commit](ctx, s.client, path, commitCompatListQuery(filters...), opts)
}

// ListCommits returns commits for a repository using RFC-compatible filters.
func (s *CommitService) ListCommits(ctx context.Context, owner, repo string, filters ...any) ([]types.Commit, error) {
	params := pathParams("owner", owner, "repo", repo)
	path := ResolvePath(commitCollectionPath, params)
	query := commitCompatListQuery(filters...)
	if pageOpts, ok := commitCompatPageOptions(filters...); ok {
		page, err := ListPage[types.Commit](ctx, s.client, path, query, pageOpts)
		if err != nil {
			return nil, err
		}
		return page.Items, nil
	}
	return ListAll[types.Commit](ctx, s.client, path, query)
}

// IterCommits returns an iterator over commits for a repository using RFC-compatible filters.
func (s *CommitService) IterCommits(ctx context.Context, owner, repo string, filters ...any) iter.Seq2[types.Commit, error] {
	params := pathParams("owner", owner, "repo", repo)
	path := ResolvePath(commitCollectionPath, params)
	query := commitCompatListQuery(filters...)
	if pageOpts, ok := commitCompatPageOptions(filters...); ok {
		return func(yield func(types.Commit, error) bool) {
			page, err := ListPage[types.Commit](ctx, s.client, path, query, pageOpts)
			if err != nil {
				yield(*new(types.Commit), err)
				return
			}
			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}
		}
	}
	return ListIter[types.Commit](ctx, s.client, path, query)
}

// GetCommit returns a single commit by SHA or ref.
func (s *CommitService) GetCommit(ctx context.Context, owner, repo, sha string) (*types.Commit, error) {
	return s.Get(ctx, pathParams("owner", owner, "repo", repo, "sha", sha))
}

// GetDiffOrPatch returns a commit diff or patch as raw bytes.
func (s *CommitService) GetDiffOrPatch(ctx context.Context, owner, repo, sha, diffType string) ([]byte, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/commits/{sha}.{diffType}", pathParams("owner", owner, "repo", repo, "sha", sha, "diffType", diffType))
	return s.client.GetRaw(ctx, path)
}

// GetPullRequest returns the pull request associated with a commit SHA.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Commits.GetPullRequest(ctx, "core", "go-forge", "abc123")
func (s *CommitService) GetPullRequest(ctx context.Context, owner, repo, sha string) (*types.PullRequest, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/commits/{sha}/pull", pathParams("owner", owner, "repo", repo, "sha", sha))
	var out types.PullRequest
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCombinedStatus returns the combined status for a given ref (branch, tag, or SHA).
func (s *CommitService) GetCombinedStatus(ctx context.Context, owner, repo, ref string) (*types.CombinedStatus, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/statuses/{ref}", pathParams("owner", owner, "repo", repo, "ref", ref))
	var out types.CombinedStatus
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCombinedStatusByRef returns the combined status for a given commit reference.
func (s *CommitService) GetCombinedStatusByRef(ctx context.Context, owner, repo, ref string) (*types.CombinedStatus, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/commits/{ref}/status", pathParams("owner", owner, "repo", repo, "ref", ref))
	var out types.CombinedStatus
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListStatuses returns all commit statuses for a given ref.
func (s *CommitService) ListStatuses(ctx context.Context, owner, repo, ref string) ([]types.CommitStatus, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/commits/{ref}/statuses", pathParams("owner", owner, "repo", repo, "ref", ref))
	var out []types.CommitStatus
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// IterStatuses returns an iterator over all commit statuses for a given ref.
func (s *CommitService) IterStatuses(ctx context.Context, owner, repo, ref string) iter.Seq2[types.CommitStatus, error] {
	return func(yield func(types.CommitStatus, error) bool) {
		statuses, err := s.ListStatuses(ctx, owner, repo, ref)
		if err != nil {
			yield(*new(types.CommitStatus), err)
			return
		}
		for _, status := range statuses {
			if !yield(status, nil) {
				return
			}
		}
	}
}

// CreateStatus creates a new commit status for the given SHA.
func (s *CommitService) CreateStatus(ctx context.Context, owner, repo, sha string, opts *types.CreateStatusOption) (*types.CommitStatus, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/statuses/{sha}", pathParams("owner", owner, "repo", repo, "sha", sha))
	var out types.CommitStatus
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetNote returns the git note for a given commit SHA.
func (s *CommitService) GetNote(ctx context.Context, owner, repo, sha string) (*types.Note, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/notes/{sha}", pathParams("owner", owner, "repo", repo, "sha", sha))
	var out types.Note
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetNote creates or updates the git note for a given commit SHA.
func (s *CommitService) SetNote(ctx context.Context, owner, repo, sha, message string) (*types.Note, error) {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/notes/{sha}", pathParams("owner", owner, "repo", repo, "sha", sha))
	var out types.Note
	if err := s.client.Post(ctx, path, types.NoteOptions{Message: message}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteNote removes the git note for a given commit SHA.
func (s *CommitService) DeleteNote(ctx context.Context, owner, repo, sha string) error {
	path := ResolvePath("/api/v1/repos/{owner}/{repo}/git/notes/{sha}", pathParams("owner", owner, "repo", repo, "sha", sha))
	return s.client.Delete(ctx, path)
}

func commitListQuery(filters ...CommitListOptions) map[string]string {
	query := make(map[string]string, len(filters))
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

func commitCompatListQuery(filters ...any) map[string]string {
	query := make(map[string]string, len(filters))
	for _, filter := range filters {
		switch v := filter.(type) {
		case CommitListOptions:
			for key, value := range v.queryParams() {
				query[key] = value
			}
		case *CommitListOptions:
			if v != nil {
				for key, value := range v.queryParams() {
					query[key] = value
				}
			}
		case types.ListCommitsOption:
			for key, value := range commitListQueryFromCompat(v) {
				query[key] = value
			}
		case *types.ListCommitsOption:
			if v != nil {
				for key, value := range commitListQueryFromCompat(*v) {
					query[key] = value
				}
			}
		}
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func commitListQueryFromCompat(filter types.ListCommitsOption) map[string]string {
	query := make(map[string]string, 6)
	if filter.Sha != "" {
		query["sha"] = filter.Sha
	}
	if filter.Path != "" {
		query["path"] = filter.Path
	}
	if filter.Stat != nil {
		query["stat"] = boolString(*filter.Stat)
	}
	if filter.Verification != nil {
		query["verification"] = boolString(*filter.Verification)
	}
	if filter.Files != nil {
		query["files"] = boolString(*filter.Files)
	}
	if filter.Not != "" {
		query["not"] = filter.Not
	}
	return query
}

func commitCompatPageOptions(filters ...any) (ListOptions, bool) {
	for _, filter := range filters {
		switch v := filter.(type) {
		case types.ListCommitsOption:
			if opts, ok := compatListOptions(v.Page, v.PageSize, v.Limit); ok {
				return opts, true
			}
		case *types.ListCommitsOption:
			if v != nil {
				if opts, ok := compatListOptions(v.Page, v.PageSize, v.Limit); ok {
					return opts, true
				}
			}
		}
	}
	return ListOptions{}, false
}
