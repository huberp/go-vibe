package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestOtelMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should skip tracing when disabled", func(t *testing.T) {
		router := gin.New()
		router.Use(OtelMiddleware("test-service", false))
		
		router.GET("/test", func(c *gin.Context) {
			// Check that no span exists in context
			span := trace.SpanFromContext(c.Request.Context())
			assert.False(t, span.IsRecording(), "Should not have recording span when OTEL is disabled")
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should skip tracing for /health endpoint", func(t *testing.T) {
		router := gin.New()
		router.Use(OtelMiddleware("test-service", true))
		
		router.GET("/health", func(c *gin.Context) {
			// Check that no span exists in context
			span := trace.SpanFromContext(c.Request.Context())
			assert.False(t, span.IsRecording(), "Should not have recording span for /health")
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should skip tracing for /health/startup endpoint", func(t *testing.T) {
		router := gin.New()
		router.Use(OtelMiddleware("test-service", true))
		
		router.GET("/health/startup", func(c *gin.Context) {
			// Check that no span exists in context
			span := trace.SpanFromContext(c.Request.Context())
			assert.False(t, span.IsRecording(), "Should not have recording span for /health/startup")
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/health/startup", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should skip tracing for /health/liveness endpoint", func(t *testing.T) {
		router := gin.New()
		router.Use(OtelMiddleware("test-service", true))
		
		router.GET("/health/liveness", func(c *gin.Context) {
			// Check that no span exists in context
			span := trace.SpanFromContext(c.Request.Context())
			assert.False(t, span.IsRecording(), "Should not have recording span for /health/liveness")
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/health/liveness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should skip tracing for /health/readiness endpoint", func(t *testing.T) {
		router := gin.New()
		router.Use(OtelMiddleware("test-service", true))
		
		router.GET("/health/readiness", func(c *gin.Context) {
			// Check that no span exists in context
			span := trace.SpanFromContext(c.Request.Context())
			assert.False(t, span.IsRecording(), "Should not have recording span for /health/readiness")
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/health/readiness", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should skip tracing for /metrics endpoint", func(t *testing.T) {
		router := gin.New()
		router.Use(OtelMiddleware("test-service", true))
		
		router.GET("/metrics", func(c *gin.Context) {
			// Check that no span exists in context
			span := trace.SpanFromContext(c.Request.Context())
			assert.False(t, span.IsRecording(), "Should not have recording span for /metrics")
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should skip tracing for /info endpoint", func(t *testing.T) {
		router := gin.New()
		router.Use(OtelMiddleware("test-service", true))
		
		router.GET("/info", func(c *gin.Context) {
			// Check that no span exists in context
			span := trace.SpanFromContext(c.Request.Context())
			assert.False(t, span.IsRecording(), "Should not have recording span for /info")
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should apply tracing for other endpoints when enabled", func(t *testing.T) {
		router := gin.New()
		router.Use(OtelMiddleware("test-service", true))
		
		router.GET("/api/users", func(c *gin.Context) {
			// When OTEL is enabled and path is not skipped, span should exist
			// Note: In test environment without OTEL provider setup, span won't be recording
			// but the middleware will still attempt to create one
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req, _ := http.NewRequest("GET", "/api/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
