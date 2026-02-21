package forge

// Forge is the top-level client for the Forgejo API.
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
}

// NewForge creates a new Forge client.
func NewForge(url, token string, opts ...Option) *Forge {
	c := NewClient(url, token, opts...)
	f := &Forge{client: c}
	f.Repos = newRepoService(c)
	f.Issues = newIssueService(c)
	f.Pulls = newPullService(c)
	// Other services initialised in their respective tasks.
	// Stub them here so tests compile:
	f.Orgs = &OrgService{}
	f.Users = &UserService{}
	f.Teams = &TeamService{}
	f.Admin = &AdminService{}
	f.Branches = &BranchService{}
	f.Releases = &ReleaseService{}
	f.Labels = &LabelService{}
	f.Webhooks = &WebhookService{}
	f.Notifications = &NotificationService{}
	f.Packages = &PackageService{}
	f.Actions = &ActionsService{}
	f.Contents = &ContentService{}
	f.Wiki = &WikiService{}
	f.Misc = &MiscService{}
	return f
}

// Client returns the underlying HTTP client.
func (f *Forge) Client() *Client { return f.client }
