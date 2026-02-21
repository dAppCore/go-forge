# Architecture

## Overview

go-forge is structured around three core ideas:

1. **Generic Resource** — a single `Resource[T, C, U]` type that provides CRUD methods for any Forgejo API resource.
2. **Code generation** — all 229 Go types are generated from `swagger.v1.json` via a custom codegen tool.
3. **Service embedding** — each service struct embeds `Resource` for free CRUD, then adds hand-written action methods for non-CRUD endpoints.

## File Layout

```
go-forge/
  client.go          HTTP client with auth, context, error handling
  pagination.go      Generic ListPage[T] and ListAll[T]
  params.go          Path variable resolution ({owner}/{repo} -> values)
  resource.go        Generic Resource[T, C, U] for CRUD
  config.go          Config resolution: flags -> env -> defaults
  forge.go           Top-level Forge client aggregating all 18 services

  repos.go           RepoService          (embeds Resource)
  issues.go          IssueService         (embeds Resource)
  pulls.go           PullService          (embeds Resource)
  orgs.go            OrgService           (embeds Resource)
  users.go           UserService          (embeds Resource)
  teams.go           TeamService          (embeds Resource)
  branches.go        BranchService        (embeds Resource)
  releases.go        ReleaseService       (embeds Resource)
  webhooks.go        WebhookService       (embeds Resource)
  admin.go           AdminService         (plain client)
  labels.go          LabelService         (plain client)
  contents.go        ContentService       (plain client)
  actions.go         ActionsService       (plain client)
  notifications.go   NotificationService  (plain client)
  packages.go        PackageService       (plain client)
  wiki.go            WikiService          (plain client)
  misc.go            MiscService          (plain client)
  commits.go         CommitService        (plain client)

  cmd/forgegen/      Code generator
    main.go          CLI entry point
    parser.go        Swagger 2.0 spec parser
    generator.go     Go source code emitter

  types/             Generated Go types (229 types across 36 files)
    generate.go      go:generate directive
    repo.go          Repository, CreateRepoOption, EditRepoOption, ...
    issue.go         Issue, CreateIssueOption, ...
    ...

  testdata/
    swagger.v1.json  Forgejo Swagger 2.0 specification
```

## Generic Resource[T, C, U]

The core abstraction. `T` is the resource type (e.g. `types.Repository`), `C` is the create options type (e.g. `types.CreateRepoOption`), `U` is the update options type (e.g. `types.EditRepoOption`).

```go
type Resource[T any, C any, U any] struct {
    client *Client
    path   string    // e.g. "/api/v1/repos/{owner}/{repo}"
}
```

Provides six methods out of the box:

| Method   | HTTP     | Description                       |
|----------|----------|-----------------------------------|
| `List`   | GET      | Single page with `ListOptions`    |
| `ListAll`| GET      | All pages (auto-pagination)       |
| `Get`    | GET      | Single resource by path params    |
| `Create` | POST     | Create with `*C` body             |
| `Update` | PATCH    | Update with `*U` body             |
| `Delete` | DELETE   | Delete by path params             |

Path parameters like `{owner}` and `{repo}` are resolved at call time via `Params` maps:

```go
repo, err := f.Repos.Get(ctx, forge.Params{"owner": "core", "repo": "go-forge"})
```

## Service Pattern

Services that fit the CRUD model embed `Resource`:

```go
type RepoService struct {
    Resource[types.Repository, types.CreateRepoOption, types.EditRepoOption]
}
```

This gives `RepoService` free `List`, `ListAll`, `Get`, `Create`, `Update`, and `Delete` methods. Action methods like `Fork`, `Transfer`, and `MirrorSync` are added as hand-written methods.

Services with heterogeneous endpoints (e.g. `AdminService`, `LabelService`) hold a plain `*Client` field instead of embedding `Resource`.

## Pagination

All list endpoints use Forgejo's `X-Total-Count` header for total item counts and `page`/`limit` query parameters.

- `ListPage[T]` — fetches a single page, returns `*PagedResult[T]` with `Items`, `TotalCount`, `Page`, `HasMore`.
- `ListAll[T]` — iterates all pages automatically, returns `[]T`.

## HTTP Client

`Client` is a low-level HTTP client that:

- Adds `Authorization: token <token>` to every request
- Sets `Accept: application/json` and `Content-Type: application/json`
- Wraps errors as `*APIError` with `StatusCode`, `Message`, `URL`
- Supports `Get`, `Post`, `Patch`, `Put`, `Delete`, `DeleteWithBody`, `GetRaw`, `PostRaw`
- Accepts `context.Context` on all methods

## Code Generation Pipeline

```
swagger.v1.json  -->  parser.go  -->  generator.go  -->  types/*.go
                      (extract)       (emit Go)         (229 types)
```

The parser reads the Swagger 2.0 spec, extracts type definitions, detects Create/Edit option pairs, and resolves Go types. The generator groups types by domain (repo, issue, PR, etc.) and writes one `.go` file per group with `// Code generated` headers.
