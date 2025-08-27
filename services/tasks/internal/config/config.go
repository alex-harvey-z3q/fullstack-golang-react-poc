package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	Env         string
}

func Load() Config {
	_ = godotenv.Load()
	c := Config{
		Port:        get("PORT", "8081"),
		DatabaseURL: get("DATABASE_URL", "postgres://app:app@localhost:5432/tasks?sslmode=disable"),
		Env:         get("ENV", "dev"),
	}
	log.Printf("env=%s", c.Env)
	return c
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
