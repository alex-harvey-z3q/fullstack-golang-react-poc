package tasks

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/config"
)

func TestRepo_List(t *testing.T) {
	t.Parallel()

	// Use compose Postgres by default; override with TEST_DATABASE_URL if needed.
	dsn := getenv("TEST_DATABASE_URL", "postgres://app:app@localhost:5432/tasks?sslmode=disable")

	ctx := context.Background()

	// Connect directly to run migration + seed (simple PoC approach).
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	// Apply migration (idempotent)
	mig := filepath.Join("..", "db", "migrate", "001_init.sql")
	if err := applySQL(ctx, pool, mig); err != nil {
		t.Fatalf("apply migration: %v", err)
	}

	// Seed a row
	_, err = pool.Exec(ctx, `INSERT INTO tasks (title, done) VALUES ('from test', false) ON CONFLICT DO NOTHING`)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Build the repo using your real constructor (same DSN).
	cfg := config.Config{DatabaseURL: dsn}
	repo := NewRepo(cfg)
	defer repo.Close()

	// Exercise the method under test
	items, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("expected at least 1 task, got 0")
	}
}

func applySQL(ctx context.Context, pool *pgxpool.Pool, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, string(b))
	return err
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
