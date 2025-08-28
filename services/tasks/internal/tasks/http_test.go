package tasks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Fake that satisfies TaskLister
type fakeSvc struct{}

func (f *fakeSvc) List(ctx context.Context) ([]Task, error) {
	return []Task{
		{ID: 1, Title: "First", Done: false},
		{ID: 2, Title: "Second", Done: true},
	}, nil
}

func TestGETTasks_ReturnsJSONList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	api := r.Group("/api")
	RegisterRoutes(api, &fakeSvc{})

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", w.Code, w.Body.String())
	}

	var got []Task
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("json unmarshal error: %v; body=%s", err, w.Body.String())
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 tasks, got %d: %#v", len(got), got)
	}
	if got[0].Title != "First" || got[1].Title != "Second" {
		t.Fatalf("unexpected titles: %#v", got)
	}
}
