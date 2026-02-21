# Configuration

go-forge resolves its Forgejo URL and API token from multiple sources with a clear priority chain.

## Priority Order

```
1. Flags (highest)     — passed directly to NewForgeFromConfig or ResolveConfig
2. Environment variables — FORGE_URL, FORGE_TOKEN
3. Defaults (lowest)   — http://localhost:3000 for URL, no token default
```

## Environment Variables

| Variable      | Description                          | Default                  |
|---------------|--------------------------------------|--------------------------|
| `FORGE_URL`   | Base URL of the Forgejo instance     | `http://localhost:3000`  |
| `FORGE_TOKEN` | API token for authentication         | _(none — required)_      |

## Using NewForgeFromConfig

The simplest way to create a configured client:

```go
// Reads from FORGE_URL and FORGE_TOKEN environment variables.
f, err := forge.NewForgeFromConfig("", "")
if err != nil {
    log.Fatal(err) // "no API token configured" if FORGE_TOKEN is unset
}
```

With flag overrides:

```go
// Flag values take priority over env vars.
f, err := forge.NewForgeFromConfig("https://forge.example.com", "my-token")
```

## Using ResolveConfig Directly

For more control, resolve the config separately:

```go
url, token, err := forge.ResolveConfig(flagURL, flagToken)
if err != nil {
    log.Fatal(err)
}
// Use url and token as needed...
f := forge.NewForge(url, token)
```

## Using NewForge Directly

If you already have the URL and token:

```go
f := forge.NewForge("https://forge.lthn.ai", "my-api-token")
```

## Client Options

`NewForge` and `NewForgeFromConfig` accept variadic options:

```go
// Custom HTTP client (e.g. with timeouts, proxies, TLS config)
f := forge.NewForge(url, token,
    forge.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)

// Custom User-Agent header
f := forge.NewForge(url, token,
    forge.WithUserAgent("my-app/1.0"),
)
```

| Option             | Description                   | Default          |
|--------------------|-------------------------------|------------------|
| `WithHTTPClient`   | Set a custom `*http.Client`   | `http.DefaultClient` |
| `WithUserAgent`    | Set the `User-Agent` header   | `go-forge/0.1`  |

## Accessing the Low-Level Client

If you need direct HTTP access for custom endpoints:

```go
client := f.Client() // returns *forge.Client
err := client.Get(ctx, "/api/v1/custom/endpoint", &result)
```
