# Services

go-forge provides 18 services accessible via the top-level `Forge` client. Each service groups related API endpoints.

```go
f := forge.NewForge("https://forge.lthn.ai", "my-token")

f.Repos         // Repository operations
f.Issues        // Issue operations
f.Pulls         // Pull request operations
f.Orgs          // Organisation operations
f.Users         // User operations
f.Teams         // Team operations
f.Admin         // Site administration
f.Branches      // Branch operations
f.Releases      // Release operations
f.Labels        // Label operations
f.Webhooks      // Webhook operations
f.Notifications // Notification operations
f.Packages      // Package registry
f.Actions       // CI/CD actions
f.Contents      // File read/write
f.Wiki          // Wiki pages
f.Misc          // Markdown, licences, gitignore, server info
f.Commits       // Commit statuses and git notes
```

---

## RepoService

Embeds `Resource[Repository, CreateRepoOption, EditRepoOption]`.
Path: `/api/v1/repos/{owner}/{repo}`

**CRUD** (via Resource): `List`, `ListAll`, `Get`, `Create`, `Update`, `Delete`

**Action methods:**

| Method           | Description                                  |
|------------------|----------------------------------------------|
| `ListOrgRepos`   | All repositories for an organisation         |
| `ListUserRepos`  | All repositories for the authenticated user  |
| `Fork`           | Fork a repository (optionally into an org)   |
| `Transfer`       | Initiate a repository transfer               |
| `AcceptTransfer` | Accept a pending transfer                    |
| `RejectTransfer` | Reject a pending transfer                    |
| `MirrorSync`     | Trigger a mirror sync                        |

---

## IssueService

Embeds `Resource[Issue, CreateIssueOption, EditIssueOption]`.
Path: `/api/v1/repos/{owner}/{repo}/issues/{index}`

**Action methods:**

| Method           | Description                        |
|------------------|------------------------------------|
| `Pin`            | Pin an issue                       |
| `Unpin`          | Unpin an issue                     |
| `SetDeadline`    | Set or update a deadline           |
| `AddReaction`    | Add a reaction emoji               |
| `DeleteReaction` | Remove a reaction emoji            |
| `StartStopwatch` | Start time tracking                |
| `StopStopwatch`  | Stop time tracking                 |
| `AddLabels`      | Add labels by ID                   |
| `RemoveLabel`    | Remove a single label              |
| `ListComments`   | List all comments on an issue      |
| `CreateComment`  | Create a comment                   |

---

## PullService

Embeds `Resource[PullRequest, CreatePullRequestOption, EditPullRequestOption]`.
Path: `/api/v1/repos/{owner}/{repo}/pulls/{index}`

**Action methods:**

| Method            | Description                            |
|-------------------|----------------------------------------|
| `Merge`           | Merge a pull request (merge, rebase, squash, etc.) |
| `Update`          | Update branch with base branch         |
| `ListReviews`     | List all reviews                       |
| `SubmitReview`    | Submit a new review                    |
| `DismissReview`   | Dismiss a review                       |
| `UndismissReview` | Undismiss a review                     |

---

## OrgService

Embeds `Resource[Organization, CreateOrgOption, EditOrgOption]`.
Path: `/api/v1/orgs/{org}`

**Action methods:**

| Method         | Description                                |
|----------------|--------------------------------------------|
| `ListMembers`  | List all members of an organisation        |
| `AddMember`    | Add a user to an organisation              |
| `RemoveMember` | Remove a user from an organisation         |
| `ListUserOrgs` | List organisations for a given user        |
| `ListMyOrgs`   | List organisations for the authenticated user |

---

## UserService

Embeds `Resource[User, struct{}, struct{}]` (read-only via CRUD).
Path: `/api/v1/users/{username}`

**Action methods:**

| Method          | Description                                |
|-----------------|--------------------------------------------|
| `GetCurrent`    | Get the authenticated user                 |
| `ListFollowers` | List followers of a user                   |
| `ListFollowing` | List users a user is following             |
| `Follow`        | Follow a user                              |
| `Unfollow`      | Unfollow a user                            |
| `ListStarred`   | List starred repositories                  |
| `Star`          | Star a repository                          |
| `Unstar`        | Unstar a repository                        |

---

## TeamService

Embeds `Resource[Team, CreateTeamOption, EditTeamOption]`.
Path: `/api/v1/teams/{id}`

**Action methods:**

| Method         | Description                        |
|----------------|------------------------------------|
| `ListMembers`  | List all team members              |
| `AddMember`    | Add a user to a team               |
| `RemoveMember` | Remove a user from a team          |
| `ListRepos`    | List repositories managed by team  |
| `AddRepo`      | Add a repository to a team         |
| `RemoveRepo`   | Remove a repository from a team    |
| `ListOrgTeams` | List all teams in an organisation  |

---

## AdminService

No Resource embedding (heterogeneous endpoints). Requires site admin privileges.

| Method                | Description                              |
|-----------------------|------------------------------------------|
| `ListUsers`           | List all users                           |
| `CreateUser`          | Create a new user                        |
| `EditUser`            | Edit an existing user                    |
| `DeleteUser`          | Delete a user                            |
| `RenameUser`          | Rename a user                            |
| `ListOrgs`            | List all organisations                   |
| `RunCron`             | Run a cron task by name                  |
| `ListCron`            | List all cron tasks                      |
| `AdoptRepo`           | Adopt an unadopted repository            |
| `GenerateRunnerToken` | Generate an actions runner registration token |

---

## BranchService

Embeds `Resource[Branch, CreateBranchRepoOption, struct{}]`.
Path: `/api/v1/repos/{owner}/{repo}/branches/{branch}`

**Action methods:**

| Method                   | Description                          |
|--------------------------|--------------------------------------|
| `ListBranchProtections`  | List all branch protection rules     |
| `GetBranchProtection`    | Get a single branch protection rule  |
| `CreateBranchProtection` | Create a new branch protection rule  |
| `EditBranchProtection`   | Update a branch protection rule      |
| `DeleteBranchProtection` | Delete a branch protection rule      |

---

## ReleaseService

Embeds `Resource[Release, CreateReleaseOption, EditReleaseOption]`.
Path: `/api/v1/repos/{owner}/{repo}/releases/{id}`

**Action methods:**

| Method        | Description                    |
|---------------|--------------------------------|
| `GetByTag`    | Get a release by tag name      |
| `DeleteByTag` | Delete a release by tag name   |
| `ListAssets`  | List assets for a release      |
| `GetAsset`    | Get a single release asset     |
| `DeleteAsset` | Delete a release asset         |

---

## LabelService

No Resource embedding (heterogeneous repo/org paths).

| Method           | Description                       |
|------------------|-----------------------------------|
| `ListRepoLabels` | List all labels in a repository   |
| `GetRepoLabel`   | Get a single repository label     |
| `CreateRepoLabel`| Create a label in a repository    |
| `EditRepoLabel`  | Update a repository label         |
| `DeleteRepoLabel`| Delete a repository label         |
| `ListOrgLabels`  | List all labels in an organisation|
| `CreateOrgLabel` | Create a label in an organisation |

---

## WebhookService

Embeds `Resource[Hook, CreateHookOption, EditHookOption]`.
Path: `/api/v1/repos/{owner}/{repo}/hooks/{id}`

**Action methods:**

| Method        | Description                            |
|---------------|----------------------------------------|
| `TestHook`    | Trigger a test delivery for a webhook  |
| `ListOrgHooks`| List all webhooks for an organisation  |

---

## ContentService

No Resource embedding (varied operations on file paths).

| Method       | Description                                |
|--------------|--------------------------------------------|
| `GetFile`    | Get metadata and content for a file        |
| `CreateFile` | Create a new file                          |
| `UpdateFile` | Update an existing file                    |
| `DeleteFile` | Delete a file (DELETE with body)           |
| `GetRawFile` | Get raw file content as bytes              |

---

## ActionsService

No Resource embedding (heterogeneous repo/org endpoints).

| Method               | Description                              |
|----------------------|------------------------------------------|
| `ListRepoSecrets`    | List secrets for a repository            |
| `CreateRepoSecret`   | Create or update a repository secret     |
| `DeleteRepoSecret`   | Delete a repository secret               |
| `ListRepoVariables`  | List action variables for a repository   |
| `CreateRepoVariable` | Create an action variable                |
| `DeleteRepoVariable` | Delete an action variable                |
| `ListOrgSecrets`     | List secrets for an organisation         |
| `ListOrgVariables`   | List action variables for an organisation|
| `DispatchWorkflow`   | Trigger a workflow run                   |

---

## NotificationService

No Resource embedding.

| Method          | Description                              |
|-----------------|------------------------------------------|
| `List`          | List all notifications                   |
| `ListRepo`      | List notifications for a repository      |
| `MarkRead`      | Mark all notifications as read           |
| `GetThread`     | Get a single notification thread         |
| `MarkThreadRead`| Mark a single thread as read             |

---

## PackageService

No Resource embedding.

| Method     | Description                                    |
|------------|------------------------------------------------|
| `List`     | List all packages for an owner                 |
| `Get`      | Get a package by owner, type, name, version    |
| `Delete`   | Delete a package                               |
| `ListFiles`| List files for a specific package version      |

---

## WikiService

No Resource embedding.

| Method       | Description                  |
|--------------|------------------------------|
| `ListPages`  | List all wiki page metadata  |
| `GetPage`    | Get a single wiki page       |
| `CreatePage` | Create a new wiki page       |
| `EditPage`   | Update an existing page      |
| `DeletePage` | Delete a wiki page           |

---

## MiscService

No Resource embedding (read-only utility endpoints).

| Method                    | Description                         |
|---------------------------|-------------------------------------|
| `RenderMarkdown`          | Render markdown text to HTML        |
| `ListLicenses`            | List available licence templates    |
| `GetLicense`              | Get a single licence template       |
| `ListGitignoreTemplates`  | List gitignore template names       |
| `GetGitignoreTemplate`    | Get a single gitignore template     |
| `GetNodeInfo`             | Get Forgejo instance NodeInfo       |
| `GetVersion`              | Get server version                  |

---

## CommitService

No Resource embedding.

| Method              | Description                                  |
|---------------------|----------------------------------------------|
| `GetCombinedStatus` | Get combined status for a ref                |
| `ListStatuses`      | List all commit statuses for a ref           |
| `CreateStatus`      | Create a commit status for a SHA             |
| `GetNote`           | Get the git note for a commit SHA            |
