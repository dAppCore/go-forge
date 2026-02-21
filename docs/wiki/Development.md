# Development

## Prerequisites

- Go 1.21+ (generics support required)
- Access to a Forgejo instance for integration testing (optional)

## Building

```bash
go build ./...
```

## Running Tests

```bash
# All tests
go test ./...

# Verbose output
go test -v ./...

# Single test
go test -v -run TestClient_Good_Get ./...

# With race detection
go test -race ./...
```

## Test Naming Convention

Tests use a `_Good`, `_Bad`, `_Ugly` suffix pattern:

| Suffix  | Purpose                              |
|---------|--------------------------------------|
| `_Good` | Happy path — expected success        |
| `_Bad`  | Expected errors — invalid input, 404 |
| `_Ugly` | Edge cases — panics, empty data      |

## Regenerating Types

When the Forgejo Swagger spec is updated:

```bash
# Using go generate
go generate ./types/...

# Or directly
go run ./cmd/forgegen/ -spec testdata/swagger.v1.json -out types/
```

Then verify:

```bash
go build ./...
go test ./...
go vet ./...
```

## Adding a New Service

1. Create `servicename.go` with the service struct:
   - If CRUD applies: embed `Resource[T, C, U]` with the appropriate types.
   - If endpoints are heterogeneous: use a plain `client *Client` field.

2. Add a constructor: `func newServiceNameService(c *Client) *ServiceNameService`

3. Add action methods for non-CRUD endpoints.

4. Create `servicename_test.go` with tests following the `_Good`/`_Bad`/`_Ugly` convention.

5. Wire the service in `forge.go`:
   - Add the field to the `Forge` struct.
   - Add `f.ServiceName = newServiceNameService(c)` in `NewForge`.

6. Run all tests: `go test ./...`

## Adding a New Action Method

Action methods are hand-written methods on service structs for endpoints that don't fit the CRUD pattern.

```go
func (s *ServiceName) ActionName(ctx context.Context, params...) (returnType, error) {
    path := fmt.Sprintf("/api/v1/path/%s/%s", param1, param2)
    var out types.ReturnType
    if err := s.client.Post(ctx, path, body, &out); err != nil {
        return nil, err
    }
    return &out, nil
}
```

All methods must:
- Accept `context.Context` as the first parameter
- Return `error` as the last return value
- Use `fmt.Sprintf` for path construction with parameters

## Project Conventions

- **UK English** in comments (organisation, colour, licence)
- **`context.Context`** as the first parameter on all methods
- Errors wrapped as `*APIError` for HTTP responses
- Generated code has `// Code generated` headers and must not be edited manually
- Commit messages use conventional commits (`feat:`, `fix:`, `docs:`)

## Project Structure

```
go-forge/
  *.go              Core library (client, pagination, resource, services)
  *_test.go         Tests for each module
  cmd/forgegen/     Code generator
  types/            Generated types (do not edit)
  testdata/         Swagger specification
  docs/             Documentation
```

## Forge Remote

```bash
git remote add forge ssh://git@forge.lthn.ai:2223/core/go-forge.git
git push -u forge main
```
