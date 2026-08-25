package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/misbakhul29/hazartgo"
)

// Logger returns a middleware that logs HTTP request details (method, path, status, latency)
func Logger() hazart.MiddlewareFunc {
	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			start := time.Now()
			next(ctx)
			latency := time.Since(start)

			log.Printf("[HAZART] %s %s %d %s",
				ctx.Method,
				ctx.Path,
				ctx.StatusCode,
				latency,
			)
		}
	}
}

// Recovery returns a middleware that recovers from panics and returns 500 Internal Server Error
func Recovery() hazart.MiddlewareFunc {
	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("[HAZART RECOVERY] Panic recovered: %v", err)
					ctx.JSON(http.StatusInternalServerError, hazart.Map{
						"error": fmt.Sprintf("Internal Server Error: %v", err),
					})
				}
			}()
			next(ctx)
		}
	}
}

// CORSOptions defines parameters for CORS middleware
type CORSOptions struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

// CORS returns a middleware for handling Cross-Origin Resource Sharing
func CORS(opts ...CORSOptions) hazart.MiddlewareFunc {
	var opt CORSOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	origins := "*"
	if len(opt.AllowOrigins) > 0 {
		origins = ""
		for i, o := range opt.AllowOrigins {
			if i > 0 {
				origins += ", "
			}
			origins += o
		}
	}

	methods := "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	if len(opt.AllowMethods) > 0 {
		methods = ""
		for i, m := range opt.AllowMethods {
			if i > 0 {
				methods += ", "
			}
			methods += m
		}
	}

	headers := "Content-Type, Authorization, X-Requested-With"
	if len(opt.AllowHeaders) > 0 {
		headers = ""
		for i, h := range opt.AllowHeaders {
			if i > 0 {
				headers += ", "
			}
			headers += h
		}
	}

	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			ctx.SetHeader("Access-Control-Allow-Origin", origins)
			ctx.SetHeader("Access-Control-Allow-Methods", methods)
			ctx.SetHeader("Access-Control-Allow-Headers", headers)

			next(ctx)
		}
	}
}

// BearerAuth returns a middleware that validates Bearer JWT token from Authorization header.
// validator receives the raw token string and returns true if token is valid.
func BearerAuth(validator func(token string) bool) hazart.MiddlewareFunc {
	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			authHeader := ctx.Request.Header.Get("Authorization")
			if authHeader == "" {
				ctx.JSON(http.StatusUnauthorized, hazart.Unauthorized("Missing Authorization header"))
				return
			}

			var token string
			_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
			if err != nil || token == "" {
				ctx.JSON(http.StatusUnauthorized, hazart.Unauthorized("Invalid Authorization header format. Expected 'Bearer <token>'"))
				return
			}

			if validator != nil && !validator(token) {
				ctx.JSON(http.StatusUnauthorized, hazart.Unauthorized("Invalid or expired Bearer token"))
				return
			}

			next(ctx)
		}
	}
}
