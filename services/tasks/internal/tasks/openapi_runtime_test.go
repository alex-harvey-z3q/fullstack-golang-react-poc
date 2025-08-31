package tasks

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
)

// distinct name to avoid clashing with fakeSvc in http_test.go
type oasFakeSvc struct{}

func (f *oasFakeSvc) List(ctx context.Context) ([]Task, error) {
	return []Task{
		{ID: 1, Title: "First", Done: false},
		{ID: 2, Title: "Second", Done: true},
	}, nil
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
