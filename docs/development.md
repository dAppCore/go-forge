---
title: Development
description: Building, testing, linting, and contributing to go-forge.
---

# Development

This guide covers everything needed to build, test, and contribute to go-forge.


## Prerequisites

- **Go 1.26** or later
- **golangci-lint** (recommended for linting)
- A Forgejo instance and API token (only needed for manual/integration testing — the test suite uses `httptest` and requires no live server)


## Building

go-forge is a library, so there is nothing to compile for normal use. The only binary in the repository is the code generator:

```bash
go build ./cmd/forgegen/
```

The `core build` CLI can also produce cross-compiled binaries of the generator for distribution. Build configuration is in `.core/build.yaml`:

```bash
core build    # Builds forgegen for all configured targets
```


## Running tests

All tests use the standard `testing` package with `net/http/httptest` for HTTP stubbing. No live Forgejo instance is required.

```bash
# Run the full suite
go test ./...

# Run a specific test by name
go test -v -run TestClient_Good_Get ./...

# Run tests with race detection
go test -race ./...

# Generate a coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Alternatively, if you have the `core` CLI installed:

```bash
core go test
core go cov          # Generate coverage
core go cov --open   # Open coverage report in browser
```

### Test naming convention

Tests follow the `_Good`, `_Bad`, `_Ugly` suffix pattern:

- **`_Good`** — Happy-path tests confirming correct behaviour.
- **`_Bad`** — Expected error conditions (e.g. 404, 500 responses).
- **`_Ugly`** — Edge cases, panics, and boundary conditions.

Examples:
```
TestClient_Good_Get
TestClient_Bad_ServerError
TestClient_Bad_NotFound
TestClient_Good_ContextCancellation
TestResource_Good_ListAll
```


## Linting

The project uses `golangci-lint` with the configuration in `.golangci.yml`:

```bash
golangci-lint run ./...
```

Or with the `core` CLI:

```bash
core go lint
```

Enabled linters: `govet`, `errcheck`, `staticcheck`, `unused`, `gosimple`, `ineffassign`, `typecheck`, `gocritic`, `gofmt`.


## Formatting and vetting

```bash
gofmt -w .
go vet ./...
```

Or:

```bash
core go fmt
core go vet
```

For a full quality-assurance pass (format, vet, lint, and test):

```bash
core go qa           # Standard checks
core go qa full      # Adds race detection, vulnerability scanning, and security checks
```


## Regenerating types

When the Forgejo API changes (after a Forgejo upgrade), regenerate the types:

1. Download the updated `swagger.v1.json` from your Forgejo instance at `/swagger.json` and place it in `testdata/swagger.v1.json`.

2. Run the generator:
   ```bash
   go generate ./types/...
   ```

   Or manually:
   ```bash
   go run ./cmd/forgegen/ -spec testdata/swagger.v1.json -out types/
   ```

3. Review the diff. The generator produces deterministic output (types are sorted alphabetically within each file), so the diff will show only genuine API changes.

4. Run the test suite to confirm nothing is broken.

The `types/generate.go` file holds the `//go:generate` directive that wires everything together:

```go
//go:generate go run ../cmd/forgegen/ -spec ../testdata/swagger.v1.json -out .
```


## Adding a new service

To add coverage for a new Forgejo API domain:

1. **Create the service file** (e.g. `topics.go`). If the API follows a standard CRUD pattern, embed `Resource[T, C, U]`:

   ```go
   type TopicService struct {
       Resource[types.Topic, types.CreateTopicOption, types.EditTopicOption]
   }

   func newTopicService(c *Client) *TopicService {
       return &TopicService{
           Resource: *NewResource[types.Topic, types.CreateTopicOption, types.EditTopicOption](
               c, "/api/v1/repos/{owner}/{repo}/topics/{topic}",
           ),
       }
   }
   ```

   If the endpoints are heterogeneous, hold a `*Client` directly instead:

   ```go
   type TopicService struct {
       client *Client
   }
   ```

2. **Add action methods** for any operations beyond standard CRUD:

   ```go
   func (s *TopicService) ListRepoTopics(ctx context.Context, owner, repo string) ([]types.Topic, error) {
       path := fmt.Sprintf("/api/v1/repos/%s/%s/topics", owner, repo)
       return ListAll[types.Topic](ctx, s.client, path, nil)
   }
   ```

3. **Wire it up** in `forge.go`:
   - Add a field to the `Forge` struct.
   - Initialise it in `NewForge()`.

4. **Write tests** in `topics_test.go` using `httptest.NewServer` to stub the HTTP responses.

5. Every list method should provide both a `ListX` (returns `[]T`) and an `IterX` (returns `iter.Seq2[T, error]`) variant.


## Coding standards

- **UK English** in all comments and documentation: organisation, colour, centre, licence (noun).
- **`context.Context`** as the first parameter of every exported method.
- **Errors** wrapped as `*APIError` with `StatusCode`, `Message`, and `URL`.
- **Conventional commits**: `feat:`, `fix:`, `docs:`, `refactor:`, `chore:`.
- **Generated code** must not be edited by hand. All changes go through the swagger spec and the generator.
- **Licence**: EUPL-1.2. All contributions are licensed under the same terms.


## Commit messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
feat: add topic service for repository topics
fix: handle empty response body in DeleteWithBody
docs: update architecture diagram for pagination
refactor: extract rate limit parsing into helper
chore: regenerate types from Forgejo 10.1 swagger spec
```

Include the co-author trailer:

```
Co-Authored-By: Virgil <virgil@lethean.io>
```


## Project structure at a glance

```
.
├── .core/
│   ├── build.yaml       Build targets for the forgegen binary
│   └── release.yaml     Release publishing configuration
├── .golangci.yml        Linter configuration
├── cmd/forgegen/        Code generator (swagger -> Go types)
├── testdata/            swagger.v1.json spec file
├── types/               229 generated types (36 files)
├── *.go                 Library source (client, services, pagination)
└── *_test.go            Tests (httptest-based, no live server needed)
```
