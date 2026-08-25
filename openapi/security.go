package openapi

// SecurityType represents strong typed security scheme name or type
type SecurityType string

const (
	SecurityBearerAuth SecurityType = "BearerAuth"
	SecurityApiKeyAuth SecurityType = "ApiKeyAuth"
	SecurityBasicAuth  SecurityType = "BasicAuth"
)

// SecurityScheme defines security auth types (Bearer JWT, API Key, Basic Auth)
type SecurityScheme struct {
	Type         string `json:"type"` // bearer, apiKey, http
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	In           string `json:"in,omitempty"` // header, query
	Name         string `json:"name,omitempty"`
}

// AddSecurityScheme registers an authentication scheme (e.g. Bearer JWT)
func (g *Generator) AddSecurityScheme(name SecurityType, scheme SecurityScheme) {
	if g.Spec.Components.SecuritySchemes == nil {
		g.Spec.Components.SecuritySchemes = make(map[string]SecurityScheme)
	}
	g.Spec.Components.SecuritySchemes[string(name)] = scheme
}

// AddGlobalSecurity applies a security scheme globally to all endpoints
func (g *Generator) AddGlobalSecurity(name SecurityType) {
	g.Spec.Security = append(g.Spec.Security, map[string][]string{
		string(name): {},
	})
}
