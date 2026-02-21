# CLAUDE.md — go-forge

## Overview

Full-coverage Go client for the Forgejo API (450 endpoints). Uses a generic `Resource[T,C,U]` pattern for CRUD operations and types generated from `swagger.v1.json`.

**Module:** `forge.lthn.ai/core/go-forge`

## Build & Test

```bash
go test ./...                          # Run all tests
go test -v -run TestName ./...         # Run single test
go generate ./types/...                # Regenerate types from swagger spec
go run ./cmd/forgegen/ -spec testdata/swagger.v1.json -out types/  # Manual codegen
```

## Architecture

- `client.go` — HTTP client with auth, context.Context, error handling
- `pagination.go` — Generic paginated requests (ListPage, ListAll)
- `params.go` — Path variable resolution ({owner}/{repo} → values)
- `resource.go` — Generic Resource[T, C, U] for CRUD (covers 91% of endpoints)
- `forge.go` — Top-level Forge client aggregating all services
- `config.go` — Config resolution: flags → env (FORGE_URL, FORGE_TOKEN) → defaults
- `types/` — Generated Go types from swagger.v1.json (229 types)
- `cmd/forgegen/` — Code generator: swagger.v1.json → types/*.go

## Services

18 service structs, each embedding Resource[T,C,U] for CRUD + hand-written action methods:

repos, issues, pulls, orgs, users, teams, admin, branches, releases, labels, webhooks, notifications, packages, actions, contents, wiki, commits, misc

## Test Naming

Tests use `_Good`, `_Bad`, `_Ugly` suffix pattern:
- `_Good`: Happy path
- `_Bad`: Expected errors
- `_Ugly`: Edge cases/panics

## Coding Standards

- All methods accept `context.Context` as first parameter
- Errors wrapped as `*APIError` with StatusCode, Message, URL
- UK English in comments (organisation, colour, etc.)
- `Co-Authored-By: Virgil <virgil@lethean.io>` in commits

## Implementation Plan

See `/Users/snider/Code/host-uk/core/docs/plans/2026-02-21-go-forge-plan.md` for the full 20-task plan.

## Forge Remote

```bash
git remote add forge ssh://git@forge.lthn.ai:2223/core/go-forge.git
```
