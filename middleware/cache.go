package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/misbakhul29/hazartgo"
)

type cacheEntry struct {
	statusCode  int
	contentType string
	body        []byte
	expiration  time.Time
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode  int
	contentType string
	body        bytes.Buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Cache returns a middleware that caches HTTP GET responses in-memory for specified TTL duration
func Cache(ttl time.Duration) hazart.MiddlewareFunc {
	var mu sync.RWMutex
	store := make(map[string]*cacheEntry)

	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			// Only cache GET requests
			if ctx.Method != http.MethodGet {
				next(ctx)
				return
			}

			key := ctx.Request.URL.RequestURI()

			mu.RLock()
			entry, found := store[key]
			mu.RUnlock()

			if found && time.Now().Before(entry.expiration) {
				ctx.SetHeader("X-Cache", "HIT")
				ctx.SetHeader("Content-Type", entry.contentType)
				ctx.Status(entry.statusCode)
				_, _ = ctx.Writer.Write(entry.body)
				return
			}

			rec := &responseRecorder{
				ResponseWriter: ctx.Writer,
				statusCode:     http.StatusOK,
			}
			ctx.Writer = rec
			ctx.SetHeader("X-Cache", "MISS")

			next(ctx)

			if rec.statusCode >= 200 && rec.statusCode < 300 {
				mu.Lock()
				store[key] = &cacheEntry{
					statusCode:  rec.statusCode,
					contentType: rec.Header().Get("Content-Type"),
					body:        rec.body.Bytes(),
					expiration:  time.Now().Add(ttl),
				}
				mu.Unlock()
			}
		}
	}
}
