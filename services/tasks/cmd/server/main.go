package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/config"
	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/tasks"
)

func main() {
	cfg := config.Load()
	r := gin.Default()

	// Basic health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// REST routes
	api := r.Group("/api")
	tasks.RegisterRoutes(api) // /api/tasks...

	log.Printf("listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Println("server error:", err)
		os.Exit(1)
	}
}
