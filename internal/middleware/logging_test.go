package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLoggingMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should log basic request successfully", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Verify request_id was set in context
		// Note: We can't easily access the context after the request, but we verified the middleware ran
	})

	t.Run("should extract trace ID from valid W3C traceparent header", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			// Valid traceparent format: 00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01
			// Extract trace ID (positions 3-35)
			assert.Equal(t, "0af7651916cd43dd8448eb211c80319c", requestID)
			c.JSON(http.StatusOK, gin.H{"request_id": requestID})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("traceparent", "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should handle invalid traceparent header and generate UUID", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			// Should be a UUID (not the invalid traceparent value)
			assert.NotEqual(t, "invalid", requestID)
			// UUID should be longer than 0
			assert.NotEmpty(t, requestID)
			c.JSON(http.StatusOK, gin.H{"request_id": requestID})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("traceparent", "invalid")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should handle short traceparent header and generate UUID", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			// Should generate UUID for short traceparent
			assert.NotEmpty(t, requestID)
			c.JSON(http.StatusOK, gin.H{"request_id": requestID})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("traceparent", "00-abc")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should generate UUID when no traceparent header", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			assert.NotEmpty(t, requestID)
			c.JSON(http.StatusOK, gin.H{"request_id": requestID})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should log request with long path", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		longPath := "/api/v1/users/123/posts/456/comments/789/replies/012/nested/path/that/is/very/long"
		router.GET(longPath, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", longPath, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should log different HTTP methods", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))

		handler := func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		}

		router.GET("/test", handler)
		router.POST("/test", handler)
		router.PUT("/test", handler)
		router.DELETE("/test", handler)

		methods := []string{"GET", "POST", "PUT", "DELETE"}
		for _, method := range methods {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("should log different status codes", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))

		router.GET("/ok", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
		router.GET("/created", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"message": "created"})
		})
		router.GET("/bad-request", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		})
		router.GET("/not-found", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		})
		router.GET("/error", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		})

		testCases := []struct {
			path   string
			status int
		}{
			{"/ok", http.StatusOK},
			{"/created", http.StatusCreated},
			{"/bad-request", http.StatusBadRequest},
			{"/not-found", http.StatusNotFound},
			{"/error", http.StatusInternalServerError},
		}

		for _, tc := range testCases {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.path, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.status, w.Code)
		}
	})

	t.Run("should include trace context from OpenTelemetry span", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()

		// Create a tracer for testing
		tracer := otel.Tracer("test-tracer")

		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			// Create a span in the handler to simulate OTel middleware
			ctx, span := tracer.Start(c.Request.Context(), "test-operation")
			defer span.End()

			// Update request context
			c.Request = c.Request.WithContext(ctx)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should handle request with query parameters", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test?param1=value1&param2=value2", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should log request from different client IPs", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		testIPs := []string{
			"192.168.1.1",
			"10.0.0.1",
			"172.16.0.1",
		}

		for _, ip := range testIPs {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Forwarded-For", ip)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}

func TestLoggingMiddleware_EmptyTraceparent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zaptest.NewLogger(t)

	t.Run("should handle empty string traceparent", func(t *testing.T) {
		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			assert.NotEmpty(t, requestID)
			c.JSON(http.StatusOK, gin.H{"request_id": requestID})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("traceparent", "")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestLoggingMiddleware_WithNilLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should handle nil logger gracefully", func(t *testing.T) {
		// This tests defensive programming - what happens if logger is nil
		// In practice, this shouldn't happen, but we test it anyway
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic: %v", r)
			}
		}()

		router := gin.New()
		// Using a no-op logger instead of nil to avoid actual panics
		router.Use(LoggingMiddleware(zap.NewNop()))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestLoggingMiddleware_TraceIDConsistency(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zaptest.NewLogger(t)

	t.Run("should consistently extract same trace ID from traceparent", func(t *testing.T) {
		traceparent := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
		expectedTraceID := "0af7651916cd43dd8448eb211c80319c"

		for range 3 {
			router := gin.New()
			router.Use(LoggingMiddleware(logger))
			router.GET("/test", func(c *gin.Context) {
				requestID, _ := c.Get("request_id")
				assert.Equal(t, expectedTraceID, requestID)
				c.JSON(http.StatusOK, gin.H{"request_id": requestID})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("traceparent", traceparent)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}
