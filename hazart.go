package hazart

import (
	"net/http"
	"reflect"

	"github.com/misbakhul29/hazartgo/binder"
	"github.com/misbakhul29/hazartgo/openapi"
)

// MiddlewareFunc defines a middleware handler wrap function
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// SecurityType defines strong typed security scheme identifiers
type SecurityType = openapi.SecurityType

const (
	SecurityBearerAuth SecurityType = openapi.SecurityBearerAuth
	SecurityApiKeyAuth SecurityType = openapi.SecurityApiKeyAuth
	SecurityBasicAuth  SecurityType = openapi.SecurityBasicAuth
)

// RouteMeta defines documentation metadata for OpenAPI generation
type RouteMeta struct {
	Summary     string
	Description string
	Tags        []string
	Security    SecurityType // e.g. hazart.SecurityBearerAuth
}

// App is the main HazartGo engine instance
type App struct {
	router      *Router
	OpenAPI     *openapi.Generator
	middlewares []MiddlewareFunc
}

// Config defines application configuration
type Config struct {
	Title          string
	Description    string
	TermsOfService string
	Contact        *openapi.Contact
	License        *openapi.License
	Version        string
}

// New creates a new HazartGo application
func New(cfg Config) *App {
	if cfg.Title == "" {
		cfg.Title = "HazartGo API"
	}
	if cfg.Version == "" {
		cfg.Version = "1.0.0"
	}

	app := &App{
		router:      newRouter(),
		OpenAPI:     openapi.NewGenerator(cfg.Title, cfg.Version),
		middlewares: make([]MiddlewareFunc, 0),
	}

	app.OpenAPI.Spec.Info.Description = cfg.Description
	app.OpenAPI.Spec.Info.TermsOfService = cfg.TermsOfService
	app.OpenAPI.Spec.Info.Contact = cfg.Contact
	app.OpenAPI.Spec.Info.License = cfg.License

	// Register Swagger / Docs routes
	app.registerDocsRoutes()

	return app
}

// Use attaches global middlewares to the application
func (a *App) Use(mw ...MiddlewareFunc) {
	a.middlewares = append(a.middlewares, mw...)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r)

	// Chain global middlewares with router handler
	var handler HandlerFunc = func(c *Context) {
		a.router.handle(c)
	}

	for i := len(a.middlewares) - 1; i >= 0; i-- {
		handler = a.middlewares[i](handler)
	}

	handler(ctx)
}

func (a *App) Listen(addr string) error {
	return http.ListenAndServe(addr, a)
}

func (a *App) registerDocsRoutes() {
	// Raw JSON Spec
	a.router.addRoute("GET", "/openapi.json", func(ctx *Context) {
		ctx.JSON(http.StatusOK, a.OpenAPI.Spec)
	})

	// Embedded Swagger UI HTML
	swaggerHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>` + a.OpenAPI.Spec.Info.Title + ` - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" crossorigin></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: '/openapi.json',
          dom_id: '#swagger-ui',
        });
      };
    </script>
</body>
</html>`

	a.router.addRoute("GET", "/docs", func(ctx *Context) {
		ctx.HTML(http.StatusOK, swaggerHTML)
	})
	a.router.addRoute("GET", "/swagger", func(ctx *Context) {
		ctx.HTML(http.StatusOK, swaggerHTML)
	})
}

// Generic Handler Registration Helper
func Register[Req any, Res any](a *App, method string, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	var reqInst Req
	var resInst Res

	reqType := reflect.TypeOf(reqInst)
	resType := reflect.TypeOf(resInst)

	var summary, description string
	var tags []string
	var security SecurityType
	if len(meta) > 0 {
		summary = meta[0].Summary
		description = meta[0].Description
		tags = meta[0].Tags
		security = meta[0].Security
	}

	// Register to OpenAPI Generator
	a.OpenAPI.RegisterRoute(method, path, reqType, resType, summary, description, tags, string(security))

	// Register HTTP Handler function
	a.router.addRoute(method, path, func(ctx *Context) {
		var req Req

		// Execute Binder & Validation
		if err := binder.BindAndValidate(ctx.Request, ctx.Params, &req); err != nil {
			ctx.JSON(http.StatusBadRequest, Map{
				"error":   "Validation / Binding Failed",
				"details": err.Error(),
			})
			return
		}

		// Call Generic Handler
		res, err := handler(ctx, &req)
		if err != nil {
			if apiErr, ok := err.(*APIError); ok {
				if apiErr.Instance == "" {
					apiErr.Instance = ctx.Path
				}
				ctx.JSON(apiErr.Status, apiErr)
				return
			}

			// Default Internal Server Error
			ctx.JSON(http.StatusInternalServerError, Map{
				"status": http.StatusInternalServerError,
				"title":  "Internal Server Error",
				"detail": err.Error(),
			})
			return
		}

		if res != nil {
			ctx.JSON(http.StatusOK, res)
		} else {
			ctx.Status(http.StatusNoContent)
		}
	})
}

// Shortcuts for HTTP Methods

func Get[Req any, Res any](a *App, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	Register(a, "GET", path, handler, meta...)
}

func Post[Req any, Res any](a *App, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	Register(a, "POST", path, handler, meta...)
}

func Put[Req any, Res any](a *App, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	Register(a, "PUT", path, handler, meta...)
}

func Delete[Req any, Res any](a *App, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	Register(a, "DELETE", path, handler, meta...)
}
