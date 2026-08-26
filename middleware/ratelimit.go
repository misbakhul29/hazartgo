package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/misbakhul29/hazartgo"
)

type clientVisitor struct {
	lastSeen time.Time
	tokens   int
}

// RateLimitConfig configures rate limiting parameters
type RateLimitConfig struct {
	RequestsPerWindow int
	Window            time.Duration
	KeyGenerator      func(ctx *hazart.Context) string
}

// RateLimiter returns a middleware that limits request frequency per client IP or key
func RateLimiter(cfg ...RateLimitConfig) hazart.MiddlewareFunc {
	config := RateLimitConfig{
		RequestsPerWindow: 60,
		Window:            time.Minute,
	}

	if len(cfg) > 0 {
		if cfg[0].RequestsPerWindow > 0 {
			config.RequestsPerWindow = cfg[0].RequestsPerWindow
		}
		if cfg[0].Window > 0 {
			config.Window = cfg[0].Window
		}
		if cfg[0].KeyGenerator != nil {
			config.KeyGenerator = cfg[0].KeyGenerator
		}
	}

	if config.KeyGenerator == nil {
		config.KeyGenerator = func(ctx *hazart.Context) string {
			ip := ctx.Request.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = ctx.Request.RemoteAddr
			}
			return ip
		}
	}

	var mu sync.Mutex
	visitors := make(map[string]*clientVisitor)

	// Clean up stale visitors periodically
	go func() {
		for {
			time.Sleep(config.Window)
			mu.Lock()
			for key, v := range visitors {
				if time.Since(v.lastSeen) > config.Window*2 {
					delete(visitors, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			key := config.KeyGenerator(ctx)

			mu.Lock()
			v, exists := visitors[key]
			now := time.Now()

			if !exists || now.Sub(v.lastSeen) > config.Window {
				visitors[key] = &clientVisitor{
					lastSeen: now,
					tokens:   config.RequestsPerWindow - 1,
				}
				mu.Unlock()
				next(ctx)
				return
			}

			if v.tokens <= 0 {
				mu.Unlock()
				ctx.JSON(http.StatusTooManyRequests, hazart.Map{
					"success": false,
					"error":   "Rate limit exceeded. Please try again later.",
				})
				return
			}

			v.tokens--
			v.lastSeen = now
			mu.Unlock()

			next(ctx)
		}
	}
}
