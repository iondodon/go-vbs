# AGENTS.md

This document captures the durable engineering principles behind `go-vbs` so they can be reused in other Go services. It is written for coding agents and humans who need a concrete operating model, not a generic style guide.

Use this guide as the default unless the target repository already has stronger local conventions.

## Core Intent

Build small Go services with:

- explicit boundaries
- thin transport code
- business logic in use-case services
- persistence behind interfaces
- generated infrastructure kept as generated
- simple local tooling
- tests centered on service behavior

Favor straightforward code over framework-heavy abstractions.

## Architecture

Organize the service around domain modules under `internal/`.

Recommended shape:

```text
cmd/<app>/                 executable entrypoint
internal/app/              dependency bootstrap and application assembly
internal/http/             transport glue: server, middleware, handler adapters
internal/<domain>/domain/  core entities and domain value types
internal/<domain>/services/ contracts and use-case interfaces
internal/<domain>/services/<usecase>/ service implementations
internal/<domain>/controller/ transport DTOs and HTTP controllers
internal/<domain>/repository/ concrete persistence adapters
internal/repository/       shared SQL layer, migrations, generated query code
pkg/                       only if you truly expose reusable library APIs
```

Treat `cmd/` and `internal/` as the default. Do not put application code in `pkg/` unless another module must import it.

## Layer Responsibilities

### Domain

- Keep domain structs small and boring.
- Domain types represent business data, not SQL rows or HTTP payloads.
- Prefer explicit named types for constrained values such as enums.
- Keep domain packages dependency-light.

### Services

- Put service contracts in `internal/<domain>/services/`.
- Define interfaces from the consumer side. The service layer states what it needs from repositories or other domains.
- Implement each use case in its own focused package under `services/...`.
- Service implementations should hold dependencies in a `Service` struct and expose a `New(...)` constructor.
- Business rules, orchestration, validation, and cross-repository workflows belong here.

### Repositories

- Repositories are adapters from generated SQL/query code into domain objects.
- Repositories should implement service-layer interfaces, not define the application contract themselves.
- Add compile-time interface assertions:

```go
var _ bookingServices.BookingRepository = (*Repository)(nil)
```

- Keep SQL in source files consumed by code generation, then map generated result types into domain models manually.
- Use transactions through repository methods when a use case requires atomic writes.

### Controllers and HTTP

- Controllers should be thin. Their job is to:
  - read path/query/body input
  - map transport data to use-case arguments
  - start/commit/rollback transactions when the HTTP request is the transaction boundary
  - serialize the result
- Keep shared HTTP wrapper logic in one place. In this project that role is `internal/http/handler`.
- Keep middleware separate from controllers.
- Register routes centrally in the server package.

### Application Wiring

- Keep dependency construction explicit in `internal/app`.
- Constructor functions are preferred over global state.
- Compile-time DI is acceptable. This project uses `wire`, but the actual principle is explicit provider graphs, not mandatory Wire usage.
- Generated wiring may be committed if it is part of the normal build path.

## Dependency Rules

Follow these rules strictly:

- `controller` may depend on `services`, transport DTOs, and standard library HTTP primitives.
- `services` may depend on domain packages, service interfaces, and repository interfaces.
- `repository` may depend on generated SQL code and domain packages.
- `domain` should not depend on controller or repository packages.
- `http` infrastructure may depend on controllers and middleware, but business rules should not live there.

Important nuance from `go-vbs`:

- The current code sometimes passes `booking/controller.DatePeriodDTO` through service and repository interfaces.
- Do not treat that coupling as the ideal pattern to reproduce.
- For new projects, promote shared business values into a domain-owned type when the same concept crosses transport and service boundaries.

## Naming and Code Style

Mirror these code-shape conventions:

- Package names are lowercase and purpose-driven.
- Cross-package aliases are explicit when they improve clarity, for example `vehicleServices`, `bookingDomain`, `authController`.
- Constructors are named `New`.
- Main structs are usually named `Controller`, `Service`, or `Repository`.
- Pass `context.Context` as the first runtime argument after the receiver.
- Return domain values or pointers, not generated SQL structs.
- Keep methods small and direct.
- Prefer explicit field mapping over reflection or magic mappers.
- Use standard library facilities first: `net/http`, `log/slog`, `database/sql`, `context`.

Write code that is easy to trace from route to controller to use case to repository to SQL.

## Error Handling

- Return errors upward instead of swallowing them.
- Let the transport boundary decide how to turn errors into HTTP responses.
- Wrap errors when extra context matters, especially around parsing or mapping failures.
- Roll back transactions on service/controller write-path failures.
- Keep error messages concrete and tied to the business action.

## Transactions

Use the HTTP boundary or application boundary to start transactions for multi-step writes.

Pattern:

- controller begins `sql.Tx`
- controller calls service with `ctx` and `tx`
- service orchestrates writes through repositories
- repository uses generated query helpers with `WithTx(tx)`
- controller commits on success and rolls back on failure

Keep transaction ownership obvious. Avoid hidden transaction creation deep in repositories.

## Persistence and SQL

The `go-vbs` model for SQL-backed services is:

- schema changes live in migration files
- query intent lives in hand-written SQL
- generated code is derived from SQL, not edited manually
- repositories map generated rows into domain structs

For projects following this pattern:

- edit `query.sql`, not `query.sql.go`
- edit migrations, not generated models
- keep SQL queries named and small
- favor repository mapping code that is explicit even if repetitive

If using `sqlc`, commit the generated output when the repository expects local builds/tests without regeneration surprises.

## Generated Code Policy

In this project, generated artifacts include:

- `internal/repository/db.go`
- `internal/repository/models.go`
- `internal/repository/query.sql.go`
- `internal/app/wire_gen.go`
- `internal/booking/services/mocks/`

Rules:

- do not hand-edit generated files
- change the source file and regenerate
- keep generator config committed
- make generation part of the standard developer workflow

## Tooling Principles

`go-vbs` uses a minimal toolchain:

- Go toolchain for build and test
- `just` as the single task runner
- `sqlc` for query generation
- `mockery` for mocks
- `goose` for migrations
- `wire` for optional compile-time DI generation

Transferable principle:

- expose one obvious command surface for developers and agents
- automate generation before build/test
- keep infra tooling lightweight and local

If another repository prefers `make`, `task`, or plain scripts, preserve the same simplicity rather than forcing `just`.

## Testing Principles

Test the service layer first.

- Unit tests should target use-case services.
- Mock repository and cross-domain dependencies through generated mocks.
- Verify behavior and orchestration, not implementation trivia.
- Use `testify/assert` and mock expecters when they improve readability.
- Keep tests close to the use case package they exercise.

Observed pattern in `go-vbs`:

- `mockery` generates mocks from interfaces in `services/`
- tests construct the service directly
- tests assert business outcomes and dependency calls

Prefer this over controller-heavy tests for most business logic.

## API and Transport Conventions

- Define routes centrally in the HTTP server package.
- Keep request DTOs near the controller that owns them.
- Use JSON tags deliberately.
- Expose a machine-readable API contract when possible. This project uses `openapi.yaml`.
- Keep a manual request file such as `requests.http` for quick local verification.
- Mount developer-only documentation endpoints only in development mode.

## Logging and Runtime

- Use structured logging with `log/slog`.
- Log lifecycle events and business-significant actions.
- Keep startup and shutdown explicit in `main`.
- Close long-lived resources on shutdown.
- Prefer graceful server shutdown with context timeouts.

## Feature Delivery Workflow

When adding a new feature in a project following this model:

1. Add or refine domain types.
2. Define or update service interfaces in the consuming domain.
3. Implement a focused use-case service under `services/...`.
4. Add or update repository interfaces needed by the service.
5. Add SQL and migrations if persistence changes.
6. Regenerate `sqlc`, mocks, and DI code if applicable.
7. Implement repository adapters and explicit row-to-domain mapping.
8. Add or update controller DTOs and HTTP handlers.
9. Register the route in the server.
10. Add service-layer tests and then run the standard task commands.

## Agent Rules

If you are an agent operating in another Go repository and asked to follow `go-vbs` principles:

- preserve layered boundaries
- keep controllers thin
- put orchestration in services
- define interfaces where they are consumed
- keep repositories as adapters, not business logic containers
- do not edit generated files directly
- prefer explicit constructors and explicit wiring
- use the repo's single task runner for generation, build, and test
- add tests at the service layer when changing behavior

## Avoid Cargo-Culting These Local Quirks

These are present in `go-vbs`, but should not be copied blindly:

- transport DTOs reused outside the transport layer
- hard-coded JWT secrets
- sparse HTTP error translation
- repository mapping code that depends on SQLite-specific generated field shapes

Carry forward the architecture and workflow. Improve incidental implementation details when the target project can support a cleaner version.

## Done Criteria

Work is aligned with this guide when:

- the new behavior fits the existing layer boundaries
- generated artifacts are in sync with their sources
- build/test commands run through the project task surface
- business logic changes are covered by service-level tests
- routes, controllers, services, repositories, and SQL remain easy to trace end-to-end
