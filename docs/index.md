---
title: go-forge
description: Full-coverage Go client for the Forgejo API with generics-based CRUD, pagination, and code-generated types.
---

# go-forge

`dappco.re/go/core/forge` is a Go client library for the [Forgejo](https://forgejo.org) REST API. It provides typed access to 18 API domains (repositories, issues, pull requests, organisations, and more) through a single top-level `Forge` client. Types are generated directly from Forgejo's `swagger.v1.json` specification, keeping the library in lockstep with the server.

**Module path:** `dappco.re/go/core/forge`
**Go version:** 1.26+
**Licence:** EUPL-1.2


## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "dappco.re/go/core/forge"
)

func main() {
    // Create a client with your Forgejo URL and API token.
    f := forge.NewForge("https://forge.lthn.ai", "your-token")

    ctx := context.Background()

    // List repositories for an organisation (first page, 50 per page).
    result, err := f.Repos.List(ctx, forge.Params{"org": "core"}, forge.DefaultList)
    if err != nil {
        log.Fatal(err)
    }
    for _, repo := range result.Items {
        fmt.Println(repo.Name)
    }

    // Get a single repository.
    repo, err := f.Repos.Get(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s — %s\n", repo.FullName, repo.Description)
}
```

### Configuration from environment

If you prefer to resolve the URL and token from environment variables rather than hard-coding them, use `NewForgeFromConfig`:

```go
// Priority: flags > env (FORGE_URL, FORGE_TOKEN) > defaults (http://localhost:3000)
f, err := forge.NewForgeFromConfig("", "", forge.WithUserAgent("my-tool/1.0"))
if err != nil {
    log.Fatal(err) // no token configured
}
```

Environment variables:

| Variable      | Purpose                              | Default                  |
|---------------|--------------------------------------|--------------------------|
| `FORGE_URL`   | Base URL of the Forgejo instance     | `http://localhost:3000`  |
| `FORGE_TOKEN` | API token for authentication         | (none -- required)       |


## Package layout

```
go-forge/
├── client.go          HTTP client, auth, error handling, rate limits
├── config.go          Config resolution: flags > env > defaults
├── forge.go           Top-level Forge struct aggregating all 18 services
├── resource.go        Generic Resource[T, C, U] for CRUD operations
├── pagination.go      ListPage, ListAll, ListIter — paginated requests
├── params.go          Path variable resolution ({owner}/{repo} -> values)
├── repos.go           RepoService — repositories, forks, transfers, mirrors
├── issues.go          IssueService — issues, comments, labels, reactions
├── pulls.go           PullService — pull requests, merges, reviews
├── orgs.go            OrgService — organisations, members
├── users.go           UserService — users, followers, stars
├── teams.go           TeamService — teams, members, repositories
├── admin.go           AdminService — site admin, cron, user management
├── branches.go        BranchService — branches, branch protections
├── releases.go        ReleaseService — releases, assets, tags
├── labels.go          LabelService — repo and org labels
├── webhooks.go        WebhookService — repo and org webhooks
├── notifications.go   NotificationService — notifications, threads
├── packages.go        PackageService — package registry
├── actions.go         ActionsService — CI/CD secrets, variables, dispatches, tasks
├── contents.go        ContentService — file read/write/delete
├── wiki.go            WikiService — wiki pages
├── commits.go         CommitService — statuses, notes
├── misc.go            MiscService — markdown, licences, gitignore, version
├── types/             229 generated Go types from swagger.v1.json
│   ├── generate.go    go:generate directive
│   ├── repo.go        Repository, CreateRepoOption, EditRepoOption, ...
│   ├── issue.go       Issue, CreateIssueOption, ...
│   ├── pr.go          PullRequest, CreatePullRequestOption, ...
│   └── ... (36 files total, grouped by domain)
├── cmd/forgegen/      Code generator: swagger spec -> types/*.go
│   ├── main.go        CLI entry point
│   ├── parser.go      Swagger spec parsing, type extraction, CRUD pair detection
│   └── generator.go   Template-based Go source file generation
└── testdata/
    └── swagger.v1.json   Forgejo API specification (input for codegen)
```


## Services

The `Forge` struct exposes 18 service fields, each handling a different API domain:

| Service         | Struct              | Embedding                        | Domain                               |
|-----------------|---------------------|----------------------------------|--------------------------------------|
| `Repos`         | `RepoService`       | `Resource[Repository, ...]`      | Repositories, forks, transfers       |
| `Issues`        | `IssueService`      | `Resource[Issue, ...]`           | Issues, comments, labels, reactions  |
| `Pulls`         | `PullService`       | `Resource[PullRequest, ...]`     | Pull requests, merges, reviews       |
| `Orgs`          | `OrgService`        | `Resource[Organization, ...]`    | Organisations, members               |
| `Users`         | `UserService`       | `Resource[User, ...]`            | Users, followers, stars              |
| `Teams`         | `TeamService`       | `Resource[Team, ...]`            | Teams, members, repos                |
| `Admin`         | `AdminService`      | (standalone)                     | Site admin, cron, user management    |
| `Branches`      | `BranchService`     | `Resource[Branch, ...]`          | Branches, protections                |
| `Releases`      | `ReleaseService`    | `Resource[Release, ...]`         | Releases, assets, tags               |
| `Labels`        | `LabelService`      | (standalone)                     | Repo and org labels                  |
| `Webhooks`      | `WebhookService`    | `Resource[Hook, ...]`            | Repo and org webhooks                |
| `Notifications` | `NotificationService` | (standalone)                   | Notifications, threads               |
| `Packages`      | `PackageService`    | (standalone)                     | Package registry                     |
| `Actions`       | `ActionsService`    | (standalone)                     | CI/CD secrets, variables, dispatches, tasks |
| `Contents`      | `ContentService`    | (standalone)                     | File read/write/delete               |
| `Wiki`          | `WikiService`       | (standalone)                     | Wiki pages                           |
| `Commits`       | `CommitService`     | (standalone)                     | Commit statuses, git notes           |
| `Misc`          | `MiscService`       | (standalone)                     | Markdown, licences, gitignore, version |

Services that embed `Resource[T, C, U]` inherit `List`, `ListAll`, `Iter`, `Get`, `Create`, `Update`, and `Delete` methods automatically. Standalone services have hand-written methods because their API endpoints are heterogeneous and do not fit a uniform CRUD pattern.


## Dependencies

This module has a small dependency set: `dappco.re/go/core` and `github.com/goccy/go-json`, plus the Go standard library (`net/http`, `context`, `iter`, etc.) where appropriate.

```
module dappco.re/go/core/forge

go 1.26.0
```
