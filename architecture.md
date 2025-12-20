# Architecture & Design

This document explains the project’s structure, layered design, core patterns, and how dependency injection is used to keep the code testable and maintainable.

## Goals

- Simple, testable service with clear boundaries.
- Layers that isolate concerns (HTTP, business logic, persistence).
- Interfaces at seams for inversion and mocking.
- Convention-driven, minimal boilerplate for CRUD.

## Project Structure

- `cmd/` – CLI entrypoints (`server`, `migrate`).
- `internal/controller/` – HTTP handlers, request validation, routes.
- `internal/service/` – Business logic; orchestrates repositories.
- `internal/repository/` – Data access abstractions with GORM implementation.
- `internal/model/` – Domain types and cross-layer validators.
- `internal/server/` – Application composition (dependency wiring), middleware.
- `internal/db/` – DB connection, migrations runner.
- `internal/utils/` – Logging and cross-cutting utilities.
- `migrations/` – SQL migrations (DDL/DML).

## Layered Architecture

- Controller
  - Binds and validates requests.
  - Translates errors to HTTP responses.
  - No business rules or persistence logic.
- Service
  - Encodes business rules and orchestration.
  - Depends on repository interfaces.
  - No HTTP or framework-specific logic.
- Repository
  - Encapsulates persistence with GORM.
  - Exposes an interface consumed by the service.
- Model
  - Domain types (`User`) and cross-layer validators (e.g., `validUserName`).
- Server
  - Composes dependencies and registers routes.
  - Sets global middleware (CORS, logging).

Dependency flow: `Repository → Service → Controller`. Composition occurs in `internal/server`, not inside handlers or services.

## Design Patterns

- Repository Pattern
  - Abstracts persistence behind an interface (`UserRepository`).
  - Allows swapping implementations and simplifies testing.
  - Keeps GORM specifics out of business logic.
- Dependency Inversion
  - Services depend on repository interfaces, not concretions or GORM.
  - High-level policy modules (services) do not depend on low-level details (ORM).
- Dependency Injection
  - Concrete implementations are created centrally and injected where needed.
  - Improves testability and reduces coupling.
- DTOs and Validation Strategy
  - Controllers use request DTOs with declarative validation tags.
  - Custom validation tags registered once and reused across endpoints.
- Error Mapping
  - Internal errors translated to typed HTTP responses (`internal/errors/codes.go`).

## Dependency Injection: Usage in This Project

DI is performed in the server layer, which is responsible for wiring repository → service → controller and initializing cross-cutting components like validation and middleware.

Example wiring (simplified) from `internal/server/api.go`:

```go
// initUserController creates and configures the user handler with its dependencies
func (s *txnAPIServer) initUserController() controller.UserHandler {
    // Repository → Service → Controller
    userRepo := repository.NewUserRepository(s.db)
    userService := service.NewUserService(userRepo)
    userController := controller.NewUserHandler(userService)
    return userController
}

// setupRoutes registers the routes and validator
func (s *txnAPIServer) setupRoutes(e *echo.Echo) {
    e.Validator = controller.NewCustomValidator()
    api := e.Group("/api/v1")
    controller.InitRoutes(api, s.initUserController())
}
```

- Repositories are constructed with `*gorm.DB`.
- Services accept repository interfaces in their constructors.
- Controllers accept service interfaces in their constructors.
- The Echo validator is created once and attached to the engine (`NewCustomValidator`).

## Validation

- Echo integrates `go-playground/validator` for struct tags.
- Custom tag `validUserName` is registered in `internal/controller/validator.go`.
- The validator delegates to model-level function references, keeping rules close to domain types.

Example registration:

```go
// internal/controller/validator.go
v := validator.New()
_ = v.RegisterValidation("validUserName", model.IsValidUserName)
```

## Testing & CI (TDD)

- Controller tests exercise binding/validation against a real DB.
- Test DB provisioning:
  - `internal/db.NewTestDB()` connects to `postgres`, auto-creates `user_test` if missing, then migrates.
- CI pipeline (`.github/workflows/test.yml`) runs `go test ./...` on PRs/commits.
- TDD workflow:
  - Write a failing test for new behavior.
  - Implement minimal code to pass the test.
  - Refactor with tests guarding behavior.

## Migrations

- GORM auto-migrates models on startup.
- SQL migrations (DDL/DML) in `migrations/` executed by `internal/db/migration.go`.
- Files are sorted and applied deterministically.

## Swagger & API

- Swagger annotations present in handlers; server can optionally host Swagger.
- Base path and schemes configured for local development.

## Logging

- Structured request logging via Echo middleware (`internal/server/log.go`).
- Global logger initialized in `internal/utils/logger.go`.

## Configuration

- `config.yaml` for runtime configuration; `config.test.yaml` for tests.
- Database DSN and ports are loaded via the CLI `server` command.

## Naming Conventions

- Files and types reflect the User domain (no `transaction` naming).
- Interfaces use clear names (`UserRepository`, `UserService`).