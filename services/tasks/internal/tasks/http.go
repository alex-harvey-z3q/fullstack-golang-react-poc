package tasks

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TaskLister is the minimal interface the HTTP layer needs.
// Your real *Service implements this, and tests can provide a fake.
type TaskLister interface {
	List(ctx context.Context) ([]Task, error)
}

// RegisterRoutes wires HTTP endpoints using the provided service.
// In main.go you'll pass the real service; in tests you can pass a fake.
func RegisterRoutes(r *gin.RouterGroup, svc TaskLister) {
	// GET /api/tasks
	r.GET("/tasks", func(c *gin.Context) {
		items, err := svc.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})
}
