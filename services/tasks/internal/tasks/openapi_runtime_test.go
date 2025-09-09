package tasks

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
)

// distinct name to avoid clashing with fakeSvc in http_test.go
type oasFakeSvc struct{}

func (f *oasFakeSvc) List(ctx context.Context) ([]Task, error) {
	now := time.Now().UTC()
	return []Task{
		{ID: 1, Title: "First", Done: false, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Title: "Second", Done: true, CreatedAt: now, UpdatedAt: now},
	}, nil
}

func (f *oasFakeSvc) Create(ctx context.Context, title string) (Task, error) {
	now := time.Now().UTC()
	return Task{ID: 3, Title: title, Done: false, CreatedAt: now, UpdatedAt: now}, nil
}

func Test_Server_GetTasks_MatchesOpenAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	api := r.Group("/api")
	RegisterRoutes(api, &oasFakeSvc{})

	specPath := filepath.Join("..", "..", "api", "openapi.yaml")
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		t.Fatalf("load openapi.yaml: %v", err)
	}
	if err := doc.Validate(context.Background()); err != nil {
		t.Fatalf("spec validation error: %v", err)
	}
	if pi := doc.Paths.Find("/api/tasks"); pi == nil || pi.Get == nil {
		t.Fatalf("GET /api/tasks not declared in spec")
	}

	// Exercise the live handler
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Route resolution against the spec
	router, err := legacyrouter.NewRouter(doc)
	if err != nil {
		t.Fatalf("build OAS router: %v", err)
	}
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		t.Fatalf("find route in spec for %s %s: %v", req.Method, req.URL.Path, err)
	}

	// Validate the real HTTP response against the OpenAPI document.
	body := io.NopCloser(bytes.NewReader(rec.Body.Bytes()))
	defer body.Close()

	in := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			Route:      route,
			Options:    &openapi3filter.Options{},
		},
		Status: rec.Code,
		Header: rec.Header(),
		Body:   body,
	}
	if err := openapi3filter.ValidateResponse(req.Context(), in); err != nil {
		t.Fatalf("response does not match OpenAPI: %v", err)
	}
}

func Test_Server_PostTasks_MatchesOpenAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	api := r.Group("/api")
	RegisterRoutes(api, &oasFakeSvc{})

	specPath := filepath.Join("..", "..", "api", "openapi.yaml")
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		t.Fatalf("load openapi.yaml: %v", err)
	}
	if err := doc.Validate(context.Background()); err != nil {
		t.Fatalf("spec validation error: %v", err)
	}
	if pi := doc.Paths.Find("/api/tasks"); pi == nil || pi.Post == nil {
		t.Fatalf("POST /api/tasks not declared in spec")
	}

	// Exercise the live handler
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(`{"title":"From test"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Route resolution against the spec
	router, err := legacyrouter.NewRouter(doc)
	if err != nil {
		t.Fatalf("build OAS router: %v", err)
	}
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		t.Fatalf("find route in spec for %s %s: %v", req.Method, req.URL.Path, err)
	}

	// Validate the real HTTP response against the OpenAPI document.
	body := io.NopCloser(bytes.NewReader(rec.Body.Bytes()))
	defer body.Close()

	in := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			Route:      route,
			Options:    &openapi3filter.Options{},
		},
		Status: rec.Code,
		Header: rec.Header(),
		Body:   body,
	}
	if err := openapi3filter.ValidateResponse(req.Context(), in); err != nil {
		t.Fatalf("response does not match OpenAPI: %v", err)
	}
}
