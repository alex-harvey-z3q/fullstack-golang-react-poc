package tasks

import (
	"context"
	"fmt"
	"time"

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

// NewRepo creates a pgx pool using a bounded context and verifies connectivity.
// It returns an error instead of exiting the process.
func NewRepo(parent context.Context, cfg config.Config) (*Repo, error) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}
	// Ensure we can reach the database now, using the same timeout.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db ping: %w", err)
	}
	return &Repo{
		db:  pool,
		qry: gen.New(pool),
	}, nil
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
