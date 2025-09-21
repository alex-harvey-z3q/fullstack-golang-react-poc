# Fullstack Monorepo — Go + Postgres + React + Angular

A learning project that demonstrates a clean, layered architecture with a Go HTTP service, PostgreSQL, and **two** frontends:

- **React** (Vite) at `services/web/react` (port **5173**)
- **Angular 18** at `services/web/angular` (port **4200**)

Both frontends call the same REST API implemented by the Go service.

## Status

- **services/tasks** — Go + Postgres
  - REST: `/api/tasks` **GET** + **POST**, validated against OpenAPI.
  - GraphQL: schema scaffolded (no `/graphql` route yet).
- **services/web/react** — React + TypeScript
  - Lists/creates tasks via REST. Vite proxy → `:8081`.
- **services/web/angular** — Angular 18 (standalone)
  - Lists/creates tasks via REST. Angular CLI dev server → `:8081` via proxy.
- **services/users** — placeholder.

## High-level Architecture

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

## Where things live

- Backend service: `services/tasks`
- OpenAPI: `services/tasks/api/openapi.yaml`
- SQL + sqlc: `services/tasks/internal/db`
- React app: `services/web/react`
- Angular app: `services/web/angular`

## Getting Started

See **QUICK_START.md** for step-by-step instructions.

## Extending the API

See **MAINTAINER_GUIDE.md** for a contract-first workflow (OpenAPI → sqlc → repo → service → HTTP) and testing pattern (spec validity + runtime OAS conformance + unit tests).

## License

MIT.
