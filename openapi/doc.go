package openapi

import "reflect"

// Doc returns a JSON Schema description for a field
func (g *Generator) generateFieldSchemaWithDoc(field reflect.StructField) *Schema {
	schema := g.typeToSchema(field.Type)

	if docTag := field.Tag.Get("doc"); docTag != "" {
		schema.Description = docTag
	}

	if exampleTag := field.Tag.Get("example"); exampleTag != "" {
		schema.Description += " (Example: " + exampleTag + ")"
	}

	return schema
}
