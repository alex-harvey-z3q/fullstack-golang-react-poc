# Fullstack Monorepo

A learning project demonstrating a fullstack architecture with Go, Postgres, and React.

## Status

- `services/tasks`: Go + Postgres (REST + GraphQL)
- `web/react`: Frontend (services/web/react): React + TypeScript client that lists tasks.
- `services/users`: Go + Mongo (TODO)

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

## Prerequisites

- Docker & Docker Compose
- Go 1.23+
- Node.js 18+ and npm

## Running the stack
```
docker compose up --build # Start Postgres + tasks service
```

Backend runs at http://localhost:8081

Health check: curl http://localhost:8081/healthz

Tasks API: curl http://localhost:8081/api/tasks

## Database

Migrations are stored in services/tasks/internal/db/migrate.

```
# Apply migrations to dev DB
make migrate

# Apply migrations to test DB (used automatically in `make test`)
make migrate-test

# Connect to dev DB
docker compose exec -it postgres psql -U app -d tasks

# Connect to test DB
docker compose exec -it postgres psql -U app -d tasks_test
```

## Backend development

```
cd services/tasks

# Run unit tests (uses tasks_test DB)
make test

# Run server locally (without Docker)
make dev
```

## Frontend development

```
cd services/web/react
npm install
npm run dev
```

- React app runs at http://localhost:5173
 (Vite default).

- Fetches data from the backend at http://localhost:8081/api/tasks.

## API Documentation

The REST API contract is defined in:
services/tasks/api/openapi.yaml

This can be used with tools like Swagger UI or Postman to generate clients or validate requests.
Currently describes `GET /api/tasks`.

## License

MIT.
