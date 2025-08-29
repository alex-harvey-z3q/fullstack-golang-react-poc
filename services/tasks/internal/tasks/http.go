package tasks

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TaskLister is a *narrow interface* that this HTTP layer depends on.
// It only requires a List method.
//
// This is deliberate: it decouples the HTTP code from the full Service struct.
// In tests, we can pass in a fake that satisfies this interface.
//
// Note that every HTTP request in Go contains a Context.
//
// A Context carries request-scoped data: cancellation, deadlines/timeouts, and key–value metadata.
// It’s passed down so DB calls and services can stop early or propagate trace/user info.
type TaskLister interface {
	List(ctx context.Context) ([]Task, error)
}

// RegisterRoutes wires up HTTP endpoints under a given router group.
// The svc argument is anything that satisfies TaskLister (usually *Service).
//
// The *gin.RouterGroup is a pointer to a gin.RouterGroup as we created
// in main.go via r.Group.
func RegisterRoutes(r *gin.RouterGroup, svc TaskLister) {
	// GET /api/tasks
	r.GET("/tasks", func(c *gin.Context) {
		// Call the service layer to fetch tasks.
		items, err := svc.List(c.Request.Context())
		if err != nil {
			// If the service returns an error (e.g., DB failure),
			// respond with HTTP 500 and the error message as JSON.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// On success, return the list of tasks as JSON with HTTP 200.
		c.JSON(http.StatusOK, items)
	})
}
