package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/config"
	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/tasks"
)

func main() {
	// Read config from services/tasks/internal/config/config.go.
	cfg := config.Load()

	// The Repository design pattern
	// (https://martinfowler.com/eaaCatalog/repository.html)
	//
	// "A Repository mediates between the domain and data mapping layers, acting like
	// an in-memory domain object collection. Client objects construct query specifications
	// declaratively and submit them to Repository for satisfaction. Objects can be added
	// to and removed from the Repository, as they can from a simple collection of objects,
	// and the mapping code encapsulated by the Repository will carry out the appropriate
	// operations behind the scenes. Conceptually, a Repository encapsulates the set of
	// objects persisted in a data store and the operations performed over them, providing
	// a more object-oriented view of the persistence layer. Repository also supports the
	// objective of achieving a clean separation and one-way dependency between the domain
	// and data mapping layers."
	//
	// This reads from services/tasks/internal/tasks/repo.go. It creates a pgxpool to
	// talk to the Postgres database.
	//
	// pgxpool is the connection pool implementation that comes with pgx, a popular Go
	// driver for PostgreSQL.
	//
	repo, err := tasks.NewRepo(context.Background(), cfg)

	if err != nil {
		log.Fatalf("repository init: %v", err)
	}

	// defer repo.Close() is Go’s way of guaranteeing that the database pool is cleaned
	// up when the program ends, avoiding resource leaks and ensuring graceful shutdowns.
	defer repo.Close()

	// Builds the domain/business layer. The service coordinates use-cases and calls into
	// the repo. (For this PoC it’s thin: it just forwards to the repo.)
	svc := tasks.NewService(repo)

	// Creates the HTTP router with default middleware (logger + recovery).
	//
	// Logger   → logs each incoming HTTP request (method, path, status code, latency, etc.).
	// Recovery → catches any panics inside the handlers, logs the stack trace, and
	//   responds with 500 Internal Server Error instead of crashing the whole process.
	//
	// A web service 'router' is the component that:
	// - Receives an incoming HTTP request (e.g. GET /api/tasks)
	// - Matches it against a set of defined rules (like “if it’s a GET and the path
	//   is /api/tasks, call this function”).
	// - Dispatches the request to the correct handler — i.e. the Go function that produces
	//   the response.
	//
	r := gin.Default()

	// Health check endpoint. Returns {"ok": true} with 200 status.
	// Useful for Docker/Kubernetes probes or for a quick curl check.
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Group all API endpoints under the /api prefix.
	api := r.Group("/api")

	// Register the task-related routes (e.g., GET /api/tasks).
	// We inject `svc` here: this allows tests to inject fakes instead of real DB.
	tasks.RegisterRoutes(api, svc)

	log.Printf("listening on :%s", cfg.Port)

	// Start the HTTP server — this blocks forever, handling requests until shutdown.
	// If the port is unavailable, Run() returns an error immediately.
	// On error, log.Fatalf prints the message and exits the process.
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
