package tasks

import (
	"context"
	"log"

	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/config"
	gen "github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/db/gen"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db  *pgxpool.Pool
	qry *gen.Queries
}

func NewRepo(cfg config.Config) *Repo {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil { log.Fatalf("db connect: %v", err) }
	return &Repo{db: pool, qry: gen.New(pool)}
}

func (r *Repo) Close() { r.db.Close() }

func (r *Repo) List(ctx context.Context) ([]Task, error) {
	rows, err := r.qry.ListTasks(ctx)
	if err != nil { return nil, err }
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
