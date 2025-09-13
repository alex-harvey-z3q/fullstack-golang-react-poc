# Fullstack Monorepo

A learning project demonstrating a fullstack architecture with Go, Postgres, and React.

## Status

- **services/tasks** — Go + Postgres
  - **REST:** Implemented (`/api/tasks` GET + POST), validated against `openapi.yaml`.
  - **GraphQL:** Schema & gqlgen config scaffolded, **not wired** into the server yet (no `/graphql` route or resolvers).

- **services/web/react** — React + TypeScript client
  - Lists tasks from the REST API; Vite dev server with `/api` proxy → `:8081`.

- **services/users** — (planned)
  - Placeholder; not implemented.

## Architecture and flow

```
Client (browser/curl)
        │  GET /api/tasks
        ▼
Gin Router (r := gin.Default)
  /api group → RegisterRoutes(api, svc)
        │  matches GET /api/tasks
        ▼
HTTP Handler (internal/tasks/http.go)
  - Extract ctx := c.Request.Context()
  - Call svc.List(ctx)
        │
        ▼
Service Layer (internal/tasks/service.go)
  - Orchestrates domain logic
  - Delegates to repo.List(ctx)
        │
        ▼
Repository (internal/tasks/repo.go)
  - Uses pgxpool (DB connection pool)
  - Calls sqlc-generated code (internal/db/gen)
        │
        ▼
PostgreSQL
  - Executes SQL from internal/db/queries/*.sql
  - Rows → sqlc models
        │
        ▼
Repository
  - Map sqlc models → domain Task
        │
        ▼
Service
  - Return []Task
        │
        ▼
HTTP Handler
  - c.JSON(200, []Task)
        │
        ▼
Client
  - Receives JSON array of tasks
```

## Quick Start
See [QUICK_START.md](QUICK_START.md) for step-by-step setup and workflow.

## Maintainers
See [MAINTAINER_GUIDE.md](MAINTAINER_GUIDE.md) for details on extending the API, running migrations, and contributing changes.

## API Documentation

The REST API contract is defined in:
services/tasks/api/openapi.yaml

This can be used with tools like Swagger UI or Postman to generate clients or validate requests.

## License

MIT.
