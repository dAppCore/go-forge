# go-forge

Full-coverage Go client for the [Forgejo](https://forgejo.org/) API.

**Module:** `forge.lthn.ai/core/go-forge`

## Features

- 18 service modules covering repositories, issues, pull requests, organisations, teams, users, admin, branches, releases, labels, webhooks, notifications, packages, actions, contents, wiki, commits, and miscellaneous endpoints
- Generic `Resource[T, C, U]` pattern providing free CRUD for most services
- 229 Go types generated from the Forgejo Swagger 2.0 specification
- Automatic pagination with `ListAll` and `ListPage` generics
- Config resolution from flags, environment variables, and defaults
- Typed error handling with `IsNotFound`, `IsForbidden`, `IsConflict` helpers

## Install

```bash
go get forge.lthn.ai/core/go-forge
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "forge.lthn.ai/core/go-forge"
)

func main() {
    // Create a client — reads FORGE_URL and FORGE_TOKEN from environment.
    f, err := forge.NewForgeFromConfig("", "")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get the authenticated user.
    me, err := f.Users.GetCurrent(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Hello, %s!\n", me.Login)

    // List repositories for an organisation.
    repos, err := f.Repos.ListOrgRepos(ctx, "core")
    if err != nil {
        log.Fatal(err)
    }
    for _, r := range repos {
        fmt.Println(r.FullName)
    }
}
```

## Documentation

- [Architecture](Architecture) — Generic Resource pattern, codegen pipeline, service design
- [Services](Services) — All 18 services with method signatures
- [Code Generation](Code-Generation) — How types are generated from the Swagger spec
- [Configuration](Configuration) — Environment variables, flags, defaults
- [Error Handling](Error-Handling) — APIError, sentinel checks, error wrapping
- [Development](Development) — Building, testing, contributing
