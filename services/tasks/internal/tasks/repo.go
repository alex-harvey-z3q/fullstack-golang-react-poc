package tasks

import (
	"context"
	"log"

	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/config"
	gen "github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/db/gen"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Defines the repository type:
//
// db holds the shared connection pool.
// qry holds sqlc’s generated query wrapper bound to that pool.
type Repo struct {
	db  *pgxpool.Pool
	qry *gen.Queries
}

// Constructor - like __init__ in Python
//
// Creates a new pgx pool using cfg.DatabaseURL.
// If connection fails, logs and exits the process (simple PoC behavior).
//
// Binds sqlc’s Queries to the pool and returns a *Repo.
func NewRepo(cfg config.Config) *Repo {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		// For this PoC we crash hard if DB connect fails.
		// In production we'd return an error instead of fatal logging.
		log.Fatalf("db connect: %v", err)
	}
	return &Repo{
		db:  pool,
		qry: gen.New(pool),
	}
}

// Close releases the connection pool resources.
// Call this at shutdown (main.go does `defer repo.Close()`).
func (r *Repo) Close() {
	r.db.Close()
}

// List fetches all tasks from the database.
// It calls the sqlc-generated query ListTasks(ctx),
// then maps the raw DB row structs into the domain Task model.
func (r *Repo) List(ctx context.Context) ([]Task, error) {
	// Run the sqlc-generated query (SELECT * FROM tasks ...).
	rows, err := r.qry.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	// Map sqlc's row structs into our domain model Task.
	out := make([]Task, 0, len(rows))
	for _, t := range rows {
		out = append(out, Task{
			ID:        t.ID,
			Title:     t.Title,
			Done:      t.Done,
			CreatedAt: t.CreatedAt.Time,
			UpdatedAt: t.UpdatedAt.Time,
		})
	}
	return out, nil
}
