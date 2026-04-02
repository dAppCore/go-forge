package forge

import "net/http"

// Forge is the top-level client for the Forgejo API.
//
// Usage:
//
//	ctx := context.Background()
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	repo, err := f.Repos.Get(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
type Forge struct {
	client *Client

	Repos         *RepoService
	Issues        *IssueService
	Pulls         *PullService
	Orgs          *OrgService
	Users         *UserService
	Teams         *TeamService
	Admin         *AdminService
	Branches      *BranchService
	Releases      *ReleaseService
	Labels        *LabelService
	Webhooks      *WebhookService
	Notifications *NotificationService
	Packages      *PackageService
	Actions       *ActionsService
	Contents      *ContentService
	Wiki          *WikiService
	Misc          *MiscService
	Commits       *CommitService
	Milestones    *MilestoneService
	ActivityPub   *ActivityPubService
}

// NewForge creates a new Forge client.
//
// Usage:
//
//	ctx := context.Background()
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	repos, err := f.Repos.ListOrgRepos(ctx, "core")
func NewForge(url, token string, opts ...Option) *Forge {
	c := NewClient(url, token, opts...)
	f := &Forge{client: c}
	f.Repos = newRepoService(c)
	f.Issues = newIssueService(c)
	f.Pulls = newPullService(c)
	f.Orgs = newOrgService(c)
	f.Users = newUserService(c)
	f.Teams = newTeamService(c)
	f.Admin = newAdminService(c)
	f.Branches = newBranchService(c)
	f.Releases = newReleaseService(c)
	f.Labels = newLabelService(c)
	f.Webhooks = newWebhookService(c)
	f.Notifications = newNotificationService(c)
	f.Packages = newPackageService(c)
	f.Actions = newActionsService(c)
	f.Contents = newContentService(c)
	f.Wiki = newWikiService(c)
	f.Misc = newMiscService(c)
	f.Commits = newCommitService(c)
	f.Milestones = newMilestoneService(c)
	f.ActivityPub = newActivityPubService(c)
	return f
}

// Client returns the underlying Forge client.
//
// Usage:
//
//	client := f.Client()
func (f *Forge) Client() *Client {
	if f == nil {
		return nil
	}
	return f.client
}

// BaseURL returns the configured Forgejo base URL.
//
// Usage:
//
//	baseURL := f.BaseURL()
func (f *Forge) BaseURL() string {
	if f == nil || f.client == nil {
		return ""
	}
	return f.client.BaseURL()
}

// RateLimit returns the last known rate limit information.
//
// Usage:
//
//	rl := f.RateLimit()
func (f *Forge) RateLimit() RateLimit {
	if f == nil || f.client == nil {
		return RateLimit{}
	}
	return f.client.RateLimit()
}

// UserAgent returns the configured User-Agent header value.
//
// Usage:
//
//	ua := f.UserAgent()
func (f *Forge) UserAgent() string {
	if f == nil || f.client == nil {
		return ""
	}
	return f.client.UserAgent()
}

// HTTPClient returns the configured underlying HTTP client.
//
// Usage:
//
//	hc := f.HTTPClient()
func (f *Forge) HTTPClient() *http.Client {
	if f == nil || f.client == nil {
		return nil
	}
	return f.client.HTTPClient()
}

// HasToken reports whether the Forge client was configured with an API token.
//
// Usage:
//
//	if f.HasToken() {
//	    _ = "authenticated"
//	}
func (f *Forge) HasToken() bool {
	if f == nil || f.client == nil {
		return false
	}
	return f.client.HasToken()
}

// String returns a safe summary of the Forge client.
//
// Usage:
//
//	s := f.String()
func (f *Forge) String() string {
	if f == nil {
		return "forge.Forge{<nil>}"
	}
	if f.client == nil {
		return "forge.Forge{client=<nil>}"
	}
	return "forge.Forge{client=" + f.client.String() + "}"
}

// GoString returns a safe Go-syntax summary of the Forge client.
//
// Usage:
//
//	s := fmt.Sprintf("%#v", f)
func (f *Forge) GoString() string { return f.String() }
