package openapi_test

import (
	"reflect"
	"testing"

	"github.com/misbakhul29/hazartgo/openapi"
)

type TestItem struct {
	ID   string `json:"id" doc:"Item ID"`
	Name string `json:"name" doc:"Item Name"`
}

func TestOpenAPI_SliceResponse(t *testing.T) {
	gen := openapi.NewGenerator("Test API", "1.0.0")

	// Register GET /items returning []TestItem
	gen.RegisterRoute("GET", "/api/v1/items", nil, reflect.TypeOf([]TestItem{}), "List Items", "Get all items", []string{"Items"}, "")

	spec := gen.Spec
	pathItem, exists := spec.Paths["/api/v1/items"]
	if !exists {
		t.Fatalf("expected path /api/v1/items in OpenAPI spec")
	}

	op := pathItem["get"]
	if op == nil {
		t.Fatalf("expected GET operation on /api/v1/items")
	}

	res200, ok := op.Responses["200"]
	if !ok {
		t.Fatalf("expected 200 response")
	}

	content, ok := res200.Content["application/json"]
	if !ok || content.Schema == nil {
		t.Fatalf("expected application/json content with schema")
	}

	if content.Schema.Type != "array" {
		t.Errorf("expected schema type 'array', got '%s'", content.Schema.Type)
	}

	if content.Schema.Items == nil {
		t.Fatalf("expected schema.items to be non-nil")
	}

	if content.Schema.Items.Ref != "#/components/schemas/TestItem" {
		t.Errorf("expected items.$ref to be '#/components/schemas/TestItem', got '%s'", content.Schema.Items.Ref)
	}

	// Verify TestItem is registered in components.schemas
	schemaItem, exists := spec.Components.Schemas["TestItem"]
	if !exists {
		t.Fatalf("expected 'TestItem' in components.schemas")
	}

	if _, ok := schemaItem.Properties["name"]; !ok {
		t.Errorf("expected property 'name' in TestItem schema")
	}
}

func TestOpenAPI_PostOmitID(t *testing.T) {
	gen := openapi.NewGenerator("Test API", "1.0.0")

	// Register POST /items taking TestItem as request body
	gen.RegisterRoute("POST", "/api/v1/items", reflect.TypeOf(TestItem{}), reflect.TypeOf(TestItem{}), "Create Item", "Create item", []string{"Items"}, "")

	pathItem := gen.Spec.Paths["/api/v1/items"]
	op := pathItem["post"]

	if op.RequestBody == nil {
		t.Fatalf("expected request body on POST operation")
	}

	reqSchema, ok := gen.Spec.Components.Schemas["TestItemRequest"]
	if !ok {
		t.Fatalf("expected TestItemRequest in components.schemas")
	}

	if _, hasID := reqSchema.Properties["id"]; hasID {
		t.Errorf("expected POST request body schema to omit 'id' property")
	}

	if _, hasName := reqSchema.Properties["name"]; !hasName {
		t.Errorf("expected POST request body schema to include 'name' property")
	}

	// Register PUT /items/:id
	gen.RegisterRoute("PUT", "/api/v1/items/:id", reflect.TypeOf(TestItem{}), reflect.TypeOf(TestItem{}), "Update Item", "Update item", []string{"Items"}, "")
	putPathItem := gen.Spec.Paths["/api/v1/items/{id}"]
	putOp := putPathItem["put"]

	if len(putOp.Parameters) == 0 || putOp.Parameters[0].Name != "id" {
		t.Errorf("expected path parameter 'id' on PUT operation")
	}

	if putOp.RequestBody == nil {
		t.Fatalf("expected request body on PUT operation")
	}
}
