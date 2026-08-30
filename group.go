package hazart

// Group represents a route group with a common path prefix and shared middlewares
type Group struct {
	prefix      string
	app         *App
	middlewares []MiddlewareFunc
}

// Group creates a new route sub-group from App
func (a *App) Group(prefix string, middlewares ...MiddlewareFunc) *Group {
	return &Group{
		prefix:      prefix,
		app:         a,
		middlewares: middlewares,
	}
}

// Group creates a nested route sub-group from an existing Group
func (g *Group) Group(prefix string, middlewares ...MiddlewareFunc) *Group {
	combinedMiddlewares := make([]MiddlewareFunc, 0, len(g.middlewares)+len(middlewares))
	combinedMiddlewares = append(combinedMiddlewares, g.middlewares...)
	combinedMiddlewares = append(combinedMiddlewares, middlewares...)

	return &Group{
		prefix:      g.prefix + prefix,
		app:         g.app,
		middlewares: combinedMiddlewares,
	}
}

// Use adds middlewares to the route group
func (g *Group) Use(mw ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, mw...)
}

// Generic Handler Registration inside Group
func GroupRegister[Req any, Res any](g *Group, method string, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	fullPath := g.prefix + path

	// Combine group middlewares with endpoint execution
	wrappedHandler := func(ctx *Context, req *Req) (*Res, error) {
		return handler(ctx, req)
	}

	_ = wrappedHandler

	// Register to underlying App router with group prefix
	Register(g.app, method, fullPath, func(ctx *Context, req *Req) (*Res, error) {
		var outerRes *Res
		var outerErr error

		// Run Group Middlewares sequentially
		var chain HandlerFunc = func(c *Context) {
			res, err := handler(c, req)
			outerRes = res
			outerErr = err
		}

		for i := len(g.middlewares) - 1; i >= 0; i-- {
			chain = g.middlewares[i](chain)
		}

		chain(ctx)
		return outerRes, outerErr
	}, meta...)
}

func GroupGet[Req any, Res any](g *Group, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	GroupRegister(g, "GET", path, handler, meta...)
}

func GroupPost[Req any, Res any](g *Group, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	GroupRegister(g, "POST", path, handler, meta...)
}

func GroupPut[Req any, Res any](g *Group, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	GroupRegister(g, "PUT", path, handler, meta...)
}

func GroupDelete[Req any, Res any](g *Group, path string, handler func(ctx *Context, req *Req) (*Res, error), meta ...RouteMeta) {
	GroupRegister(g, "DELETE", path, handler, meta...)
}
