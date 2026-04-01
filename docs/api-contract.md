# API Contract Inventory

`CODEX.md` was not present in `/workspace`, so this inventory follows the repository’s existing Go/doc/test conventions.

Coverage notes: rows list direct tests when a symbol is named in test names or referenced explicitly in test code. Services that embed `Resource[...]` are documented by their declared methods; promoted CRUD methods are covered under the `Resource` rows instead of being duplicated for every service.

## `forge`

| Kind | Name | Signature | Description | Test Coverage |
| --- | --- | --- | --- | --- |
| type | APIError | `type APIError struct` | APIError represents an error response from the Forgejo API. | `TestAPIError_Good_Error`, `TestClient_Bad_ServerError`, `TestIsConflict_Bad_NotConflict` (+2 more) |
| type | ActionsService | `type ActionsService struct` | ActionsService handles CI/CD actions operations across repositories and organisations — secrets, variables, and workflow dispatches. No Resource embedding — heterogeneous endpoints across repo and org levels. | `TestActionsService_Bad_NotFound`, `TestActionsService_Good_CreateRepoSecret`, `TestActionsService_Good_CreateRepoVariable` (+7 more) |
| type | AdminService | `type AdminService struct` | AdminService handles site administration operations. Unlike other services, AdminService does not embed Resource[T,C,U] because admin endpoints are heterogeneous. | `TestAdminService_Bad_CreateUser_Forbidden`, `TestAdminService_Bad_DeleteUser_NotFound`, `TestAdminService_Good_AdoptRepo` (+9 more) |
| type | BranchService | `type BranchService struct` | BranchService handles branch operations within a repository. | `TestBranchService_Good_CreateProtection`, `TestBranchService_Good_Get`, `TestBranchService_Good_List` |
| type | Client | `type Client struct` | Client is a low-level HTTP client for the Forgejo API. | `TestClient_Bad_Conflict`, `TestClient_Bad_Forbidden`, `TestClient_Bad_NotFound` (+9 more) |
| type | CommitService | `type CommitService struct` | CommitService handles commit-related operations such as commit statuses and git notes. No Resource embedding — collection and item commit paths differ, and the remaining endpoints are heterogeneous across status and note paths. | `TestCommitService_Bad_NotFound`, `TestCommitService_Good_CreateStatus`, `TestCommitService_Good_Get` (+4 more) |
| type | ContentService | `type ContentService struct` | ContentService handles file read/write operations via the Forgejo API. No Resource embedding — paths vary by operation. | `TestContentService_Bad_GetRawNotFound`, `TestContentService_Bad_NotFound`, `TestContentService_Good_CreateFile` (+4 more) |
| type | Forge | `type Forge struct` | Forge is the top-level client for the Forgejo API. | `TestForge_Good_Client`, `TestForge_Good_NewForge` |
| type | IssueService | `type IssueService struct` | IssueService handles issue operations within a repository. | `TestIssueService_Bad_List`, `TestIssueService_Good_Create`, `TestIssueService_Good_CreateComment` (+6 more) |
| type | LabelService | `type LabelService struct` | LabelService handles repository labels, organisation labels, and issue labels. No Resource embedding — paths are heterogeneous. | `TestLabelService_Bad_NotFound`, `TestLabelService_Good_CreateOrgLabel`, `TestLabelService_Good_CreateRepoLabel` (+5 more) |
| type | ListOptions | `type ListOptions struct` | ListOptions controls pagination. | `TestListPage_Good_QueryParams` |
| type | MilestoneService | `type MilestoneService struct` | MilestoneService handles repository milestones. | No direct tests. |
| type | MiscService | `type MiscService struct` | MiscService handles miscellaneous Forgejo API endpoints such as markdown rendering, licence templates, gitignore templates, and server metadata. No Resource embedding — heterogeneous read-only endpoints. | `TestMiscService_Bad_NotFound`, `TestMiscService_Good_GetGitignoreTemplate`, `TestMiscService_Good_GetLicense` (+5 more) |
| type | NotificationService | `type NotificationService struct` | NotificationService handles notification operations via the Forgejo API. No Resource embedding — varied endpoint shapes. | `TestNotificationService_Bad_NotFound`, `TestNotificationService_Good_GetThread`, `TestNotificationService_Good_List` (+3 more) |
| type | Option | `type Option func(*Client)` | Option configures the Client. | No direct tests. |
| type | OrgService | `type OrgService struct` | OrgService handles organisation operations. | `TestOrgService_Good_Get`, `TestOrgService_Good_List`, `TestOrgService_Good_ListMembers` |
| type | PackageService | `type PackageService struct` | PackageService handles package registry operations via the Forgejo API. No Resource embedding — paths vary by operation. | `TestPackageService_Bad_NotFound`, `TestPackageService_Good_Delete`, `TestPackageService_Good_Get` (+2 more) |
| type | PagedResult | `type PagedResult[T any] struct` | PagedResult holds a single page of results with metadata. | No direct tests. |
| type | Params | `type Params map[string]string` | Params maps path variable names to values. Example: Params{"owner": "core", "repo": "go-forge"} | `TestBranchService_Good_Get`, `TestBranchService_Good_List`, `TestCommitService_Good_Get` (+32 more) |
| type | PullService | `type PullService struct` | PullService handles pull request operations within a repository. | `TestPullService_Bad_Merge`, `TestPullService_Good_Create`, `TestPullService_Good_Get` (+2 more) |
| type | RateLimit | `type RateLimit struct` | RateLimit represents the rate limit information from the Forgejo API. | `TestClient_Good_RateLimit` |
| type | ReleaseService | `type ReleaseService struct` | ReleaseService handles release operations within a repository. | `TestReleaseService_Good_Get`, `TestReleaseService_Good_GetByTag`, `TestReleaseService_Good_List` |
| type | RepoService | `type RepoService struct` | RepoService handles repository operations. | `TestRepoService_Bad_Get`, `TestRepoService_Good_Delete`, `TestRepoService_Good_Fork` (+3 more) |
| type | Resource | `type Resource[T any, C any, U any] struct` | Resource provides generic CRUD operations for a Forgejo API resource. T is the resource type, C is the create options type, U is the update options type. | `TestResource_Bad_IterError`, `TestResource_Good_Create`, `TestResource_Good_Delete` (+6 more) |
| type | TeamService | `type TeamService struct` | TeamService handles team operations. | `TestTeamService_Good_AddMember`, `TestTeamService_Good_Get`, `TestTeamService_Good_ListMembers` |
| type | UserService | `type UserService struct` | UserService handles user operations. | `TestUserService_Good_Get`, `TestUserService_Good_GetCurrent`, `TestUserService_Good_ListFollowers` |
| type | WebhookService | `type WebhookService struct` | WebhookService handles webhook (hook) operations within a repository. Embeds Resource for standard CRUD on /api/v1/repos/{owner}/{repo}/hooks/{id}. | `TestWebhookService_Bad_NotFound`, `TestWebhookService_Good_Create`, `TestWebhookService_Good_Get` (+3 more) |
| type | WikiService | `type WikiService struct` | WikiService handles wiki page operations for a repository. No Resource embedding — custom endpoints for wiki CRUD. | `TestWikiService_Bad_NotFound`, `TestWikiService_Good_CreatePage`, `TestWikiService_Good_DeletePage` (+3 more) |
| function | IsConflict | `func IsConflict(err error) bool` | IsConflict returns true if the error is a 409 response. | `TestClient_Bad_Conflict`, `TestIsConflict_Bad_NotConflict`, `TestIsConflict_Good` (+1 more) |
| function | IsForbidden | `func IsForbidden(err error) bool` | IsForbidden returns true if the error is a 403 response. | `TestAdminService_Bad_CreateUser_Forbidden`, `TestClient_Bad_Forbidden`, `TestIsForbidden_Bad_NotForbidden` |
| function | IsNotFound | `func IsNotFound(err error) bool` | IsNotFound returns true if the error is a 404 response. | `TestActionsService_Bad_NotFound`, `TestAdminService_Bad_DeleteUser_NotFound`, `TestClient_Bad_NotFound` (+10 more) |
| function | ListAll | `func ListAll[T any](ctx context.Context, c *Client, path string, query map[string]string) ([]T, error)` | ListAll fetches all pages of results. | `TestOrgService_Good_List`, `TestPagination_Bad_ServerError`, `TestPagination_Good_EmptyResult` (+3 more) |
| function | ListIter | `func ListIter[T any](ctx context.Context, c *Client, path string, query map[string]string) iter.Seq2[T, error]` | ListIter returns an iterator over all resources across all pages. | `TestPagination_Good_Iter` |
| function | ListPage | `func ListPage[T any](ctx context.Context, c *Client, path string, query map[string]string, opts ListOptions) (*PagedResult[T], error)` | ListPage fetches a single page of results. Extra query params can be passed via the query map. | `TestListPage_Good_QueryParams` |
| function | NewClient | `func NewClient(url, token string, opts ...Option) *Client` | NewClient creates a new Forgejo API client. | `TestClient_Bad_Conflict`, `TestClient_Bad_Forbidden`, `TestClient_Bad_NotFound` (+23 more) |
| function | NewForge | `func NewForge(url, token string, opts ...Option) *Forge` | NewForge creates a new Forge client. | `TestActionsService_Bad_NotFound`, `TestActionsService_Good_CreateRepoSecret`, `TestActionsService_Good_CreateRepoVariable` (+109 more) |
| function | NewForgeFromConfig | `func NewForgeFromConfig(flagURL, flagToken string, opts ...Option) (*Forge, error)` | NewForgeFromConfig creates a new Forge client using resolved configuration. It returns an error if no API token is available from flags or environment. | `TestNewForgeFromConfig_Bad_NoToken` |
| function | NewResource | `func NewResource[T any, C any, U any](c *Client, path string) *Resource[T, C, U]` | NewResource creates a new Resource for the given path pattern. The path should be the item path (e.g., /repos/{owner}/{repo}/issues/{index}). The collection path is derived by stripping the last /{placeholder} segment. | `TestResource_Bad_IterError`, `TestResource_Good_Create`, `TestResource_Good_Delete` (+6 more) |
| function | ResolveConfig | `func ResolveConfig(flagURL, flagToken string) (url, token string, err error)` | ResolveConfig resolves the Forgejo URL and API token from flags, environment variables, and built-in defaults. Priority order: flags > env > defaults. Environment variables: - FORGE_URL — base URL of the Forgejo instance - FORGE_TOKEN — API token for authentication | `TestResolveConfig_Good_DefaultURL`, `TestResolveConfig_Good_EnvOverrides`, `TestResolveConfig_Good_FlagOverridesEnv` |
| function | ResolvePath | `func ResolvePath(path string, params Params) string` | ResolvePath substitutes {placeholders} in path with values from params. | `TestResolvePath_Good_NoParams`, `TestResolvePath_Good_Simple`, `TestResolvePath_Good_URLEncoding` (+1 more) |
| function | WithHTTPClient | `func WithHTTPClient(hc *http.Client) Option` | WithHTTPClient sets a custom http.Client. | `TestClient_Good_WithHTTPClient` |
| function | WithUserAgent | `func WithUserAgent(ua string) Option` | WithUserAgent sets the User-Agent header. | `TestClient_Good_Options` |
| method | APIError.Error | `func (e *APIError) Error() string` | No doc comment. | `TestAPIError_Good_Error` |
| method | ActionsService.CreateRepoSecret | `func (s *ActionsService) CreateRepoSecret(ctx context.Context, owner, repo, name string, data string) error` | CreateRepoSecret creates or updates a secret in a repository. Forgejo expects a PUT with {"data": "secret-value"} body. | `TestActionsService_Good_CreateRepoSecret` |
| method | ActionsService.CreateRepoVariable | `func (s *ActionsService) CreateRepoVariable(ctx context.Context, owner, repo, name, value string) error` | CreateRepoVariable creates a new action variable in a repository. Forgejo expects a POST with {"value": "var-value"} body. | `TestActionsService_Good_CreateRepoVariable` |
| method | ActionsService.DeleteRepoSecret | `func (s *ActionsService) DeleteRepoSecret(ctx context.Context, owner, repo, name string) error` | DeleteRepoSecret removes a secret from a repository. | `TestActionsService_Good_DeleteRepoSecret` |
| method | ActionsService.DeleteRepoVariable | `func (s *ActionsService) DeleteRepoVariable(ctx context.Context, owner, repo, name string) error` | DeleteRepoVariable removes an action variable from a repository. | `TestActionsService_Good_DeleteRepoVariable` |
| method | ActionsService.DispatchWorkflow | `func (s *ActionsService) DispatchWorkflow(ctx context.Context, owner, repo, workflow string, opts map[string]any) error` | DispatchWorkflow triggers a workflow run. | `TestActionsService_Good_DispatchWorkflow` |
| method | ActionsService.IterOrgSecrets | `func (s *ActionsService) IterOrgSecrets(ctx context.Context, org string) iter.Seq2[types.Secret, error]` | IterOrgSecrets returns an iterator over all secrets for an organisation. | No direct tests. |
| method | ActionsService.IterOrgVariables | `func (s *ActionsService) IterOrgVariables(ctx context.Context, org string) iter.Seq2[types.ActionVariable, error]` | IterOrgVariables returns an iterator over all action variables for an organisation. | No direct tests. |
| method | ActionsService.IterRepoSecrets | `func (s *ActionsService) IterRepoSecrets(ctx context.Context, owner, repo string) iter.Seq2[types.Secret, error]` | IterRepoSecrets returns an iterator over all secrets for a repository. | No direct tests. |
| method | ActionsService.IterRepoVariables | `func (s *ActionsService) IterRepoVariables(ctx context.Context, owner, repo string) iter.Seq2[types.ActionVariable, error]` | IterRepoVariables returns an iterator over all action variables for a repository. | No direct tests. |
| method | ActionsService.ListOrgSecrets | `func (s *ActionsService) ListOrgSecrets(ctx context.Context, org string) ([]types.Secret, error)` | ListOrgSecrets returns all secrets for an organisation. | `TestActionsService_Good_ListOrgSecrets` |
| method | ActionsService.ListOrgVariables | `func (s *ActionsService) ListOrgVariables(ctx context.Context, org string) ([]types.ActionVariable, error)` | ListOrgVariables returns all action variables for an organisation. | `TestActionsService_Good_ListOrgVariables` |
| method | ActionsService.ListRepoSecrets | `func (s *ActionsService) ListRepoSecrets(ctx context.Context, owner, repo string) ([]types.Secret, error)` | ListRepoSecrets returns all secrets for a repository. | `TestActionsService_Bad_NotFound`, `TestActionsService_Good_ListRepoSecrets` |
| method | ActionsService.ListRepoVariables | `func (s *ActionsService) ListRepoVariables(ctx context.Context, owner, repo string) ([]types.ActionVariable, error)` | ListRepoVariables returns all action variables for a repository. | `TestActionsService_Good_ListRepoVariables` |
| method | AdminService.AdoptRepo | `func (s *AdminService) AdoptRepo(ctx context.Context, owner, repo string) error` | AdoptRepo adopts an unadopted repository (admin only). | `TestAdminService_Good_AdoptRepo` |
| method | AdminService.CreateUser | `func (s *AdminService) CreateUser(ctx context.Context, opts *types.CreateUserOption) (*types.User, error)` | CreateUser creates a new user (admin only). | `TestAdminService_Bad_CreateUser_Forbidden`, `TestAdminService_Good_CreateUser` |
| method | AdminService.DeleteUser | `func (s *AdminService) DeleteUser(ctx context.Context, username string) error` | DeleteUser deletes a user (admin only). | `TestAdminService_Bad_DeleteUser_NotFound`, `TestAdminService_Good_DeleteUser` |
| method | AdminService.EditUser | `func (s *AdminService) EditUser(ctx context.Context, username string, opts map[string]any) error` | EditUser edits an existing user (admin only). | `TestAdminService_Good_EditUser` |
| method | AdminService.GenerateRunnerToken | `func (s *AdminService) GenerateRunnerToken(ctx context.Context) (string, error)` | GenerateRunnerToken generates an actions runner registration token. | `TestAdminService_Good_GenerateRunnerToken` |
| method | AdminService.IterCron | `func (s *AdminService) IterCron(ctx context.Context) iter.Seq2[types.Cron, error]` | IterCron returns an iterator over all cron tasks (admin only). | No direct tests. |
| method | AdminService.IterOrgs | `func (s *AdminService) IterOrgs(ctx context.Context) iter.Seq2[types.Organization, error]` | IterOrgs returns an iterator over all organisations (admin only). | No direct tests. |
| method | AdminService.IterUsers | `func (s *AdminService) IterUsers(ctx context.Context) iter.Seq2[types.User, error]` | IterUsers returns an iterator over all users (admin only). | No direct tests. |
| method | AdminService.ListCron | `func (s *AdminService) ListCron(ctx context.Context) ([]types.Cron, error)` | ListCron returns all cron tasks (admin only). | `TestAdminService_Good_ListCron` |
| method | AdminService.ListOrgs | `func (s *AdminService) ListOrgs(ctx context.Context) ([]types.Organization, error)` | ListOrgs returns all organisations (admin only). | `TestAdminService_Good_ListOrgs` |
| method | AdminService.ListUsers | `func (s *AdminService) ListUsers(ctx context.Context) ([]types.User, error)` | ListUsers returns all users (admin only). | `TestAdminService_Good_ListUsers` |
| method | AdminService.RenameUser | `func (s *AdminService) RenameUser(ctx context.Context, username, newName string) error` | RenameUser renames a user (admin only). | `TestAdminService_Good_RenameUser` |
| method | AdminService.RunCron | `func (s *AdminService) RunCron(ctx context.Context, task string) error` | RunCron runs a cron task by name (admin only). | `TestAdminService_Good_RunCron` |
| method | BranchService.CreateBranchProtection | `func (s *BranchService) CreateBranchProtection(ctx context.Context, owner, repo string, opts *types.CreateBranchProtectionOption) (*types.BranchProtection, error)` | CreateBranchProtection creates a new branch protection rule. | `TestBranchService_Good_CreateProtection` |
| method | BranchService.DeleteBranchProtection | `func (s *BranchService) DeleteBranchProtection(ctx context.Context, owner, repo, name string) error` | DeleteBranchProtection deletes a branch protection rule. | No direct tests. |
| method | BranchService.EditBranchProtection | `func (s *BranchService) EditBranchProtection(ctx context.Context, owner, repo, name string, opts *types.EditBranchProtectionOption) (*types.BranchProtection, error)` | EditBranchProtection updates an existing branch protection rule. | No direct tests. |
| method | BranchService.GetBranchProtection | `func (s *BranchService) GetBranchProtection(ctx context.Context, owner, repo, name string) (*types.BranchProtection, error)` | GetBranchProtection returns a single branch protection by name. | No direct tests. |
| method | BranchService.IterBranchProtections | `func (s *BranchService) IterBranchProtections(ctx context.Context, owner, repo string) iter.Seq2[types.BranchProtection, error]` | IterBranchProtections returns an iterator over all branch protections for a repository. | No direct tests. |
| method | BranchService.ListBranchProtections | `func (s *BranchService) ListBranchProtections(ctx context.Context, owner, repo string) ([]types.BranchProtection, error)` | ListBranchProtections returns all branch protections for a repository. | No direct tests. |
| method | Client.Delete | `func (c *Client) Delete(ctx context.Context, path string) error` | Delete performs a DELETE request. | `TestClient_Good_Delete` |
| method | Client.DeleteWithBody | `func (c *Client) DeleteWithBody(ctx context.Context, path string, body any) error` | DeleteWithBody performs a DELETE request with a JSON body. | No direct tests. |
| method | Client.Get | `func (c *Client) Get(ctx context.Context, path string, out any) error` | Get performs a GET request. | `TestClient_Bad_Forbidden`, `TestClient_Bad_NotFound`, `TestClient_Bad_ServerError` (+3 more) |
| method | Client.GetRaw | `func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error)` | GetRaw performs a GET request and returns the raw response body as bytes instead of JSON-decoding. Useful for endpoints that return raw file content. | No direct tests. |
| method | Client.Patch | `func (c *Client) Patch(ctx context.Context, path string, body, out any) error` | Patch performs a PATCH request. | No direct tests. |
| method | Client.Post | `func (c *Client) Post(ctx context.Context, path string, body, out any) error` | Post performs a POST request. | `TestClient_Bad_Conflict`, `TestClient_Good_Post` |
| method | Client.PostRaw | `func (c *Client) PostRaw(ctx context.Context, path string, body any) ([]byte, error)` | PostRaw performs a POST request with a JSON body and returns the raw response body as bytes instead of JSON-decoding. Useful for endpoints such as /markdown that return raw HTML text. | No direct tests. |
| method | Client.Put | `func (c *Client) Put(ctx context.Context, path string, body, out any) error` | Put performs a PUT request. | No direct tests. |
| method | Client.RateLimit | `func (c *Client) RateLimit() RateLimit` | RateLimit returns the last known rate limit information. | `TestClient_Good_RateLimit` |
| method | CommitService.CreateStatus | `func (s *CommitService) CreateStatus(ctx context.Context, owner, repo, sha string, opts *types.CreateStatusOption) (*types.CommitStatus, error)` | CreateStatus creates a new commit status for the given SHA. | `TestCommitService_Good_CreateStatus` |
| method | CommitService.Get | `func (s *CommitService) Get(ctx context.Context, params Params) (*types.Commit, error)` | Get returns a single commit by SHA or ref. | `TestCommitService_Good_Get`, `TestCommitService_Good_List` |
| method | CommitService.GetCombinedStatus | `func (s *CommitService) GetCombinedStatus(ctx context.Context, owner, repo, ref string) (*types.CombinedStatus, error)` | GetCombinedStatus returns the combined status for a given ref (branch, tag, or SHA). | `TestCommitService_Good_GetCombinedStatus` |
| method | CommitService.GetNote | `func (s *CommitService) GetNote(ctx context.Context, owner, repo, sha string) (*types.Note, error)` | GetNote returns the git note for a given commit SHA. | `TestCommitService_Bad_NotFound`, `TestCommitService_Good_GetNote` |
| method | CommitService.Iter | `func (s *CommitService) Iter(ctx context.Context, params Params) iter.Seq2[types.Commit, error]` | Iter returns an iterator over all commits for a repository. | No direct tests. |
| method | CommitService.List | `func (s *CommitService) List(ctx context.Context, params Params, opts ListOptions) (*PagedResult[types.Commit], error)` | List returns a single page of commits for a repository. | `TestCommitService_Good_List` |
| method | CommitService.ListAll | `func (s *CommitService) ListAll(ctx context.Context, params Params) ([]types.Commit, error)` | ListAll returns all commits for a repository. | No direct tests. |
| method | CommitService.ListStatuses | `func (s *CommitService) ListStatuses(ctx context.Context, owner, repo, ref string) ([]types.CommitStatus, error)` | ListStatuses returns all commit statuses for a given ref. | `TestCommitService_Good_ListStatuses` |
| method | ContentService.CreateFile | `func (s *ContentService) CreateFile(ctx context.Context, owner, repo, filepath string, opts *types.CreateFileOptions) (*types.FileResponse, error)` | CreateFile creates a new file in a repository. | `TestContentService_Good_CreateFile` |
| method | ContentService.DeleteFile | `func (s *ContentService) DeleteFile(ctx context.Context, owner, repo, filepath string, opts *types.DeleteFileOptions) error` | DeleteFile deletes a file from a repository. Uses DELETE with a JSON body. | `TestContentService_Good_DeleteFile` |
| method | ContentService.GetFile | `func (s *ContentService) GetFile(ctx context.Context, owner, repo, filepath string) (*types.ContentsResponse, error)` | GetFile returns metadata and content for a file in a repository. | `TestContentService_Bad_NotFound`, `TestContentService_Good_GetFile` |
| method | ContentService.GetRawFile | `func (s *ContentService) GetRawFile(ctx context.Context, owner, repo, filepath string) ([]byte, error)` | GetRawFile returns the raw file content as bytes. | `TestContentService_Bad_GetRawNotFound`, `TestContentService_Good_GetRawFile` |
| method | ContentService.UpdateFile | `func (s *ContentService) UpdateFile(ctx context.Context, owner, repo, filepath string, opts *types.UpdateFileOptions) (*types.FileResponse, error)` | UpdateFile updates an existing file in a repository. | `TestContentService_Good_UpdateFile` |
| method | Forge.Client | `func (f *Forge) Client() *Client` | Client returns the underlying HTTP client. | `TestForge_Good_Client` |
| method | IssueService.AddLabels | `func (s *IssueService) AddLabels(ctx context.Context, owner, repo string, index int64, labelIDs []int64) error` | AddLabels adds labels to an issue. | No direct tests. |
| method | IssueService.AddReaction | `func (s *IssueService) AddReaction(ctx context.Context, owner, repo string, index int64, reaction string) error` | AddReaction adds a reaction to an issue. | No direct tests. |
| method | IssueService.CreateComment | `func (s *IssueService) CreateComment(ctx context.Context, owner, repo string, index int64, body string) (*types.Comment, error)` | CreateComment creates a comment on an issue. | `TestIssueService_Good_CreateComment` |
| method | IssueService.DeleteReaction | `func (s *IssueService) DeleteReaction(ctx context.Context, owner, repo string, index int64, reaction string) error` | DeleteReaction removes a reaction from an issue. | No direct tests. |
| method | IssueService.IterComments | `func (s *IssueService) IterComments(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.Comment, error]` | IterComments returns an iterator over all comments on an issue. | No direct tests. |
| method | IssueService.ListComments | `func (s *IssueService) ListComments(ctx context.Context, owner, repo string, index int64) ([]types.Comment, error)` | ListComments returns all comments on an issue. | No direct tests. |
| method | IssueService.Pin | `func (s *IssueService) Pin(ctx context.Context, owner, repo string, index int64) error` | Pin pins an issue. | `TestIssueService_Good_Pin` |
| method | IssueService.RemoveLabel | `func (s *IssueService) RemoveLabel(ctx context.Context, owner, repo string, index int64, labelID int64) error` | RemoveLabel removes a single label from an issue. | No direct tests. |
| method | IssueService.SetDeadline | `func (s *IssueService) SetDeadline(ctx context.Context, owner, repo string, index int64, deadline string) error` | SetDeadline sets or updates the deadline on an issue. | No direct tests. |
| method | IssueService.StartStopwatch | `func (s *IssueService) StartStopwatch(ctx context.Context, owner, repo string, index int64) error` | StartStopwatch starts the stopwatch on an issue. | No direct tests. |
| method | IssueService.StopStopwatch | `func (s *IssueService) StopStopwatch(ctx context.Context, owner, repo string, index int64) error` | StopStopwatch stops the stopwatch on an issue. | No direct tests. |
| method | IssueService.Unpin | `func (s *IssueService) Unpin(ctx context.Context, owner, repo string, index int64) error` | Unpin unpins an issue. | No direct tests. |
| method | LabelService.CreateOrgLabel | `func (s *LabelService) CreateOrgLabel(ctx context.Context, org string, opts *types.CreateLabelOption) (*types.Label, error)` | CreateOrgLabel creates a new label in an organisation. | `TestLabelService_Good_CreateOrgLabel` |
| method | LabelService.CreateRepoLabel | `func (s *LabelService) CreateRepoLabel(ctx context.Context, owner, repo string, opts *types.CreateLabelOption) (*types.Label, error)` | CreateRepoLabel creates a new label in a repository. | `TestLabelService_Good_CreateRepoLabel` |
| method | LabelService.DeleteRepoLabel | `func (s *LabelService) DeleteRepoLabel(ctx context.Context, owner, repo string, id int64) error` | DeleteRepoLabel deletes a label from a repository. | `TestLabelService_Good_DeleteRepoLabel` |
| method | LabelService.EditRepoLabel | `func (s *LabelService) EditRepoLabel(ctx context.Context, owner, repo string, id int64, opts *types.EditLabelOption) (*types.Label, error)` | EditRepoLabel updates an existing label in a repository. | `TestLabelService_Good_EditRepoLabel` |
| method | LabelService.GetRepoLabel | `func (s *LabelService) GetRepoLabel(ctx context.Context, owner, repo string, id int64) (*types.Label, error)` | GetRepoLabel returns a single label by ID. | `TestLabelService_Bad_NotFound`, `TestLabelService_Good_GetRepoLabel` |
| method | LabelService.IterOrgLabels | `func (s *LabelService) IterOrgLabels(ctx context.Context, org string) iter.Seq2[types.Label, error]` | IterOrgLabels returns an iterator over all labels for an organisation. | No direct tests. |
| method | LabelService.IterRepoLabels | `func (s *LabelService) IterRepoLabels(ctx context.Context, owner, repo string) iter.Seq2[types.Label, error]` | IterRepoLabels returns an iterator over all labels for a repository. | No direct tests. |
| method | LabelService.ListOrgLabels | `func (s *LabelService) ListOrgLabels(ctx context.Context, org string) ([]types.Label, error)` | ListOrgLabels returns all labels for an organisation. | `TestLabelService_Good_ListOrgLabels` |
| method | LabelService.ListRepoLabels | `func (s *LabelService) ListRepoLabels(ctx context.Context, owner, repo string) ([]types.Label, error)` | ListRepoLabels returns all labels for a repository. | `TestLabelService_Good_ListRepoLabels` |
| method | MilestoneService.Create | `func (s *MilestoneService) Create(ctx context.Context, owner, repo string, opts *types.CreateMilestoneOption) (*types.Milestone, error)` | Create creates a new milestone. | No direct tests. |
| method | MilestoneService.Get | `func (s *MilestoneService) Get(ctx context.Context, owner, repo string, id int64) (*types.Milestone, error)` | Get returns a single milestone by ID. | No direct tests. |
| method | MilestoneService.ListAll | `func (s *MilestoneService) ListAll(ctx context.Context, params Params) ([]types.Milestone, error)` | ListAll returns all milestones for a repository. | No direct tests. |
| method | MiscService.GetGitignoreTemplate | `func (s *MiscService) GetGitignoreTemplate(ctx context.Context, name string) (*types.GitignoreTemplateInfo, error)` | GetGitignoreTemplate returns a single gitignore template by name. | `TestMiscService_Good_GetGitignoreTemplate` |
| method | MiscService.GetLicense | `func (s *MiscService) GetLicense(ctx context.Context, name string) (*types.LicenseTemplateInfo, error)` | GetLicense returns a single licence template by name. | `TestMiscService_Bad_NotFound`, `TestMiscService_Good_GetLicense` |
| method | MiscService.GetNodeInfo | `func (s *MiscService) GetNodeInfo(ctx context.Context) (*types.NodeInfo, error)` | GetNodeInfo returns the NodeInfo metadata for the Forgejo instance. | `TestMiscService_Good_GetNodeInfo` |
| method | MiscService.GetVersion | `func (s *MiscService) GetVersion(ctx context.Context) (*types.ServerVersion, error)` | GetVersion returns the server version. | `TestMiscService_Good_GetVersion` |
| method | MiscService.ListGitignoreTemplates | `func (s *MiscService) ListGitignoreTemplates(ctx context.Context) ([]string, error)` | ListGitignoreTemplates returns all available gitignore template names. | `TestMiscService_Good_ListGitignoreTemplates` |
| method | MiscService.ListLicenses | `func (s *MiscService) ListLicenses(ctx context.Context) ([]types.LicensesTemplateListEntry, error)` | ListLicenses returns all available licence templates. | `TestMiscService_Good_ListLicenses` |
| method | MiscService.RenderMarkdown | `func (s *MiscService) RenderMarkdown(ctx context.Context, text, mode string) (string, error)` | RenderMarkdown renders markdown text to HTML. The response is raw HTML text, not JSON. | `TestMiscService_Good_RenderMarkdown` |
| method | MiscService.RenderMarkdownRaw | `func (s *MiscService) RenderMarkdownRaw(ctx context.Context, text string) (string, error)` | RenderMarkdownRaw renders raw markdown text to HTML. The request body is sent as text/plain and the response is raw HTML text, not JSON. | `TestMiscService_RenderMarkdownRaw_Good` |
| method | NotificationService.GetThread | `func (s *NotificationService) GetThread(ctx context.Context, id int64) (*types.NotificationThread, error)` | GetThread returns a single notification thread by ID. | `TestNotificationService_Bad_NotFound`, `TestNotificationService_Good_GetThread` |
| method | NotificationService.Iter | `func (s *NotificationService) Iter(ctx context.Context) iter.Seq2[types.NotificationThread, error]` | Iter returns an iterator over all notifications for the authenticated user. | No direct tests. |
| method | NotificationService.IterRepo | `func (s *NotificationService) IterRepo(ctx context.Context, owner, repo string) iter.Seq2[types.NotificationThread, error]` | IterRepo returns an iterator over all notifications for a specific repository. | No direct tests. |
| method | NotificationService.List | `func (s *NotificationService) List(ctx context.Context) ([]types.NotificationThread, error)` | List returns all notifications for the authenticated user. | `TestNotificationService_Good_List` |
| method | NotificationService.ListRepo | `func (s *NotificationService) ListRepo(ctx context.Context, owner, repo string) ([]types.NotificationThread, error)` | ListRepo returns all notifications for a specific repository. | `TestNotificationService_Good_ListRepo` |
| method | NotificationService.MarkRead | `func (s *NotificationService) MarkRead(ctx context.Context) error` | MarkRead marks all notifications as read. | `TestNotificationService_Good_MarkRead` |
| method | NotificationService.MarkThreadRead | `func (s *NotificationService) MarkThreadRead(ctx context.Context, id int64) error` | MarkThreadRead marks a single notification thread as read. | `TestNotificationService_Good_MarkThreadRead` |
| method | OrgService.AddMember | `func (s *OrgService) AddMember(ctx context.Context, org, username string) error` | AddMember adds a user to an organisation. | No direct tests. |
| method | OrgService.IterMembers | `func (s *OrgService) IterMembers(ctx context.Context, org string) iter.Seq2[types.User, error]` | IterMembers returns an iterator over all members of an organisation. | No direct tests. |
| method | OrgService.IterMyOrgs | `func (s *OrgService) IterMyOrgs(ctx context.Context) iter.Seq2[types.Organization, error]` | IterMyOrgs returns an iterator over all organisations for the authenticated user. | No direct tests. |
| method | OrgService.IterUserOrgs | `func (s *OrgService) IterUserOrgs(ctx context.Context, username string) iter.Seq2[types.Organization, error]` | IterUserOrgs returns an iterator over all organisations for a user. | No direct tests. |
| method | OrgService.ListMembers | `func (s *OrgService) ListMembers(ctx context.Context, org string) ([]types.User, error)` | ListMembers returns all members of an organisation. | `TestOrgService_Good_ListMembers` |
| method | OrgService.ListMyOrgs | `func (s *OrgService) ListMyOrgs(ctx context.Context) ([]types.Organization, error)` | ListMyOrgs returns all organisations for the authenticated user. | No direct tests. |
| method | OrgService.ListUserOrgs | `func (s *OrgService) ListUserOrgs(ctx context.Context, username string) ([]types.Organization, error)` | ListUserOrgs returns all organisations for a user. | No direct tests. |
| method | OrgService.RemoveMember | `func (s *OrgService) RemoveMember(ctx context.Context, org, username string) error` | RemoveMember removes a user from an organisation. | No direct tests. |
| method | PackageService.Delete | `func (s *PackageService) Delete(ctx context.Context, owner, pkgType, name, version string) error` | Delete removes a package by owner, type, name, and version. | `TestPackageService_Good_Delete` |
| method | PackageService.Get | `func (s *PackageService) Get(ctx context.Context, owner, pkgType, name, version string) (*types.Package, error)` | Get returns a single package by owner, type, name, and version. | `TestPackageService_Bad_NotFound`, `TestPackageService_Good_Get` |
| method | PackageService.Iter | `func (s *PackageService) Iter(ctx context.Context, owner string) iter.Seq2[types.Package, error]` | Iter returns an iterator over all packages for a given owner. | No direct tests. |
| method | PackageService.IterFiles | `func (s *PackageService) IterFiles(ctx context.Context, owner, pkgType, name, version string) iter.Seq2[types.PackageFile, error]` | IterFiles returns an iterator over all files for a specific package version. | No direct tests. |
| method | PackageService.List | `func (s *PackageService) List(ctx context.Context, owner string) ([]types.Package, error)` | List returns all packages for a given owner. | `TestPackageService_Good_List` |
| method | PackageService.ListFiles | `func (s *PackageService) ListFiles(ctx context.Context, owner, pkgType, name, version string) ([]types.PackageFile, error)` | ListFiles returns all files for a specific package version. | `TestPackageService_Good_ListFiles` |
| method | PullService.DismissReview | `func (s *PullService) DismissReview(ctx context.Context, owner, repo string, index, reviewID int64, msg string) error` | DismissReview dismisses a pull request review. | No direct tests. |
| method | PullService.IterReviews | `func (s *PullService) IterReviews(ctx context.Context, owner, repo string, index int64) iter.Seq2[types.PullReview, error]` | IterReviews returns an iterator over all reviews on a pull request. | No direct tests. |
| method | PullService.ListReviews | `func (s *PullService) ListReviews(ctx context.Context, owner, repo string, index int64) ([]types.PullReview, error)` | ListReviews returns all reviews on a pull request. | No direct tests. |
| method | PullService.Merge | `func (s *PullService) Merge(ctx context.Context, owner, repo string, index int64, method string) error` | Merge merges a pull request. Method is one of "merge", "rebase", "rebase-merge", "squash", "fast-forward-only", "manually-merged". | `TestPullService_Bad_Merge`, `TestPullService_Good_Merge` |
| method | PullService.SubmitReview | `func (s *PullService) SubmitReview(ctx context.Context, owner, repo string, index int64, review map[string]any) (*types.PullReview, error)` | SubmitReview creates a new review on a pull request. | No direct tests. |
| method | PullService.UndismissReview | `func (s *PullService) UndismissReview(ctx context.Context, owner, repo string, index, reviewID int64) error` | UndismissReview undismisses a pull request review. | No direct tests. |
| method | PullService.Update | `func (s *PullService) Update(ctx context.Context, owner, repo string, index int64) error` | Update updates a pull request branch with the base branch. | No direct tests. |
| method | ReleaseService.DeleteAsset | `func (s *ReleaseService) DeleteAsset(ctx context.Context, owner, repo string, releaseID, assetID int64) error` | DeleteAsset deletes a single asset from a release. | No direct tests. |
| method | ReleaseService.DeleteByTag | `func (s *ReleaseService) DeleteByTag(ctx context.Context, owner, repo, tag string) error` | DeleteByTag deletes a release by its tag name. | No direct tests. |
| method | ReleaseService.GetAsset | `func (s *ReleaseService) GetAsset(ctx context.Context, owner, repo string, releaseID, assetID int64) (*types.Attachment, error)` | GetAsset returns a single asset for a release. | No direct tests. |
| method | ReleaseService.GetByTag | `func (s *ReleaseService) GetByTag(ctx context.Context, owner, repo, tag string) (*types.Release, error)` | GetByTag returns a release by its tag name. | `TestReleaseService_Good_GetByTag` |
| method | ReleaseService.IterAssets | `func (s *ReleaseService) IterAssets(ctx context.Context, owner, repo string, releaseID int64) iter.Seq2[types.Attachment, error]` | IterAssets returns an iterator over all assets for a release. | No direct tests. |
| method | ReleaseService.ListAssets | `func (s *ReleaseService) ListAssets(ctx context.Context, owner, repo string, releaseID int64) ([]types.Attachment, error)` | ListAssets returns all assets for a release. | No direct tests. |
| method | RepoService.AcceptTransfer | `func (s *RepoService) AcceptTransfer(ctx context.Context, owner, repo string) error` | AcceptTransfer accepts a pending repository transfer. | No direct tests. |
| method | RepoService.Fork | `func (s *RepoService) Fork(ctx context.Context, owner, repo, org string) (*types.Repository, error)` | Fork forks a repository. If org is non-empty, forks into that organisation. | `TestRepoService_Good_Fork` |
| method | RepoService.IterOrgRepos | `func (s *RepoService) IterOrgRepos(ctx context.Context, org string) iter.Seq2[types.Repository, error]` | IterOrgRepos returns an iterator over all repositories for an organisation. | No direct tests. |
| method | RepoService.IterUserRepos | `func (s *RepoService) IterUserRepos(ctx context.Context) iter.Seq2[types.Repository, error]` | IterUserRepos returns an iterator over all repositories for the authenticated user. | No direct tests. |
| method | RepoService.ListOrgRepos | `func (s *RepoService) ListOrgRepos(ctx context.Context, org string) ([]types.Repository, error)` | ListOrgRepos returns all repositories for an organisation. | `TestRepoService_Good_ListOrgRepos` |
| method | RepoService.ListUserRepos | `func (s *RepoService) ListUserRepos(ctx context.Context) ([]types.Repository, error)` | ListUserRepos returns all repositories for the authenticated user. | No direct tests. |
| method | RepoService.MirrorSync | `func (s *RepoService) MirrorSync(ctx context.Context, owner, repo string) error` | MirrorSync triggers a mirror sync. | No direct tests. |
| method | RepoService.RejectTransfer | `func (s *RepoService) RejectTransfer(ctx context.Context, owner, repo string) error` | RejectTransfer rejects a pending repository transfer. | No direct tests. |
| method | RepoService.Transfer | `func (s *RepoService) Transfer(ctx context.Context, owner, repo string, opts map[string]any) error` | Transfer initiates a repository transfer. | No direct tests. |
| method | Resource.Create | `func (r *Resource[T, C, U]) Create(ctx context.Context, params Params, body *C) (*T, error)` | Create creates a new resource. | `TestResource_Good_Create` |
| method | Resource.Delete | `func (r *Resource[T, C, U]) Delete(ctx context.Context, params Params) error` | Delete removes a resource. | `TestResource_Good_Delete` |
| method | Resource.Get | `func (r *Resource[T, C, U]) Get(ctx context.Context, params Params) (*T, error)` | Get returns a single resource by appending id to the path. | `TestResource_Good_Get` |
| method | Resource.Iter | `func (r *Resource[T, C, U]) Iter(ctx context.Context, params Params) iter.Seq2[T, error]` | Iter returns an iterator over all resources across all pages. | `TestResource_Bad_IterError`, `TestResource_Good_Iter`, `TestResource_Good_IterBreakEarly` |
| method | Resource.List | `func (r *Resource[T, C, U]) List(ctx context.Context, params Params, opts ListOptions) (*PagedResult[T], error)` | List returns a single page of resources. | `TestResource_Good_List` |
| method | Resource.ListAll | `func (r *Resource[T, C, U]) ListAll(ctx context.Context, params Params) ([]T, error)` | ListAll returns all resources across all pages. | `TestResource_Good_ListAll` |
| method | Resource.Update | `func (r *Resource[T, C, U]) Update(ctx context.Context, params Params, body *U) (*T, error)` | Update modifies an existing resource. | `TestResource_Good_Update` |
| method | TeamService.AddMember | `func (s *TeamService) AddMember(ctx context.Context, teamID int64, username string) error` | AddMember adds a user to a team. | `TestTeamService_Good_AddMember` |
| method | TeamService.AddRepo | `func (s *TeamService) AddRepo(ctx context.Context, teamID int64, org, repo string) error` | AddRepo adds a repository to a team. | No direct tests. |
| method | TeamService.IterMembers | `func (s *TeamService) IterMembers(ctx context.Context, teamID int64) iter.Seq2[types.User, error]` | IterMembers returns an iterator over all members of a team. | No direct tests. |
| method | TeamService.IterOrgTeams | `func (s *TeamService) IterOrgTeams(ctx context.Context, org string) iter.Seq2[types.Team, error]` | IterOrgTeams returns an iterator over all teams in an organisation. | No direct tests. |
| method | TeamService.IterRepos | `func (s *TeamService) IterRepos(ctx context.Context, teamID int64) iter.Seq2[types.Repository, error]` | IterRepos returns an iterator over all repositories managed by a team. | No direct tests. |
| method | TeamService.ListMembers | `func (s *TeamService) ListMembers(ctx context.Context, teamID int64) ([]types.User, error)` | ListMembers returns all members of a team. | `TestTeamService_Good_ListMembers` |
| method | TeamService.ListOrgTeams | `func (s *TeamService) ListOrgTeams(ctx context.Context, org string) ([]types.Team, error)` | ListOrgTeams returns all teams in an organisation. | No direct tests. |
| method | TeamService.ListRepos | `func (s *TeamService) ListRepos(ctx context.Context, teamID int64) ([]types.Repository, error)` | ListRepos returns all repositories managed by a team. | No direct tests. |
| method | TeamService.RemoveMember | `func (s *TeamService) RemoveMember(ctx context.Context, teamID int64, username string) error` | RemoveMember removes a user from a team. | No direct tests. |
| method | TeamService.RemoveRepo | `func (s *TeamService) RemoveRepo(ctx context.Context, teamID int64, org, repo string) error` | RemoveRepo removes a repository from a team. | No direct tests. |
| method | UserService.Follow | `func (s *UserService) Follow(ctx context.Context, username string) error` | Follow follows a user as the authenticated user. | No direct tests. |
| method | UserService.GetCurrent | `func (s *UserService) GetCurrent(ctx context.Context) (*types.User, error)` | GetCurrent returns the authenticated user. | `TestUserService_Good_GetCurrent` |
| method | UserService.IterFollowers | `func (s *UserService) IterFollowers(ctx context.Context, username string) iter.Seq2[types.User, error]` | IterFollowers returns an iterator over all followers of a user. | No direct tests. |
| method | UserService.IterFollowing | `func (s *UserService) IterFollowing(ctx context.Context, username string) iter.Seq2[types.User, error]` | IterFollowing returns an iterator over all users that a user is following. | No direct tests. |
| method | UserService.IterStarred | `func (s *UserService) IterStarred(ctx context.Context, username string) iter.Seq2[types.Repository, error]` | IterStarred returns an iterator over all repositories starred by a user. | No direct tests. |
| method | UserService.ListFollowers | `func (s *UserService) ListFollowers(ctx context.Context, username string) ([]types.User, error)` | ListFollowers returns all followers of a user. | `TestUserService_Good_ListFollowers` |
| method | UserService.ListFollowing | `func (s *UserService) ListFollowing(ctx context.Context, username string) ([]types.User, error)` | ListFollowing returns all users that a user is following. | No direct tests. |
| method | UserService.ListStarred | `func (s *UserService) ListStarred(ctx context.Context, username string) ([]types.Repository, error)` | ListStarred returns all repositories starred by a user. | No direct tests. |
| method | UserService.Star | `func (s *UserService) Star(ctx context.Context, owner, repo string) error` | Star stars a repository as the authenticated user. | No direct tests. |
| method | UserService.Unfollow | `func (s *UserService) Unfollow(ctx context.Context, username string) error` | Unfollow unfollows a user as the authenticated user. | No direct tests. |
| method | UserService.Unstar | `func (s *UserService) Unstar(ctx context.Context, owner, repo string) error` | Unstar unstars a repository as the authenticated user. | No direct tests. |
| method | WebhookService.IterOrgHooks | `func (s *WebhookService) IterOrgHooks(ctx context.Context, org string) iter.Seq2[types.Hook, error]` | IterOrgHooks returns an iterator over all webhooks for an organisation. | No direct tests. |
| method | WebhookService.ListOrgHooks | `func (s *WebhookService) ListOrgHooks(ctx context.Context, org string) ([]types.Hook, error)` | ListOrgHooks returns all webhooks for an organisation. | `TestWebhookService_Good_ListOrgHooks` |
| method | WebhookService.TestHook | `func (s *WebhookService) TestHook(ctx context.Context, owner, repo string, id int64) error` | TestHook triggers a test delivery for a webhook. | `TestWebhookService_Good_TestHook` |
| method | WikiService.CreatePage | `func (s *WikiService) CreatePage(ctx context.Context, owner, repo string, opts *types.CreateWikiPageOptions) (*types.WikiPage, error)` | CreatePage creates a new wiki page. | `TestWikiService_Good_CreatePage` |
| method | WikiService.DeletePage | `func (s *WikiService) DeletePage(ctx context.Context, owner, repo, pageName string) error` | DeletePage removes a wiki page. | `TestWikiService_Good_DeletePage` |
| method | WikiService.EditPage | `func (s *WikiService) EditPage(ctx context.Context, owner, repo, pageName string, opts *types.CreateWikiPageOptions) (*types.WikiPage, error)` | EditPage updates an existing wiki page. | `TestWikiService_Good_EditPage` |
| method | WikiService.GetPage | `func (s *WikiService) GetPage(ctx context.Context, owner, repo, pageName string) (*types.WikiPage, error)` | GetPage returns a single wiki page by name. | `TestWikiService_Bad_NotFound`, `TestWikiService_Good_GetPage` |
| method | WikiService.ListPages | `func (s *WikiService) ListPages(ctx context.Context, owner, repo string) ([]types.WikiPageMetaData, error)` | ListPages returns all wiki page metadata for a repository. | `TestWikiService_Good_ListPages` |

## `forge/types`

| Kind | Name | Signature | Description | Test Coverage |
| --- | --- | --- | --- | --- |
| type | APIError | `type APIError struct` | APIError is an api error with a message | `TestAPIError_Good_Error`, `TestClient_Bad_ServerError`, `TestIsConflict_Bad_NotConflict` (+2 more) |
| type | APIForbiddenError | `type APIForbiddenError struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | APIInvalidTopicsError | `type APIInvalidTopicsError struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | APINotFound | `type APINotFound struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | APIRepoArchivedError | `type APIRepoArchivedError struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | APIUnauthorizedError | `type APIUnauthorizedError struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | APIValidationError | `type APIValidationError struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | AccessToken | `type AccessToken struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | ActionTask | `type ActionTask struct` | ActionTask represents a ActionTask | No direct tests; only indirect coverage via callers. |
| type | ActionTaskResponse | `type ActionTaskResponse struct` | ActionTaskResponse returns a ActionTask | No direct tests; only indirect coverage via callers. |
| type | ActionVariable | `type ActionVariable struct` | ActionVariable return value of the query API | `TestActionsService_Good_ListOrgVariables`, `TestActionsService_Good_ListRepoVariables` |
| type | Activity | `type Activity struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | ActivityPub | `type ActivityPub struct` | ActivityPub type | No direct tests; only indirect coverage via callers. |
| type | AddCollaboratorOption | `type AddCollaboratorOption struct` | AddCollaboratorOption options when adding a user as a collaborator of a repository | No direct tests; only indirect coverage via callers. |
| type | AddTimeOption | `type AddTimeOption struct` | AddTimeOption options for adding time to an issue | No direct tests; only indirect coverage via callers. |
| type | AnnotatedTag | `type AnnotatedTag struct` | AnnotatedTag represents an annotated tag | No direct tests; only indirect coverage via callers. |
| type | AnnotatedTagObject | `type AnnotatedTagObject struct` | AnnotatedTagObject contains meta information of the tag object | No direct tests; only indirect coverage via callers. |
| type | Attachment | `type Attachment struct` | Attachment a generic attachment | No direct tests; only indirect coverage via callers. |
| type | BlockedUser | `type BlockedUser struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | Branch | `type Branch struct` | Branch represents a repository branch | `TestBranchService_Good_Get`, `TestBranchService_Good_List` |
| type | BranchProtection | `type BranchProtection struct` | BranchProtection represents a branch protection for a repository | `TestBranchService_Good_CreateProtection` |
| type | ChangeFileOperation | `type ChangeFileOperation struct` | ChangeFileOperation for creating, updating or deleting a file | No direct tests; only indirect coverage via callers. |
| type | ChangeFilesOptions | `type ChangeFilesOptions struct` | ChangeFilesOptions options for creating, updating or deleting multiple files Note: `author` and `committer` are optional (if only one is given, it will be used for the other, otherwise the authenticated user will be used) | No direct tests; only indirect coverage via callers. |
| type | ChangedFile | `type ChangedFile struct` | ChangedFile store information about files affected by the pull request | No direct tests; only indirect coverage via callers. |
| type | CombinedStatus | `type CombinedStatus struct` | CombinedStatus holds the combined state of several statuses for a single commit | `TestCommitService_Good_GetCombinedStatus` |
| type | Comment | `type Comment struct` | Comment represents a comment on a commit or issue | `TestIssueService_Good_CreateComment` |
| type | Commit | `type Commit struct` | No doc comment. | `TestCommitService_Good_Get`, `TestCommitService_Good_GetNote`, `TestCommitService_Good_List` (+1 more) |
| type | CommitAffectedFiles | `type CommitAffectedFiles struct` | CommitAffectedFiles store information about files affected by the commit | No direct tests; only indirect coverage via callers. |
| type | CommitDateOptions | `type CommitDateOptions struct` | CommitDateOptions store dates for GIT_AUTHOR_DATE and GIT_COMMITTER_DATE | No direct tests; only indirect coverage via callers. |
| type | CommitMeta | `type CommitMeta struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | CommitStats | `type CommitStats struct` | CommitStats is statistics for a RepoCommit | No direct tests; only indirect coverage via callers. |
| type | CommitStatus | `type CommitStatus struct` | CommitStatus holds a single status of a single Commit | `TestCommitService_Good_CreateStatus`, `TestCommitService_Good_GetCombinedStatus`, `TestCommitService_Good_ListStatuses` |
| type | CommitStatusState | `type CommitStatusState struct` | CommitStatusState holds the state of a CommitStatus It can be "pending", "success", "error" and "failure" CommitStatusState has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | CommitUser | `type CommitUser struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | Compare | `type Compare struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | ContentsResponse | `type ContentsResponse struct` | ContentsResponse contains information about a repo's entry's (dir, file, symlink, submodule) metadata and content | `TestContentService_Good_CreateFile`, `TestContentService_Good_GetFile`, `TestContentService_Good_UpdateFile` |
| type | CreateAccessTokenOption | `type CreateAccessTokenOption struct` | CreateAccessTokenOption options when create access token | No direct tests; only indirect coverage via callers. |
| type | CreateBranchProtectionOption | `type CreateBranchProtectionOption struct` | CreateBranchProtectionOption options for creating a branch protection | `TestBranchService_Good_CreateProtection` |
| type | CreateBranchRepoOption | `type CreateBranchRepoOption struct` | CreateBranchRepoOption options when creating a branch in a repository | No direct tests; only indirect coverage via callers. |
| type | CreateEmailOption | `type CreateEmailOption struct` | CreateEmailOption options when creating email addresses | No direct tests; only indirect coverage via callers. |
| type | CreateFileOptions | `type CreateFileOptions struct` | CreateFileOptions options for creating files Note: `author` and `committer` are optional (if only one is given, it will be used for the other, otherwise the authenticated user will be used) | `TestContentService_Good_CreateFile` |
| type | CreateForkOption | `type CreateForkOption struct` | CreateForkOption options for creating a fork | No direct tests; only indirect coverage via callers. |
| type | CreateGPGKeyOption | `type CreateGPGKeyOption struct` | CreateGPGKeyOption options create user GPG key | No direct tests; only indirect coverage via callers. |
| type | CreateHookOption | `type CreateHookOption struct` | CreateHookOption options when create a hook | `TestWebhookService_Good_Create` |
| type | CreateHookOptionConfig | `type CreateHookOptionConfig struct` | CreateHookOptionConfig has all config options in it required are "content_type" and "url" Required CreateHookOptionConfig has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | CreateIssueCommentOption | `type CreateIssueCommentOption struct` | CreateIssueCommentOption options for creating a comment on an issue | `TestIssueService_Good_CreateComment` |
| type | CreateIssueOption | `type CreateIssueOption struct` | CreateIssueOption options to create one issue | `TestIssueService_Good_Create` |
| type | CreateKeyOption | `type CreateKeyOption struct` | CreateKeyOption options when creating a key | No direct tests; only indirect coverage via callers. |
| type | CreateLabelOption | `type CreateLabelOption struct` | CreateLabelOption options for creating a label | `TestLabelService_Good_CreateOrgLabel`, `TestLabelService_Good_CreateRepoLabel` |
| type | CreateMilestoneOption | `type CreateMilestoneOption struct` | CreateMilestoneOption options for creating a milestone | No direct tests; only indirect coverage via callers. |
| type | CreateOAuth2ApplicationOptions | `type CreateOAuth2ApplicationOptions struct` | CreateOAuth2ApplicationOptions holds options to create an oauth2 application | No direct tests; only indirect coverage via callers. |
| type | CreateOrUpdateSecretOption | `type CreateOrUpdateSecretOption struct` | CreateOrUpdateSecretOption options when creating or updating secret | No direct tests; only indirect coverage via callers. |
| type | CreateOrgOption | `type CreateOrgOption struct` | CreateOrgOption options for creating an organization | No direct tests; only indirect coverage via callers. |
| type | CreatePullRequestOption | `type CreatePullRequestOption struct` | CreatePullRequestOption options when creating a pull request | `TestPullService_Good_Create` |
| type | CreatePullReviewComment | `type CreatePullReviewComment struct` | CreatePullReviewComment represent a review comment for creation api | No direct tests; only indirect coverage via callers. |
| type | CreatePullReviewCommentOptions | `type CreatePullReviewCommentOptions struct` | CreatePullReviewCommentOptions are options to create a pull review comment CreatePullReviewCommentOptions has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | CreatePullReviewOptions | `type CreatePullReviewOptions struct` | CreatePullReviewOptions are options to create a pull review | No direct tests; only indirect coverage via callers. |
| type | CreatePushMirrorOption | `type CreatePushMirrorOption struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | CreateQuotaGroupOptions | `type CreateQuotaGroupOptions struct` | CreateQutaGroupOptions represents the options for creating a quota group | No direct tests; only indirect coverage via callers. |
| type | CreateQuotaRuleOptions | `type CreateQuotaRuleOptions struct` | CreateQuotaRuleOptions represents the options for creating a quota rule | No direct tests; only indirect coverage via callers. |
| type | CreateReleaseOption | `type CreateReleaseOption struct` | CreateReleaseOption options when creating a release | No direct tests; only indirect coverage via callers. |
| type | CreateRepoOption | `type CreateRepoOption struct` | CreateRepoOption options when creating repository | No direct tests; only indirect coverage via callers. |
| type | CreateStatusOption | `type CreateStatusOption struct` | CreateStatusOption holds the information needed to create a new CommitStatus for a Commit | `TestCommitService_Good_CreateStatus` |
| type | CreateTagOption | `type CreateTagOption struct` | CreateTagOption options when creating a tag | No direct tests; only indirect coverage via callers. |
| type | CreateTagProtectionOption | `type CreateTagProtectionOption struct` | CreateTagProtectionOption options for creating a tag protection | No direct tests; only indirect coverage via callers. |
| type | CreateTeamOption | `type CreateTeamOption struct` | CreateTeamOption options for creating a team | No direct tests; only indirect coverage via callers. |
| type | CreateUserOption | `type CreateUserOption struct` | CreateUserOption create user options | `TestAdminService_Bad_CreateUser_Forbidden`, `TestAdminService_Good_CreateUser` |
| type | CreateVariableOption | `type CreateVariableOption struct` | CreateVariableOption the option when creating variable | `TestActionsService_Good_CreateRepoVariable` |
| type | CreateWikiPageOptions | `type CreateWikiPageOptions struct` | CreateWikiPageOptions form for creating wiki | `TestWikiService_Good_CreatePage`, `TestWikiService_Good_EditPage` |
| type | Cron | `type Cron struct` | Cron represents a Cron task | `TestAdminService_Good_ListCron` |
| type | DeleteEmailOption | `type DeleteEmailOption struct` | DeleteEmailOption options when deleting email addresses | No direct tests; only indirect coverage via callers. |
| type | DeleteFileOptions | `type DeleteFileOptions struct` | DeleteFileOptions options for deleting files (used for other File structs below) Note: `author` and `committer` are optional (if only one is given, it will be used for the other, otherwise the authenticated user will be used) | `TestContentService_Good_DeleteFile` |
| type | DeleteLabelsOption | `type DeleteLabelsOption struct` | DeleteLabelOption options for deleting a label | No direct tests; only indirect coverage via callers. |
| type | DeployKey | `type DeployKey struct` | DeployKey a deploy key | No direct tests; only indirect coverage via callers. |
| type | DismissPullReviewOptions | `type DismissPullReviewOptions struct` | DismissPullReviewOptions are options to dismiss a pull review | No direct tests; only indirect coverage via callers. |
| type | DispatchWorkflowOption | `type DispatchWorkflowOption struct` | DispatchWorkflowOption options when dispatching a workflow | No direct tests; only indirect coverage via callers. |
| type | EditAttachmentOptions | `type EditAttachmentOptions struct` | EditAttachmentOptions options for editing attachments | No direct tests; only indirect coverage via callers. |
| type | EditBranchProtectionOption | `type EditBranchProtectionOption struct` | EditBranchProtectionOption options for editing a branch protection | No direct tests; only indirect coverage via callers. |
| type | EditDeadlineOption | `type EditDeadlineOption struct` | EditDeadlineOption options for creating a deadline | No direct tests; only indirect coverage via callers. |
| type | EditGitHookOption | `type EditGitHookOption struct` | EditGitHookOption options when modifying one Git hook | No direct tests; only indirect coverage via callers. |
| type | EditHookOption | `type EditHookOption struct` | EditHookOption options when modify one hook | No direct tests; only indirect coverage via callers. |
| type | EditIssueCommentOption | `type EditIssueCommentOption struct` | EditIssueCommentOption options for editing a comment | No direct tests; only indirect coverage via callers. |
| type | EditIssueOption | `type EditIssueOption struct` | EditIssueOption options for editing an issue | `TestIssueService_Good_Update` |
| type | EditLabelOption | `type EditLabelOption struct` | EditLabelOption options for editing a label | `TestLabelService_Good_EditRepoLabel` |
| type | EditMilestoneOption | `type EditMilestoneOption struct` | EditMilestoneOption options for editing a milestone | No direct tests; only indirect coverage via callers. |
| type | EditOrgOption | `type EditOrgOption struct` | EditOrgOption options for editing an organization | No direct tests; only indirect coverage via callers. |
| type | EditPullRequestOption | `type EditPullRequestOption struct` | EditPullRequestOption options when modify pull request | No direct tests; only indirect coverage via callers. |
| type | EditQuotaRuleOptions | `type EditQuotaRuleOptions struct` | EditQuotaRuleOptions represents the options for editing a quota rule | No direct tests; only indirect coverage via callers. |
| type | EditReactionOption | `type EditReactionOption struct` | EditReactionOption contain the reaction type | No direct tests; only indirect coverage via callers. |
| type | EditReleaseOption | `type EditReleaseOption struct` | EditReleaseOption options when editing a release | No direct tests; only indirect coverage via callers. |
| type | EditRepoOption | `type EditRepoOption struct` | EditRepoOption options when editing a repository's properties | `TestRepoService_Good_Update` |
| type | EditTagProtectionOption | `type EditTagProtectionOption struct` | EditTagProtectionOption options for editing a tag protection | No direct tests; only indirect coverage via callers. |
| type | EditTeamOption | `type EditTeamOption struct` | EditTeamOption options for editing a team | No direct tests; only indirect coverage via callers. |
| type | EditUserOption | `type EditUserOption struct` | EditUserOption edit user options | No direct tests; only indirect coverage via callers. |
| type | Email | `type Email struct` | Email an email address belonging to a user | `TestAdminService_Bad_CreateUser_Forbidden`, `TestAdminService_Good_CreateUser` |
| type | ExternalTracker | `type ExternalTracker struct` | ExternalTracker represents settings for external tracker | No direct tests; only indirect coverage via callers. |
| type | ExternalWiki | `type ExternalWiki struct` | ExternalWiki represents setting for external wiki | No direct tests; only indirect coverage via callers. |
| type | FileCommitResponse | `type FileCommitResponse struct` | No doc comment. | `TestContentService_Good_CreateFile` |
| type | FileDeleteResponse | `type FileDeleteResponse struct` | FileDeleteResponse contains information about a repo's file that was deleted | `TestContentService_Good_DeleteFile` |
| type | FileLinksResponse | `type FileLinksResponse struct` | FileLinksResponse contains the links for a repo's file | No direct tests; only indirect coverage via callers. |
| type | FileResponse | `type FileResponse struct` | FileResponse contains information about a repo's file | `TestContentService_Good_CreateFile`, `TestContentService_Good_UpdateFile` |
| type | FilesResponse | `type FilesResponse struct` | FilesResponse contains information about multiple files from a repo | No direct tests; only indirect coverage via callers. |
| type | ForgeLike | `type ForgeLike struct` | ForgeLike activity data type ForgeLike has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | GPGKey | `type GPGKey struct` | GPGKey a user GPG key to sign commit and tag in repository | No direct tests; only indirect coverage via callers. |
| type | GPGKeyEmail | `type GPGKeyEmail struct` | GPGKeyEmail an email attached to a GPGKey | No direct tests; only indirect coverage via callers. |
| type | GeneralAPISettings | `type GeneralAPISettings struct` | GeneralAPISettings contains global api settings exposed by it | No direct tests; only indirect coverage via callers. |
| type | GeneralAttachmentSettings | `type GeneralAttachmentSettings struct` | GeneralAttachmentSettings contains global Attachment settings exposed by API | No direct tests; only indirect coverage via callers. |
| type | GeneralRepoSettings | `type GeneralRepoSettings struct` | GeneralRepoSettings contains global repository settings exposed by API | No direct tests; only indirect coverage via callers. |
| type | GeneralUISettings | `type GeneralUISettings struct` | GeneralUISettings contains global ui settings exposed by API | No direct tests; only indirect coverage via callers. |
| type | GenerateRepoOption | `type GenerateRepoOption struct` | GenerateRepoOption options when creating repository using a template | No direct tests; only indirect coverage via callers. |
| type | GitBlobResponse | `type GitBlobResponse struct` | GitBlobResponse represents a git blob | No direct tests; only indirect coverage via callers. |
| type | GitEntry | `type GitEntry struct` | GitEntry represents a git tree | No direct tests; only indirect coverage via callers. |
| type | GitHook | `type GitHook struct` | GitHook represents a Git repository hook | No direct tests; only indirect coverage via callers. |
| type | GitObject | `type GitObject struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | GitTreeResponse | `type GitTreeResponse struct` | GitTreeResponse returns a git tree | No direct tests; only indirect coverage via callers. |
| type | GitignoreTemplateInfo | `type GitignoreTemplateInfo struct` | GitignoreTemplateInfo name and text of a gitignore template | `TestMiscService_Good_GetGitignoreTemplate` |
| type | Hook | `type Hook struct` | Hook a hook is a web hook when one repository changed | `TestWebhookService_Good_Create`, `TestWebhookService_Good_Get`, `TestWebhookService_Good_List` (+2 more) |
| type | Identity | `type Identity struct` | Identity for a person's identity like an author or committer | No direct tests; only indirect coverage via callers. |
| type | InternalTracker | `type InternalTracker struct` | InternalTracker represents settings for internal tracker | No direct tests; only indirect coverage via callers. |
| type | Issue | `type Issue struct` | Issue represents an issue in a repository | `TestIssueService_Good_Create`, `TestIssueService_Good_Get`, `TestIssueService_Good_List` (+2 more) |
| type | IssueConfig | `type IssueConfig struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | IssueConfigContactLink | `type IssueConfigContactLink struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | IssueConfigValidation | `type IssueConfigValidation struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | IssueDeadline | `type IssueDeadline struct` | IssueDeadline represents an issue deadline | No direct tests; only indirect coverage via callers. |
| type | IssueFormField | `type IssueFormField struct` | IssueFormField represents a form field | No direct tests; only indirect coverage via callers. |
| type | IssueFormFieldType | `type IssueFormFieldType struct` | IssueFormFieldType has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | IssueFormFieldVisible | `type IssueFormFieldVisible struct` | IssueFormFieldVisible defines issue form field visible IssueFormFieldVisible has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | IssueLabelsOption | `type IssueLabelsOption struct` | IssueLabelsOption a collection of labels | No direct tests; only indirect coverage via callers. |
| type | IssueMeta | `type IssueMeta struct` | IssueMeta basic issue information | No direct tests; only indirect coverage via callers. |
| type | IssueTemplate | `type IssueTemplate struct` | IssueTemplate represents an issue template for a repository | No direct tests; only indirect coverage via callers. |
| type | IssueTemplateLabels | `type IssueTemplateLabels struct` | IssueTemplateLabels has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | Label | `type Label struct` | Label a label to an issue or a pr | `TestLabelService_Good_CreateOrgLabel`, `TestLabelService_Good_CreateRepoLabel`, `TestLabelService_Good_EditRepoLabel` (+3 more) |
| type | LabelTemplate | `type LabelTemplate struct` | LabelTemplate info of a Label template | No direct tests; only indirect coverage via callers. |
| type | LicenseTemplateInfo | `type LicenseTemplateInfo struct` | LicensesInfo contains information about a License | `TestMiscService_Good_GetLicense` |
| type | LicensesTemplateListEntry | `type LicensesTemplateListEntry struct` | LicensesListEntry is used for the API | `TestMiscService_Good_ListLicenses` |
| type | MarkdownOption | `type MarkdownOption struct` | MarkdownOption markdown options | `TestMiscService_Good_RenderMarkdown` |
| type | MarkupOption | `type MarkupOption struct` | MarkupOption markup options | No direct tests; only indirect coverage via callers. |
| type | MergePullRequestOption | `type MergePullRequestOption struct` | MergePullRequestForm form for merging Pull Request | No direct tests; only indirect coverage via callers. |
| type | MigrateRepoOptions | `type MigrateRepoOptions struct` | MigrateRepoOptions options for migrating repository's this is used to interact with api v1 | No direct tests; only indirect coverage via callers. |
| type | Milestone | `type Milestone struct` | Milestone milestone is a collection of issues on one repository | No direct tests; only indirect coverage via callers. |
| type | NewIssuePinsAllowed | `type NewIssuePinsAllowed struct` | NewIssuePinsAllowed represents an API response that says if new Issue Pins are allowed | No direct tests; only indirect coverage via callers. |
| type | NodeInfo | `type NodeInfo struct` | NodeInfo contains standardized way of exposing metadata about a server running one of the distributed social networks | `TestMiscService_Good_GetNodeInfo` |
| type | NodeInfoServices | `type NodeInfoServices struct` | NodeInfoServices contains the third party sites this server can connect to via their application API | No direct tests; only indirect coverage via callers. |
| type | NodeInfoSoftware | `type NodeInfoSoftware struct` | NodeInfoSoftware contains Metadata about server software in use | `TestMiscService_Good_GetNodeInfo` |
| type | NodeInfoUsage | `type NodeInfoUsage struct` | NodeInfoUsage contains usage statistics for this server | No direct tests; only indirect coverage via callers. |
| type | NodeInfoUsageUsers | `type NodeInfoUsageUsers struct` | NodeInfoUsageUsers contains statistics about the users of this server | No direct tests; only indirect coverage via callers. |
| type | Note | `type Note struct` | Note contains information related to a git note | `TestCommitService_Good_GetNote` |
| type | NoteOptions | `type NoteOptions struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | NotificationCount | `type NotificationCount struct` | NotificationCount number of unread notifications | No direct tests; only indirect coverage via callers. |
| type | NotificationSubject | `type NotificationSubject struct` | NotificationSubject contains the notification subject (Issue/Pull/Commit) | `TestNotificationService_Good_GetThread`, `TestNotificationService_Good_List`, `TestNotificationService_Good_ListRepo` |
| type | NotificationThread | `type NotificationThread struct` | NotificationThread expose Notification on API | `TestNotificationService_Good_GetThread`, `TestNotificationService_Good_List`, `TestNotificationService_Good_ListRepo` |
| type | NotifySubjectType | `type NotifySubjectType struct` | NotifySubjectType represent type of notification subject NotifySubjectType has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | OAuth2Application | `type OAuth2Application struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | Organization | `type Organization struct` | Organization represents an organization | `TestAdminService_Good_ListOrgs`, `TestOrgService_Good_Get`, `TestOrgService_Good_List` |
| type | OrganizationPermissions | `type OrganizationPermissions struct` | OrganizationPermissions list different users permissions on an organization | No direct tests; only indirect coverage via callers. |
| type | PRBranchInfo | `type PRBranchInfo struct` | PRBranchInfo information about a branch | No direct tests; only indirect coverage via callers. |
| type | Package | `type Package struct` | Package represents a package | `TestPackageService_Good_Get`, `TestPackageService_Good_List` |
| type | PackageFile | `type PackageFile struct` | PackageFile represents a package file | `TestPackageService_Good_ListFiles` |
| type | PayloadCommit | `type PayloadCommit struct` | PayloadCommit represents a commit | No direct tests; only indirect coverage via callers. |
| type | PayloadCommitVerification | `type PayloadCommitVerification struct` | PayloadCommitVerification represents the GPG verification of a commit | No direct tests; only indirect coverage via callers. |
| type | PayloadUser | `type PayloadUser struct` | PayloadUser represents the author or committer of a commit | No direct tests; only indirect coverage via callers. |
| type | Permission | `type Permission struct` | Permission represents a set of permissions | No direct tests; only indirect coverage via callers. |
| type | PublicKey | `type PublicKey struct` | PublicKey publickey is a user key to push code to repository | No direct tests; only indirect coverage via callers. |
| type | PullRequest | `type PullRequest struct` | PullRequest represents a pull request | `TestPullService_Good_Create`, `TestPullService_Good_Get`, `TestPullService_Good_List` |
| type | PullRequestMeta | `type PullRequestMeta struct` | PullRequestMeta PR info if an issue is a PR | No direct tests; only indirect coverage via callers. |
| type | PullReview | `type PullReview struct` | PullReview represents a pull request review | No direct tests; only indirect coverage via callers. |
| type | PullReviewComment | `type PullReviewComment struct` | PullReviewComment represents a comment on a pull request review | No direct tests; only indirect coverage via callers. |
| type | PullReviewRequestOptions | `type PullReviewRequestOptions struct` | PullReviewRequestOptions are options to add or remove pull review requests | No direct tests; only indirect coverage via callers. |
| type | PushMirror | `type PushMirror struct` | PushMirror represents information of a push mirror | No direct tests; only indirect coverage via callers. |
| type | QuotaGroup | `type QuotaGroup struct` | QuotaGroup represents a quota group | No direct tests; only indirect coverage via callers. |
| type | QuotaGroupList | `type QuotaGroupList struct` | QuotaGroupList represents a list of quota groups QuotaGroupList has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | QuotaInfo | `type QuotaInfo struct` | QuotaInfo represents information about a user's quota | No direct tests; only indirect coverage via callers. |
| type | QuotaRuleInfo | `type QuotaRuleInfo struct` | QuotaRuleInfo contains information about a quota rule | No direct tests; only indirect coverage via callers. |
| type | QuotaUsed | `type QuotaUsed struct` | QuotaUsed represents the quota usage of a user | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedArtifact | `type QuotaUsedArtifact struct` | QuotaUsedArtifact represents an artifact counting towards a user's quota | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedArtifactList | `type QuotaUsedArtifactList struct` | QuotaUsedArtifactList represents a list of artifacts counting towards a user's quota QuotaUsedArtifactList has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedAttachment | `type QuotaUsedAttachment struct` | QuotaUsedAttachment represents an attachment counting towards a user's quota | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedAttachmentList | `type QuotaUsedAttachmentList struct` | QuotaUsedAttachmentList represents a list of attachment counting towards a user's quota QuotaUsedAttachmentList has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedPackage | `type QuotaUsedPackage struct` | QuotaUsedPackage represents a package counting towards a user's quota | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedPackageList | `type QuotaUsedPackageList struct` | QuotaUsedPackageList represents a list of packages counting towards a user's quota QuotaUsedPackageList has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedSize | `type QuotaUsedSize struct` | QuotaUsedSize represents the size-based quota usage of a user | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedSizeAssets | `type QuotaUsedSizeAssets struct` | QuotaUsedSizeAssets represents the size-based asset usage of a user | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedSizeAssetsAttachments | `type QuotaUsedSizeAssetsAttachments struct` | QuotaUsedSizeAssetsAttachments represents the size-based attachment quota usage of a user | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedSizeAssetsPackages | `type QuotaUsedSizeAssetsPackages struct` | QuotaUsedSizeAssetsPackages represents the size-based package quota usage of a user | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedSizeGit | `type QuotaUsedSizeGit struct` | QuotaUsedSizeGit represents the size-based git (lfs) quota usage of a user | No direct tests; only indirect coverage via callers. |
| type | QuotaUsedSizeRepos | `type QuotaUsedSizeRepos struct` | QuotaUsedSizeRepos represents the size-based repository quota usage of a user | No direct tests; only indirect coverage via callers. |
| type | Reaction | `type Reaction struct` | Reaction contain one reaction | No direct tests; only indirect coverage via callers. |
| type | Reference | `type Reference struct` | No doc comment. | No direct tests; only indirect coverage via callers. |
| type | Release | `type Release struct` | Release represents a repository release | `TestReleaseService_Good_Get`, `TestReleaseService_Good_GetByTag`, `TestReleaseService_Good_List` |
| type | RenameUserOption | `type RenameUserOption struct` | RenameUserOption options when renaming a user | `TestAdminService_Good_RenameUser` |
| type | ReplaceFlagsOption | `type ReplaceFlagsOption struct` | ReplaceFlagsOption options when replacing the flags of a repository | No direct tests; only indirect coverage via callers. |
| type | RepoCollaboratorPermission | `type RepoCollaboratorPermission struct` | RepoCollaboratorPermission to get repository permission for a collaborator | No direct tests; only indirect coverage via callers. |
| type | RepoCommit | `type RepoCommit struct` | No doc comment. | `TestCommitService_Good_Get`, `TestCommitService_Good_List` |
| type | RepoTopicOptions | `type RepoTopicOptions struct` | RepoTopicOptions a collection of repo topic names | No direct tests; only indirect coverage via callers. |
| type | RepoTransfer | `type RepoTransfer struct` | RepoTransfer represents a pending repo transfer | No direct tests; only indirect coverage via callers. |
| type | Repository | `type Repository struct` | Repository represents a repository | `TestRepoService_Good_Fork`, `TestRepoService_Good_Get`, `TestRepoService_Good_ListOrgRepos` (+1 more) |
| type | RepositoryMeta | `type RepositoryMeta struct` | RepositoryMeta basic repository information | No direct tests; only indirect coverage via callers. |
| type | ReviewStateType | `type ReviewStateType struct` | ReviewStateType review state type ReviewStateType has no fields in the swagger spec. | No direct tests; only indirect coverage via callers. |
| type | SearchResults | `type SearchResults struct` | SearchResults results of a successful search | No direct tests; only indirect coverage via callers. |
| type | Secret | `type Secret struct` | Secret represents a secret | `TestActionsService_Good_ListOrgSecrets`, `TestActionsService_Good_ListRepoSecrets` |
| type | ServerVersion | `type ServerVersion struct` | ServerVersion wraps the version of the server | `TestMiscService_Good_GetVersion` |
| type | SetUserQuotaGroupsOptions | `type SetUserQuotaGroupsOptions struct` | SetUserQuotaGroupsOptions represents the quota groups of a user | No direct tests; only indirect coverage via callers. |
| type | StateType | `type StateType string` | StateType is the state of an issue or PR: "open", "closed". | No direct tests; only indirect coverage via callers. |
| type | StopWatch | `type StopWatch struct` | StopWatch represent a running stopwatch | No direct tests; only indirect coverage via callers. |
| type | SubmitPullReviewOptions | `type SubmitPullReviewOptions struct` | SubmitPullReviewOptions are options to submit a pending pull review | No direct tests; only indirect coverage via callers. |
| type | Tag | `type Tag struct` | Tag represents a repository tag | No direct tests; only indirect coverage via callers. |
| type | TagArchiveDownloadCount | `type TagArchiveDownloadCount struct` | TagArchiveDownloadCount counts how many times a archive was downloaded | No direct tests; only indirect coverage via callers. |
| type | TagProtection | `type TagProtection struct` | TagProtection represents a tag protection | No direct tests; only indirect coverage via callers. |
| type | Team | `type Team struct` | Team represents a team in an organization | `TestTeamService_Good_Get` |
| type | TimeStamp | `type TimeStamp string` | TimeStamp is a Forgejo timestamp string. | No direct tests; only indirect coverage via callers. |
| type | TimelineComment | `type TimelineComment struct` | TimelineComment represents a timeline comment (comment of any type) on a commit or issue | No direct tests; only indirect coverage via callers. |
| type | TopicName | `type TopicName struct` | TopicName a list of repo topic names | No direct tests; only indirect coverage via callers. |
| type | TopicResponse | `type TopicResponse struct` | TopicResponse for returning topics | No direct tests; only indirect coverage via callers. |
| type | TrackedTime | `type TrackedTime struct` | TrackedTime worked time for an issue / pr | No direct tests; only indirect coverage via callers. |
| type | TransferRepoOption | `type TransferRepoOption struct` | TransferRepoOption options when transfer a repository's ownership | No direct tests; only indirect coverage via callers. |
| type | UpdateBranchRepoOption | `type UpdateBranchRepoOption struct` | UpdateBranchRepoOption options when updating a branch in a repository | No direct tests; only indirect coverage via callers. |
| type | UpdateFileOptions | `type UpdateFileOptions struct` | UpdateFileOptions options for updating files Note: `author` and `committer` are optional (if only one is given, it will be used for the other, otherwise the authenticated user will be used) | `TestContentService_Good_UpdateFile` |
| type | UpdateRepoAvatarOption | `type UpdateRepoAvatarOption struct` | UpdateRepoAvatarUserOption options when updating the repo avatar | No direct tests; only indirect coverage via callers. |
| type | UpdateUserAvatarOption | `type UpdateUserAvatarOption struct` | UpdateUserAvatarUserOption options when updating the user avatar | No direct tests; only indirect coverage via callers. |
| type | UpdateVariableOption | `type UpdateVariableOption struct` | UpdateVariableOption the option when updating variable | No direct tests; only indirect coverage via callers. |
| type | User | `type User struct` | User represents a user | `TestAdminService_Good_CreateUser`, `TestAdminService_Good_ListUsers`, `TestOrgService_Good_ListMembers` (+4 more) |
| type | UserHeatmapData | `type UserHeatmapData struct` | UserHeatmapData represents the data needed to create a heatmap | No direct tests; only indirect coverage via callers. |
| type | UserSettings | `type UserSettings struct` | UserSettings represents user settings | No direct tests; only indirect coverage via callers. |
| type | UserSettingsOptions | `type UserSettingsOptions struct` | UserSettingsOptions represents options to change user settings | No direct tests; only indirect coverage via callers. |
| type | WatchInfo | `type WatchInfo struct` | WatchInfo represents an API watch status of one repository | No direct tests; only indirect coverage via callers. |
| type | WikiCommit | `type WikiCommit struct` | WikiCommit page commit/revision | No direct tests; only indirect coverage via callers. |
| type | WikiCommitList | `type WikiCommitList struct` | WikiCommitList commit/revision list | No direct tests; only indirect coverage via callers. |
| type | WikiPage | `type WikiPage struct` | WikiPage a wiki page | `TestWikiService_Good_CreatePage`, `TestWikiService_Good_EditPage`, `TestWikiService_Good_GetPage` |
| type | WikiPageMetaData | `type WikiPageMetaData struct` | WikiPageMetaData wiki page meta information | `TestWikiService_Good_ListPages` |

## `cmd/forgegen`

| Kind | Name | Signature | Description | Test Coverage |
| --- | --- | --- | --- | --- |
| type | CRUDPair | `type CRUDPair struct` | CRUDPair groups a base type with its corresponding Create and Edit option types. | No direct tests. |
| type | GoField | `type GoField struct` | GoField is the intermediate representation for a single struct field. | No direct tests. |
| type | GoType | `type GoType struct` | GoType is the intermediate representation for a Go type to be generated. | `TestParser_Good_FieldTypes` |
| type | SchemaDefinition | `type SchemaDefinition struct` | SchemaDefinition represents a single type definition in the swagger spec. | No direct tests. |
| type | SchemaProperty | `type SchemaProperty struct` | SchemaProperty represents a single property within a schema definition. | No direct tests. |
| type | Spec | `type Spec struct` | Spec represents a Swagger 2.0 specification document. | No direct tests. |
| type | SpecInfo | `type SpecInfo struct` | SpecInfo holds metadata about the API specification. | No direct tests. |
| function | DetectCRUDPairs | `func DetectCRUDPairs(spec *Spec) []CRUDPair` | DetectCRUDPairs finds Create*Option / Edit*Option pairs in the swagger definitions and maps them back to the base type name. | `TestGenerate_Good_CreatesFiles`, `TestGenerate_Good_RepositoryType`, `TestGenerate_Good_TimeImport` (+2 more) |
| function | ExtractTypes | `func ExtractTypes(spec *Spec) map[string]*GoType` | ExtractTypes converts all swagger definitions into Go type intermediate representations. | `TestGenerate_Good_CreatesFiles`, `TestGenerate_Good_RepositoryType`, `TestGenerate_Good_TimeImport` (+3 more) |
| function | Generate | `func Generate(types map[string]*GoType, pairs []CRUDPair, outDir string) error` | Generate writes Go source files for the extracted types, grouped by logical domain. | `TestGenerate_Good_CreatesFiles`, `TestGenerate_Good_RepositoryType`, `TestGenerate_Good_TimeImport` (+1 more) |
| function | LoadSpec | `func LoadSpec(path string) (*Spec, error)` | LoadSpec reads and parses a Swagger 2.0 JSON file from the given path. | `TestGenerate_Good_CreatesFiles`, `TestGenerate_Good_RepositoryType`, `TestGenerate_Good_TimeImport` (+5 more) |
