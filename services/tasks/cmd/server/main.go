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

	// The Repository design pattern is a structural pattern that abstracts data access,
	// providing a centralised way to manage data operations. By separating the data layer
	// from business logic, it enhances code maintainability, testability, and flexibility,
	// making it easier to work with various data sources in an application.
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
	defer repo.Close()

	// Builds the domain/business layer. The service coordinates use-cases and calls into
	// the repo. (Right now it’s thin: it just forwards to the repo.)
	svc := tasks.NewService(repo)

	// Creates the HTTP router with default middleware (logger + recovery).
	// Logger → logs each incoming HTTP request (method, path, status code, latency, etc.).
	// Recovery → catches any panics inside your handlers, logs the stack trace, and
	// responds with 500 Internal Server Error instead of crashing the whole process.
	//
	// A web service 'router' is the component that:
	// Receives an incoming HTTP request (e.g. GET /api/tasks)
	// Matches it against a set of rules you’ve defined (like “if it’s a GET and the path
	// is /api/tasks, call this function”).
	// Dispatches the request to the correct handler — your Go function that produces the
	// response.
	r := gin.Default()

	// Health check endpoint. Returns {"ok": true} with 200 status.
	// Useful for Docker/Kubernetes probes or for a quick curl check.
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// Group all API endpoints under the /api prefix.
	api := r.Group("/api")

	// Register your task-related routes (e.g., GET /api/tasks).
	// Notice we inject `svc` here: this allows tests to inject fakes instead of real DB.
	tasks.RegisterRoutes(api, svc)

	log.Printf("listening on :%s", cfg.Port)

	// Start the HTTP server — this blocks forever, handling requests until shutdown.
	// If the port is unavailable, Run() returns an error immediately.
	// On error, log.Fatalf prints the message and exits the process.
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
