# Repository Guidelines

This repository is a Go (1.24.x) backend using Echo v4, oapi-codegen, PostgreSQL (with optional SQLite), structured logging (slog), and golang-migrate.

## Communication Language

**IMPORTANT: Always respond in Japanese (日本語) when working in this repository.** All explanations, error messages, and communication should be in Japanese to match the project's primary language and documentation.

## Project Structure & Module Organization
- `cmd/api/`: Application entrypoint.
- `internal/api/openapi.yaml`: OpenAPI spec (human‑reviewed; do not change casually).
- `internal/api/gen/`: oapi-codegen output (do not edit by hand).
- `internal/api/handler/`: Echo handlers (HTTP → service orchestration).
- `internal/domain/{model,service,repository}/`: Domain entities, business logic, repository interfaces.
- `internal/infra/{db,repository}/`: DB connections/Tx and SQL implementations.
- `internal/policy/`: Scope and access decisions.
- `internal/pkg/{errors,response}/`: Error types and Problem response mapping.
- `internal/{logger,config}/`: Logging and configuration.
- `migrations/`: SQL migrations. `scripts/`: helper scripts.

## Build, Test, and Development Commands
- Codegen: `go generate && git diff --exit-code internal/api/gen` (keeps generated code in sync).
- Build: `go build ./cmd/api`  Run: `go run ./cmd/api` (set config envs first).
- Test: `go test ./...`  Vet: `go vet ./...`
- Migrate (example): `migrate -path migrations -database "$DSN" up`
  - `DSN` starting with `postgres://` uses PostgreSQL; otherwise SQLite.

## Coding Style & Naming Conventions
- Use standard Go formatting: `go fmt ./...`; prefer simple, readable code.
- Receivers: one letter (`r`, `s`, `h`). Interfaces: `NameService`, `NameRepository`.
- First arg is `ctx context.Context`. Avoid `panic`; return errors and map to HTTP in handlers.
- Logging keys: `req_id`, `path`, `status`, `latency_ms`. Never log secrets.
- SQL: always use placeholders; no string concatenation for conditions.

## Testing Guidelines
- Focus on service/policy branches (target ≥60%); CRUD for repositories; handlers: success + 400/404 paths.
- Prefer table-driven tests; use `t.Helper()` and fail-fast assertions when appropriate.
- Run: `go test ./...`; add DB-backed tests only where valuable and deterministic.

## Commit & Pull Request Guidelines
- Commits: imperative mood, concise scope, reference issues (e.g., `Fix: scope check for RUNTIME_REQUIRED`).
- PRs: clear description, linked issue, reproduction/verification steps, and note schema or OpenAPI changes.
- Include a short change summary, e.g.:
  
  ```
  # Change Summary
  - Added: handler/ListOssComponents pagination
  - Modified: repository tag filter
  - DB: No schema change
  - OpenAPI: No change
  ```

## Security & Configuration Tips
- Configure via env (no secrets in code). Respect DSN behavior (Postgres vs SQLite). Arrays: PG `TEXT[]`, SQLite JSON string.
- Basic auth only; plan for future LDAP/JWT. Keep logs free of passwords/tokens.

