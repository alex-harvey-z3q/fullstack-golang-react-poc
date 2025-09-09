package tasks

import (
	"context"
	"net/http"
	"strings"

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

// taskCreator is an optional capability discovered at runtime.
// If the provided service also implements this, POST /api/tasks will succeed.
// Otherwise the handler returns 501 Not Implemented.
type taskCreator interface {
	Create(ctx context.Context, title string) (Task, error)
}

// RegisterRoutes wires up HTTP endpoints under a given router group.
// The svc argument only needs to satisfy TaskLister; if it also implements
// taskCreator, POST /api/tasks will be enabled.
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

	// POST /api/tasks
	type createReq struct {
		Title string `json:"title"`
	}
	r.POST("/tasks", func(c *gin.Context) {
		// Discover create capability at runtime.
		cr, ok := svc.(taskCreator)
		if !ok {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "create not supported"})
			return
		}

		var req createReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			return
		}
		title := strings.TrimSpace(req.Title)
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
			return
		}

		t, err := cr.Create(c.Request.Context(), title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, t)
	})
}
