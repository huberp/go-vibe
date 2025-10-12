package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// LoggingMiddleware logs HTTP requests with structured logging and W3C trace context support
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID or use existing trace ID from W3C trace context
		var requestID string
		
		// Check for W3C traceparent header
		traceparent := c.GetHeader("traceparent")
		if traceparent != "" {
			// Extract trace ID from W3C traceparent header (format: version-trace_id-parent_id-flags)
			// We'll use the trace ID as request ID for consistency
			if len(traceparent) >= 55 {
				requestID = traceparent[3:35] // Extract 32-char trace ID
			}
		}
		
		// If no traceparent or extraction failed, generate new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		c.Set("request_id", requestID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request with trace context if available
		duration := time.Since(start)
		logFields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
		}

		// Add trace context information if available
		if span := trace.SpanFromContext(c.Request.Context()); span.SpanContext().IsValid() {
			logFields = append(logFields,
				zap.String("trace_id", span.SpanContext().TraceID().String()),
				zap.String("span_id", span.SpanContext().SpanID().String()),
			)
		}

		logger.Info("http request", logFields...)
	}
}
