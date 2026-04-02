# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Full-coverage Go client for the Forgejo API (~450 endpoints). Uses a generic `Resource[T,C,U]` pattern for CRUD operations and types generated from `swagger.v1.json`.

**Module:** `dappco.re/go/core/forge` (Go 1.26, depends on `go-io` and `go-log`)

## Build & Test

```bash
go test ./...                          # Run all tests
go test -v -run TestName ./...         # Run single test
go generate ./types/...                # Regenerate types from swagger spec
go run ./cmd/forgegen/ -spec testdata/swagger.v1.json -out types/  # Manual codegen
```

No linter or formatter is configured beyond standard `go vet`.

## Architecture

The library is a flat package (`package forge`) with a layered design:

1. **`client.go`** — Low-level HTTP client (`Client`). All requests go through `doJSON()`. Handles auth (token header), JSON marshalling, rate-limit tracking, and error parsing into `*APIError`. Also provides `GetRaw`/`PostRaw` for non-JSON responses.

2. **`resource.go`** — Generic `Resource[T, C, U]` struct (T=response type, C=create options, U=update options). Provides List, ListAll, Iter, Get, Create, Update, Delete. Covers ~91% of endpoints.

3. **`pagination.go`** — Generic `ListPage[T]`, `ListAll[T]`, `ListIter[T]` functions. Uses `X-Total-Count` header for pagination. `ListIter` returns `iter.Seq2[T, error]` for Go 1.26 range-over-func.

4. **`params.go`** — `Params` map type and `ResolvePath()` for substituting `{placeholders}` in API paths.

5. **`config.go`** — Config resolution: flags > env (`FORGE_URL`, `FORGE_TOKEN`) > defaults (`http://localhost:3000`).

6. **`forge.go`** — Top-level `Forge` struct aggregating all 19 service fields. Created via `NewForge(url, token)` or `NewForgeFromConfig(flagURL, flagToken)`.

7. **Service files** (`repos.go`, `issues.go`, etc.) — Each service struct embeds `Resource[T,C,U]` for standard CRUD, then adds hand-written action methods (e.g. `Fork`, `Pin`, `MirrorSync`). 19 services total: repos, issues, pulls, orgs, users, teams, admin, branches, releases, labels, webhooks, notifications, packages, actions, contents, wiki, commits, misc, activitypub.

8. **`types/`** — Generated Go types from `swagger.v1.json` (229 types). The `//go:generate` directive lives in `types/generate.go`. **Do not hand-edit generated type files** — modify `cmd/forgegen/` instead.

9. **`cmd/forgegen/`** — Code generator: parses swagger spec (`parser.go`), extracts types and CRUD pairs, generates Go files (`generator.go`).

## Test Patterns

- Tests use `httptest.NewServer` with inline handlers — no mocks or external services
- Test naming uses `_Good`, `_Bad`, `_Ugly` suffixes (happy path, expected errors, edge cases/panics)
- Service-level tests live in `*_test.go` files alongside each service

## Coding Standards

- All methods accept `context.Context` as first parameter
- Errors wrapped as `*APIError` with StatusCode, Message, URL; use `IsNotFound()`, `IsForbidden()`, `IsConflict()` helpers
- Internal errors use `coreerr.E()` from `go-log` (imported as `coreerr "dappco.re/go/core/log"`), never `fmt.Errorf` or `errors.New`
- File I/O uses `go-io` (imported as `coreio "dappco.re/go/core/io"`), never `os.ReadFile`/`os.WriteFile`/`os.MkdirAll`
- UK English in comments (organisation, colour, etc.)
- `Co-Authored-By: Virgil <virgil@lethean.io>` in commits
- `Client` uses functional options pattern (`WithHTTPClient`, `WithUserAgent`)

## Adding a New Service

1. Create `newservice.go` with a struct embedding `Resource[T, C, U]`
2. Add action methods that call `s.client.Get/Post/Patch/Put/Delete` directly
3. Wire it into `Forge` struct in `forge.go` and `NewForge()`
4. Add tests in `newservice_test.go`

## Forge Remote

```bash
git remote add forge ssh://git@forge.lthn.ai:2223/core/go-forge.git
```
