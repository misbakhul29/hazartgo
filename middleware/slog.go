package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/misbakhul29/hazartgo"
)

// StructuredLogger returns a middleware using Go log/slog structured JSON logger
func StructuredLogger(logger *slog.Logger) hazart.MiddlewareFunc {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			start := time.Now()
			next(ctx)
			latency := time.Since(start)

			traceID, _ := ctx.Get("trace_id")
			if traceID == nil {
				traceID = ""
			}

			logger.Info("HTTP Request",
				slog.String("method", ctx.Method),
				slog.String("path", ctx.Path),
				slog.Int("status", ctx.StatusCode),
				slog.Duration("latency", latency),
				slog.Any("trace_id", traceID),
				slog.String("client_ip", ctx.Request.RemoteAddr),
			)
		}
	}
}
