package hazart

import (
	"net/http"
	"strings"
)

// HandlerFunc defines the standard HTTP handler function used inside HazartGo router.
type HandlerFunc func(ctx *Context)

// node represents a single node in the Radix Tree router.
type node struct {
	path     string
	part     string
	children []*node
	isWild   bool
	handlers map[string]HandlerFunc // map[HTTPMethod]HandlerFunc
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int, method string, handler HandlerFunc) {
	if len(parts) == height {
		n.path = pattern
		if n.handlers == nil {
			n.handlers = make(map[string]HandlerFunc)
		}
		n.handlers[method] = handler
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1, method, handler)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// Router handles path matching and route registration.
type Router struct {
	root *node
}

func newRouter() *Router {
	return &Router{
		root: &node{},
	}
}

func parsePath(path string) []string {
	vs := strings.Split(path, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Router) addRoute(method string, path string, handler HandlerFunc) {
	parts := parsePath(path)
	r.root.insert(path, parts, 0, method, handler)
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)
	n := r.root.search(searchParts, 0)

	if n != nil {
		parts := parsePath(n.path)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *Router) handle(ctx *Context) {
	n, params := r.getRoute(ctx.Method, ctx.Path)
	if n != nil && n.handlers[ctx.Method] != nil {
		ctx.Params = params
		n.handlers[ctx.Method](ctx)
	} else {
		ctx.JSON(http.StatusNotFound, Map{"error": "404 Not Found", "path": ctx.Path})
	}
}
