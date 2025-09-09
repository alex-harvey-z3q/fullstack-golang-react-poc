package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// Fake that satisfies TaskLister
type fakeSvc struct{}

func (f *fakeSvc) List(ctx context.Context) ([]Task, error) {
	now := time.Now().UTC()
	return []Task{
		{ID: 1, Title: "First", Done: false, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Title: "Second", Done: true, CreatedAt: now, UpdatedAt: now},
	}, nil
}

func (f *fakeSvc) Create(ctx context.Context, title string) (Task, error) {
	now := time.Now().UTC()
	return Task{ID: 3, Title: title, Done: false, CreatedAt: now, UpdatedAt: now}, nil
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

func TestPOSTTasks_CreatesAndReturnsJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	RegisterRoutes(api, &fakeSvc{})

	body := bytes.NewBufferString(`{"title":"From UI"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body=%s", w.Code, w.Body.String())
	}
	var got Task
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("json: %v; body=%s", err, w.Body.String())
	}
	if got.Title != "From UI" || got.ID == 0 {
		t.Fatalf("unexpected response: %#v", got)
	}
}
