# Error Handling

All go-forge methods return errors wrapped as `*APIError` for HTTP error responses. This page covers the error types and how to handle them.

## APIError

When the Forgejo API returns an HTTP status >= 400, go-forge wraps the response as:

```go
type APIError struct {
    StatusCode int    // HTTP status code (e.g. 404, 403, 409)
    Message    string // Error message from the API response body
    URL        string // The request path
}
```

The `Error()` method formats as: `forge: /api/v1/repos/core/missing 404: not found`

## Sentinel Checks

Three helper functions test for common HTTP error statuses:

```go
repo, err := f.Repos.Get(ctx, forge.Params{"owner": "core", "repo": "missing"})
if err != nil {
    if forge.IsNotFound(err) {
        // 404 — resource does not exist
        fmt.Println("Repository not found")
    } else if forge.IsForbidden(err) {
        // 403 — insufficient permissions
        fmt.Println("Access denied")
    } else if forge.IsConflict(err) {
        // 409 — conflict (e.g. resource already exists)
        fmt.Println("Conflict")
    } else {
        // Other error (network, 500, etc.)
        log.Fatal(err)
    }
}
```

| Function       | Status Code | Typical Cause                     |
|----------------|-------------|-----------------------------------|
| `IsNotFound`   | 404         | Resource does not exist           |
| `IsForbidden`  | 403         | Insufficient permissions or scope |
| `IsConflict`   | 409         | Resource already exists           |

## Extracting APIError

Use `errors.As` to access the full `APIError`:

```go
var apiErr *forge.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("Status: %d\n", apiErr.StatusCode)
    fmt.Printf("Message: %s\n", apiErr.Message)
    fmt.Printf("URL: %s\n", apiErr.URL)
}
```

## Non-API Errors

Errors that are not HTTP responses (network failures, JSON decode errors, request construction failures) are returned as plain `error` values with descriptive prefixes:

| Prefix                      | Cause                              |
|-----------------------------|------------------------------------|
| `forge: marshal body:`      | Failed to JSON-encode request body |
| `forge: create request:`    | Failed to build HTTP request       |
| `forge: request GET ...:`   | Network error during request       |
| `forge: decode response:`   | Failed to JSON-decode response     |
| `forge: parse url:`         | Invalid URL in pagination          |
| `forge: read response body:`| Failed to read raw response body   |

## Context Cancellation

All methods accept `context.Context`. When the context is cancelled or times out, the underlying `http.Client` returns a context error:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

repos, err := f.Repos.ListOrgRepos(ctx, "core")
if err != nil {
    // Could be context.DeadlineExceeded or context.Canceled
    log.Fatal(err)
}
```
