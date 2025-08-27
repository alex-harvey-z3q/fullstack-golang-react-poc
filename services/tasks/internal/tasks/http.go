package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/alex-harvey-z3q/fullstack-golang-react-poc/services/tasks/internal/config"
)

func RegisterRoutes(r *gin.RouterGroup) {
	cfg := config.Load()
	repo := NewRepo(cfg)

	svc := NewService(repo)

	r.GET("/tasks", func(c *gin.Context) {
		items, err := svc.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	// TODO: POST /tasks, GET/PUT/DELETE /tasks/:id using repo methods
}
