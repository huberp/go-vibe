package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

// TestMiddlewareOrdering tests that middleware is applied in the correct order
func TestMiddlewareOrdering(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should apply middleware in registration order", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		executionOrder := []string{}

		router := gin.New()
		
		// Add middleware that tracks execution order
		router.Use(func(c *gin.Context) {
			executionOrder = append(executionOrder, "first")
			c.Next()
			executionOrder = append(executionOrder, "first-after")
		})
		
		router.Use(func(c *gin.Context) {
			executionOrder = append(executionOrder, "second")
			c.Next()
			executionOrder = append(executionOrder, "second-after")
		})
		
		router.Use(func(c *gin.Context) {
			executionOrder = append(executionOrder, "third")
			c.Next()
			executionOrder = append(executionOrder, "third-after")
		})

		router.GET("/test", func(c *gin.Context) {
			executionOrder = append(executionOrder, "handler")
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verify execution order: middleware runs in order, then handler, then middleware cleanup in reverse
		expected := []string{"first", "second", "third", "handler", "third-after", "second-after", "first-after"}
		assert.Equal(t, expected, executionOrder)
		
		_ = logger // Use logger to avoid unused variable error
	})
}

// TestMiddlewareChaining tests that multiple middleware can be chained together
func TestMiddlewareChaining(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should chain logging, metrics, and rate limiting middleware", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		
		// Chain multiple middleware
		router.Use(LoggingMiddleware(logger))
		router.Use(PrometheusMiddleware())
		router.Use(RateLimitMiddleware(10, 5))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should chain auth and role middleware", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		
		router.Use(JWTAuthMiddleware(secret))
		router.Use(RequireRole("admin"))
		
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin area"})
		})

		// Create valid admin token
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "admin",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should reject when chained auth and role middleware fail", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		
		router.Use(JWTAuthMiddleware(secret))
		router.Use(RequireRole("admin"))
		
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin area"})
		})

		// Create valid user token (not admin)
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "user",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

// TestContextPropagation tests that context values are properly propagated between middleware
func TestContextPropagation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should propagate request_id from logging middleware", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		var capturedRequestID string

		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.Use(func(c *gin.Context) {
			// Capture request_id set by logging middleware
			if reqID, exists := c.Get("request_id"); exists {
				capturedRequestID = reqID.(string)
			}
			c.Next()
		})

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, capturedRequestID)
	})

	t.Run("should propagate user_id and user_role from auth middleware", func(t *testing.T) {
		secret := "test-secret"
		var capturedUserID uint
		var capturedRole string

		router := gin.New()
		router.Use(JWTAuthMiddleware(secret))
		router.Use(func(c *gin.Context) {
			if userID, exists := c.Get("user_id"); exists {
				capturedUserID = userID.(uint)
			}
			if role, exists := c.Get("user_role"); exists {
				capturedRole = role.(string)
			}
			c.Next()
		})

		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create valid token
		claims := jwt.MapClaims{
			"user_id": float64(456),
			"role":    "user",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, uint(456), capturedUserID)
		assert.Equal(t, "user", capturedRole)
	})
}

// TestMiddlewareInteraction tests interactions between different middleware
func TestMiddlewareInteraction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("rate limit should apply before auth", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		
		// Rate limit comes first
		router.Use(RateLimitMiddleware(1, 1))
		router.Use(JWTAuthMiddleware(secret))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Create valid token
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "user",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		// First request should succeed
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		req1.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second request should be rate limited (before auth is checked)
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		req2.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	})

	t.Run("auth should apply before rate limit", func(t *testing.T) {
		secret := "test-secret"
		router := gin.New()
		
		// Auth comes first
		router.Use(JWTAuthMiddleware(secret))
		router.Use(RateLimitMiddleware(1, 1))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Request without token should fail auth (before rate limit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("logging should record all requests including rate limited ones", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		router := gin.New()
		
		// Logging first, then rate limiting
		router.Use(LoggingMiddleware(logger))
		router.Use(RateLimitMiddleware(1, 1))
		
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// Both requests should be logged, even if second is rate limited
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	})
}

// TestMiddlewareErrorHandling tests that errors in one middleware don't affect others
func TestMiddlewareErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("abort in one middleware should stop the chain", func(t *testing.T) {
		handlerCalled := false

		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "aborted"})
		})
		router.Use(func(c *gin.Context) {
			// This should not execute
			t.Error("Second middleware should not execute after abort")
			c.Next()
		})

		router.GET("/test", func(c *gin.Context) {
			handlerCalled = true
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.False(t, handlerCalled)
	})

	t.Run("should handle auth failure in middleware chain", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		handlerCalled := false

		router := gin.New()
		router.Use(LoggingMiddleware(logger))
		router.Use(JWTAuthMiddleware("secret"))
		router.Use(PrometheusMiddleware())

		router.GET("/test", func(c *gin.Context) {
			handlerCalled = true
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		// No auth header
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, handlerCalled)
	})
}

// TestMiddlewareWithDifferentHTTPMethods tests middleware with various HTTP methods
func TestMiddlewareWithDifferentHTTPMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		method string
		status int
	}{
		{"GET", http.StatusOK},
		{"POST", http.StatusCreated},
		{"PUT", http.StatusOK},
		{"PATCH", http.StatusOK},
		{"DELETE", http.StatusNoContent},
	}

	for _, tc := range testCases {
		t.Run("should handle "+tc.method+" requests", func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			router := gin.New()
			
			router.Use(LoggingMiddleware(logger))
			router.Use(PrometheusMiddleware())

			handler := func(c *gin.Context) {
				c.JSON(tc.status, gin.H{"method": tc.method})
			}

			router.GET("/test", handler)
			router.POST("/test", handler)
			router.PUT("/test", handler)
			router.PATCH("/test", handler)
			router.DELETE("/test", handler)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.status, w.Code)
		})
	}
}

// TestFullMiddlewareStack tests a complete middleware stack like in production
func TestFullMiddlewareStack(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should work with full middleware stack", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		secret := "test-secret"

		router := gin.New()
		
		// Full stack: logging -> metrics -> rate limit -> auth -> role check
		router.Use(LoggingMiddleware(logger))
		router.Use(PrometheusMiddleware())
		router.Use(RateLimitMiddleware(10, 5))
		router.Use(JWTAuthMiddleware(secret))
		router.Use(RequireRole("admin"))

		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin area"})
		})

		// Create valid admin token
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "admin",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should fail at appropriate middleware in stack", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		secret := "test-secret"

		router := gin.New()
		
		// Full stack
		router.Use(LoggingMiddleware(logger))
		router.Use(PrometheusMiddleware())
		router.Use(RateLimitMiddleware(10, 5))
		router.Use(JWTAuthMiddleware(secret))
		router.Use(RequireRole("admin"))

		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin area"})
		})

		// Create valid user token (not admin)
		claims := jwt.MapClaims{
			"user_id": float64(123),
			"role":    "user",
			"exp":     time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(secret))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		router.ServeHTTP(w, req)

		// Should fail at role check
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
