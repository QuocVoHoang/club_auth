# Go Base Monolith

Base Go monolith service using Gin, clean architecture boundaries, PostgreSQL, Gorm, migrations, transaction manager, request IDs, and health/readiness endpoints.

## Development Setup

1. Create `.env` from `.env.example`.

2. Start PostgreSQL with Docker Compose: `docker compose up -d`

3. Install dev tools: `make tool`

4. Run the command `make migrate-up`

5. Run the command `make run`

## Useful Commands

- `make run`: Run the application locally.
- `make test`: Run all Go tests.
- `make test-race`: Run tests with the Go race detector.
- `make fmt`: Format all Go files with `gofmt`.
- `make vet`: Run `go vet` static checks.
- `make lint`: Run `golangci-lint`.
- `make migrate-up`: Apply all pending database migrations.
- `make migrate-down`: Roll back the latest database migration.
- `make migrate-create NAME=create_users_table`: Create a new migration file pair.
- `make migrate-version`: Show the current database migration version.
- `make migrate-force VERSION=1`: Force the migration version after a dirty migration state.
- `make up`: Start Docker Compose services in the background.
- `make down`: Stop Docker Compose services.

## Health Checks

- Liveness: `GET /api/healthz`
- Readiness: `GET /api/readyz`

`/api/readyz` verifies PostgreSQL connectivity.

## Swagger/OpenAPI
This service supports Swagger UI out of the box.

How to use:
- Start the service: `make run` (or via Docker Compose: `docker compose up`)
- Open your browser and navigate to: `http://localhost:<PORT>/api/v1/doc/index.html`
  - Replace `<PORT>` with the value from your `.env` (for example, 8080)

Regenerate API docs:
- If you add or change route annotations (Swaggo comments) in the code, regenerate the docs with
  - `make gen`

Notes:
- The Swagger endpoint is registered at `/api/v1/doc/*any` in `internal/framework/route/route.go` using `github.com/swaggo/gin-swagger`.
- Generated artifacts are stored in `internal/docs` (docs.go, swagger.json, swagger.yaml). Do not edit these manually; regenerate with `make gen`.
- Ensure `ENV`, `PORT` and database settings in `.env` are correct before running.
