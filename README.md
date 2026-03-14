# go-vbs

Go implementation of a small vehicle booking system API.

The service exposes endpoints for:

- issuing JWT access and refresh tokens
- fetching a vehicle by UUID
- creating bookings for a date range
- listing bookings behind JWT auth

It uses SQLite for persistence, `sqlc` for query generation, `mockery` for test mocks, and `just` as the project task runner.

## Requirements

- Go 1.23+
- `just`
- GCC or another C toolchain compatible with `github.com/mattn/go-sqlite3`

Optional tools used by project tasks:

- `sqlc`
- `mockery`
- `goose`
- `wire`

You can install the Go-based tools with:

```bash
just install-tools
```

## Quick Start

Build the binary:

```bash
just build
```

Run the server:

```bash
just run
```

The server listens on `127.0.0.1:8000`.

Run the test suite:

```bash
just test
```

## Task Runner

The project uses `just` only.

Common tasks:

```bash
just build
just run
just test
just test-coverage
just sqlc
just mocks
just migrate-up
just migrate-status
```

See all available recipes with:

```bash
just --list
```

## Generated Code

Two parts of the repository are generated:

- `internal/repository/db.go`, `internal/repository/models.go`, `internal/repository/query.sql.go` from `internal/repository/query.sql` via `sqlc`
- `internal/booking/business/mocks/` from `internal/booking/business` interfaces via `mockery`

Normal `just build` runs `sqlc` first.
Normal `just test` runs `sqlc` and `mocks` first.

## API

Current routes:

- `GET /login`
- `GET /refresh`
- `GET /vehicles/{uuid}`
- `POST /bookings`
- `GET /bookings`

OpenAPI spec:

- `openapi.yaml`

Manual request examples:

- `requests.http`

Swagger UI is available in development mode after downloading the static assets:

```bash
just swagger-ui
GO_ENV=development ./go-vbs
```

Then open:

```text
http://localhost:8000/docs
```

## Database

The app uses a local SQLite database file:

- `vbs.db`

Schema and seed data live in:

- `internal/repository/migrations/`

Migration commands:

```bash
just migrate-up
just migrate-down
just migrate-status
just migrate-create add_something
```

## Project Layout

```text
cmd/go-vbs/       application entrypoint
internal/app/     dependency bootstrap
internal/auth/    authentication controller
internal/booking/ booking business logic, controllers, repositories, mocks
internal/customer/ customer domain and repository
internal/vehicle/ vehicle domain, use cases, repository
internal/http/    HTTP server, middleware, and handler glue
internal/repository/ SQLite connection, migrations, sqlc input and generated queries
```

## Architecture Notes

- Runtime wiring currently happens in `internal/app/application.go`.
- The `wire` recipe is still available, but it is not part of the normal build or test flow.
- Persistence uses `github.com/mattn/go-sqlite3`, so CGO support is required for builds.

## Status

The repository is set up so that:

- `just build` builds the application
- `just test` passes
- generated query code and mocks are regenerated through `just`
