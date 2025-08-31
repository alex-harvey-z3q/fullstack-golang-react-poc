# Fullstack Monorepo

- `services/tasks`: Go + Postgres (REST + GraphQL)
- `services/users`: Go + Mongo (TODO)
- `web/react`: React + TypeScript (TODO)

# Flow

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

# Development

To login to the database and test database:

```
% docker compose exec -it postgres psql -U app -d tasks
% docker compose exec -it postgres psql -U app -d tasks_test
```
