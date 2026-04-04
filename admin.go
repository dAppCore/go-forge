package forge

import (
	"context"
	"iter"
	"net/http"
	"net/url"
	"strconv"

	core "dappco.re/go/core"
	"dappco.re/go/core/forge/types"
)

// AdminService handles site administration operations.
// Unlike other services, AdminService does not embed Resource[T,C,U]
// because admin endpoints are heterogeneous.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.Admin.ListUsers(ctx)
type AdminService struct {
	client *Client
}

// AdminActionsRunListOptions controls filtering for admin Actions run listings.
//
// Usage:
//
//	opts := forge.AdminActionsRunListOptions{Event: "push", Status: "success"}
type AdminActionsRunListOptions struct {
	Event   string
	Branch  string
	Status  string
	Actor   string
	HeadSHA string
}

// String returns a safe summary of the admin Actions run filters.
func (o AdminActionsRunListOptions) String() string {
	return optionString("forge.AdminActionsRunListOptions",
		"event", o.Event,
		"branch", o.Branch,
		"status", o.Status,
		"actor", o.Actor,
		"head_sha", o.HeadSHA,
	)
}

// GoString returns a safe Go-syntax summary of the admin Actions run filters.
func (o AdminActionsRunListOptions) GoString() string { return o.String() }

func (o AdminActionsRunListOptions) queryParams() map[string]string {
	query := make(map[string]string, 5)
	if o.Event != "" {
		query["event"] = o.Event
	}
	if o.Branch != "" {
		query["branch"] = o.Branch
	}
	if o.Status != "" {
		query["status"] = o.Status
	}
	if o.Actor != "" {
		query["actor"] = o.Actor
	}
	if o.HeadSHA != "" {
		query["head_sha"] = o.HeadSHA
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

// AdminUnadoptedListOptions controls filtering for unadopted repository listings.
//
// Usage:
//
//	opts := forge.AdminUnadoptedListOptions{Pattern: "core/*"}
type AdminUnadoptedListOptions struct {
	Pattern string
}

// String returns a safe summary of the unadopted repository filters.
func (o AdminUnadoptedListOptions) String() string {
	return optionString("forge.AdminUnadoptedListOptions", "pattern", o.Pattern)
}

// GoString returns a safe Go-syntax summary of the unadopted repository filters.
func (o AdminUnadoptedListOptions) GoString() string { return o.String() }

func (o AdminUnadoptedListOptions) queryParams() map[string]string {
	if o.Pattern == "" {
		return nil
	}
	return map[string]string{"pattern": o.Pattern}
}

func newAdminService(c *Client) *AdminService {
	return &AdminService{client: c}
}

// ListUsers returns all users (admin only).
func (s *AdminService) ListUsers(ctx context.Context) ([]types.User, error) {
	return ListAll[types.User](ctx, s.client, "/api/v1/admin/users", nil)
}

// IterUsers returns an iterator over all users (admin only).
func (s *AdminService) IterUsers(ctx context.Context) iter.Seq2[types.User, error] {
	return ListIter[types.User](ctx, s.client, "/api/v1/admin/users", nil)
}

// CreateUser creates a new user (admin only).
func (s *AdminService) CreateUser(ctx context.Context, opts *types.CreateUserOption) (*types.User, error) {
	var out types.User
	if err := s.client.Post(ctx, "/api/v1/admin/users", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateUserKey adds a public key on behalf of a user.
func (s *AdminService) CreateUserKey(ctx context.Context, username string, opts *types.CreateKeyOption) (*types.PublicKey, error) {
	path := ResolvePath("/api/v1/admin/users/{username}/keys", Params{"username": username})
	var out types.PublicKey
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteUserKey deletes a user's public key.
func (s *AdminService) DeleteUserKey(ctx context.Context, username string, id int64) error {
	path := ResolvePath("/api/v1/admin/users/{username}/keys/{id}", Params{"username": username, "id": int64String(id)})
	return s.client.Delete(ctx, path)
}

// CreateUserOrg creates an organisation on behalf of a user.
func (s *AdminService) CreateUserOrg(ctx context.Context, username string, opts *types.CreateOrgOption) (*types.Organization, error) {
	path := ResolvePath("/api/v1/admin/users/{username}/orgs", Params{"username": username})
	var out types.Organization
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUserQuota returns a user's quota information.
func (s *AdminService) GetUserQuota(ctx context.Context, username string) (*types.QuotaInfo, error) {
	path := ResolvePath("/api/v1/admin/users/{username}/quota", Params{"username": username})
	var out types.QuotaInfo
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SetUserQuotaGroups sets the user's quota groups to a given list.
func (s *AdminService) SetUserQuotaGroups(ctx context.Context, username string, opts *types.SetUserQuotaGroupsOptions) error {
	path := ResolvePath("/api/v1/admin/users/{username}/quota/groups", Params{"username": username})
	return s.client.Post(ctx, path, opts, nil)
}

// CreateUserRepo creates a repository on behalf of a user.
func (s *AdminService) CreateUserRepo(ctx context.Context, username string, opts *types.CreateRepoOption) (*types.Repository, error) {
	path := ResolvePath("/api/v1/admin/users/{username}/repos", Params{"username": username})
	var out types.Repository
	if err := s.client.Post(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditUser edits an existing user (admin only).
func (s *AdminService) EditUser(ctx context.Context, username string, opts map[string]any) error {
	path := ResolvePath("/api/v1/admin/users/{username}", Params{"username": username})
	return s.client.Patch(ctx, path, opts, nil)
}

// DeleteUser deletes a user (admin only).
func (s *AdminService) DeleteUser(ctx context.Context, username string) error {
	path := ResolvePath("/api/v1/admin/users/{username}", Params{"username": username})
	return s.client.Delete(ctx, path)
}

// RenameUser renames a user (admin only).
func (s *AdminService) RenameUser(ctx context.Context, username, newName string) error {
	path := ResolvePath("/api/v1/admin/users/{username}/rename", Params{"username": username})
	return s.client.Post(ctx, path, &types.RenameUserOption{NewName: newName}, nil)
}

// ListOrgs returns all organisations (admin only).
func (s *AdminService) ListOrgs(ctx context.Context) ([]types.Organization, error) {
	return ListAll[types.Organization](ctx, s.client, "/api/v1/admin/orgs", nil)
}

// IterOrgs returns an iterator over all organisations (admin only).
func (s *AdminService) IterOrgs(ctx context.Context) iter.Seq2[types.Organization, error] {
	return ListIter[types.Organization](ctx, s.client, "/api/v1/admin/orgs", nil)
}

// ListEmails returns all email addresses (admin only).
func (s *AdminService) ListEmails(ctx context.Context) ([]types.Email, error) {
	return ListAll[types.Email](ctx, s.client, "/api/v1/admin/emails", nil)
}

// IterEmails returns an iterator over all email addresses (admin only).
func (s *AdminService) IterEmails(ctx context.Context) iter.Seq2[types.Email, error] {
	return ListIter[types.Email](ctx, s.client, "/api/v1/admin/emails", nil)
}

// ListHooks returns all global hooks (admin only).
func (s *AdminService) ListHooks(ctx context.Context) ([]types.Hook, error) {
	return ListAll[types.Hook](ctx, s.client, "/api/v1/admin/hooks", nil)
}

// IterHooks returns an iterator over all global hooks (admin only).
func (s *AdminService) IterHooks(ctx context.Context) iter.Seq2[types.Hook, error] {
	return ListIter[types.Hook](ctx, s.client, "/api/v1/admin/hooks", nil)
}

// GetHook returns a single global hook by ID (admin only).
func (s *AdminService) GetHook(ctx context.Context, id int64) (*types.Hook, error) {
	path := ResolvePath("/api/v1/admin/hooks/{id}", Params{"id": int64String(id)})
	var out types.Hook
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateHook creates a new global hook (admin only).
func (s *AdminService) CreateHook(ctx context.Context, opts *types.CreateHookOption) (*types.Hook, error) {
	var out types.Hook
	if err := s.client.Post(ctx, "/api/v1/admin/hooks", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditHook updates an existing global hook (admin only).
func (s *AdminService) EditHook(ctx context.Context, id int64, opts *types.EditHookOption) (*types.Hook, error) {
	path := ResolvePath("/api/v1/admin/hooks/{id}", Params{"id": int64String(id)})
	var out types.Hook
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteHook deletes a global hook (admin only).
func (s *AdminService) DeleteHook(ctx context.Context, id int64) error {
	path := ResolvePath("/api/v1/admin/hooks/{id}", Params{"id": int64String(id)})
	return s.client.Delete(ctx, path)
}

// ListQuotaGroups returns all available quota groups.
func (s *AdminService) ListQuotaGroups(ctx context.Context) ([]types.QuotaGroup, error) {
	return ListAll[types.QuotaGroup](ctx, s.client, "/api/v1/admin/quota/groups", nil)
}

// IterQuotaGroups returns an iterator over all available quota groups.
func (s *AdminService) IterQuotaGroups(ctx context.Context) iter.Seq2[types.QuotaGroup, error] {
	return func(yield func(types.QuotaGroup, error) bool) {
		groups, err := s.ListQuotaGroups(ctx)
		if err != nil {
			yield(*new(types.QuotaGroup), err)
			return
		}
		for _, group := range groups {
			if !yield(group, nil) {
				return
			}
		}
	}
}

// CreateQuotaGroup creates a new quota group.
func (s *AdminService) CreateQuotaGroup(ctx context.Context, opts *types.CreateQuotaGroupOptions) (*types.QuotaGroup, error) {
	var out types.QuotaGroup
	if err := s.client.Post(ctx, "/api/v1/admin/quota/groups", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetQuotaGroup returns information about a quota group.
func (s *AdminService) GetQuotaGroup(ctx context.Context, quotagroup string) (*types.QuotaGroup, error) {
	path := ResolvePath("/api/v1/admin/quota/groups/{quotagroup}", Params{"quotagroup": quotagroup})
	var out types.QuotaGroup
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteQuotaGroup deletes a quota group.
func (s *AdminService) DeleteQuotaGroup(ctx context.Context, quotagroup string) error {
	path := ResolvePath("/api/v1/admin/quota/groups/{quotagroup}", Params{"quotagroup": quotagroup})
	return s.client.Delete(ctx, path)
}

// AddQuotaGroupRule adds a quota rule to a quota group.
func (s *AdminService) AddQuotaGroupRule(ctx context.Context, quotagroup, quotarule string) error {
	path := ResolvePath("/api/v1/admin/quota/groups/{quotagroup}/rules/{quotarule}", Params{"quotagroup": quotagroup, "quotarule": quotarule})
	return s.client.Put(ctx, path, nil, nil)
}

// RemoveQuotaGroupRule removes a quota rule from a quota group.
func (s *AdminService) RemoveQuotaGroupRule(ctx context.Context, quotagroup, quotarule string) error {
	path := ResolvePath("/api/v1/admin/quota/groups/{quotagroup}/rules/{quotarule}", Params{"quotagroup": quotagroup, "quotarule": quotarule})
	return s.client.Delete(ctx, path)
}

// ListQuotaGroupUsers returns all users in a quota group.
func (s *AdminService) ListQuotaGroupUsers(ctx context.Context, quotagroup string) ([]types.User, error) {
	path := ResolvePath("/api/v1/admin/quota/groups/{quotagroup}/users", Params{"quotagroup": quotagroup})
	return ListAll[types.User](ctx, s.client, path, nil)
}

// IterQuotaGroupUsers returns an iterator over all users in a quota group.
func (s *AdminService) IterQuotaGroupUsers(ctx context.Context, quotagroup string) iter.Seq2[types.User, error] {
	path := ResolvePath("/api/v1/admin/quota/groups/{quotagroup}/users", Params{"quotagroup": quotagroup})
	return ListIter[types.User](ctx, s.client, path, nil)
}

// AddQuotaGroupUser adds a user to a quota group.
func (s *AdminService) AddQuotaGroupUser(ctx context.Context, quotagroup, username string) error {
	path := ResolvePath("/api/v1/admin/quota/groups/{quotagroup}/users/{username}", Params{"quotagroup": quotagroup, "username": username})
	return s.client.Put(ctx, path, nil, nil)
}

// RemoveQuotaGroupUser removes a user from a quota group.
func (s *AdminService) RemoveQuotaGroupUser(ctx context.Context, quotagroup, username string) error {
	path := ResolvePath("/api/v1/admin/quota/groups/{quotagroup}/users/{username}", Params{"quotagroup": quotagroup, "username": username})
	return s.client.Delete(ctx, path)
}

// ListQuotaRules returns all available quota rules.
func (s *AdminService) ListQuotaRules(ctx context.Context) ([]types.QuotaRuleInfo, error) {
	return ListAll[types.QuotaRuleInfo](ctx, s.client, "/api/v1/admin/quota/rules", nil)
}

// IterQuotaRules returns an iterator over all available quota rules.
func (s *AdminService) IterQuotaRules(ctx context.Context) iter.Seq2[types.QuotaRuleInfo, error] {
	return func(yield func(types.QuotaRuleInfo, error) bool) {
		rules, err := s.ListQuotaRules(ctx)
		if err != nil {
			yield(*new(types.QuotaRuleInfo), err)
			return
		}
		for _, rule := range rules {
			if !yield(rule, nil) {
				return
			}
		}
	}
}

// CreateQuotaRule creates a new quota rule.
func (s *AdminService) CreateQuotaRule(ctx context.Context, opts *types.CreateQuotaRuleOptions) (*types.QuotaRuleInfo, error) {
	var out types.QuotaRuleInfo
	if err := s.client.Post(ctx, "/api/v1/admin/quota/rules", opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetQuotaRule returns information about a quota rule.
func (s *AdminService) GetQuotaRule(ctx context.Context, quotarule string) (*types.QuotaRuleInfo, error) {
	path := ResolvePath("/api/v1/admin/quota/rules/{quotarule}", Params{"quotarule": quotarule})
	var out types.QuotaRuleInfo
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// EditQuotaRule updates an existing quota rule.
func (s *AdminService) EditQuotaRule(ctx context.Context, quotarule string, opts *types.EditQuotaRuleOptions) (*types.QuotaRuleInfo, error) {
	path := ResolvePath("/api/v1/admin/quota/rules/{quotarule}", Params{"quotarule": quotarule})
	var out types.QuotaRuleInfo
	if err := s.client.Patch(ctx, path, opts, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteQuotaRule deletes a quota rule.
func (s *AdminService) DeleteQuotaRule(ctx context.Context, quotarule string) error {
	path := ResolvePath("/api/v1/admin/quota/rules/{quotarule}", Params{"quotarule": quotarule})
	return s.client.Delete(ctx, path)
}

// ListUnadoptedRepos returns all unadopted repositories on the instance.
func (s *AdminService) ListUnadoptedRepos(ctx context.Context, filters ...AdminUnadoptedListOptions) ([]string, error) {
	return ListAll[string](ctx, s.client, "/api/v1/admin/unadopted", adminUnadoptedQuery(filters...))
}

// IterUnadoptedRepos returns an iterator over all unadopted repositories on the instance.
func (s *AdminService) IterUnadoptedRepos(ctx context.Context, filters ...AdminUnadoptedListOptions) iter.Seq2[string, error] {
	return ListIter[string](ctx, s.client, "/api/v1/admin/unadopted", adminUnadoptedQuery(filters...))
}

// SearchEmails searches all email addresses by keyword (admin only).
func (s *AdminService) SearchEmails(ctx context.Context, q string) ([]types.Email, error) {
	return ListAll[types.Email](ctx, s.client, "/api/v1/admin/emails/search", map[string]string{"q": q})
}

// IterSearchEmails returns an iterator over all email addresses matching a keyword (admin only).
func (s *AdminService) IterSearchEmails(ctx context.Context, q string) iter.Seq2[types.Email, error] {
	return ListIter[types.Email](ctx, s.client, "/api/v1/admin/emails/search", map[string]string{"q": q})
}

// RunCron runs a cron task by name (admin only).
func (s *AdminService) RunCron(ctx context.Context, task string) error {
	path := ResolvePath("/api/v1/admin/cron/{task}", Params{"task": task})
	return s.client.Post(ctx, path, nil, nil)
}

// ListCron returns all cron tasks (admin only).
func (s *AdminService) ListCron(ctx context.Context) ([]types.Cron, error) {
	return ListAll[types.Cron](ctx, s.client, "/api/v1/admin/cron", nil)
}

// IterCron returns an iterator over all cron tasks (admin only).
func (s *AdminService) IterCron(ctx context.Context) iter.Seq2[types.Cron, error] {
	return ListIter[types.Cron](ctx, s.client, "/api/v1/admin/cron", nil)
}

// ListActionsRuns returns a single page of Actions workflow runs across the instance.
func (s *AdminService) ListActionsRuns(ctx context.Context, filters AdminActionsRunListOptions, opts ListOptions) (*PagedResult[types.ActionTask], error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 {
		opts.Limit = 50
	}

	u, err := url.Parse("/api/v1/admin/actions/runs")
	if err != nil {
		return nil, core.E("AdminService.ListActionsRuns", "forge: parse path", err)
	}

	q := u.Query()
	for key, value := range filters.queryParams() {
		q.Set(key, value)
	}
	q.Set("page", strconv.Itoa(opts.Page))
	q.Set("limit", strconv.Itoa(opts.Limit))
	u.RawQuery = q.Encode()

	var out types.ActionTaskResponse
	resp, err := s.client.doJSON(ctx, http.MethodGet, u.String(), nil, &out)
	if err != nil {
		return nil, err
	}

	totalCount, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))
	items := make([]types.ActionTask, 0, len(out.Entries))
	for _, run := range out.Entries {
		if run != nil {
			items = append(items, *run)
		}
	}

	return &PagedResult[types.ActionTask]{
		Items:      items,
		TotalCount: totalCount,
		Page:       opts.Page,
		HasMore: (totalCount > 0 && (opts.Page-1)*opts.Limit+len(items) < totalCount) ||
			(totalCount == 0 && len(items) >= opts.Limit),
	}, nil
}

// IterActionsRuns returns an iterator over all Actions workflow runs across the instance.
func (s *AdminService) IterActionsRuns(ctx context.Context, filters AdminActionsRunListOptions) iter.Seq2[types.ActionTask, error] {
	return func(yield func(types.ActionTask, error) bool) {
		page := 1
		for {
			result, err := s.ListActionsRuns(ctx, filters, ListOptions{Page: page, Limit: 50})
			if err != nil {
				yield(*new(types.ActionTask), err)
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

// AdoptRepo adopts an unadopted repository (admin only).
func (s *AdminService) AdoptRepo(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/admin/unadopted/{owner}/{repo}", Params{"owner": owner, "repo": repo})
	return s.client.Post(ctx, path, nil, nil)
}

// DeleteUnadoptedRepo deletes an unadopted repository's files.
func (s *AdminService) DeleteUnadoptedRepo(ctx context.Context, owner, repo string) error {
	path := ResolvePath("/api/v1/admin/unadopted/{owner}/{repo}", Params{"owner": owner, "repo": repo})
	return s.client.Delete(ctx, path)
}

func adminUnadoptedQuery(filters ...AdminUnadoptedListOptions) map[string]string {
	if len(filters) == 0 {
		return nil
	}

	query := make(map[string]string, 1)
	for _, filter := range filters {
		if filter.Pattern != "" {
			query["pattern"] = filter.Pattern
		}
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

// GenerateRunnerToken generates an actions runner registration token.
func (s *AdminService) GenerateRunnerToken(ctx context.Context) (string, error) {
	var out struct {
		Token string `json:"token"`
	}
	if err := s.client.Get(ctx, "/api/v1/admin/runners/registration-token", &out); err != nil {
		return "", err
	}
	return out.Token, nil
}
