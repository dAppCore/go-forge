package forge

// Forge is the top-level client for the Forgejo API.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_ = f.Repos
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
}

// NewForge creates a new Forge client.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_ = f
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
	return f
}

// Client returns the underlying HTTP client.
func (f *Forge) Client() *Client { return f.client }
