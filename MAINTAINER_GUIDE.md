# Maintainer Guide — Extending the Tasks API (with React & Angular clients)

This guide explains how to add new endpoints and features while staying consistent with the current architecture and the OpenAPI contract. It also notes where React and Angular touch the API.

- Backend: `services/tasks` (Go, Gin, Postgres, sqlc)
- OpenAPI: `services/tasks/api/openapi.yaml`
- Frontends:
  - React: `services/web/react`
  - Angular: `services/web/angular`

---

## Checklist: Add a New Endpoint

1. **Design the contract first**
   - Edit `services/tasks/api/openapi.yaml`.
   - Declare paths, methods, request/response schemas, and status codes.
   - Keep error responses consistent: `{"error": string}` via the shared `Error` schema.
   - Validate the spec (tests will also validate it).

2. **Write/adjust SQL**
   - Add queries under `services/tasks/internal/db/queries/*.sql`.
   - Use sqlc annotations (`-- name:`, `:one`, `:many`, `:exec`).
   - Generate code:
     ```bash
     cd services/tasks
     make sqlc
     ```
   - If schema changes are needed, add a migration under `internal/db/migrate/` and apply it:
     ```bash
     make migrate
     make migrate-test
     ```

3. **Repository layer**
   - Implement repo methods in `services/tasks/internal/tasks/repo.go` using `gen.Queries`.
   - Map sqlc models to domain models in `services/tasks/internal/tasks/model.go`.

4. **Service layer**
   - Add orchestration methods in `services/tasks/internal/tasks/service.go`.
   - Keep the interface between HTTP and service **narrow** (capability interfaces like `TaskLister`, `taskCreator`, etc.).

5. **HTTP layer**
   - Extend `RegisterRoutes` in `services/tasks/internal/tasks/http.go`.
   - Introduce a small, focused interface for each capability (e.g., `TaskGetter`, `TaskUpdater`, `TaskDeleter`).
   - Validate inputs, propagate `c.Request.Context()`, and return JSON errors as `{"error":"…"}` with proper status codes.

6. **Tests (required)**
   - **Spec validity test:** update `services/tasks/api/openapi_spec_test.go` to assert the new path/method exists.
   - **Runtime OAS tests:** add tests in `services/tasks/internal/tasks/openapi_runtime_test.go` that:
     1. hit the live handler via `httptest`, and
     2. validate the real response against `openapi.yaml` using `openapi3filter.ValidateResponse`.
   - **Handler tests:** add unit tests in `services/tasks/internal/tasks/http_test.go` using fakes that satisfy the new capability interface.
   - **Repo tests:** seed the test DB and exercise new repo methods (see `repo_test.go` for the pattern).

7. **Frontends (React & Angular)**
   - **React:** update `services/web/react/src/api.ts` with a new wrapper and use it from components.
   - **Angular:** update `services/web/angular/src/app/task.service.ts` with a new method and use it from components.
   - Surface errors consistently (HTTP status + body) as done for `createTask`.

8. **Housekeeping**
   - `go mod tidy`
   - End-to-end check:
     ```bash
     docker compose up --build
     ```

---

## Patterns & Examples

### Example: GET `/api/tasks/{id}`

**OpenAPI (snippet):**
```yaml
paths:
  /api/tasks/{id}:
    get:
      summary: Get task by id
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: integer, format: int32 }
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema: { $ref: "#/components/schemas/Task" }
        "400":
          description: Bad Request
          content:
            application/json:
              schema: { $ref: "#/components/schemas/Error" }
        "404":
          description: Not Found
          content:
            application/json:
              schema: { $ref: "#/components/schemas/Error" }
        "500":
          description: Internal Server Error
          content:
            application/json:
              schema: { $ref: "#/components/schemas/Error" }
```

**SQL (`internal/db/queries/tasks.sql`):**
```sql
-- name: GetTask :one
SELECT id, title, done, created_at, updated_at
FROM tasks
WHERE id = $1;
```

**Repo (`internal/tasks/repo.go`):**
```go
func (r *Repo) Get(ctx context.Context, id int32) (Task, error) {
  row, err := r.qry.GetTask(ctx, id)
  if err != nil {
    return Task{}, err
  }
  return Task{
    ID:        row.ID,
    Title:     row.Title,
    Done:      row.Done,
    CreatedAt: row.CreatedAt.Time,
    UpdatedAt: row.UpdatedAt.Time,
  }, nil
}
```

**Service (`internal/tasks/service.go`):**
```go
func (s *Service) Get(ctx context.Context, id int32) (Task, error) {
  return s.repo.Get(ctx, id)
}
```

**HTTP (`internal/tasks/http.go`):**
```go
type TaskGetter interface {
  Get(ctx context.Context, id int32) (Task, error)
}

r.GET("/tasks/:id", func(c *gin.Context) {
  g, ok := svc.(TaskGetter)
  if !ok {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "get not supported"})
    return
  }

  id64, err := strconv.ParseInt(c.Param("id"), 10, 32)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
    return
  }

  t, err := g.Get(c.Request.Context(), int32(id64))
  if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
      c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
      return
    }
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, t)
})
```

**Handler test (pattern):**
```go
type fakeGetter struct{}
func (f *fakeGetter) List(ctx context.Context) ([]Task, error) { return nil, nil } // to satisfy TaskLister required by RegisterRoutes
func (f *fakeGetter) Get(ctx context.Context, id int32) (Task, error) {
  now := time.Now()
  return Task{ID: id, Title: "One", Done: false, CreatedAt: now, UpdatedAt: now}, nil
}

// In test:
r := gin.New()
api := r.Group("/api")
RegisterRoutes(api, &fakeGetter{})
req := httptest.NewRequest(http.MethodGet, "/api/tasks/1", nil)
rec := httptest.NewRecorder()
r.ServeHTTP(rec, req)
// assert 200, decode body, etc.
```

---

## Useful Commands

```bash
# Run the stack
make up           # or: make up-d
make down
make logs

# Server locally
cd services/tasks && make dev

# DB tasks
cd services/tasks
make migrate
make migrate-test
make reset-db

# Codegen
cd services/tasks && make sqlc

# Tests
make lint
make test-backend
```
