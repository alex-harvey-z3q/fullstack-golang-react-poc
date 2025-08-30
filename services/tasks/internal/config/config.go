package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all environment-driven settings for this service.
// Port → HTTP listen port
// DatabaseURL → connection string for Postgres
// Env → environment name (e.g. "dev", "prod") to switch behavior
type Config struct {
	Port        string
	DatabaseURL string
	Env         string
}

// Load reads environment variables into a Config struct.
// It also supports loading from a local .env file for developer convenience.
func Load() Config {
	// godotenv.Load() will read a .env file if present
	// (useful locally; in Docker/K8s env vars come from the runtime).
	_ = godotenv.Load()

	c := Config{
		// Use get() helper to fetch an env var or fall back to a default.
		Port:        get("PORT", "8081"),
		DatabaseURL: get("DATABASE_URL", "postgres://app:app@localhost:5432/tasks?sslmode=disable"),
		Env:         get("ENV", "dev"),
	}

	// Log the environment for visibility at startup.
	log.Printf("env=%s", c.Env)
	return c
}

// get is a tiny helper: returns the env var if set, otherwise the provided default.
// Keeps Load() cleaner and avoids repeating the same logic.
func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
