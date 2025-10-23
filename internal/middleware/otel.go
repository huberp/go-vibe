package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// SkipPaths defines paths that should be excluded from OpenTelemetry tracing.
// These typically include health checks, metrics, and informational endpoints.
var SkipPaths = []string{
	"/health",
	"/health/",
	"/metrics",
	"/info",
}

// shouldTraceRequest determines if a request should be traced based on its path.
// Returns true if the request should be traced, false if it should be skipped.
func shouldTraceRequest(req *http.Request) bool {
	path := req.URL.Path

	// Return false to filter out (skip) these paths
	for _, skipPath := range SkipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath+"/") {
			return false
		}
	}

	return true
}

// OtelMiddleware returns a middleware that conditionally applies OpenTelemetry tracing.
// It skips tracing for specified paths like health checks, metrics, and info endpoints.
func OtelMiddleware(serviceName string, enabled bool) gin.HandlerFunc {
	if !enabled {
		// Return a no-op middleware if OTEL is disabled
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Create otelgin middleware with filter
	return otelgin.Middleware(serviceName, otelgin.WithFilter(shouldTraceRequest))
}
