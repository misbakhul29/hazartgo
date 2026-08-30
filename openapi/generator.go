package openapi

import (
	"reflect"
	"strings"
)

// Spec represents OpenAPI 3.1 structure
type Spec struct {
	OpenAPI    string                 `json:"openapi"`
	Info       Info                   `json:"info"`
	Paths      map[string]PathItem    `json:"paths"`
	Components Components             `json:"components"`
	Security   []map[string][]string  `json:"security,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type Info struct {
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version"`
}

type PathItem map[string]*Operation

type Operation struct {
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	OperationID string                 `json:"operationId,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Parameters  []Parameter            `json:"parameters,omitempty"`
	RequestBody *RequestBody           `json:"requestBody,omitempty"`
	Responses   map[string]Response    `json:"responses"`
	Security    []map[string][]string  `json:"security,omitempty"`
}

type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"` // path, query, header
	Required    bool    `json:"required"`
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema"`
}

type RequestBody struct {
	Required string               `json:"required,omitempty"`
	Content  map[string]MediaType `json:"content"`
}

type MediaType struct {
	Schema *Schema `json:"schema"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type Schema struct {
	Type                 string            `json:"type,omitempty"`
	Ref                  string            `json:"$ref,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string          `json:"required,omitempty"`
	Description          string            `json:"description,omitempty"`
	Items                *Schema           `json:"items,omitempty"`
}

type Components struct {
	Schemas         map[string]*Schema        `json:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
}

// Generator manages OpenAPI spec generation
type Generator struct {
	Spec *Spec
}

func NewGenerator(title, version string) *Generator {
	return &Generator{
		Spec: &Spec{
			OpenAPI: "3.1.0",
			Info: Info{
				Title:   title,
				Version: version,
			},
			Paths: make(map[string]PathItem),
			Components: Components{
				Schemas: make(map[string]*Schema),
			},
		},
	}
}

// RegisterRoute registers a route with its request/response struct types
func (g *Generator) RegisterRoute(method, path string, reqType reflect.Type, resType reflect.Type, summary, description string, tags []string, security string) {
	openapiPath := convertPathToOpenAPI(path)
	pathItem, exists := g.Spec.Paths[openapiPath]
	if !exists {
		pathItem = make(PathItem)
	}

	op := &Operation{
		Summary:     summary,
		Description: description,
		Tags:        tags,
		Responses:   make(map[string]Response),
	}

	if security != "" {
		op.Security = []map[string][]string{
			{security: {}},
		}
	}

	// Process Request Struct (Query, Path, Body)
	if reqType != nil {
		if reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}
		if reqType.Kind() == reflect.Struct {
			g.processRequestStruct(op, reqType)
		}
	}

	// Process Response Struct
	if resType != nil && resType.Kind() == reflect.Ptr {
		resType = resType.Elem()
	}

	if resType != nil && resType.Kind() == reflect.Struct && resType.NumField() == 0 {
		resType = nil
	}

	if resType != nil {
		schemaRef := g.registerSchema(resType)
		status := "200"
		if strings.EqualFold(method, "DELETE") {
			status = "204"
		}
		op.Responses[status] = Response{
			Description: "Successful Operation",
			Content: map[string]MediaType{
				"application/json": {
					Schema: schemaRef,
				},
			},
		}
	} else {
		if strings.EqualFold(method, "DELETE") {
			op.Responses["204"] = Response{Description: "Resource deleted successfully"}
		} else {
			op.Responses["200"] = Response{Description: "Successful Operation"}
		}
	}

	pathItem[strings.ToLower(method)] = op
	g.Spec.Paths[openapiPath] = pathItem
}

func convertPathToOpenAPI(path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if strings.HasPrefix(p, ":") {
			parts[i] = "{" + p[1:] + "}"
		}
	}
	return strings.Join(parts, "/")
}

func (g *Generator) processRequestStruct(op *Operation, t reflect.Type) {
	var bodyProperties = make(map[string]*Schema)
	var bodyRequired []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		doc := field.Tag.Get("doc")
		validateTag := field.Tag.Get("validate")
		isRequired := strings.Contains(validateTag, "required")

		if pathTag := field.Tag.Get("path"); pathTag != "" {
			op.Parameters = append(op.Parameters, Parameter{
				Name:        pathTag,
				In:          "path",
				Required:    true,
				Description: doc,
				Schema:      g.typeToSchema(field.Type),
			})
		} else if queryTag := field.Tag.Get("query"); queryTag != "" {
			op.Parameters = append(op.Parameters, Parameter{
				Name:        queryTag,
				In:          "query",
				Required:    isRequired,
				Description: doc,
				Schema:      g.typeToSchema(field.Type),
			})
		} else if headerTag := field.Tag.Get("header"); headerTag != "" {
			op.Parameters = append(op.Parameters, Parameter{
				Name:        headerTag,
				In:          "header",
				Required:    isRequired,
				Description: doc,
				Schema:      g.typeToSchema(field.Type),
			})
		} else if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			jsonName := strings.Split(jsonTag, ",")[0]
			bodyProperties[jsonName] = g.typeToSchema(field.Type)
			if isRequired {
				bodyRequired = append(bodyRequired, jsonName)
			}
		}
	}

	if len(bodyProperties) > 0 {
		schemaName := t.Name() + "Request"
		g.Spec.Components.Schemas[schemaName] = &Schema{
			Type:        "object",
			Properties:  bodyProperties,
			Required:    bodyRequired,
			Description: "Request body payload",
		}
		op.RequestBody = &RequestBody{
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{Ref: "#/components/schemas/" + schemaName},
				},
			},
		}
	}
}

func (g *Generator) registerSchema(t reflect.Type) *Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		return &Schema{
			Type:  "array",
			Items: g.registerSchema(t.Elem()),
		}
	}

	if t.Kind() != reflect.Struct {
		return g.typeToSchema(t)
	}

	schemaName := t.Name()
	if schemaName == "" {
		return &Schema{Type: "object"}
	}

	if _, exists := g.Spec.Components.Schemas[schemaName]; !exists {
		props := make(map[string]*Schema)
		var required []string

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			jsonName := strings.Split(jsonTag, ",")[0]
			validateTag := field.Tag.Get("validate")

			fieldSchema := g.registerSchema(field.Type)
			docTag := field.Tag.Get("doc")
			if docTag != "" {
				fieldSchema.Description = docTag
			}

			props[jsonName] = fieldSchema
			if strings.Contains(validateTag, "required") {
				required = append(required, jsonName)
			}
		}

		g.Spec.Components.Schemas[schemaName] = &Schema{
			Type:       "object",
			Properties: props,
			Required:   required,
		}
	}

	return &Schema{Ref: "#/components/schemas/" + schemaName}
}

func (g *Generator) typeToSchema(t reflect.Type) *Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return &Schema{Type: "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Schema{Type: "integer"}
	case reflect.Float32, reflect.Float64:
		return &Schema{Type: "number"}
	case reflect.Bool:
		return &Schema{Type: "boolean"}
	case reflect.Slice, reflect.Array:
		return &Schema{
			Type:  "array",
			Items: g.registerSchema(t.Elem()),
		}
	case reflect.Struct:
		return g.registerSchema(t)
	default:
		return &Schema{Type: "object"}
	}
}
