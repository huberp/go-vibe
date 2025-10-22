package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// OtelMiddleware returns a middleware that conditionally applies OpenTelemetry tracing.
// It skips tracing for specified paths like health checks, metrics, and info endpoints.
func OtelMiddleware(serviceName string, enabled bool) gin.HandlerFunc {
	if !enabled {
		// Return a no-op middleware if OTEL is disabled
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Paths to skip from tracing
	skipPaths := []string{
		"/health",
		"/health/",
		"/metrics",
		"/info",
	}

	// Create the base otelgin middleware
	otelMiddleware := otelgin.Middleware(serviceName)

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Check if path should be skipped
		shouldSkip := false
		for _, skipPath := range skipPaths {
			if path == skipPath || strings.HasPrefix(path, skipPath+"/") {
				shouldSkip = true
				break
			}
		}

		if shouldSkip {
			// Skip OTEL tracing for this path
			c.Next()
			return
		}

		// Apply OTEL middleware
		otelMiddleware(c)
	}
}
