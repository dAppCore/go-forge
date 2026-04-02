---
title: Architecture
description: Internals of go-forge — the HTTP client, generic Resource pattern, pagination, code generation, and error handling.
---

# Architecture

This document explains how go-forge is structured internally. It covers the layered design, key types, data flow for a typical API call, and how types are generated from the Forgejo swagger specification.


## Layered design

go-forge is organised in three layers, each building on the one below:

```
┌─────────────────────────────────────────────────┐
│  Forge (top-level client)                       │
│  Aggregates 20 service structs                  │
├─────────────────────────────────────────────────┤
│  Service layer                                  │
│  RepoService, IssueService, PullService, ...    │
│  Embed Resource[T,C,U] or hold a *Client        │
├─────────────────────────────────────────────────┤
│  Foundation layer                               │
│  Client (HTTP), Resource[T,C,U] (generics),     │
│  Pagination, Params, Config                     │
└─────────────────────────────────────────────────┘
```


## Client — the HTTP layer

`Client` (`client.go`) is the lowest-level component. It handles:

- **Authentication** — every request includes an `Authorization: token <token>` header.
- **JSON marshalling** — request bodies are marshalled to JSON; responses are decoded from JSON.
- **Error parsing** — HTTP 4xx/5xx responses are converted to `*APIError` with status code, message, and URL.
- **Rate limit tracking** — `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `X-RateLimit-Reset` headers are captured after each request.
- **Context propagation** — all methods accept `context.Context` as their first parameter.

### Key methods

```go
func (c *Client) Get(ctx context.Context, path string, out any) error
func (c *Client) Post(ctx context.Context, path string, body, out any) error
func (c *Client) Patch(ctx context.Context, path string, body, out any) error
func (c *Client) Put(ctx context.Context, path string, body, out any) error
func (c *Client) Delete(ctx context.Context, path string) error
func (c *Client) DeleteWithBody(ctx context.Context, path string, body any) error
func (c *Client) GetRaw(ctx context.Context, path string) ([]byte, error)
func (c *Client) PostRaw(ctx context.Context, path string, body any) ([]byte, error)
```

The `Raw` variants return the response body as `[]byte` instead of decoding JSON. This is used by endpoints that return non-JSON content (e.g. the markdown rendering endpoint returns raw HTML).

### Options

The client supports functional options:

```go
f := forge.NewForge("https://forge.lthn.ai", "token",
    forge.WithHTTPClient(customHTTPClient),
    forge.WithUserAgent("my-agent/1.0"),
)
```


## APIError — structured error handling

All API errors are returned as `*APIError`:

```go
type APIError struct {
    StatusCode int
    Message    string
    URL        string
}
```

Three helper functions allow classification without type-asserting:

```go
forge.IsNotFound(err)  // 404
forge.IsForbidden(err) // 403
forge.IsConflict(err)  // 409
```

These use `errors.As` internally, so they work correctly with wrapped errors.


## Params — path variable resolution

API paths contain `{placeholders}` (e.g. `/api/v1/repos/{owner}/{repo}`). The `Params` type is a `map[string]string` that resolves these:

```go
type Params map[string]string

path := forge.ResolvePath("/api/v1/repos/{owner}/{repo}", forge.Params{
    "owner": "core",
    "repo":  "go-forge",
})
// Result: "/api/v1/repos/core/go-forge"
```

Values are URL-path-escaped to handle special characters safely.


## Resource[T, C, U] — the generic CRUD core

The heart of the library is `Resource[T, C, U]` (`resource.go`), a generic struct parameterised on three types:

- **T** — the resource type (e.g. `types.Repository`)
- **C** — the create-options type (e.g. `types.CreateRepoOption`)
- **U** — the update-options type (e.g. `types.EditRepoOption`)

It provides seven methods that map directly to REST operations:

| Method    | HTTP verb | Description                              |
|-----------|-----------|------------------------------------------|
| `List`    | GET       | Single page of results with metadata     |
| `ListAll` | GET       | All results across all pages             |
| `Iter`    | GET       | `iter.Seq2[T, error]` iterator over all  |
| `Get`     | GET       | Single resource by path params           |
| `Create`  | POST      | Create a new resource                    |
| `Update`  | PATCH     | Modify an existing resource              |
| `Delete`  | DELETE    | Remove a resource                        |

Services embed this struct to inherit CRUD for free:

```go
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
```

Services that do not fit the CRUD pattern (e.g. `AdminService`, `LabelService`, `NotificationService`) hold a `*Client` directly and implement methods by hand.


## Pagination

`pagination.go` provides three generic functions for paginated endpoints:

### ListPage — single page

```go
func ListPage[T any](ctx context.Context, c *Client, path string, query map[string]string, opts ListOptions) (*PagedResult[T], error)
```

Returns a `PagedResult[T]` containing:
- `Items []T` — the results on this page
- `TotalCount int` — from the `X-Total-Count` response header
- `Page int` — the current page number
- `HasMore bool` — whether more pages exist

`ListOptions` controls the page number (1-based) and items per page:

```go
var DefaultList = ListOptions{Page: 1, Limit: 50}
```

### ListAll — all pages

```go
func ListAll[T any](ctx context.Context, c *Client, path string, query map[string]string) ([]T, error)
```

Fetches every page sequentially and returns the concatenated slice. Uses a page size of 50.

### ListIter — range-over-func iterator

```go
func ListIter[T any](ctx context.Context, c *Client, path string, query map[string]string) iter.Seq2[T, error]
```

Returns a Go 1.23+ range-over-func iterator that lazily fetches pages as you consume items:

```go
for repo, err := range f.Repos.IterOrgRepos(ctx, "core") {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(repo.Name)
}
```


## Data flow of a typical API call

Here is the path a call like `f.Repos.Get(ctx, params)` takes:

```
1. Caller invokes f.Repos.Get(ctx, Params{"owner":"core", "repo":"go-forge"})
2. Resource.Get calls ResolvePath(r.path, params)
       "/api/v1/repos/{owner}/{repo}" -> "/api/v1/repos/core/go-forge"
3. Resource.Get calls Client.Get(ctx, resolvedPath, &out)
4. Client.doJSON builds http.Request:
       - Method: GET
       - URL: baseURL + resolvedPath
       - Headers: Authorization, Accept, User-Agent
5. Client.httpClient.Do(req) sends the request
6. Client reads response:
       - Updates rate limit from headers
       - If status >= 400: parseError -> return *APIError
       - If status < 400: json.Decode into &out
7. Resource.Get returns (*T, error) to caller
```


## Code generation — types/

The `types/` package contains 229 Go types generated from Forgejo's `swagger.v1.json` specification. The code generator lives at `cmd/forgegen/`.

### Pipeline

```
swagger.v1.json  -->  parser.go  -->  GoType/GoField IR  -->  generator.go  -->  types/*.go
                      LoadSpec()      ExtractTypes()          Generate()
                                      DetectCRUDPairs()
```

1. **parser.go** — `LoadSpec()` reads the JSON spec. `ExtractTypes()` converts swagger definitions into an intermediate representation (`GoType` with `GoField` children). `DetectCRUDPairs()` finds matching `Create*Option`/`Edit*Option` pairs.

2. **generator.go** — `Generate()` groups types by domain using prefix matching (the `typeGrouping` map), then renders each group into a separate `.go` file using `text/template`.

### Running the generator

```bash
go generate ./types/...
```

Or manually:

```bash
go run ./cmd/forgegen/ -spec testdata/swagger.v1.json -out types/
```

The `generate.go` file in `types/` contains the `//go:generate` directive:

```go
//go:generate go run ../cmd/forgegen/ -spec ../testdata/swagger.v1.json -out .
```

### Type grouping

Types are distributed across 36 files based on a name-prefix mapping. For example, types starting with `Repository` or `Repo` go into `repo.go`, types starting with `Issue` go into `issue.go`, and so on. The `classifyType()` function also strips `Create`/`Edit`/`Delete` prefixes and `Option` suffixes before matching, so `CreateRepoOption` lands in `repo.go` alongside `Repository`.

### Type mapping from swagger to Go

| Swagger type  | Swagger format | Go type     |
|---------------|----------------|-------------|
| string        | (none)         | `string`    |
| string        | date-time      | `time.Time` |
| string        | binary         | `[]byte`    |
| integer       | int64          | `int64`     |
| integer       | int32          | `int32`     |
| integer       | (none)         | `int`       |
| number        | float          | `float32`   |
| number        | (none)         | `float64`   |
| boolean       | —              | `bool`      |
| array         | —              | `[]T`       |
| object        | —              | `map[string]any` |
| `$ref`        | —              | `*RefType`  |


## Config resolution

`config.go` provides `ResolveConfig()` which resolves the Forgejo URL and API token with the following priority:

1. Explicit flag values (passed as function arguments)
2. Environment variables (`FORGE_URL`, `FORGE_TOKEN`)
3. Built-in defaults (`http://localhost:3000` for URL; no default for token)

`NewForgeFromConfig()` wraps this into a one-liner that returns an error if no token is available from any source.
