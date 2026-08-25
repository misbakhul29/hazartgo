package hazart

import (
	"encoding/json"
	"fmt"
	"net/http"
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
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
		Params:  make(map[string]string),
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
