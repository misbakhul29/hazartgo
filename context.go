package hazart

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/misbakhul29/hazartgo/binder"
)

// Map is a shortcut for map[string]any.
type Map map[string]any

// Context holds request and response state for a single HTTP handler execution.
type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	store      map[string]any
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
		Params:  make(map[string]string),
		store:   make(map[string]any),
	}
}

// Param returns URL route path parameter value by key name.
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// Query returns URL query parameter value by key name.
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// SetHeader sets a response HTTP header key-value.
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// Set stores a key-value pair in context store
func (c *Context) Set(key string, val any) {
	if c.store == nil {
		c.store = make(map[string]any)
	}
	c.store[key] = val
}

// Get retrieves a stored value from context store
func (c *Context) Get(key string) (any, bool) {
	if c.store == nil {
		return nil, false
	}
	val, ok := c.store[key]
	return val, ok
}

// Status sets HTTP response status code.
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// JSON sends a JSON response with status code.
func (c *Context) JSON(code int, v any) error {
	c.StatusCode = code
	c.SetHeader("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	return json.NewEncoder(c.Writer).Encode(v)
}

// Success sends a standardized JSON success response wrapper
func (c *Context) Success(code int, data any, message ...string) error {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}
	return c.JSON(code, Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

// Error sends a standardized JSON error response wrapper
func (c *Context) Error(code int, message string, details ...any) error {
	res := Map{
		"success": false,
		"error":   message,
	}
	if len(details) > 0 {
		res["details"] = details[0]
	}
	return c.JSON(code, res)
}

// String sends a plain text response.
func (c *Context) String(code int, format string, values ...any) error {
	c.StatusCode = code
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	_, err := fmt.Fprintf(c.Writer, format, values...)
	return err
}

// HTML sends an HTML response.
func (c *Context) HTML(code int, html string) error {
	c.SetHeader("Content-Type", "text/html")
	c.Writer.WriteHeader(code)
	c.StatusCode = code
	_, err := c.Writer.Write([]byte(html))
	return err
}

// BindAndValidate is a generic helper to bind request params/body to struct T and validate via struct tags.
func BindAndValidate[T any](c *Context) (*T, error) {
	var target T
	if err := binder.BindAndValidate(c.Request, c.Params, &target); err != nil {
		return nil, err
	}
	return &target, nil
}
