package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/config"
	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/tasks"
)

func main() {
	cfg := config.Load()

	// Real repo + service for production/runtime
	repo := tasks.NewRepo(cfg)
	defer repo.Close()
	svc := tasks.NewService(repo)

	r := gin.Default()

	// Healthcheck
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// API routes with injected service (NEW signature)
	api := r.Group("/api")
	tasks.RegisterRoutes(api, svc)

	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
