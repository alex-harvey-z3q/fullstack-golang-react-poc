package api_test

import (
	"context"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func Test_OpenAPI_SpecIsValid_And_DeclaresGETTasks(t *testing.T) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("openapi.yaml")
	if err != nil {
		t.Fatalf("load openapi.yaml: %v", err)
	}
	if err := doc.Validate(context.Background()); err != nil {
		t.Fatalf("spec validation error: %v", err)
	}
	// Path + GET and POST must exist; deeper response checks are exercised by the runtime test.
	if pi := doc.Paths.Find("/api/tasks"); pi == nil || pi.Get == nil {
		t.Fatalf("GET /api/tasks not declared in openapi.yaml")
	}
	if pi := doc.Paths.Find("/api/tasks"); pi == nil || pi.Post == nil {
		t.Fatalf("POST /api/tasks not declared in openapi.yaml")
	}
}
