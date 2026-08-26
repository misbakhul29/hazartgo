package hazart

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SSEStream represents a Server-Sent Events stream connection handler
type SSEStream struct {
	ctx *Context
}

// SSEvent sends a structured Server-Sent Event to client
func (c *Context) SSEvent(event string, data any) error {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported by response writer")
	}

	c.SetHeader("Content-Type", "text/event-stream")
	c.SetHeader("Cache-Control", "no-cache")
	c.SetHeader("Connection", "keep-alive")

	var payload string
	switch v := data.(type) {
	case string:
		payload = v
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		payload = string(bytes)
	}

	if event != "" {
		fmt.Fprintf(c.Writer, "event: %s\n", event)
	}
	fmt.Fprintf(c.Writer, "data: %s\n\n", payload)

	flusher.Flush()
	return nil
}

// SSE Handler helper for continuous event streams
func SSE(fn func(c *Context) error) HandlerFunc {
	return func(c *Context) {
		c.SetHeader("Content-Type", "text/event-stream")
		c.SetHeader("Cache-Control", "no-cache")
		c.SetHeader("Connection", "keep-alive")
		_ = fn(c)
	}
}
