package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/misbakhul29/hazartgo"
)

// TraceHeader key name for distributed tracing correlation ID
const TraceHeader = "X-Trace-ID"

func generateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(bytes)
}

// Tracing returns a middleware that injects or propagates a unique X-Trace-ID header in context
func Tracing() hazart.MiddlewareFunc {
	return func(next hazart.HandlerFunc) hazart.HandlerFunc {
		return func(ctx *hazart.Context) {
			traceID := ctx.Request.Header.Get(TraceHeader)
			if traceID == "" {
				traceID = generateTraceID()
			}

			ctx.Set("trace_id", traceID)
			ctx.SetHeader(TraceHeader, traceID)

			next(ctx)
		}
	}
}
