package middleware

import (
	"net/http"
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

	// Use otelgin's built-in filter to skip health checks, metrics, and info endpoints
	filter := func(req *http.Request) bool {
		path := req.URL.Path
		
		// Return false to filter out (skip) these paths
		skipPaths := []string{
			"/health",
			"/health/",
			"/metrics",
			"/info",
		}
		
		for _, skipPath := range skipPaths {
			if path == skipPath || strings.HasPrefix(path, skipPath+"/") {
				return false
			}
		}
		
		return true
	}

	// Create otelgin middleware with filter
	return otelgin.Middleware(serviceName, otelgin.WithFilter(filter))
}
